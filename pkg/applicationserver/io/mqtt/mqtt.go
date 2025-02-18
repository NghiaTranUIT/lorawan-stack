// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package mqtt implements the MQTT frontend.
package mqtt

import (
	"context"
	"fmt"
	stdio "io"
	"net"

	"github.com/TheThingsIndustries/mystique/pkg/auth"
	mqttlog "github.com/TheThingsIndustries/mystique/pkg/log"
	mqttnet "github.com/TheThingsIndustries/mystique/pkg/net"
	"github.com/TheThingsIndustries/mystique/pkg/packet"
	"github.com/TheThingsIndustries/mystique/pkg/session"
	"github.com/TheThingsIndustries/mystique/pkg/topic"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	ttsauth "go.thethings.network/lorawan-stack/v3/pkg/auth"
	"go.thethings.network/lorawan-stack/v3/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/v3/pkg/errorcontext"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/mqtt"
	"go.thethings.network/lorawan-stack/v3/pkg/ratelimit"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"google.golang.org/grpc/metadata"
)

const qosUpstream byte = 0

type srv struct {
	ctx    context.Context
	server io.Server
	format Format
	lis    mqttnet.Listener
}

// Serve serves the MQTT frontend.
func Serve(ctx context.Context, server io.Server, listener net.Listener, format Format, protocol string) error {
	ctx = log.NewContextWithField(ctx, "namespace", "applicationserver/io/mqtt")
	ctx = mqttlog.NewContext(ctx, mqtt.Logger(log.FromContext(ctx)))
	s := &srv{ctx, server, format, mqttnet.NewListener(listener, protocol)}
	go func() {
		<-ctx.Done()
		s.lis.Close()
	}()
	return s.accept()
}

func (s *srv) accept() error {
	for {
		mqttConn, err := s.lis.Accept()
		if err != nil {
			if s.ctx.Err() == nil {
				log.FromContext(s.ctx).WithError(err).Warn("Accept failed")
			}
			return err
		}

		remoteAddr := mqttConn.RemoteAddr().String()
		ctx := log.NewContextWithFields(s.ctx, log.Fields("remote_addr", remoteAddr))

		resource := ratelimit.ApplicationAcceptMQTTConnectionResource(remoteAddr)
		if err := ratelimit.Require(s.server.RateLimiter(), resource); err != nil {
			if err := mqttConn.Close(); err != nil {
				log.FromContext(ctx).WithError(err).Warn("Close connection failed")
			}
			log.FromContext(ctx).WithError(err).Debug("Drop connection")
			continue
		}

		go func() {
			conn := &connection{server: s.server, mqtt: mqttConn, format: s.format}
			if err := conn.setup(ctx); err != nil {
				switch err {
				case stdio.EOF, stdio.ErrUnexpectedEOF:
				default:
					log.FromContext(ctx).WithError(err).Warn("Failed to setup connection")
				}
				mqttConn.Close()
				return
			}
		}()
	}
}

type connection struct {
	format  Format
	server  io.Server
	mqtt    mqttnet.Conn
	session session.Session
	io      *io.Subscription

	resource ratelimit.Resource
}

func (c *connection) setup(ctx context.Context) error {
	ctx = auth.NewContextWithInterface(ctx, c)
	ctx, cancel := errorcontext.New(ctx)
	c.session = session.New(ctx, c.mqtt, c.deliver)
	if err := c.session.ReadConnect(); err != nil {
		cancel(err)
		return err
	}
	ctx = c.io.Context()

	logger := log.FromContext(ctx)
	controlCh := make(chan packet.ControlPacket)

	// Read control packets
	go func() {
		for {
			pkt, err := c.session.ReadPacket()
			if err != nil {
				if err != stdio.EOF {
					logger.WithError(err).Warn("Error when reading packet")
				}
				cancel(err)
				return
			}
			if pkt != nil {
				logger.Debugf("Schedule %s packet", packet.Name[pkt.PacketType()])
				select {
				case <-ctx.Done():
					return
				case controlCh <- pkt:
				}
			}
		}
	}()

	// Publish upstream
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case up := <-c.io.Up():
				logger := logger.WithField("device_uid", unique.ID(up.Context, up.EndDeviceIdentifiers))
				var topicParts []string
				switch up.Up.(type) {
				case *ttnpb.ApplicationUp_UplinkMessage:
					topicParts = c.format.UplinkTopic(unique.ID(up.Context, c.io.ApplicationIDs()), up.DeviceId)
				case *ttnpb.ApplicationUp_JoinAccept:
					topicParts = c.format.JoinAcceptTopic(unique.ID(up.Context, c.io.ApplicationIDs()), up.DeviceId)
				case *ttnpb.ApplicationUp_DownlinkAck:
					topicParts = c.format.DownlinkAckTopic(unique.ID(up.Context, c.io.ApplicationIDs()), up.DeviceId)
				case *ttnpb.ApplicationUp_DownlinkNack:
					topicParts = c.format.DownlinkNackTopic(unique.ID(up.Context, c.io.ApplicationIDs()), up.DeviceId)
				case *ttnpb.ApplicationUp_DownlinkSent:
					topicParts = c.format.DownlinkSentTopic(unique.ID(up.Context, c.io.ApplicationIDs()), up.DeviceId)
				case *ttnpb.ApplicationUp_DownlinkFailed:
					topicParts = c.format.DownlinkFailedTopic(unique.ID(up.Context, c.io.ApplicationIDs()), up.DeviceId)
				case *ttnpb.ApplicationUp_DownlinkQueued:
					topicParts = c.format.DownlinkQueuedTopic(unique.ID(up.Context, c.io.ApplicationIDs()), up.DeviceId)
				case *ttnpb.ApplicationUp_DownlinkQueueInvalidated:
					topicParts = c.format.DownlinkQueueInvalidatedTopic(unique.ID(up.Context, c.io.ApplicationIDs()), up.DeviceId)
				case *ttnpb.ApplicationUp_LocationSolved:
					topicParts = c.format.LocationSolvedTopic(unique.ID(up.Context, c.io.ApplicationIDs()), up.DeviceId)
				case *ttnpb.ApplicationUp_ServiceData:
					topicParts = c.format.ServiceDataTopic(unique.ID(up.Context, c.io.ApplicationIDs()), up.DeviceId)
				}
				if topicParts == nil {
					continue
				}
				buf, err := c.format.FromUp(up.ApplicationUp)
				if err != nil {
					logger.WithError(err).Warn("Failed to marshal upstream message")
					continue
				}
				logger.Debug("Publish upstream message")
				c.session.Publish(&packet.PublishPacket{
					TopicName:  topic.Join(topicParts),
					TopicParts: topicParts,
					QoS:        qosUpstream,
					Message:    buf,
				})
			}
		}
	}()

	// Write packets
	go func() {
		for {
			var err error
			select {
			case <-ctx.Done():
				return
			case pkt := <-controlCh:
				err = c.mqtt.Send(pkt)
			case pkt, ok := <-c.session.PublishChan():
				if !ok {
					return
				}
				logger.Debug("Write publish packet")
				err = c.mqtt.Send(pkt)
			}
			if err != nil {
				cancel(err)
				return
			}
		}
	}()

	// Close connection on context closure
	go func() {
		select {
		case <-ctx.Done():
			logger.WithError(ctx.Err()).Info("Disconnected")
			c.session.Close()
			c.mqtt.Close()
		}
	}()

	logger.Info("Connected")
	return nil
}

type topicAccess struct {
	appUID string
	reads  [][]string
	writes [][]string
}

func (c *connection) Connect(ctx context.Context, info *auth.Info) (_ context.Context, err error) {
	ids := ttnpb.ApplicationIdentifiers{
		ApplicationId: info.Username,
	}
	if err := ids.ValidateContext(ctx); err != nil {
		return nil, err
	}

	md := metadata.New(map[string]string{
		"id":            ids.ApplicationId,
		"authorization": fmt.Sprintf("Bearer %s", info.Password),
	})
	if ctxMd, ok := metadata.FromIncomingContext(ctx); ok {
		md = metadata.Join(ctxMd, md)
	}
	ctx = metadata.NewIncomingContext(ctx, md)

	ctx = c.server.FillContext(ctx)
	uid := unique.ID(ctx, ids)
	ctx = log.NewContextWithField(ctx, "application_uid", uid)

	defer func() {
		if err != nil {
			registerConnectFail(ctx, ids, err)
		}
		switch {
		case errors.IsPermissionDenied(err):
			err = packet.ConnectNotAuthorized
		case errors.IsResourceExhausted(err):
			err = packet.ConnectServerUnavailable
		}
	}()

	if err := rights.RequireApplication(ctx, ids); err != nil {
		return nil, err
	}

	c.io, err = c.server.Subscribe(ctx, "mqtt", &ids, true)
	if err != nil {
		return nil, err
	}
	ctx = c.io.Context()

	authTokenID := ""
	if _, v, _, err := ttsauth.SplitToken(string(info.Password)); err == nil && v != "" {
		authTokenID = v
	}
	c.resource = ratelimit.ApplicationMQTTDownResource(ctx, ids, authTokenID)

	access := topicAccess{
		appUID: uid,
	}
	if err := rights.RequireApplication(ctx, ids, ttnpb.RIGHT_APPLICATION_TRAFFIC_READ); err == nil {
		access.reads = append(access.reads,
			c.format.UplinkTopic(uid, topic.PartWildcard),
			c.format.JoinAcceptTopic(uid, topic.PartWildcard),
			c.format.DownlinkAckTopic(uid, topic.PartWildcard),
			c.format.DownlinkNackTopic(uid, topic.PartWildcard),
			c.format.DownlinkSentTopic(uid, topic.PartWildcard),
			c.format.DownlinkFailedTopic(uid, topic.PartWildcard),
			c.format.DownlinkQueuedTopic(uid, topic.PartWildcard),
			c.format.DownlinkQueueInvalidatedTopic(uid, topic.PartWildcard),
			c.format.LocationSolvedTopic(uid, topic.PartWildcard),
			c.format.ServiceDataTopic(uid, topic.PartWildcard),
		)
	}
	if err := rights.RequireApplication(ctx, ids, ttnpb.RIGHT_APPLICATION_TRAFFIC_DOWN_WRITE); err == nil {
		access.writes = append(access.writes,
			c.format.DownlinkPushTopic(uid, topic.PartWildcard),
			c.format.DownlinkReplaceTopic(uid, topic.PartWildcard),
		)
	}
	info.Metadata = access
	info.Interface = c
	return ctx, nil
}

var errNotAuthorized = errors.DefinePermissionDenied("not_authorized", "not authorized")

func (c *connection) Subscribe(info *auth.Info, requestedTopic string, requestedQoS byte) (acceptedTopic string, acceptedQoS byte, err error) {
	access := info.Metadata.(topicAccess)
	accepted, ok := c.format.AcceptedTopic(access.appUID, topic.Split(requestedTopic))
	if !ok {
		return "", 0, errNotAuthorized.New()
	}
	return topic.Join(accepted), requestedQoS, nil
}

func (c *connection) CanRead(info *auth.Info, topicParts ...string) bool {
	access := info.Metadata.(topicAccess)
	for _, reads := range access.reads {
		if topic.MatchPath(topicParts, reads) {
			return true
		}
	}
	return false
}

func (c *connection) CanWrite(info *auth.Info, topicParts ...string) bool {
	access := info.Metadata.(topicAccess)
	for _, writes := range access.writes {
		if topic.MatchPath(topicParts, writes) {
			return true
		}
	}
	return false
}

func (c *connection) deliver(pkt *packet.PublishPacket) {
	logger := log.FromContext(c.io.Context()).WithField("topic", pkt.TopicName)

	if err := ratelimit.Require(c.server.RateLimiter(), c.resource); err != nil {
		logger.WithError(err).Warn("Terminate connection")
		c.io.Disconnect(err)
		return
	}

	var deviceID string
	var op func(io.Server, context.Context, ttnpb.EndDeviceIdentifiers, []*ttnpb.ApplicationDownlink) error
	switch {
	case c.format.IsDownlinkPushTopic(pkt.TopicParts):
		deviceID = c.format.ParseDownlinkPushTopic(pkt.TopicParts)
		op = io.Server.DownlinkQueuePush
	case c.format.IsDownlinkReplaceTopic(pkt.TopicParts):
		deviceID = c.format.ParseDownlinkReplaceTopic(pkt.TopicParts)
		op = io.Server.DownlinkQueueReplace
	default:
		logger.Error("Invalid topic path")
		return
	}
	items, err := c.format.ToDownlinks(pkt.Message)
	if err != nil {
		logger.WithError(err).Warn("Failed to decode downlink messages")
		return
	}
	ids := ttnpb.EndDeviceIdentifiers{
		ApplicationIdentifiers: *c.io.ApplicationIDs(),
		DeviceId:               deviceID,
	}
	if err := ids.ValidateContext(c.io.Context()); err != nil {
		logger.WithError(err).Warn("Failed to validate message identifiers")
		return
	}
	logger.WithFields(log.Fields(
		"device_uid", unique.ID(c.io.Context(), ids),
		"count", len(items.Downlinks),
	)).Debug("Handle downlink messages")
	if err := op(c.server, c.io.Context(), ids, items.Downlinks); err != nil {
		logger.WithError(err).Warn("Failed to handle downlink messages")
	}
}

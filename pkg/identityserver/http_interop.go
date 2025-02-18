// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
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

package identityserver

import (
	"context"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/v3/pkg/interop"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
)

type interopServer struct {
	*IdentityServer
	interop.Authorizer
}

func (srv *interopServer) HomeNSRequest(ctx context.Context, in *interop.HomeNSReq) (*interop.TTIHomeNSAns, error) {
	ctx = log.NewContextWithField(ctx, "namespace", "identityserver/interop")
	if err := srv.RequireAuthorized(ctx); err != nil {
		return nil, err
	}

	ids := &ttnpb.EndDeviceIdentifiers{
		JoinEui: (*types.EUI64)(&in.ReceiverID),
		DevEui:  (*types.EUI64)(&in.DevEUI),
	}

	var dev *ttnpb.EndDevice
	err := srv.withDatabase(ctx, func(db *gorm.DB) (err error) {
		dev, err = store.GetEndDeviceStore(db).GetEndDevice(ctx, ids, &pbtypes.FieldMask{
			Paths: []string{"network_server_address"},
		})
		return err
	})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, interop.ErrUnknownDevEUI.WithCause(err)
		}
		return nil, err
	}

	var (
		conf  = srv.configFromContext(ctx)
		hNSID *types.EUI64
	)

	header, err := in.AnswerHeader()
	if err != nil {
		return nil, interop.ErrMalformedMessage.WithCause(err)
	}
	ans := &interop.TTIHomeNSAns{
		HomeNSAns: interop.HomeNSAns{
			JsNsMessageHeader: interop.JsNsMessageHeader{
				MessageHeader: header,
				SenderID:      in.ReceiverID,
				ReceiverID:    in.SenderID,
				ReceiverNSID:  in.SenderNSID,
			},
			Result: interop.Result{
				ResultCode: interop.ResultSuccess,
			},
			HNetID: interop.NetID(conf.Network.NetID),
		},
		HNSAddress: dev.NetworkServerAddress,
		HTenantID:  conf.Network.TenantID,
	}
	if in.ProtocolVersion.SupportsNSID() {
		ans.HNSID = (*interop.EUI64)(hNSID)
	}
	return ans, nil
}

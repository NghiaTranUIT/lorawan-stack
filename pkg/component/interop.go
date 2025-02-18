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

package component

import (
	"crypto/tls"
	"io/ioutil"
	stdlog "log"
	"net"
	"net/http"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/interop"
)

// RegisterInterop registers an interop subsystem to the component.
func (c *Component) RegisterInterop(s interop.Registerer) {
	c.interopSubsystems = append(c.interopSubsystems, s)
}

func (c *Component) serveInterop(lis net.Listener) error {
	srv := http.Server{
		Handler:           c.interop,
		ReadTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ErrorLog:          stdlog.New(ioutil.Discard, "", 0),
	}
	go func() {
		<-c.Context().Done()
		srv.Close()
	}()
	return srv.Serve(lis)
}

func (c *Component) interopEndpoints() []Endpoint {
	return []Endpoint{
		// TODO: Enable TCP endpoint (https://github.com/TheThingsNetwork/lorawan-stack/issues/717)
		NewTLSEndpoint(c.config.Interop.ListenTLS, "Interop",
			WithTLSClientAuth(tls.VerifyClientCertIfGiven, c.interop.ClientCAPool(), nil),
			WithNextProtos("h2", "http/1.1"),
		),
	}
}

func (c *Component) listenInterop() error {
	return c.serveOnEndpoints(c.interopEndpoints(), (*Component).serveInterop, "interop")
}

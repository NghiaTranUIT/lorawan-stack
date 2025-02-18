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

// Package rpcclient contains the default options for TTN gRPC clients.
package rpcclient

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"go.opencensus.io/plugin/ocgrpc"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/metrics"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcmiddleware"
	_ "go.thethings.network/lorawan-stack/v3/pkg/rpcmiddleware/discover" // Register service discovery resolvers
	"go.thethings.network/lorawan-stack/v3/pkg/rpcmiddleware/rpclog"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcmiddleware/warning"
	"go.thethings.network/lorawan-stack/v3/pkg/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// DefaultDialOptions for gRPC clients
func DefaultDialOptions(ctx context.Context) []grpc.DialOption {
	streamInterceptors := []grpc.StreamClientInterceptor{
		metrics.StreamClientInterceptor,
		grpc_opentracing.StreamClientInterceptor(),
		rpclog.StreamClientInterceptor(ctx), // Gets logger from global context
		warning.StreamClientInterceptor,
		errors.StreamClientInterceptor(),
	}

	unaryInterceptors := []grpc.UnaryClientInterceptor{
		metrics.UnaryClientInterceptor,
		grpc_opentracing.UnaryClientInterceptor(),
		rpclog.UnaryClientInterceptor(ctx), // Gets logger from global context
		warning.UnaryClientInterceptor,
		errors.UnaryClientInterceptor(),
	}

	return []grpc.DialOption{
		grpc.WithStatsHandler(rpcmiddleware.StatsHandlers{new(ocgrpc.ClientHandler), metrics.StatsHandler}),
		grpc.WithUserAgent(fmt.Sprintf(
			"%s go/%s ttn/%s",
			filepath.Base(os.Args[0]),
			strings.TrimPrefix(runtime.Version(), "go"),
			version.String(),
		)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamInterceptors...)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInterceptors...)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5 * time.Minute,
			Timeout:             20 * time.Second,
			PermitWithoutStream: false,
		}),
	}
}

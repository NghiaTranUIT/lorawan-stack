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

package mqtt

import (
	"context"

	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

var evtConnectFail = events.Define(
	"as.mqtt.connect.fail", "fail to connect to MQTT",
	events.WithVisibility(ttnpb.RIGHT_APPLICATION_TRAFFIC_READ),
	events.WithErrorDataType(),
)

func registerConnectFail(ctx context.Context, ids ttnpb.ApplicationIdentifiers, err error) {
	events.Publish(evtConnectFail.NewWithIdentifiersAndData(ctx, &ids, err))
}

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

import { connect } from 'react-redux'

import withRequest from '@ttn-lw/lib/components/with-request'

import { selectJsConfig } from '@ttn-lw/lib/selectors/env'

import { getJoinEUIPrefixes } from '@console/store/actions/join-server'

import { selectJoinEUIPrefixesFetching } from '@console/store/selectors/join-server'

const mapStateToProps = state => {
  const { enabled: jsEnabled } = selectJsConfig()

  return {
    fetching: selectJoinEUIPrefixesFetching(state) && jsEnabled,
    jsEnabled,
  }
}
const mapDispatchToProps = { getPrefixes: getJoinEUIPrefixes }

export default DeviceAdd =>
  connect(
    mapStateToProps,
    mapDispatchToProps,
  )(
    withRequest(
      ({ getPrefixes, jsEnabled }) => {
        if (jsEnabled) {
          return getPrefixes()
        }
      },
      ({ fetching }) => fetching,
    )(DeviceAdd),
  )

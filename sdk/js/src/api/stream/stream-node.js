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

import axios from 'axios'

import Token from '../../util/token'

import { notify, EVENTS } from './shared'

/**
 * Opens a new stream.
 *
 * @async
 * @param {object} payload -  - The body of the initial request.
 * @param {string} url - The stream endpoint.
 *
 * @example
 * (async () => {
 *    const stream = await stream(
 *      { identifiers: [{ application_ids: { application_id: 'my-app' }}]},
 *      'http://localhost:1885/api/v3/events',
 *    )
 *
 *    // Add listeners to the stream.
 *    stream
 *      .on('start', () => console.log('conn opened'))
 *      .on('chunk', chunk => console.log('received chunk', chunk))
 *      .on('error', error => console.log(error))
 *      .on('close', wasClientRequest => console.log(wasClientRequest ? 'conn closed by client' : 'conn closed by server'))
 *
 *    // Start the stream after attaching the listeners.
 *    stream.open()
 *
 *    // Close the stream after 20 s.
 *    setTimeout(() => stream.close(), 20000)
 * })()
 *
 * @returns {object} The stream subscription object with the `on` function for
 * attaching listeners and the `close` function to close the stream.
 */
export default async (payload, url) => {
  let listeners = Object.values(EVENTS).reduce((acc, curr) => ({ ...acc, [curr]: null }), {})
  let reader = null
  let closeRequested = false

  const token = new Token().get()

  let Authorization = null
  if (typeof token === 'function') {
    Authorization = `Bearer ${(await token()).access_token}`
  } else {
    Authorization = `Bearer ${token}`
  }

  let buffer = ''
  axios({
    url,
    data: JSON.stringify(payload),
    method: 'POST',
    responseType: 'stream',
    headers: {
      Authorization,
      Accept: 'text/event-stream',
    },
  })
    .then(response => response.data)
    .then(stream => {
      reader = stream
    })

  return {
    open: () => {
      notify(listeners[EVENTS.START])

      reader.on('data', data => {
        const parsed = data.toString('utf8')
        buffer += parsed
        const lines = buffer.split(/\n\n/)
        buffer = lines.pop()
        for (const line of lines) {
          notify(listeners[EVENTS.CHUNK], JSON.parse(line).result)
        }
      })
      reader.on('end', () => {
        notify(listeners[EVENTS.CLOSE], closeRequested)
        listeners = {}
      })
      reader.on('error', error => {
        notify(listeners[EVENTS.ERROR], error)
        listeners = {}
      })
    },
    on(eventName, callback) {
      if (listeners[eventName] === undefined) {
        throw new Error(
          `${eventName} event is not supported. Should be one of: start, error, chunk or close`,
        )
      }

      listeners[eventName] = callback

      return this
    },
    close: () => {
      if (reader) {
        closeRequested = true
        reader.cancel()
      }
    },
  }
}

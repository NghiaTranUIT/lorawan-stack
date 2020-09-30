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

const selectEventsStore = (state, entityId) => state[entityId]

export const createEventsSelector = entity =>
  function(state, entityId) {
    const store = selectEventsStore(state.events[entity], entityId)

    return store ? store.events : []
  }

export const createEventsStatusSelector = entity =>
  function(state, entityId) {
    const store = selectEventsStore(state.events[entity], entityId)

    return store ? store.status : 'unknown'
  }

export const createEventsPausedSelector = entity =>
  function(state, entityId) {
    const store = selectEventsStore(state.events[entity], entityId)

    return Boolean(store.paused)
  }

export const createEventsInterruptedSelector = entity =>
  function(state, entityId) {
    const store = selectEventsStore(state.events[entity], entityId)

    return Boolean(store.interrupted)
  }

export const createEventsErrorSelector = entity =>
  function(state, entityId) {
    const store = selectEventsStore(state.events[entity], entityId)

    return store ? store.error : undefined
  }

export const createEventsTruncatedSelector = entity =>
  function(state, entityId) {
    const store = selectEventsStore(state.events[entity], entityId)

    return Boolean(store.truncated)
  }

export const createLatestEventSelector = function(entity) {
  const eventsSelector = createEventsSelector(entity)

  return function selectLatestEvent(state, entityId) {
    const events = eventsSelector(state, entityId)

    return events[0]
  }
}

export const createInterruptedStreamsSelector = entity => state => {
  const eventsStore = state.events[entity]

  return Object.keys(eventsStore).reduce((acc, id) => {
    if (eventsStore[id].interrupted) {
      acc[id] = eventsStore[id]
    }

    return acc
  }, {})
}

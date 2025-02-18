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

/* eslint-disable no-invalid-this */

import traverse from 'traverse'

/**
 * Takes registry responses from different components and merges them into a
 * single entity record.
 *
 * @param {object} parts - An object containing the device record responded from
 * the registry and the paths that were requested from the component.
 * Shape: { device: …, paths: … }.
 * @param {string} base - An optional base device record, that the merge will
 * take as base.
 * @param {object} minimum - Paths that will always be merged for all records.
 * @returns {object} The merged device record.
 */
export default (parts, base = {}, minimum = [['ids'], ['created_at'], ['updated_at']]) => {
  const result = base

  // Cycle through all responses.
  for (const part of parts) {
    for (const path of part.paths ? [...minimum, ...part.paths] : []) {
      // For each path requested, get the corresponding value of the device
      // record.
      const val = traverse(part.device).get(path)

      // Consider also falsy boolean values, for example.
      const isBoolean = typeof val === 'boolean'
      if (val || isBoolean) {
        if (typeof val === 'object') {
          // In case of a whole sub-object being selected, write each leaf node
          // explicitly to achieve a deep merge instead of whole object
          // overrides.
          if (Object.keys(val).length === 0) {
            const currentValue = traverse(result).get(path)

            // Set empty object values only if value is unset, otherwise they might
            // override legitimate values.
            if (typeof currentValue === 'undefined') {
              traverse(result).set(path, val)
            }

            continue
          }

          traverse(val).forEach(function (e) {
            if (Array.isArray(e) && e.length > 0) {
              traverse(result).set(path, val)

              return
            }

            if (this.isLeaf) {
              // TODO: Remove this once https://github.com/TheThingsNetwork/lorawan-stack/issues/4766 is resolved.
              if (this.key.endsWith('_at') && e === '0001-01-01T00:00:00Z') {
                return
              }
              // Write the sub object leaf into the result.
              traverse(result).set([...path, ...this.path], e)
            }
          })
        } else {
          // In case of a simple leaf, just write it into the result.
          traverse(result).set(path, val)
        }
      }
    }
  }

  return result
}

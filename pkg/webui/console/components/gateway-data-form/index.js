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

import React from 'react'
import bind from 'autobind-decorator'

import delay from '@console/constants/delays'
import frequencyPlans from '@console/constants/frequency-plans'

import Form from '@ttn-lw/components/form'
import Input from '@ttn-lw/components/input'
import Checkbox from '@ttn-lw/components/checkbox'
import SubmitBar from '@ttn-lw/components/submit-bar'
import UnitInput from '@ttn-lw/components/unit-input'
import KeyValueMap from '@ttn-lw/components/key-value-map'

import Message from '@ttn-lw/lib/components/message'

import { GsFrequencyPlansSelect } from '@console/containers/freq-plans-select'
import OwnersSelect from '@console/containers/owners-select'

import tooltipIds from '@ttn-lw/lib/constants/tooltip-ids'
import PropTypes from '@ttn-lw/lib/prop-types'
import sharedMessages from '@ttn-lw/lib/shared-messages'

import { unit as unitRegexp, emptyDuration as emptyDurationRegexp } from '@console/lib/regexp'

import validationSchema from './validation-schema'

const isEmptyFrequencyPlan = value => value === frequencyPlans.EMPTY_FREQ_PLAN

class GatewayDataForm extends React.Component {
  static propTypes = {
    /** The SubmitBar content. */
    children: PropTypes.node.isRequired,
    error: PropTypes.error,
    gsEnabled: PropTypes.bool.isRequired,
    initialValues: PropTypes.gateway,
    onSubmit: PropTypes.func.isRequired,
  }

  static defaultProps = {
    error: '',
    initialValues: validationSchema.cast({}),
  }

  constructor(props) {
    super(props)

    this.state = {
      shouldDisplayWarning: this.isNotValidDuration(props.initialValues.schedule_anytime_delay),
      showFrequencyPlanWarning: isEmptyFrequencyPlan(props.initialValues.frequency_plan_id),
    }
  }

  @bind
  onScheduleAnytimeDelayChange(value) {
    this.setState({ shouldDisplayWarning: this.isNotValidDuration(value) })
  }

  @bind
  handleFrequencyPlanChange(freqPlan) {
    this.setState({ showFrequencyPlanWarning: isEmptyFrequencyPlan(freqPlan.value) })
  }

  @bind
  onSubmit(values, helpers) {
    const { onSubmit, gsEnabled } = this.props

    const castedValues = validationSchema.cast(
      isEmptyFrequencyPlan(values.frequency_plan_id)
        ? { ...values, frequency_plan_id: '' }
        : values,
      { context: { gsEnabled } },
    )

    onSubmit(castedValues, helpers)
  }

  decodeDelayValue(value) {
    if (emptyDurationRegexp.test(value)) {
      return {
        duration: undefined,
        unit: value,
      }
    }
    const duration = value.split(unitRegexp)[0]
    const unit = value.split(duration)[1]
    return {
      duration: duration ? Number(duration) : undefined,
      unit,
    }
  }

  isNotValidDuration(value) {
    const { duration, unit } = this.decodeDelayValue(value)
    switch (unit) {
      case 'ms':
        return duration < delay.MINIMUM_GATEWAY_SCHEDULE_ANYTIME_DELAY
      case 's':
        return duration < delay.MINIMUM_GATEWAY_SCHEDULE_ANYTIME_DELAY / 1000
      case 'm':
        return duration < delay.MINIMUM_GATEWAY_SCHEDULE_ANYTIME_DELAY / 60000
      case 'h':
        return duration < delay.MINIMUM_GATEWAY_SCHEDULE_ANYTIME_DELAY / 3600000
    }
  }

  render() {
    const { error, initialValues, children, gsEnabled } = this.props
    const { shouldDisplayWarning, showFrequencyPlanWarning } = this.state

    return (
      <Form
        error={error}
        onSubmit={this.onSubmit}
        initialValues={initialValues}
        validationSchema={validationSchema}
        validationContext={{ gsEnabled }}
      >
        <Form.SubTitle title={sharedMessages.generalSettings} />
        <OwnersSelect name="owner_id" required autoFocus />
        <Form.Field
          title={sharedMessages.gatewayID}
          name="ids.gateway_id"
          placeholder={sharedMessages.gatewayIdPlaceholder}
          required
          component={Input}
          tooltipId={tooltipIds.GATEWAY_ID}
        />
        <Form.Field
          title={sharedMessages.gatewayEUI}
          name="ids.eui"
          type="byte"
          min={8}
          max={8}
          placeholder={sharedMessages.gatewayEUI}
          component={Input}
          tooltipId={tooltipIds.GATEWAY_EUI}
        />
        <Form.Field
          title={sharedMessages.gatewayName}
          placeholder={sharedMessages.gatewayNamePlaceholder}
          name="name"
          component={Input}
          tooltipId={tooltipIds.GATEWAY_NAME}
        />
        <Form.Field
          title={sharedMessages.gatewayDescription}
          description={sharedMessages.gatewayDescDescription}
          placeholder={sharedMessages.gatewayDescPlaceholder}
          name="description"
          type="textarea"
          component={Input}
          tooltipId={tooltipIds.GATEWAY_DESCRIPTION}
        />
        <Form.Field
          title={sharedMessages.gatewayServerAddress}
          description={sharedMessages.gsServerAddressDescription}
          placeholder={sharedMessages.addressPlaceholder}
          name="gateway_server_address"
          component={Input}
        />
        <Form.Field
          title={sharedMessages.requireAuthenticatedConnection}
          name="require_authenticated_connection"
          component={Checkbox}
          label={sharedMessages.enabled}
          description={sharedMessages.requireAuthenticatedConnectionDescription}
          tooltipId={tooltipIds.REQUIRE_AUTHENTICATED_CONNECTION}
        />
        <Form.Field
          title={sharedMessages.gatewayStatus}
          name="status_public"
          component={Checkbox}
          label={sharedMessages.public}
          description={sharedMessages.statusDescription}
          tooltipId={tooltipIds.GATEWAY_STATUS}
        />
        <Form.Field
          title={sharedMessages.gatewayLocation}
          name="location_public"
          component={Checkbox}
          label={sharedMessages.public}
          description={sharedMessages.locationDescription}
          tooltipId={tooltipIds.GATEWAY_LOCATION}
        />
        <Form.Field
          name="attributes"
          title={sharedMessages.attributes}
          keyPlaceholder={sharedMessages.key}
          valuePlaceholder={sharedMessages.value}
          addMessage={sharedMessages.addAttributes}
          component={KeyValueMap}
          description={sharedMessages.attributeDescription}
          tooltipId={tooltipIds.GATEWAY_ATTRIBUTES}
        />
        <Message component="h4" content={sharedMessages.lorawanOptions} />
        {gsEnabled && (
          <GsFrequencyPlansSelect
            name="frequency_plan_id"
            menuPlacement="top"
            onChange={this.handleFrequencyPlanChange}
            tooltipId={tooltipIds.FREQUENCY_PLAN}
            warning={showFrequencyPlanWarning ? sharedMessages.frequencyPlanWarning : undefined}
            required
          />
        )}
        <Form.Field
          title={sharedMessages.gatewayScheduleDownlinkLate}
          name="schedule_downlink_late"
          component={Checkbox}
          description={sharedMessages.scheduleDownlinkLateDescription}
          tooltipId={tooltipIds.SCHEDULE_DOWNLINK_LATE}
        />
        <Form.Field
          title={sharedMessages.enforceDutyCycle}
          name="enforce_duty_cycle"
          component={Checkbox}
          label={sharedMessages.enabled}
          description={sharedMessages.enforceDutyCycleDescription}
          tooltipId={tooltipIds.ENFORCE_DUTY_CYCLE}
        />
        <Form.Field
          title={sharedMessages.scheduleAnyTimeDelay}
          name="schedule_anytime_delay"
          component={UnitInput}
          inputWidth="s"
          description={{
            ...sharedMessages.scheduleAnyTimeDescription,
            values: {
              minimumValue: delay.MINIMUM_GATEWAY_SCHEDULE_ANYTIME_DELAY,
              defaultValue: delay.DEFAULT_GATEWAY_SCHEDULE_ANYTIME_DELAY,
            },
          }}
          units={[
            { label: sharedMessages.milliseconds, value: 'ms' },
            { label: sharedMessages.seconds, value: 's' },
            { label: sharedMessages.minutes, value: 'm' },
            { label: sharedMessages.hours, value: 'h' },
          ]}
          onChange={this.onScheduleAnytimeDelayChange}
          warning={
            shouldDisplayWarning
              ? {
                  ...sharedMessages.delayWarning,
                  values: { minimumValue: delay.MINIMUM_GATEWAY_SCHEDULE_ANYTIME_DELAY },
                }
              : undefined
          }
          required
          tooltipId={tooltipIds.SCHEDULE_ANYTIME_DELAY}
        />
        <Form.SubTitle title={sharedMessages.gatewayUpdateOptions} />
        <Form.Field
          title={sharedMessages.automaticUpdates}
          name="auto_update"
          component={Checkbox}
          description={sharedMessages.autoUpdateDescription}
        />
        <Form.Field
          title={sharedMessages.channel}
          description={sharedMessages.updateChannelDescription}
          placeholder={sharedMessages.stable}
          name="update_channel"
          component={Input}
          autoComplete="on"
        />
        <SubmitBar>{children}</SubmitBar>
      </Form>
    )
  }
}

export default GatewayDataForm

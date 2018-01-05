// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package test

import (
	"fmt"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/errors"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/smartystreets/assertions"
)

func settings(in interface{}) (*ttnpb.IdentityServerSettings, error) {
	if s, ok := in.(ttnpb.IdentityServerSettings); ok {
		return &s, nil
	}

	if s, ok := in.(*ttnpb.IdentityServerSettings); ok {
		return s, nil
	}

	return nil, errors.Errorf("Expected: '%v' to be of type ttnpb.IdentityServerSettings but it was not", in)
}

// ShouldBeSettings checks if two Settings resemble each other.
func ShouldBeSettings(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("Expected: one settings to match but got %v", len(expected))
	}

	a, err := settings(actual)
	if err != nil {
		return err.Error()
	}

	b, err := settings(expected[0])
	if err != nil {
		return err.Error()
	}

	return all(
		ShouldBeSettingsIgnoringAutoFields(a, b),
		assertions.ShouldHappenWithin(a.UpdatedAt, time.Millisecond, b.UpdatedAt),
	)
}

// ShouldBeSettingsIgnoringAutoFields checks if two Settings resemble each other
// without looking at fields that are generated by the database: UpdatedAt.
func ShouldBeSettingsIgnoringAutoFields(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("Expected: one settings to match but got %v", len(expected))
	}

	a, err := settings(actual)
	if err != nil {
		return err.Error()
	}

	b, err := settings(expected[0])
	if err != nil {
		return err.Error()
	}

	return all(
		assertions.ShouldResemble(a.BlacklistedIDs, b.BlacklistedIDs),
		assertions.ShouldEqual(a.SkipValidation, b.SkipValidation),
		assertions.ShouldEqual(a.SelfRegistration, b.SelfRegistration),
		assertions.ShouldEqual(a.AdminApproval, b.AdminApproval),
		assertions.ShouldEqual(a.ValidationTokenTTL, b.ValidationTokenTTL),
		assertions.ShouldResemble(a.AllowedEmails, b.AllowedEmails),
	)
}

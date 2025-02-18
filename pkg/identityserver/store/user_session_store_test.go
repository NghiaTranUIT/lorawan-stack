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

package store

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
)

func TestUserSessionStore(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()

	WithDB(t, func(t *testing.T, db *gorm.DB) {
		prepareTest(db, &Account{}, &User{}, &UserSession{})

		user := &User{
			Account: Account{
				UID: "test",
			},
			Name: "Test User",
		}

		userIDs := &ttnpb.UserIdentifiers{UserId: "test"}
		doesNotExistIDs := &ttnpb.UserIdentifiers{UserId: "does_not_exist"}

		if err := newStore(db).createEntity(ctx, user); err != nil {
			panic(err)
		}

		store := GetUserSessionStore(db)

		_, err := store.CreateSession(ctx, &ttnpb.UserSession{UserIds: doesNotExistIDs})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}

		created, err := store.CreateSession(ctx, &ttnpb.UserSession{
			UserIds:       userIDs,
			SessionSecret: "123412341234123412341234",
		})

		a.So(err, should.BeNil)
		if a.So(created, should.NotBeNil) {
			a.So(created.SessionId, should.NotBeEmpty)
			a.So(created.SessionSecret, should.Equal, "123412341234123412341234")
			a.So(*created.CreatedAt, should.HappenAfter, time.Now().Add(-1*time.Hour))
			a.So(*created.UpdatedAt, should.HappenAfter, time.Now().Add(-1*time.Hour))
			a.So(created.ExpiresAt, should.BeNil)
		}

		_, err = store.GetSession(ctx, doesNotExistIDs, created.SessionId)

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}

		got, err := store.GetSession(ctx, userIDs, created.SessionId)

		a.So(err, should.BeNil)
		if a.So(got, should.NotBeNil) {
			a.So(got.CreatedAt, should.Resemble, created.CreatedAt)
			a.So(got.UpdatedAt, should.Resemble, created.UpdatedAt)
			a.So(got.ExpiresAt, should.BeNil)
		}

		got, err = store.GetSessionByID(ctx, created.SessionId)

		a.So(err, should.BeNil)
		if a.So(got, should.NotBeNil) {
			a.So(got.CreatedAt, should.Resemble, created.CreatedAt)
			a.So(got.UpdatedAt, should.Resemble, created.UpdatedAt)
			a.So(got.ExpiresAt, should.BeNil)
		}

		later := time.Now().Add(time.Hour)
		updated, err := store.UpdateSession(ctx, &ttnpb.UserSession{
			UserIds:   userIDs,
			SessionId: created.SessionId,
			ExpiresAt: &later,
		})

		a.So(err, should.BeNil)
		if a.So(updated, should.NotBeNil) {
			a.So(updated.CreatedAt, should.Resemble, created.CreatedAt)
			a.So(updated.UpdatedAt, should.NotResemble, created.UpdatedAt)
			a.So(updated.ExpiresAt, should.NotBeNil)
		}

		_, err = store.UpdateSession(ctx, &ttnpb.UserSession{
			UserIds: &ttnpb.UserIdentifiers{UserId: "does_not_exist"},
		})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}

		_, err = store.UpdateSession(ctx, &ttnpb.UserSession{UserIds: userIDs, SessionId: "00000000-0000-0000-0000-000000000000"})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}

		_, err = store.FindSessions(ctx, doesNotExistIDs)

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}

		list, err := store.FindSessions(ctx, userIDs)

		a.So(err, should.BeNil)
		if a.So(list, should.HaveLength, 1) {
			a.So(list[0].CreatedAt, should.Resemble, created.CreatedAt)
			a.So(list[0].UpdatedAt, should.Resemble, updated.UpdatedAt)
			a.So(list[0].ExpiresAt, should.Resemble, updated.ExpiresAt)
		}

		err = store.DeleteSession(ctx, doesNotExistIDs, created.SessionId)

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}

		err = store.DeleteSession(ctx, userIDs, created.SessionId)

		a.So(err, should.BeNil)

		_, err = store.GetSession(ctx, userIDs, created.SessionId)

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}

		_, err = store.GetSessionByID(ctx, created.SessionId)

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}

		list, err = store.FindSessions(ctx, userIDs)

		a.So(err, should.BeNil)
		a.So(list, should.BeEmpty)

		for _, sessionSecret := range []string{
			"123412341234123412341234",
			"12341234123412341234123121",
			"12341234123412341234132143124",
			"111123124321543453456652532154",
		} {
			_, err := store.CreateSession(ctx, &ttnpb.UserSession{
				UserIds:       userIDs,
				SessionSecret: sessionSecret,
			})
			a.So(err, should.BeNil)
		}
		list, err = store.FindSessions(ctx, userIDs)

		a.So(err, should.BeNil)
		a.So(list, should.HaveLength, 4)

		err = store.DeleteAllUserSessions(ctx, userIDs)

		list, err = store.FindSessions(ctx, userIDs)

		a.So(err, should.BeNil)
		a.So(list, should.BeEmpty)
	})
}

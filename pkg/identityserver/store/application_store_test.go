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

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
)

func TestApplicationStore(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()

	WithDB(t, func(t *testing.T, db *gorm.DB) {
		prepareTest(db, &Application{}, &Attribute{})
		store := GetApplicationStore(db)
		s := newStore(db)

		created, err := store.CreateApplication(ctx, &ttnpb.Application{
			Ids:         &ttnpb.ApplicationIdentifiers{ApplicationId: "foo"},
			Name:        "Foo Application",
			Description: "The Amazing Foo Application",
			Attributes: map[string]string{
				"foo": "bar",
				"bar": "baz",
				"baz": "qux",
			},
		})

		a.So(err, should.BeNil)
		if a.So(created, should.NotBeNil) {
			a.So(created.GetIds().GetApplicationId(), should.Equal, "foo")
			a.So(created.Name, should.Equal, "Foo Application")
			a.So(created.Description, should.Equal, "The Amazing Foo Application")
			a.So(created.Attributes, should.HaveLength, 3)
			a.So(*created.CreatedAt, should.HappenAfter, time.Now().Add(-1*time.Hour))
			a.So(*created.UpdatedAt, should.HappenAfter, time.Now().Add(-1*time.Hour))
		}

		got, err := store.GetApplication(ctx, &ttnpb.ApplicationIdentifiers{ApplicationId: "foo"}, &pbtypes.FieldMask{Paths: []string{"name", "attributes"}})

		a.So(err, should.BeNil)
		if a.So(got, should.NotBeNil) {
			a.So(got.GetIds().GetApplicationId(), should.Equal, "foo")
			a.So(got.Name, should.Equal, "Foo Application")
			a.So(got.Description, should.BeEmpty)
			a.So(got.Attributes, should.HaveLength, 3)
			a.So(got.CreatedAt, should.Resemble, created.CreatedAt)
			a.So(got.UpdatedAt, should.Resemble, created.UpdatedAt)
		}

		_, err = store.UpdateApplication(ctx, &ttnpb.Application{
			Ids: &ttnpb.ApplicationIdentifiers{ApplicationId: "bar"},
		}, nil)

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}

		updated, err := store.UpdateApplication(ctx, &ttnpb.Application{
			Ids:         &ttnpb.ApplicationIdentifiers{ApplicationId: "foo"},
			Name:        "Foobar Application",
			Description: "The Amazing Foobar Application",
			Attributes: map[string]string{
				"foo": "bar",
				"baz": "baz",
				"qux": "foo",
			},
		}, &pbtypes.FieldMask{Paths: []string{"description", "attributes"}})

		a.So(err, should.BeNil)
		if a.So(updated, should.NotBeNil) {
			a.So(updated.Description, should.Equal, "The Amazing Foobar Application")
			a.So(updated.Attributes, should.HaveLength, 3)
			a.So(updated.CreatedAt, should.Resemble, created.CreatedAt)
			a.So(*updated.UpdatedAt, should.HappenAfter, *created.CreatedAt)
		}

		got, err = store.GetApplication(ctx, &ttnpb.ApplicationIdentifiers{ApplicationId: "foo"}, nil)

		a.So(err, should.BeNil)
		if a.So(got, should.NotBeNil) {
			a.So(got.GetIds().GetApplicationId(), should.Equal, created.GetIds().GetApplicationId())
			a.So(got.Name, should.Equal, created.Name)
			a.So(got.Description, should.Equal, updated.Description)
			a.So(got.Attributes, should.Resemble, updated.Attributes)
			a.So(got.CreatedAt, should.Resemble, created.CreatedAt)
			a.So(got.UpdatedAt, should.Resemble, updated.UpdatedAt)
		}

		list, err := store.FindApplications(ctx, nil, &pbtypes.FieldMask{Paths: []string{"name"}})

		a.So(err, should.BeNil)
		if a.So(list, should.HaveLength, 1) {
			a.So(list[0].Name, should.EndWith, got.Name)
		}

		err = store.DeleteApplication(ctx, &ttnpb.ApplicationIdentifiers{ApplicationId: "foo"})

		a.So(err, should.BeNil)

		got, err = store.GetApplication(ctx, &ttnpb.ApplicationIdentifiers{ApplicationId: "foo"}, nil)

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}

		err = store.RestoreApplication(ctx, &ttnpb.ApplicationIdentifiers{ApplicationId: "foo"})

		a.So(err, should.BeNil)

		got, err = store.GetApplication(ctx, &ttnpb.ApplicationIdentifiers{ApplicationId: "foo"}, nil)

		a.So(err, should.BeNil)

		err = store.DeleteApplication(ctx, &ttnpb.ApplicationIdentifiers{ApplicationId: "foo"})

		a.So(err, should.BeNil)

		list, err = store.FindApplications(ctx, nil, nil)

		a.So(err, should.BeNil)
		a.So(list, should.BeEmpty)

		list, err = store.FindApplications(WithSoftDeleted(ctx, false), nil, nil)

		a.So(err, should.BeNil)
		a.So(list, should.NotBeEmpty)

		entity, _ := s.findDeletedEntity(ctx, &ttnpb.ApplicationIdentifiers{ApplicationId: "foo"}, "id")

		err = store.PurgeApplication(ctx, &ttnpb.ApplicationIdentifiers{ApplicationId: "foo"})

		a.So(err, should.BeNil)

		var attribute []Attribute
		s.query(ctx, Attribute{}).Where(&Attribute{
			EntityID:   entity.PrimaryKey(),
			EntityType: "application",
		}).Find(&attribute)

		a.So(attribute, should.HaveLength, 0)

		// Check that application ids are released after purge
		_, err = store.CreateApplication(ctx, &ttnpb.Application{
			Ids: &ttnpb.ApplicationIdentifiers{ApplicationId: "foo"},
		})

		a.So(err, should.BeNil)
	})
}

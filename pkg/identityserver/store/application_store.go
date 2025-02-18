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
	"context"
	"fmt"
	"reflect"
	"runtime/trace"
	"strings"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcmiddleware/warning"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// GetApplicationStore returns an ApplicationStore on the given db (or transaction).
func GetApplicationStore(db *gorm.DB) ApplicationStore {
	return &applicationStore{store: newStore(db)}
}

type applicationStore struct {
	*store
}

// selectApplicationFields selects relevant fields (based on fieldMask) and preloads details if needed.
func selectApplicationFields(ctx context.Context, query *gorm.DB, fieldMask *pbtypes.FieldMask) *gorm.DB {
	if len(fieldMask.GetPaths()) == 0 {
		return query.Preload("Attributes")
	}
	var applicationColumns []string
	var notFoundPaths []string
	for _, path := range ttnpb.TopLevelFields(fieldMask.GetPaths()) {
		switch path {
		case "ids", "created_at", "updated_at", "deleted_at":
			// always selected
		case attributesField:
			query = query.Preload("Attributes")
		default:
			if columns, ok := applicationColumnNames[path]; ok {
				applicationColumns = append(applicationColumns, columns...)
			} else {
				notFoundPaths = append(notFoundPaths, path)
			}
		}
	}
	if len(notFoundPaths) > 0 {
		warning.Add(ctx, fmt.Sprintf("unsupported field mask paths: %s", strings.Join(notFoundPaths, ", ")))
	}
	return query.Select(cleanFields(append(append(modelColumns, "deleted_at", "application_id"), applicationColumns...)...))
}

func (s *applicationStore) CreateApplication(ctx context.Context, app *ttnpb.Application) (*ttnpb.Application, error) {
	defer trace.StartRegion(ctx, "create application").End()
	appModel := Application{
		ApplicationID: app.GetIds().GetApplicationId(), // The ID is not mutated by fromPB.
	}
	appModel.fromPB(app, nil)
	if err := s.createEntity(ctx, &appModel); err != nil {
		return nil, err
	}
	var appProto ttnpb.Application
	appModel.toPB(&appProto, nil)
	return &appProto, nil
}

func (s *applicationStore) FindApplications(ctx context.Context, ids []*ttnpb.ApplicationIdentifiers, fieldMask *pbtypes.FieldMask) ([]*ttnpb.Application, error) {
	defer trace.StartRegion(ctx, "find applications").End()
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.GetApplicationId()
	}
	query := s.query(ctx, Application{}, withApplicationID(idStrings...))
	query = selectApplicationFields(ctx, query, fieldMask)
	query = query.Order(orderFromContext(ctx, "applications", "application_id", "ASC"))
	if limit, offset := limitAndOffsetFromContext(ctx); limit != 0 {
		countTotal(ctx, query.Model(&Application{}))
		query = query.Limit(limit).Offset(offset)
	}
	if onlyExpired, expireThreshold := expiredFromContext(ctx); onlyExpired {
		query = query.Scopes(withExpiredEntities(expireThreshold))
	}
	var appModels []Application
	query = query.Find(&appModels)
	setTotal(ctx, uint64(len(appModels)))
	if query.Error != nil {
		return nil, query.Error
	}
	appProtos := make([]*ttnpb.Application, len(appModels))
	for i, appModel := range appModels {
		appProto := &ttnpb.Application{}
		appModel.toPB(appProto, fieldMask)
		appProtos[i] = appProto
	}
	return appProtos, nil
}

func (s *applicationStore) GetApplication(ctx context.Context, id *ttnpb.ApplicationIdentifiers, fieldMask *pbtypes.FieldMask) (*ttnpb.Application, error) {
	defer trace.StartRegion(ctx, "get application").End()
	query := s.query(ctx, Application{}, withApplicationID(id.GetApplicationId()))
	query = selectApplicationFields(ctx, query, fieldMask)
	var appModel Application
	if err := query.First(&appModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errNotFoundForID(id)
		}
		return nil, err
	}
	appProto := &ttnpb.Application{}
	appModel.toPB(appProto, fieldMask)
	return appProto, nil
}

func (s *applicationStore) UpdateApplication(ctx context.Context, app *ttnpb.Application, fieldMask *pbtypes.FieldMask) (updated *ttnpb.Application, err error) {
	defer trace.StartRegion(ctx, "update application").End()
	query := s.query(ctx, Application{}, withApplicationID(app.GetIds().GetApplicationId()))
	query = selectApplicationFields(ctx, query, fieldMask)
	var appModel Application
	if err = query.First(&appModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errNotFoundForID(app.GetIds())
		}
		return nil, err
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return nil, err
	}
	oldAttributes := appModel.Attributes
	columns := appModel.fromPB(app, fieldMask)
	if err = s.updateEntity(ctx, &appModel, columns...); err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(oldAttributes, appModel.Attributes) {
		if err = s.replaceAttributes(ctx, "application", appModel.ID, oldAttributes, appModel.Attributes); err != nil {
			return nil, err
		}
	}
	updated = &ttnpb.Application{}
	appModel.toPB(updated, fieldMask)
	return updated, nil
}

func (s *applicationStore) DeleteApplication(ctx context.Context, id *ttnpb.ApplicationIdentifiers) error {
	defer trace.StartRegion(ctx, "delete application").End()
	return s.deleteEntity(ctx, id)
}

func (s *applicationStore) RestoreApplication(ctx context.Context, id *ttnpb.ApplicationIdentifiers) error {
	defer trace.StartRegion(ctx, "restore application").End()
	return s.restoreEntity(ctx, id)
}

func (s *applicationStore) PurgeApplication(ctx context.Context, id *ttnpb.ApplicationIdentifiers) error {
	defer trace.StartRegion(ctx, "purge application").End()
	query := s.query(ctx, Application{}, withSoftDeleted(), withApplicationID(id.GetApplicationId()))
	query = selectApplicationFields(ctx, query, nil)
	var appModel Application
	if err := query.First(&appModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errNotFoundForID(id)
		}
		return err
	}
	// delete application attributes before purging
	if len(appModel.Attributes) > 0 {
		if err := s.replaceAttributes(ctx, "application", appModel.ID, appModel.Attributes, nil); err != nil {
			return err
		}
	}
	return s.purgeEntity(ctx, id)
}

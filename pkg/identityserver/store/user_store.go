// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
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

// GetUserStore returns an UserStore on the given db (or transaction).
func GetUserStore(db *gorm.DB) UserStore {
	return &userStore{store: newStore(db)}
}

type userStore struct {
	*store
}

// selectUserFields selects relevant fields (based on fieldMask) and preloads details if needed.
func selectUserFields(ctx context.Context, query *gorm.DB, fieldMask *pbtypes.FieldMask) *gorm.DB {
	if len(fieldMask.GetPaths()) == 0 {
		return query.Preload("Attributes").Preload("ProfilePicture").Select([]string{"accounts.uid", "users.*"})
	}
	var userColumns []string
	var notFoundPaths []string
	userColumns = append(userColumns, "users.deleted_at", "accounts.uid")
	for _, column := range modelColumns {
		userColumns = append(userColumns, "users."+column)
	}
	for _, path := range ttnpb.TopLevelFields(fieldMask.GetPaths()) {
		switch path {
		case "ids", "created_at", "updated_at", "deleted_at":
			// always selected
		case attributesField:
			query = query.Preload("Attributes")
		case profilePictureField:
			userColumns = append(userColumns, "profile_picture_id")
			query = query.Preload("ProfilePicture")
		default:
			if columns, ok := userColumnNames[path]; ok {
				userColumns = append(userColumns, columns...)
			} else {
				notFoundPaths = append(notFoundPaths, path)
			}
		}
	}
	if len(notFoundPaths) > 0 {
		warning.Add(ctx, fmt.Sprintf("unsupported field mask paths: %s", strings.Join(notFoundPaths, ", ")))
	}
	return query.Select(userColumns)
}

func (s *userStore) CreateUser(ctx context.Context, usr *ttnpb.User) (*ttnpb.User, error) {
	defer trace.StartRegion(ctx, "create user").End()
	userModel := User{
		Account: Account{UID: usr.GetIds().GetUserId()}, // The ID is not mutated by fromPB.
	}
	fieldMask := &pbtypes.FieldMask{Paths: append(defaultUserFieldMask.GetPaths(), passwordField)}
	userModel.fromPB(usr, fieldMask)
	if err := s.createEntity(ctx, &userModel); err != nil {
		return nil, err
	}
	var userProto ttnpb.User
	userModel.toPB(&userProto, nil)
	return &userProto, nil
}

func (s *userStore) FindUsers(ctx context.Context, ids []*ttnpb.UserIdentifiers, fieldMask *pbtypes.FieldMask) ([]*ttnpb.User, error) {
	defer trace.StartRegion(ctx, "find users").End()
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.GetUserId()
	}
	query := s.query(ctx, User{}, withUserID(idStrings...))
	query = selectUserFields(ctx, query, fieldMask)
	query = query.Order(orderFromContext(ctx, "users", `"accounts"."uid"`, "ASC"))
	if limit, offset := limitAndOffsetFromContext(ctx); limit != 0 {
		countTotal(ctx, query.Model(User{}))
		query = query.Limit(limit).Offset(offset)
	}
	if onlyExpired, expireThreshold := expiredFromContext(ctx); onlyExpired {
		query = query.Scopes(withExpiredEntities(expireThreshold))
	}
	var userModels []userWithUID
	query = query.Find(&userModels)
	setTotal(ctx, uint64(len(userModels)))
	if query.Error != nil {
		return nil, query.Error
	}
	userProtos := make([]*ttnpb.User, len(userModels))
	for i, userModel := range userModels {
		userProto := &ttnpb.User{}
		userModel.toPB(userProto, fieldMask)
		userProtos[i] = userProto
	}
	return userProtos, nil
}

func (s *userStore) ListAdmins(ctx context.Context, fieldMask *pbtypes.FieldMask) ([]*ttnpb.User, error) {
	defer trace.StartRegion(ctx, "list admins").End()

	query := s.query(ctx, User{}, withUserID()).Where(&User{Admin: true})
	query = selectUserFields(ctx, query, fieldMask)
	query = query.Order(orderFromContext(ctx, "users", `"accounts"."uid"`, "ASC"))
	if limit, offset := limitAndOffsetFromContext(ctx); limit != 0 {
		countTotal(ctx, query.Model(User{}))
		query = query.Limit(limit).Offset(offset)
	}
	var userModels []userWithUID
	query = query.Find(&userModels)
	setTotal(ctx, uint64(len(userModels)))
	if query.Error != nil {
		return nil, query.Error
	}
	userProtos := make([]*ttnpb.User, len(userModels))
	for i, userModel := range userModels {
		userProto := &ttnpb.User{}
		userModel.toPB(userProto, fieldMask)
		userProtos[i] = userProto
	}
	return userProtos, nil
}

func (s *userStore) GetUser(ctx context.Context, id *ttnpb.UserIdentifiers, fieldMask *pbtypes.FieldMask) (*ttnpb.User, error) {
	defer trace.StartRegion(ctx, "get user").End()
	query := s.query(ctx, User{}, withUserID(id.GetUserId()))
	query = selectUserFields(ctx, query, fieldMask)
	var userModel userWithUID
	if err := query.First(&userModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errNotFoundForID(id)
		}
		return nil, err
	}
	userProto := &ttnpb.User{}
	userModel.toPB(userProto, fieldMask)
	return userProto, nil
}

func (s *userStore) GetUserByPrimaryEmailAddress(ctx context.Context, email string, fieldMask *pbtypes.FieldMask) (*ttnpb.User, error) {
	defer trace.StartRegion(ctx, "get user by primary email address").End()
	query := s.query(ctx, User{}, withPrimaryEmailAddress(email))
	query = query.Joins("LEFT JOIN accounts ON accounts.account_type = ? AND accounts.account_id = users.id", "user")
	query = selectUserFields(ctx, query, fieldMask)
	var userModel userWithUID
	if err := query.First(&userModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errUserNotFound.WithAttributes("user_id", email)
		}
		return nil, err
	}
	userProto := &ttnpb.User{}
	userModel.toPB(userProto, fieldMask)
	return userProto, nil
}

func (s *userStore) UpdateUser(ctx context.Context, usr *ttnpb.User, fieldMask *pbtypes.FieldMask) (updated *ttnpb.User, err error) {
	defer trace.StartRegion(ctx, "update user").End()
	query := s.query(ctx, User{}, withUserID(usr.GetIds().GetUserId()))
	query = selectUserFields(ctx, query, fieldMask)
	var userModel userWithUID
	if err = query.First(&userModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errNotFoundForID(usr.GetIds())
		}
		return nil, err
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return nil, err
	}
	oldAttributes, oldProfilePicture := userModel.Attributes, userModel.ProfilePicture
	columns := userModel.fromPB(usr, fieldMask)
	newProfilePicture := userModel.ProfilePicture
	if newProfilePicture != oldProfilePicture {
		if oldProfilePicture != nil {
			if err = s.query(ctx, Picture{}).Delete(oldProfilePicture).Error; err != nil {
				return nil, err
			}
		}
		if newProfilePicture != nil {
			if err = s.createEntity(ctx, &newProfilePicture); err != nil {
				return nil, err
			}
			userModel.ProfilePictureID, userModel.ProfilePicture = &newProfilePicture.ID, nil
			columns = append(columns, "profile_picture_id")
		}
	}
	if err = s.updateEntity(ctx, &userModel.User, columns...); err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(oldAttributes, userModel.Attributes) {
		if err = s.replaceAttributes(ctx, "user", userModel.ID, oldAttributes, userModel.Attributes); err != nil {
			return nil, err
		}
	}
	userModel.ProfilePicture = newProfilePicture
	updated = &ttnpb.User{}
	userModel.toPB(updated, fieldMask)
	return updated, nil
}

func (s *userStore) DeleteUser(ctx context.Context, id *ttnpb.UserIdentifiers) (err error) {
	defer trace.StartRegion(ctx, "delete user").End()
	return s.deleteEntity(ctx, id)
}

func (s *userStore) RestoreUser(ctx context.Context, id *ttnpb.UserIdentifiers) (err error) {
	defer trace.StartRegion(ctx, "restore user").End()
	return s.restoreEntity(ctx, id)
}

func (s *userStore) PurgeUser(ctx context.Context, id *ttnpb.UserIdentifiers) (err error) {
	defer trace.StartRegion(ctx, "purge user").End()
	query := s.query(ctx, User{}, withSoftDeleted(), withUserID(id.GetUserId()))
	query = selectUserFields(ctx, query, nil)
	var userModel userWithUID
	if err = query.First(&userModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errNotFoundForID(id)
		}
		return err
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return err
	}
	if len(userModel.Attributes) > 0 {
		if err := s.replaceAttributes(ctx, "user", userModel.ID, userModel.Attributes, nil); err != nil {
			return err
		}
	}
	if userModel.ProfilePicture != nil {
		if err = s.query(ctx, Picture{}).Delete(userModel.ProfilePicture).Error; err != nil {
			return err
		}
	}

	err = s.purgeEntity(ctx, id)
	if err != nil {
		return err
	}
	// Purge account after purging user because it is necessary for user query
	return s.query(ctx, Account{}, withSoftDeleted()).Where(Account{
		UID:         id.IDString(),
		AccountType: id.EntityType(),
	}).Delete(Account{}).Error
}

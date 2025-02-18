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
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"runtime/trace"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

func newStore(db *gorm.DB) *store { return &store{DB: db} }

type store struct {
	DB *gorm.DB
}

func (s *store) query(ctx context.Context, model interface{}, funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	query := s.DB.Model(model).Scopes(withContext(ctx), withSoftDeletedIfRequested(ctx))
	if len(funcs) > 0 {
		query = query.Scopes(funcs...)
	}
	return query
}

func (s *store) findEntity(ctx context.Context, entityID ttnpb.IDStringer, fields ...string) (modelInterface, error) {
	model := modelForID(entityID)
	query := s.query(ctx, model, withID(entityID))
	if len(fields) == 1 && fields[0] == "id" {
		fields[0] = s.DB.NewScope(model).TableName() + ".id"
	}
	if len(fields) > 0 {
		query = query.Select(fields)
	}
	if err := query.First(model).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errNotFoundForID(entityID)
		}
		return nil, convertError(err)
	}
	return model, nil
}

func (s *store) findDeletedEntity(ctx context.Context, entityID ttnpb.IDStringer, fields ...string) (modelInterface, error) {
	return s.findEntity(WithSoftDeleted(ctx, false), entityID, fields...)
}

func (s *store) createEntity(ctx context.Context, model interface{}) error {
	if model, ok := model.(modelInterface); ok {
		model.SetContext(ctx)
	}
	return s.DB.Create(model).Error
}

func (s *store) updateEntity(ctx context.Context, model interface{}, columns ...string) error {
	query := s.query(ctx, model)
	query = query.Select(append(columns, "updated_at"))
	return query.Save(model).Error
}

func (s *store) deleteEntity(ctx context.Context, entityID ttnpb.IDStringer) error {
	model, err := s.findEntity(ctx, entityID, "id")
	if err != nil {
		return err
	}
	if err = s.DB.Delete(model).Error; err != nil {
		return err
	}
	switch entityType := entityID.EntityType(); entityType {
	case "user", "organization":
		err = s.DB.Where(Account{
			AccountType: entityType,
			AccountID:   model.PrimaryKey(),
		}).Delete(Account{}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *store) restoreEntity(ctx context.Context, entityID ttnpb.IDStringer) error {
	model, err := s.findDeletedEntity(ctx, entityID, "id")
	if err != nil {
		return err
	}
	switch entityType := entityID.EntityType(); entityType {
	case "user", "organization":
		err := s.DB.Unscoped().Model(Account{}).Where(Account{
			AccountType: entityType,
			AccountID:   model.PrimaryKey(),
		}).UpdateColumn("deleted_at", gorm.Expr("NULL")).Error
		if err != nil {
			return err
		}
	}
	return s.DB.Unscoped().Model(model).UpdateColumn("deleted_at", gorm.Expr("NULL")).Error
}

func (s *store) purgeEntity(ctx context.Context, entityID ttnpb.IDStringer) error {
	model, err := s.findDeletedEntity(ctx, entityID, "id")
	if err != nil {
		return err
	}
	return s.DB.Unscoped().Delete(model).Error
}

var (
	errDatabase      = errors.DefineInternal("database", "database error")
	errAlreadyExists = errors.DefineAlreadyExists("already_exists", "entity already exists")

	// ErrIDTaken is returned when an entity can not be created because the ID is already taken.
	ErrIDTaken = errors.DefineAlreadyExists("id_taken", "ID already taken")
	// ErrEUITaken is returned when an entity can not be created because the EUI is already taken.
	ErrEUITaken = errors.DefineAlreadyExists("eui_taken", "EUI already taken")
)

var uniqueViolationRegex = regexp.MustCompile(`duplicate key value( .+)? violates unique constraint "([a-z_]+)"`)

func convertError(err error) error {
	switch err {
	case nil, context.Canceled, context.DeadlineExceeded:
		return err
	}
	if ttnErr, ok := errors.From(err); ok {
		return ttnErr
	}
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code.Name() {
		case "unique_violation":
			if match := uniqueViolationRegex.FindStringSubmatch(pqErr.Message); match != nil {
				switch {
				case strings.HasSuffix(match[2], "_id_index"):
					return ErrIDTaken.WithCause(err)
				case strings.HasSuffix(match[2], "_eui_index"):
					return ErrEUITaken.WithCause(err)
				default:
					return errAlreadyExists.WithCause(err).WithAttributes("index", match[2])
				}
			}
			return errAlreadyExists.WithCause(err)
		default:
			return errDatabase.WithCause(err).WithAttributes("code", pqErr.Code.Name())
		}
	}
	return errDatabase.WithCause(err)
}

// Open opens a new database connection.
func Open(ctx context.Context, dsn string) (*gorm.DB, error) {
	dbURI, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	dbName := strings.TrimPrefix(dbURI.Path, "/")
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db = db.Set("db:name", dbName)
	var dbVersion string
	err = db.Raw("SELECT version()").Row().Scan(&dbVersion)
	if err != nil {
		return nil, err
	}
	db = db.Set("db:version", dbVersion)
	switch {
	case strings.Contains(dbVersion, "CockroachDB"):
		db = db.Set("db:kind", "CockroachDB")
	case strings.Contains(dbVersion, "PostgreSQL"):
		db = db.Set("db:kind", "PostgreSQL")
	}
	SetLogger(db, log.FromContext(ctx))
	return db, nil
}

// Initialize initializes the database.
func Initialize(db *gorm.DB) error {
	if dbKind, ok := db.Get("db:kind"); ok {
		switch dbKind {
		case "CockroachDB":
			if dbName, ok := db.Get("db:name"); ok {
				if err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName)).Error; err != nil {
					return err
				}
			}
		case "PostgreSQL":
			if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// ErrTransactionRecovered is returned when a panic is caught from a SQL transaction.
var ErrTransactionRecovered = errors.DefineInternal("transaction_recovered", "Internal Server Error")

// Transact executes f in a db transaction.
func Transact(ctx context.Context, db *gorm.DB, f func(db *gorm.DB) error) (err error) {
	defer trace.StartRegion(ctx, "database transaction").End()
	tx := db.Begin()
	if tx.Error != nil {
		return convertError(tx.Error)
	}
	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintln(os.Stderr, p)
			os.Stderr.Write(debug.Stack())
			if pErr, ok := p.(error); ok {
				switch pErr {
				case context.Canceled, context.DeadlineExceeded:
					err = pErr
				default:
					err = ErrTransactionRecovered.WithCause(pErr)
				}
			} else {
				err = ErrTransactionRecovered.WithAttributes("panic", p)
			}
			log.FromContext(ctx).WithError(err).Error("Transaction panicked")
		}
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
		err = convertError(err)
	}()
	SetLogger(tx, log.FromContext(ctx).WithField("namespace", "db"))
	return f(tx)
}

func entityTypeForID(id ttnpb.IDStringer) string {
	return strings.Replace(id.EntityType(), " ", "_", -1)
}

func modelForEntityType(entityType string) modelInterface {
	switch entityType {
	case "application":
		return &Application{}
	case "client":
		return &Client{}
	case "end_device":
		return &EndDevice{}
	case "gateway":
		return &Gateway{}
	case "organization":
		return &Organization{}
	case "user":
		return &User{}
	default:
		panic(fmt.Sprintf("can't find model for entity type %s", entityType))
	}
}

func modelForID(id ttnpb.IDStringer) modelInterface {
	return modelForEntityType(entityTypeForID(id))
}

var (
	errApplicationNotFound  = errors.DefineNotFound("application_not_found", "application `{application_id}` not found")
	errClientNotFound       = errors.DefineNotFound("client_not_found", "client `{client_id}` not found")
	errGatewayNotFound      = errors.DefineNotFound("gateway_not_found", "gateway `{gateway_id}` not found")
	errEndDeviceNotFound    = errors.DefineNotFound("end_device_not_found", "end device `{application_id}:{device_id}` not found")
	errOrganizationNotFound = errors.DefineNotFound("organization_not_found", "organization `{organization_id}` not found")
	errUserNotFound         = errors.DefineNotFound("user_not_found", "user `{user_id}` not found")
	errSessionNotFound      = errors.DefineNotFound("session_not_found", "session `{session_id}` for user `{user_id}` not found")

	errAuthorizationNotFound     = errors.DefineNotFound("authorization_not_found", "authorization of `{user_id}` for `{client_id}` not found")
	errAuthorizationCodeNotFound = errors.DefineNotFound("authorization_code_not_found", "authorization code not found")
	errAccessTokenNotFound       = errors.DefineNotFound("access_token_not_found", "access token `{access_token_id}` not found")

	errAPIKeyNotFound = errors.DefineNotFound("api_key_not_found", "API key not found")

	errMigrationNotFound = errors.DefineNotFound("migration_not_found", "migration not found")
)

func errNotFoundForID(id ttnpb.IDStringer) error {
	switch t := entityTypeForID(id); t {
	case "application":
		return errApplicationNotFound.WithAttributes("application_id", id.IDString())
	case "client":
		return errClientNotFound.WithAttributes("client_id", id.IDString())
	case "end_device":
		appID, devID := splitEndDeviceIDString(id.IDString())
		return errEndDeviceNotFound.WithAttributes("application_id", appID, "device_id", devID)
	case "gateway":
		return errGatewayNotFound.WithAttributes("gateway_id", id.IDString())
	case "organization":
		return errOrganizationNotFound.WithAttributes("organization_id", id.IDString())
	case "user":
		return errUserNotFound.WithAttributes("user_id", id.IDString())
	default:
		panic(fmt.Sprintf("can't find errNotFound for entity type %s", t))
	}
}

// SetLogger sets the database logger.
func SetLogger(db *gorm.DB, log log.Interface) {
	db.SetLogger(logger{Interface: log})
}

type logger struct {
	log.Interface
}

// Print implements the gorm.logger interface.
func (l logger) Print(v ...interface{}) {
	if len(v) < 3 {
		l.Error(fmt.Sprint(v...))
		return
	}
	logger := l.Interface
	source, ok := v[1].(string)
	if !ok {
		l.Error(fmt.Sprint(v...))
		return
	}
	logger = logger.WithField("source", filepath.Base(source))
	switch v[0] {
	case "log", "error":
		if err, ok := v[2].(error); ok {
			err = convertError(err)
			if errors.Resemble(err, errDatabase) {
				logger.WithError(err).Error("Database error")
			}
			return
		}
		logger.Error(fmt.Sprint(v[2:]...))
		return
	case "sql":
		if len(v) != 6 {
			return
		}
		duration, _ := v[2].(time.Duration)
		query, _ := v[3].(string)
		values, _ := v[4].([]interface{})
		rows, _ := v[5].(int64)
		logger.WithFields(log.Fields(
			"duration", duration,
			"query", query,
			"values", values,
			"rows", rows,
		)).Debug("Run database query")
	default:
		l.Error(fmt.Sprint(v...))
	}
}

func cleanFields(fields ...string) []string {
	seen := make(map[string]struct{}, len(fields))
	out := make([]string, 0, len(fields))
	for _, field := range fields {
		if _, ok := seen[field]; ok {
			continue
		}
		seen[field] = struct{}{}
		out = append(out, field)
	}
	return out
}

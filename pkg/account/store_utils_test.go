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

package account_test

import (
	"context"

	"github.com/gogo/protobuf/types"
	account_store "go.thethings.network/lorawan-stack/v3/pkg/account/store"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type mockStoreContents struct {
	calls []string
	req   struct {
		ctx       context.Context
		fieldMask *types.FieldMask
		session   *ttnpb.UserSession
		sessionID string
		userIDs   *ttnpb.UserIdentifiers
		token     string
	}
	res struct {
		session    *ttnpb.UserSession
		user       *ttnpb.User
		loginToken *ttnpb.LoginToken
	}
	err struct {
		getUser       error
		createSession error
		getSession    error
		deleteSession error
		loginToken    error
	}
}

type mockStore struct {
	store.UserStore
	store.LoginTokenStore
	store.UserSessionStore

	mockStoreContents
}

func (s *mockStore) reset() {
	s.mockStoreContents = mockStoreContents{}
}

var (
	mockErrUnauthenticated = grpc.Errorf(codes.Unauthenticated, "Unauthenticated")
	mockErrNotFound        = grpc.Errorf(codes.NotFound, "NotFound")
)

func (s *mockStore) GetUser(ctx context.Context, id *ttnpb.UserIdentifiers, fieldMask *types.FieldMask) (*ttnpb.User, error) {
	s.req.ctx, s.req.userIDs, s.req.fieldMask = ctx, id, fieldMask
	s.calls = append(s.calls, "GetUser")
	return s.res.user, s.err.getUser
}

func (s *mockStore) CreateSession(ctx context.Context, sess *ttnpb.UserSession) (*ttnpb.UserSession, error) {
	s.req.ctx, s.req.session = ctx, sess
	s.calls = append(s.calls, "CreateSession")
	return s.res.session, s.err.createSession
}

func (s *mockStore) GetSession(ctx context.Context, userIDs *ttnpb.UserIdentifiers, sessionID string) (*ttnpb.UserSession, error) {
	s.req.ctx, s.req.userIDs, s.req.sessionID = ctx, userIDs, sessionID
	s.calls = append(s.calls, "GetSession")
	return s.res.session, s.err.getSession
}

func (s *mockStore) DeleteSession(ctx context.Context, userIDs *ttnpb.UserIdentifiers, sessionID string) error {
	s.req.ctx, s.req.userIDs, s.req.sessionID = ctx, userIDs, sessionID
	s.calls = append(s.calls, "DeleteSession")
	return s.err.deleteSession
}

func (s *mockStore) ConsumeLoginToken(ctx context.Context, token string) (*ttnpb.LoginToken, error) {
	s.req.ctx, s.req.token = ctx, token
	s.calls = append(s.calls, "ConsumeLoginToken")
	return s.res.loginToken, s.err.loginToken
}

func (s *mockStore) WithSoftDeleted(ctx context.Context, b bool) context.Context {
	return ctx
}

func (s *mockStore) Transact(ctx context.Context, f func(context.Context, account_store.Interface) error) error {
	return f(ctx, s)
}

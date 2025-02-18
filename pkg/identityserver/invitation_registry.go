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

package identityserver

import (
	"context"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/v3/pkg/auth"
	"go.thethings.network/lorawan-stack/v3/pkg/email"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/emails"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

var evtCreateInvitation = events.Define(
	"invitation.create", "create invitation",
	events.WithAuthFromContext(),
	events.WithClientInfoFromContext(),
)

var errNoInviteRights = errors.DefinePermissionDenied(
	"no_invite_rights",
	"no rights for inviting users",
)

func (is *IdentityServer) sendInvitation(ctx context.Context, in *ttnpb.SendInvitationRequest) (*ttnpb.Invitation, error) {
	authInfo, err := is.authInfo(ctx)
	if err != nil {
		return nil, err
	}
	if !authInfo.GetUniversalRights().IncludesAll(ttnpb.RIGHT_SEND_INVITES) {
		return nil, errNoInviteRights.New()
	}
	token, err := auth.GenerateKey(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	ttl := is.configFromContext(ctx).UserRegistration.Invitation.TokenTTL
	expires := now.Add(ttl)
	invitation := &ttnpb.Invitation{
		Email:     in.Email,
		Token:     token,
		ExpiresAt: &expires,
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		invitation, err = store.GetInvitationStore(db).CreateInvitation(ctx, invitation)
		return err
	})
	if err != nil {
		return nil, err
	}
	events.Publish(evtCreateInvitation.NewWithIdentifiersAndData(ctx, nil, invitation))
	err = is.SendEmail(ctx, func(data emails.Data) email.MessageData {
		data.User.Email = in.Email
		return &emails.Invitation{
			Data:            data,
			InvitationToken: invitation.Token,
			TTL:             ttl,
		}
	})
	if err != nil {
		log.FromContext(ctx).WithError(err).Error("Could not send invitation email")
	}
	return invitation, nil
}

func (is *IdentityServer) listInvitations(ctx context.Context, req *ttnpb.ListInvitationsRequest) (invitations *ttnpb.Invitations, err error) {
	authInfo, err := is.authInfo(ctx)
	if err != nil {
		return nil, err
	}
	if !authInfo.GetUniversalRights().IncludesAll(ttnpb.RIGHT_SEND_INVITES) {
		return nil, errNoInviteRights.New()
	}
	invitations = &ttnpb.Invitations{}
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		invitations.Invitations, err = store.GetInvitationStore(db).FindInvitations(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

func (is *IdentityServer) deleteInvitation(ctx context.Context, in *ttnpb.DeleteInvitationRequest) (*pbtypes.Empty, error) {
	authInfo, err := is.authInfo(ctx)
	if err != nil {
		return nil, err
	}
	if !authInfo.GetUniversalRights().IncludesAll(ttnpb.RIGHT_SEND_INVITES) {
		return nil, errNoInviteRights.New()
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) error {
		return store.GetInvitationStore(db).DeleteInvitation(ctx, in.Email)
	})
	if err != nil {
		return nil, err
	}
	return ttnpb.Empty, nil
}

type invitationRegistry struct {
	*IdentityServer
}

func (ir *invitationRegistry) Send(ctx context.Context, req *ttnpb.SendInvitationRequest) (*ttnpb.Invitation, error) {
	return ir.sendInvitation(ctx, req)
}

func (ir *invitationRegistry) List(ctx context.Context, req *ttnpb.ListInvitationsRequest) (*ttnpb.Invitations, error) {
	return ir.listInvitations(ctx, req)
}

func (ir *invitationRegistry) Delete(ctx context.Context, req *ttnpb.DeleteInvitationRequest) (*pbtypes.Empty, error) {
	return ir.deleteInvitation(ctx, req)
}

// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"strings"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/require"
)

func TestInvitationProfile(t *testing.T) {
	for _, tc := range []struct{ first, last, want string }{
		{"José", "García", "jose.garcia"},
		{"María José", "De la Peña", "mariajose.delapena"},
		{"Anne", "O'Neill-Smith", "anne.oneillsmith"},
		{"李", "王", "user.user"},
		{"123", "Smith", "u123.smith"},
	} {
		t.Run(tc.want, func(t *testing.T) {
			name := invitationUsername(tc.first, tc.last)
			require.Equal(t, tc.want, name)
			require.True(t, model.IsValidUsername(name))
			require.Equal(t, name+"2", invitationUsernameCandidate(name, 2))
		})
	}
	user := &model.User{FirstName: " José ", LastName: " García "}
	require.Nil(t, prepareInvitationProfile(user))
	require.Equal(t, "José", user.FirstName)
	require.Equal(t, "García", user.LastName)
	for _, user := range []*model.User{
		{FirstName: " ", LastName: "Smith"},
		{FirstName: "\uFEFF", LastName: "Smith"},
		{FirstName: "John", LastName: " "},
		{FirstName: strings.Repeat("a", 65), LastName: "Smith"},
	} {
		require.NotNil(t, prepareInvitationProfile(user))
	}
	base := invitationUsername(strings.Repeat("a", 64), strings.Repeat("b", 64))
	require.True(t, model.IsValidUsername(invitationUsernameCandidate(base, 1000)))
	require.LessOrEqual(t, len(invitationUsernameCandidate(base, 1000)), model.UserNameMaxLength)
}

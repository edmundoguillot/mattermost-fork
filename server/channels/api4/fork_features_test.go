// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestForkFeaturesCustomSchemesKeepAuthorization(t *testing.T) {
	original := model.BuildForkFeatures
	model.BuildForkFeatures = "true"
	defer func() { model.BuildForkFeatures = original }()
	th := Setup(t).InitBasic(t)
	require.NoError(t, th.App.SetPhase2PermissionsMigrationStatus(true))

	// The local profile must survive clearing or reloading an uploaded license.
	th.App.Srv().SetLicense(nil)
	th.App.Srv().LoadLicense()
	scheme := &model.Scheme{
		Name:        model.NewId(),
		DisplayName: "Local permissions",
		Scope:       model.SchemeScopeTeam,
	}
	created, _, err := th.SystemAdminClient.CreateScheme(context.Background(), scheme)
	require.NoError(t, err)
	require.NotEmpty(t, created.Id)

	role, _, err := th.SystemAdminClient.GetRoleByName(context.Background(), created.DefaultTeamUserRole)
	require.NoError(t, err)
	_, _, err = th.SystemAdminClient.PatchRole(context.Background(), role.Id, &model.RolePatch{Permissions: &role.Permissions})
	require.NoError(t, err)

	scheme.Name = model.NewId()
	_, response, err := th.Client.CreateScheme(context.Background(), scheme)
	require.Error(t, err)
	CheckForbiddenStatus(t, response)
}

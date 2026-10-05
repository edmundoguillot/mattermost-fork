// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package platform

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestForkFeaturesExposeMatchingServerAndClientEntitlements(t *testing.T) {
	original := model.BuildForkFeatures
	model.BuildForkFeatures = "true"
	t.Cleanup(func() { model.BuildForkFeatures = original })
	ps := &PlatformService{}
	require.True(t, *ps.License().Features.CustomPermissionsSchemes)
	client := ps.ClientLicense()
	require.Equal(t, "true", client["CustomPermissionsSchemes"])
	require.Equal(t, "true", client["GuestAccountsPermissions"])
	require.Equal(t, "false", client["Cloud"])
	require.Equal(t, "false", client["LDAP"])
	require.Equal(t, "Local fork", client["SkuName"])
	// Stored licenses cannot replace local entitlements in a fork build.
	ps.licenseValue.Store(model.NewTestLicenseWithFalseDefaults())
	require.True(t, *ps.License().Features.CustomPermissionsSchemes)
}

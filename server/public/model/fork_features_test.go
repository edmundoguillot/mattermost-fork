// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestForkFeatureLicense(t *testing.T) {
	l := NewForkFeatureLicense()
	require.True(t, MinimumEnterpriseAdvancedLicense(l))
	require.True(t, l.IsStarted())
	require.False(t, l.IsExpired())
	require.False(t, l.IsCloud())
	require.False(t, l.IsTrialLicense())
	require.False(t, l.IsSeatCountEnforced)
	for _, feature := range []*bool{l.Features.CustomPermissionsSchemes, l.Features.GuestAccounts, l.Features.GuestAccountsPermissions, l.Features.CustomTermsOfService, l.Features.Metrics, l.Features.Elasticsearch, l.Features.AdvancedLogging, l.Features.SharedChannels} {
		require.True(t, *feature)
	}
	for _, feature := range []*bool{l.Features.LDAP, l.Features.SAML, l.Features.Cluster, l.Features.DataRetention, l.Features.Compliance, l.Features.MHPNS, l.Features.Cloud, l.Features.GoogleOAuth, l.Features.Office365OAuth, l.Features.OpenId, l.Features.OutgoingOAuthConnections, l.Features.AutoTranslation, l.Features.EnterprisePlugins, l.Features.FutureFeatures} {
		require.False(t, *feature)
	}
	// Feature objects are independent: callers cannot mutate future entitlements.
	*l.Features.CustomPermissionsSchemes = false
	require.True(t, *NewForkFeatureLicense().Features.CustomPermissionsSchemes)
}

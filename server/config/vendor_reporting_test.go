// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestVendorReportingDisabledOnLoadAndSave(t *testing.T) {
	t.Setenv("MM_LOGSETTINGS_ENABLEDIAGNOSTICS", "true")
	t.Setenv("MM_LOGSETTINGS_ENABLESENTRY", "true")
	t.Setenv("MM_SERVICESETTINGS_ENABLESECURITYFIXALERT", "true")
	t.Setenv("MM_SERVICESETTINGS_SPLITKEY", "vendor-key")
	t.Setenv("MM_ANNOUNCEMENTSETTINGS_ADMINNOTICESENABLED", "true")
	t.Setenv("MM_ANNOUNCEMENTSETTINGS_USERNOTICESENABLED", "true")
	t.Setenv("MM_PLUGINSETTINGS_ENABLEREMOTEMARKETPLACE", "true")

	cfg := &model.Config{}
	cfg.SetDefaults()
	backing, err := NewMemoryStoreWithOptions(&MemoryStoreOptions{InitialConfig: cfg})
	require.NoError(t, err)
	store, err := NewStoreFromBacking(backing, nil, false)
	require.NoError(t, err)
	t.Cleanup(func() { store.Close() })

	assertDisabled := func(cfg *model.Config) {
		t.Helper()
		require.False(t, *cfg.LogSettings.EnableDiagnostics)
		require.False(t, *cfg.LogSettings.EnableSentry)
		require.False(t, *cfg.ServiceSettings.EnableSecurityFixAlert)
		require.Empty(t, *cfg.ServiceSettings.SplitKey)
		require.False(t, *cfg.ServiceSettings.EnableAWSMetering)
		require.False(t, *cfg.AnnouncementSettings.AdminNoticesEnabled)
		require.False(t, *cfg.AnnouncementSettings.UserNoticesEnabled)
		require.False(t, *cfg.PluginSettings.EnableRemoteMarketplace)
		require.False(t, cfg.PluginSettings.PluginStates[model.PluginIdNPS].Enable)
		require.Empty(t, *cfg.CloudSettings.PreviewModalBucketURL)
	}
	assertDisabled(store.Get())

	// An admin save must not re-enable reporting, even with environment overrides.
	cfg.PluginSettings.PluginStates[model.PluginIdNPS] = &model.PluginState{Enable: true}
	*cfg.ServiceSettings.EnableAWSMetering = true
	*cfg.CloudSettings.PreviewModalBucketURL = "https://example.com/vendor-content"
	_, saved, err := store.Set(cfg)
	require.NoError(t, err)
	assertDisabled(saved)
	require.NoError(t, store.Load())
	assertDisabled(store.Get())
}

func TestVendorPushProxyDisabled(t *testing.T) {
	for _, endpoint := range []string{"https://push-test.mattermost.com", "https://PUSH.mattermost.com:443", "https://push.example.com"} {
		t.Run(endpoint, func(t *testing.T) {
			cfg := &model.Config{}
			cfg.SetDefaults()
			*cfg.EmailSettings.SendPushNotifications = true
			*cfg.EmailSettings.PushNotificationServer = endpoint
			disableVendorReporting(cfg)
			if endpoint == "https://push.example.com" {
				require.True(t, *cfg.EmailSettings.SendPushNotifications)
				require.Equal(t, endpoint, *cfg.EmailSettings.PushNotificationServer)
			} else {
				require.False(t, *cfg.EmailSettings.SendPushNotifications)
				require.Empty(t, *cfg.EmailSettings.PushNotificationServer)
			}
		})
	}
}

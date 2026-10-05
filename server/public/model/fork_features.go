// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

// BuildForkFeatures enables local entitlements for implementations included in
// this fork. The Makefile sets it for fork builds; ordinary tests retain the
// upstream licensing behavior. This is not a signed Mattermost license.
var BuildForkFeatures = "false"

// NewForkFeatureLicense adapts local entitlements to the existing feature checks.
// Closed-source services and vendor-hosted services deliberately remain disabled.
func NewForkFeatureLicense() *License {
	features := &Features{
		Users:                     new(2147483647),
		FutureFeatures:            new(false),
		MFA:                       new(true),
		Metrics:                   new(true),
		Elasticsearch:             new(true),
		Announcement:              new(true),
		ThemeManagement:           new(true),
		EmailNotificationContents: new(true),
		CustomPermissionsSchemes:  new(true),
		CustomTermsOfService:      new(true),
		GuestAccounts:             new(true),
		GuestAccountsPermissions:  new(true),
		LockTeammateNameDisplay:   new(true),
		AdvancedLogging:           new(true),
		SharedChannels:            new(true),
		RemoteClusterService:      new(true),
		AutoTranslation:           new(false),
	}
	features.SetDefaults()
	return &License{
		Id:           "00000000000000000000000001",
		IssuedAt:     1,
		StartsAt:     1,
		ExpiresAt:    253402300799000,
		Customer:     &Customer{Name: "Local fork"},
		Features:     features,
		SkuName:      "Local fork",
		SkuShortName: LicenseShortSkuEnterpriseAdvanced,
	}
}

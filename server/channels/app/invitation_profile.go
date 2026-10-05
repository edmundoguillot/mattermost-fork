// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"

	"github.com/mattermost/mattermost/server/public/model"
)

func prepareInvitationProfile(user *model.User) *model.AppError {
	user.FirstName = strings.TrimSpace(model.SanitizeUnicode(user.FirstName))
	user.LastName = strings.TrimSpace(model.SanitizeUnicode(user.LastName))
	if user.FirstName == "" || utf8.RuneCountInString(user.FirstName) > model.UserFirstNameMaxRunes {
		return model.InvalidUserError("first_name", user.Id, user.FirstName)
	}
	if user.LastName == "" || utf8.RuneCountInString(user.LastName) > model.UserLastNameMaxRunes {
		return model.InvalidUserError("last_name", user.Id, user.LastName)
	}
	return nil
}

// Normalize accents and discard spaces and punctuation. Preserve the actual
// names in the profile. Names without Latin characters get a stable fallback.
func invitationUsername(firstName, lastName string) string {
	part := func(name string) string {
		var result strings.Builder
		for _, r := range norm.NFD.String(strings.ToLower(name)) {
			if unicode.Is(unicode.Mn, r) {
				continue
			}
			if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
				result.WriteRune(r)
			}
		}
		if result.Len() == 0 {
			return "user"
		}
		return result.String()
	}
	// Keep both name components even for very long names.
	first := part(firstName)
	last := part(lastName)
	if first[0] < 'a' || first[0] > 'z' {
		first = "u" + first
	}
	if len(first) > 30 {
		first = first[:30]
	}
	if len(last) > 30 {
		last = last[:30]
	}
	return first + "." + last
}

func invitationUsernameCandidate(base string, attempt int) string {
	suffix := ""
	if attempt > 1 {
		suffix = strconv.Itoa(attempt)
	}
	limit := model.UserNameMaxLength - len(suffix)
	if len(base) > limit {
		base = base[:limit]
	}
	return base + suffix
}

package models

import "strings"

const (
	allUsers              = "allUsers"
	allAuthenticatedUsers = "allAuthenticatedUsers"
	prefixDeleted         = "deleted:"
	prefixUser            = "user:"
	prefixGroup           = "group:"
	prefixServiceAccount  = "serviceAccount:"
)

type User string

func (u *User) IsDeleted() bool {
	return strings.HasPrefix(string(*u), "deleted:")
}

func (u *User) IsAllUsers() bool {
	return string(*u) == allUsers
}

func (u *User) IsAllAuthenticatedUsers() bool {
	return string(*u) == allAuthenticatedUsers
}

func (u *User) IsUser() bool {
	return strings.HasPrefix(string(*u), prefixDeleted+prefixUser) || strings.HasPrefix(string(*u), prefixUser)
}

func (u *User) IsGroup() bool {
	return strings.HasPrefix(string(*u), prefixDeleted+prefixGroup) || strings.HasPrefix(string(*u), prefixGroup)
}

func (u *User) IsServiceAccount() bool {
	return strings.HasPrefix(string(*u), prefixDeleted+prefixServiceAccount) || strings.HasPrefix(string(*u), prefixServiceAccount)
}

func (u *User) GetID() string {
	id := string(*u)
	if id == allUsers || id == allAuthenticatedUsers {
		return id
	}

	if strings.HasPrefix(id, prefixDeleted) {
		id = strings.TrimPrefix(id, prefixDeleted)
	}

	if strings.HasPrefix(id, prefixUser) {
		id = strings.TrimPrefix(id, prefixUser)
	}

	if strings.HasPrefix(id, prefixGroup) {
		id = strings.TrimPrefix(id, prefixGroup)
	}

	if strings.HasPrefix(id, prefixServiceAccount) {
		id = strings.TrimPrefix(id, prefixServiceAccount)
	}

	idx := strings.Index(id, "?uid=")
	if idx >= 0 {
		id = id[:idx]
	}

	return id
}

type Users []User

func (u *Users) GetActiveUsers() Users {
	m := map[User]bool{}
	for _, us := range *u {
		if m[us] {
			continue
		}

		if !us.IsAllUsers() && !us.IsAllAuthenticatedUsers() && !us.IsDeleted() {
			m[us] = true
		}
	}

	var result Users
	for us := range m {
		result = append(result, us)
	}

	return result
}

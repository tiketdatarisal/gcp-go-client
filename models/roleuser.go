package models

type RoleUser struct {
	Role  string `json:"role"`
	Users Users  `json:"users"`
}

func (r *RoleUser) GetActiveUsers() Users {
	return r.Users.GetActiveUsers()
}

type RoleUsers []RoleUser

func (r *RoleUsers) GetActiveUsers() Users {
	var users Users
	for _, ru := range *r {
		users = append(users, ru.Users...)
	}

	if users == nil {
		return nil
	}

	return users.GetActiveUsers();
}

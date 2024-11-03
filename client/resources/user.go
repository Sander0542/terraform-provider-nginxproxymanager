package resources

type User struct {
	resource
	Name        string          `json:"name"`
	Nickname    string          `json:"nickname"`
	Email       string          `json:"email"`
	Avatar      string          `json:"avatar"`
	IsDisabled  bool            `json:"is_disabled"`
	Roles       []string        `json:"roles"`
	Permissions UserPermissions `json:"permissions"`
}

type UserCollection []User

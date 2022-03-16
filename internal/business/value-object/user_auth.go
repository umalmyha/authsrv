package valueobj

type UserAuth struct {
	roles  []string
	scopes []string
}

func NewUserAuth(roles, scopes []string) UserAuth {
	var auth UserAuth

	auth.roles = roles
	if auth.roles == nil {
		auth.roles = make([]string, 0)
	}

	auth.scopes = scopes
	if auth.scopes == nil {
		auth.scopes = make([]string, 0)
	}

	return auth
}

func (auth UserAuth) Roles() []string {
	return auth.roles
}

func (auth UserAuth) Scopes() []string {
	return auth.scopes
}

package dto

type NewUser struct {
	Username        string  `json:"username"`
	Email           *string `json:"email"`
	Password        string  `json:"password"`
	ConfirmPassword string  `json:"confirmPassword"`
	IsSuperuser     bool    `json:"-"`
	FirstName       *string `json:"firstName"`
	LastName        *string `json:"lastName"`
	MiddleName      *string `json:"middleName"`
}

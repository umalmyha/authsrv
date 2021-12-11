package dto

type UpdateUser struct {
	Email      *string `json:"email"`
	FirstName  *string `json:"firstName"`
	LastName   *string `json:"lastName"`
	MiddleName *string `json:"middleName"`
}

package scope

type ScopeDto struct {
	Id          string  `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description"`
}

func (s ScopeDto) IsPresent() bool {
	return s.Id != ""
}

type NewScopeDto struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

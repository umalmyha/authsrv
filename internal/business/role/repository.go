package role

type Repository struct {
	uow *unitOfWork
}

func NewRepository(u *unitOfWork) *Repository {
	return &Repository{
		uow: u,
	}
}

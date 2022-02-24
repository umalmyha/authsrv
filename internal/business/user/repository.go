package user

type Repository struct {
	uow *unitOfWork
}

func NewRepository(u *unitOfWork) *Repository {
	return &Repository{
		uow: u,
	}
}

func (repo *Repository) Add(user *User) error {
	return repo.uow.RegisterNew(user)
}

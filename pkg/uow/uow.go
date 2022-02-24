package uow

import "context"

type UnitOfWork[E any] interface {
	RegisterClean(E) error
	RegisterNew(E) error
	RegisterAmended(E) error
	RegisterDeleted(E) error
	Flush(ctx context.Context) error
	Dispose() error
}

type Entitier[E any] interface {
	Key() string
	IsPresent() bool
	IsTheSameAs(E) bool
	Clone() E
}

type UnitOfWorkRepository[E any] interface {
	ById(entity E) error
	Add(entity E) error
	Update(entity E) error
	Remove(entity E) error
}

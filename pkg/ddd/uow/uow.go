package uow

import "context"

type UnitOfWork[E Entitier[E]] interface {
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
	Equal(E) bool
	Clone() E
}

type UnitOfWorkRepository[E Entitier[E], K any] interface {
	ById(K) (E, error)
	Add(entity E) error
	Update(entity E) error
	Remove(entity E) error
}

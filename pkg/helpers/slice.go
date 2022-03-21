package helpers

import (
	"container/list"
	"errors"

	"golang.org/x/exp/constraints"
)

type mapperFn[E any, V any] func(E, int, []E) V
type reducerFn[E any, V any] func(V, E, int, []E) V
type matchFn[E any] func(E, int, []E) bool
type groupKeyFn[E any, V any, K constraints.Ordered] func(E, int, []E) (K, V)
type fromListReducerFn[E any, R any] func(E) R

func Map[E any, V any](in []E, mapper mapperFn[E, V]) []V {
	out := make([]V, 0)
	for i, elem := range in {
		out = append(out, mapper(elem, i, in))
	}
	return out
}

func Reduce[E any, V any](in []E, reducer reducerFn[E, V], out V) V {
	for i, elem := range in {
		out = reducer(out, elem, i, in)
	}
	return out
}

func FindIndex[E any](in []E, matcher matchFn[E]) int {
	index := -1
	for i, elem := range in {
		if matcher(elem, i, in) {
			index = i
			break
		}
	}
	return index
}

func GroupBy[E any, V any, K constraints.Ordered](in []E, keyGrouper groupKeyFn[E, V, K]) map[K][]V {
	groups := make(map[K][]V)
	for i, elem := range in {
		key, value := keyGrouper(elem, i, in)
		if _, found := groups[key]; found {
			groups[key] = append(groups[key], value)
		} else {
			groups[key] = []V{value}
		}
	}
	return groups
}

func FromList[E any](lst *list.List) []E {
	var elems = make([]E, 0)
	for e := lst.Front(); e != nil; e = e.Next() {
		elem, ok := e.Value.(E)
		if !ok {
			panic(errors.New("list elements must have the same type to be converted to a slice"))
		}
		elems = append(elems, elem)
	}
	return elems
}

func FromListWithReducer[E any, R any](lst *list.List, reducerFn fromListReducerFn[E, R]) []R {
	var elems = make([]R, 0)
	for e := lst.Front(); e != nil; e = e.Next() {
		elem, ok := e.Value.(E)
		if !ok {
			panic(errors.New("list elements must have the same type to be converted to a slice"))
		}
		elems = append(elems, reducerFn(elem))
	}
	return elems
}

func ToList[E any](sl []E) *list.List {
	lst := list.New()
	for _, elem := range sl {
		lst.PushBack(elem)
	}
	return lst
}

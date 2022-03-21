package uow

type EntityFinderFn[E Entitier[E]] func() (E, error)

type valueReceiver[E Entitier[E]] struct {
	onReceived func(E)
	value      E
}

func (vr *valueReceiver[E]) Receive() E {
	if vr.isValuePresent(vr.value) {
		vr.onReceived(vr.value)
	}
	return vr.value
}

func (vr *valueReceiver[E]) IfNotPresent(finderFn EntityFinderFn[E]) (E, error) {
	if vr.isValuePresent(vr.value) {
		return vr.value, nil
	}

	value, err := finderFn()
	if err != nil {
		return value, err
	}

	if vr.isValuePresent(value) {
		vr.onReceived(value)
	}
	return value, nil
}

func (vr *valueReceiver[E]) isValuePresent(v E) bool {
	if isNilPtr(v) {
		return false
	}
	return v.IsPresent()
}

package uow

type EntityFinderFn[E Entitier[E]] func() (E, error)

type valueReceiver[E Entitier[E]] struct {
	onReceived func(E)
	value      E
}

func (vr *valueReceiver[E]) Receive() E {
	if vr.isValuePresent() {
		vr.onReceived(vr.value)
	}
	return vr.value
}

func (vr *valueReceiver[E]) IfNotPresent(finderFn EntityFinderFn[E]) (E, error) {
	if vr.isValuePresent() {
		return vr.value, nil
	}

	value, err := finderFn()
	if err != nil {
		return value, err
	}

	vr.onReceived(value)
	return value, nil
}

func (vr *valueReceiver[E]) isValuePresent() bool {
	if isPtrToNil(vr.value) {
		return false
	}
	return vr.value.IsPresent()
}

package args

type iterator struct {
	args      []string
	currIndex int
}

func (iter *iterator) HasNext() bool {
	return iter.currIndex < len(iter.args)
}

func (iter *iterator) Next() string {
	if iter.HasNext() {
		arg := iter.args[iter.currIndex]
		iter.currIndex++
		return arg
	}
	return ""
}

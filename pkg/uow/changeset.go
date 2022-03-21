package uow

import "fmt"

type EntityMatcherFn[E Entitier[E]] func(E) bool

type ChangeSet[E Entitier[E]] struct {
	unchanged map[string]E
	created   map[string]E
	updated   map[string]E
	deleted   map[string]E
}

func NewChangeSet[E Entitier[E]]() *ChangeSet[E] {
	return &ChangeSet[E]{
		unchanged: make(map[string]E),
		created:   make(map[string]E),
		updated:   make(map[string]E),
		deleted:   make(map[string]E),
	}
}

func (es *ChangeSet[E]) Attach(entity E) {
	id := entity.Key()

	if _, found := es.deleted[id]; found {
		return
	}

	if _, found := es.created[id]; found {
		return
	}

	if _, found := es.updated[id]; found {
		return
	}

	es.unchanged[id] = entity.Clone()
}

func (es *ChangeSet[E]) AttachRange(entities ...E) {
	for _, entity := range entities {
		es.Attach(entity)
	}
}

func (es *ChangeSet[E]) Add(entity E) error {
	id := entity.Key()

	if _, found := es.created[id]; found {
		return fmt.Errorf("entity with id %s already exists - was created recently", id)
	}

	if _, found := es.updated[id]; found {
		return fmt.Errorf("entity with id %s already exists - was updated recently", id)
	}

	if _, found := es.unchanged[id]; found {
		return fmt.Errorf("entity with id %s already exists - was found in storage recently", id)
	}

	if entry, found := es.deleted[id]; found {
		delete(es.deleted, id)
		if !entity.Equal(entry) {
			es.updated[id] = entity.Clone()
		} else {
			es.unchanged[id] = entity.Clone()
		}
		return nil
	}

	es.created[id] = entity.Clone()
	return nil
}

func (es *ChangeSet[E]) AddRange(entities ...E) error {
	for _, entity := range entities {
		if err := es.Add(entity); err != nil {
			return err
		}
	}
	return nil
}

func (es *ChangeSet[E]) Remove(entity E) error {
	id := entity.Key()

	if _, found := es.deleted[id]; found {
		return fmt.Errorf("entity with id %s is already removed", id)
	}

	if _, found := es.created[id]; found {
		delete(es.created, id)
		return nil
	}

	if _, found := es.updated[id]; found {
		delete(es.updated, id)
		es.deleted[id] = entity.Clone()
		return nil
	}

	if _, found := es.unchanged[id]; !found {
		return fmt.Errorf("entity with id %s doesn't exist", id)
	}
	delete(es.unchanged, id)

	es.deleted[id] = entity.Clone()
	return nil
}

func (es *ChangeSet[E]) RemoveRange(entities ...E) error {
	for _, entity := range entities {
		if err := es.Remove(entity); err != nil {
			return err
		}
	}
	return nil
}

func (es *ChangeSet[E]) Update(entity E) error {
	id := entity.Key()

	if _, found := es.deleted[id]; found {
		return fmt.Errorf("entity with id %s is already removed", id)
	}

	if entry, found := es.created[id]; found {
		if !entity.Equal(entry) {
			es.created[id] = entity.Clone()
		}
		return nil
	}

	if entry, found := es.updated[id]; found {
		if !entity.Equal(entry) {
			es.updated[id] = entity.Clone()
		}
		return nil
	}

	if entry, found := es.unchanged[id]; found {
		if !entity.Equal(entry) {
			delete(es.unchanged, id)
			es.updated[id] = entity.Clone()
		}
		return nil
	}

	return fmt.Errorf("entity with id %s doesn't exist", id)
}

func (es *ChangeSet[E]) UpdateRange(entities ...E) error {
	for _, entity := range entities {
		if err := es.Update(entity); err != nil {
			return err
		}
	}
	return nil
}

func (es *ChangeSet[E]) FindByKey(entity E) *valueReceiver[E] {
	id := entity.Key()

	if entry, found := es.created[id]; found {
		return es.receiverWithValue(entry)
	}

	if entry, found := es.updated[id]; found {
		return es.receiverWithValue(entry)
	}

	if entry, found := es.unchanged[id]; found {
		return es.receiverWithValue(entry)
	}

	return es.emptyReceiver()
}

func (es *ChangeSet[E]) Find(matcherFn EntityMatcherFn[E]) *valueReceiver[E] {
	if entry, found := es.findMatchedIn(matcherFn, es.created); found {
		return es.receiverWithValue(entry)
	}

	if entry, found := es.findMatchedIn(matcherFn, es.updated); found {
		return es.receiverWithValue(entry)
	}

	if entry, found := es.findMatchedIn(matcherFn, es.unchanged); found {
		return es.receiverWithValue(entry)
	}

	return es.emptyReceiver()
}

func (es *ChangeSet[E]) Filter(matcherFn EntityMatcherFn[E]) []E {
	filtered := make([]E, 0)

	if entries := es.filterMatchedIn(matcherFn, es.created); len(entries) > 0 {
		filtered = append(filtered, entries...)
	}

	if entries := es.filterMatchedIn(matcherFn, es.updated); len(entries) > 0 {
		filtered = append(filtered, entries...)
	}

	if entries := es.filterMatchedIn(matcherFn, es.unchanged); len(entries) > 0 {
		filtered = append(filtered, entries...)
	}

	return filtered
}

func (es *ChangeSet[E]) IsCreated(entity E) bool {
	if _, found := es.created[entity.Key()]; found {
		return true
	}
	return false
}

func (es *ChangeSet[E]) IsUpdated(entity E) bool {
	if _, found := es.updated[entity.Key()]; found {
		return true
	}
	return false
}

func (es *ChangeSet[E]) IsRemoved(entity E) bool {
	if _, found := es.deleted[entity.Key()]; found {
		return true
	}
	return false
}

func (es *ChangeSet[E]) IsUnchanged(entity E) bool {
	if _, found := es.unchanged[entity.Key()]; found {
		return true
	}
	return false
}

func (es *ChangeSet[E]) Exists(entity E) bool {
	return es.IsCreated(entity) || es.IsUpdated(entity) || es.IsUnchanged(entity)
}

func (es *ChangeSet[E]) Created() []E {
	created := make([]E, 0)

	for _, entry := range es.created {
		created = append(created, entry.Clone())
	}

	return created
}

func (es *ChangeSet[E]) Updated() []E {
	changed := make([]E, 0)

	for _, entry := range es.updated {
		changed = append(changed, entry.Clone())
	}

	return changed
}

func (es *ChangeSet[E]) Deleted() []E {
	deleted := make([]E, 0)

	for _, entry := range es.deleted {
		deleted = append(deleted, entry.Clone())
	}

	return deleted
}

func (es *ChangeSet[E]) Clean() []E {
	clean := make([]E, 0)

	for _, entry := range es.unchanged {
		clean = append(clean, entry.Clone())
	}

	return clean
}

func (es *ChangeSet[E]) Delta() ([]E, []E, []E) {
	return es.Created(), es.Updated(), es.Deleted()
}

func (es *ChangeSet[E]) DeltaWithMatched(entities []E, matcherFn EntityMatcherFn[E]) ([]E, []E, []E) {
	matched := es.Filter(matcherFn)
	matchedMap := make(map[string]E)
	for _, entry := range matched {
		matchedMap[entry.Key()] = entry
	}

	created := make([]E, 0)
	updated := make([]E, 0)
	deleted := make([]E, 0)

	for _, entity := range entities {
		if matchedEntry, found := matchedMap[entity.Key()]; found {
			if !entity.Equal(matchedEntry) {
				updated = append(updated, entity.Clone())
			}
			delete(matchedMap, entity.Key())
		} else {
			created = append(created, entity.Clone())
		}
	}

	for _, rmEntry := range matchedMap {
		deleted = append(deleted, rmEntry.Clone())
	}

	return created, updated, deleted
}

func (es *ChangeSet[E]) All() []E {
	created := es.Created()
	updated := es.Updated()
	clean := es.Clean()

	all := append(created, updated...)
	all = append(all, clean...)
	return all
}

func (es *ChangeSet[E]) Cleanup() {
	es.unchanged = make(map[string]E)
	es.created = make(map[string]E)
	es.updated = make(map[string]E)
	es.deleted = make(map[string]E)
}

func (es *ChangeSet[E]) receiverWithValue(value E) *valueReceiver[E] {
	return &valueReceiver[E]{
		value:      value.Clone(),
		onReceived: es.onValueReceived,
	}
}

func (es *ChangeSet[E]) emptyReceiver() *valueReceiver[E] {
	return &valueReceiver[E]{
		onReceived: es.onValueReceived,
	}
}

func (es *ChangeSet[E]) findMatchedIn(matcherFn EntityMatcherFn[E], in map[string]E) (E, bool) {
	for _, elem := range in {
		if matcherFn(elem) {
			return elem, true
		}
	}
	var entry E
	return entry, false
}

func (es *ChangeSet[E]) filterMatchedIn(matcherFn EntityMatcherFn[E], in map[string]E) []E {
	matched := make([]E, 0)
	for _, elem := range in {
		if matcherFn(elem) {
			matched = append(matched, elem)
		}
	}
	return matched
}

func (es *ChangeSet[E]) onValueReceived(entity E) {
	es.unchanged[entity.Key()] = entity
}

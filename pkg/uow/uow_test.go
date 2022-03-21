package uow

import (
	"testing"

	"golang.org/x/exp/slices"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

type Person struct {
	id   string
	name string
	age  int
}

func (p Person) Key() string {
	return p.id
}

func (p Person) IsPresent() bool {
	return p.id != ""
}

func (p Person) Equal(other Person) bool {
	return p.id == other.id &&
		p.name == other.name &&
		p.age == other.age
}

func (p Person) Clone() Person {
	return p
}

// 1 clean
// 3 created
// 2 updated
// 2 deleted
func basicChangesetData(fatalFn func(error)) *ChangeSet[Person] {
	cleanPersons := []Person{
		{id: "ebd28805-b914-4f30-aebd-e3caa9f2a5ee", name: "Henry", age: 40},
	}

	createdPersons := []Person{
		{id: "9a8c5f3a-2b3b-4841-8b3f-633735efd349", name: "Liam", age: 20},
		{id: "8e54313d-6943-4767-a67b-fdf25d5526aa", name: "Noah", age: 25},
		{id: "068f4c71-c0ab-4f36-96d0-8736498857cc", name: "Oliver", age: 20},
	}

	updatedPersons := []Person{
		{id: "2bba0cd7-e1bf-4f42-9c83-6879113ef918", name: "Elijah", age: 18},
		{id: "5cbeb127-55bd-4514-af69-16aa07b720dc", name: "Lucas", age: 18},
	}

	deletedPersons := []Person{
		{id: "ddfb2e65-cd80-4489-b526-d1484d7e9f67", name: "Joseph", age: 34},
		{id: "19280583-d137-4094-b605-3a7375d3e791", name: "Thomas", age: 22},
	}

	chs := NewChangeSet[Person]()

	chs.AttachRange(cleanPersons...)

	if err := chs.AddRange(createdPersons...); err != nil {
		fatalFn(err)
	}

	chs.AttachRange(updatedPersons...)
	updatedPersons[0].age = 35
	updatedPersons[1].age = 35
	if err := chs.UpdateRange(updatedPersons...); err != nil {
		fatalFn(err)
	}

	chs.AttachRange(deletedPersons...)
	if err := chs.RemoveRange(deletedPersons...); err != nil {
		fatalFn(err)
	}

	return chs
}

func TestChangeSetUpdate(t *testing.T) {
	person := Person{id: "f897829d-4d01-45f5-b510-1f730eeffd1f", name: "Antony", age: 32}
	newPerson := Person{id: "03f52452-5bef-4328-861d-864bc3c35482", name: "David", age: 20}
	deletePerson := Person{id: "36fea0d3-13e3-4011-b4b7-248ee35a82d8", name: "James", age: 25}
	createdAndModifiedPerson := Person{id: "83036130-8869-4cd3-b5e7-94ac2ee9737d", name: "Casey", age: 45}

	chs := NewChangeSet[Person]()
	t.Log("Given the need to test ChangeSet behaviour on update")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen Person is not modified", testId)
		{
			chs.Attach(person)
			if err := chs.Update(person); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on update: %v", failed, err)
			}

			if len(chs.Updated()) > 0 {
				t.Fatalf("\t%s\tPerson wasn't updated, but moved to updated state", failed)
			}
			t.Logf("\t%s\tPerson mustn't move to updated state", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is modified", testId)
		{
			person.age = 33
			if err := chs.Update(person); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on update: %v", failed, err)
			}

			if len(chs.Updated()) == 0 {
				t.Fatalf("\t%s\tPerson was updated, but wan't moved to updated state", failed)
			}
			t.Logf("\t%s\tPerson must move to updated state", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is not modified, but already in updated state", testId)
		{
			if err := chs.Update(person); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on update: %v", failed, err)
			}

			if len(chs.Updated()) > 1 {
				t.Fatalf("\t%s\tPerson wasn't updated recently, but duplicated in updated state", failed)
			}
			t.Logf("\t%s\tPerson must be in updated state in singular instance", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is modified and already in updated state", testId)
		{
			stateBeforeUpdate := chs.FindByKey(person).Receive()

			person.age = 90
			if err := chs.Update(person); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on update: %v", failed, err)
			}

			stateAfterUpdate := chs.FindByKey(person).Receive()

			if stateBeforeUpdate.Equal(stateAfterUpdate) {
				t.Fatalf("\t%s\tPerson state has been changed since the last update call, but wasn't modified in ChangeSet", failed)
			}

			t.Logf("\t%s\tPerson must be modified in updated state, since there are changes happened", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person doesn't exist in ChangeSet", testId)
		{
			err := chs.Update(newPerson)
			if err == nil {
				t.Fatalf("\t%s\tNew non-existing entity is requested for update, but error is not raised", failed)
			}
			t.Logf("\t%s\tCorresponding error is raised: %v", success, err)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is deleted, but trying to update", testId)
		{
			chs.Attach(deletePerson)
			if err := chs.Remove(deletePerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on remove: %v", failed, err)
			}

			err := chs.Update(deletePerson)
			if err == nil {
				t.Fatalf("\t%s\tDeleted entity is requested for update, but error is not raised", failed)
			}
			t.Logf("\t%s\tCorresponding error is raised: %v", success, err)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is created and modified in the meantime", testId)
		{
			if err := chs.Add(createdAndModifiedPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on add: %v", failed, err)
			}

			stateBeforeUpdate := chs.FindByKey(createdAndModifiedPerson).Receive()

			createdAndModifiedPerson.age = 19
			if err := chs.Update(createdAndModifiedPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on add: %v", failed, err)
			}

			stateAfterUpdate := chs.FindByKey(createdAndModifiedPerson).Receive()

			if stateBeforeUpdate.Equal(stateAfterUpdate) {
				t.Fatalf("\t%s\tState hasn't changed after update", failed)
			}

			if !chs.IsCreated(createdAndModifiedPerson) {
				t.Fatalf("\t%s\tPerson was created and modified in the meantime, so still must stay in created state", failed)
			}

			t.Logf("\t%s\tPerson must be modified in created state", success)
		}
	}
}

func TestChangeSetAdd(t *testing.T) {
	createdPerson := Person{id: "d12c3047-7b6e-4994-9b80-7ec36d38983c", name: "Steven", age: 25}
	updatedPerson := Person{id: "20f12177-d125-45fc-a721-3df6294e13de", name: "James", age: 40}
	cleanPerson := Person{id: "13407a69-c067-4a98-9998-8da260a87395", name: "Lester", age: 18}
	deletedPerson := Person{id: "87ed3005-0012-4158-8941-d983ac49a14a", name: "Kevin", age: 23}
	delAndUpdPerson := Person{id: "433fda1d-09d3-407e-8b16-25c60f057a11", name: "Arnold", age: 70}

	chs := NewChangeSet[Person]()
	t.Log("Given the need to test ChangeSet behaviour on add")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen Person is new", testId)
		{
			if err := chs.Add(createdPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on add: %v", failed, err)
			}

			if len(chs.Created()) == 0 {
				t.Fatalf("\t%s\tPerson is new, but not moved to created state", failed)
			}
			t.Logf("\t%s\tPerson must be moved to created state", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is already in created state", testId)
		{
			err := chs.Add(createdPerson)
			if err == nil {
				t.Fatalf("\t%s\tPerson is already in created state, but no error raised", failed)
			}
			t.Logf("\t%s\tCorresponding error is raised: %v", success, err)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is already in updated state", testId)
		{
			chs.Attach(updatedPerson)

			updatedPerson.age = 50
			if err := chs.Update(updatedPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on update: %v", failed, err)
			}

			err := chs.Add(updatedPerson)
			if err == nil {
				t.Fatalf("\t%s\tPerson is already in updated state, but no error raised", failed)
			}
			t.Logf("\t%s\tCorresponding error is raised: %v", success, err)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is in unchanged state", testId)
		{
			chs.Attach(cleanPerson)
			err := chs.Add(cleanPerson)
			if err == nil {
				t.Fatalf("\t%s\tPerson is exist in unchanged state, but no error raised", failed)
			}
			t.Logf("\t%s\tCorresponding error is raised: %v", success, err)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is in deleted state, but was added afterwards", testId)
		{
			chs.Attach(deletedPerson)
			if err := chs.Remove(deletedPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on delete: %v", failed, err)
			}

			if err := chs.Add(deletedPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on delete: %v", failed, err)
			}

			if len(chs.Deleted()) != 0 || !chs.IsUnchanged(deletedPerson) {
				t.Fatalf("\t%s\tPerson was removed and added aftrewards with no modification, so must be in unchanged state", failed)
			}
			t.Logf("\t%s\tPerson must be in unchanged state", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is in deleted state, but was modified and added afterwards", testId)
		{
			chs.Attach(delAndUpdPerson)
			if err := chs.Remove(delAndUpdPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on delete: %v", failed, err)
			}

			delAndUpdPerson.age = 80
			if err := chs.Add(delAndUpdPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on delete: %v", failed, err)
			}

			if len(chs.Deleted()) != 0 || !chs.IsUpdated(delAndUpdPerson) {
				t.Fatalf("\t%s\tPerson was removed, modified and added aftrewards, so must be in updated state", failed)
			}
			t.Logf("\t%s\tPerson must be in updated state", success)
		}
	}
}

func TestChangesetRemove(t *testing.T) {
	deletePerson := Person{id: "5a90ed6a-73c4-4c03-ae33-b6a25ce1918b", name: "William", age: 15}
	createdPerson := Person{id: "697972bb-19a9-42a6-a939-6c85333f9511", name: "Ryan", age: 20}
	cleanPerson := Person{id: "e7332ef5-050a-4bdb-a9ef-8e16020cb154", name: "John", age: 24}
	updatedPerson := Person{id: "87b23ea2-a07e-441c-a962-ca1d308179e9", name: "Albert", age: 33}

	chs := NewChangeSet[Person]()

	t.Log("Given the need to test ChangeSet behaviour on remove")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen Person is in unchanged state", testId)
		{
			chs.Attach(deletePerson)
			if err := chs.Remove(deletePerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on delete: %v", failed, err)
			}

			if len(chs.Deleted()) == 0 {
				t.Fatalf("\t%s\tPerson was removed, but wasn't moved to deleted state", failed)
			}
			t.Logf("\t%s\tPerson must be in deleted state", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is already in deleted state", testId)
		{
			err := chs.Remove(deletePerson)
			if err == nil {
				t.Fatalf("\t%s\tPerson already in deleted state, but error wasn't raised", failed)
			}
			t.Logf("\t%s\tCorresponding error is raised: %v", success, err)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is already in created state", testId)
		{
			if err := chs.Add(createdPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on add: %v", failed, err)
			}

			if err := chs.Remove(createdPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on remove: %v", failed, err)
			}

			if chs.IsCreated(createdPerson) || chs.IsRemoved(createdPerson) {
				t.Fatalf("\t%s\tPerson was created and removed immediately, so is not tracked anymore", failed)
			}
			t.Logf("\t%s\tPerson must not be tracked anymore", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is not tracked", testId)
		{
			err := chs.Remove(cleanPerson)
			if err == nil {
				t.Fatalf("\t%s\tPerson is untracked, but no error raised", failed)
			}
			t.Logf("\t%s\tCorresponding error is raised: %v", success, err)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is in updated state", testId)
		{
			chs.Attach(updatedPerson)
			updatedPerson.age = 40
			if err := chs.Update(updatedPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on update: %v", failed, err)
			}

			if err := chs.Remove(updatedPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on remove: %v", failed, err)
			}

			if chs.IsUpdated(updatedPerson) || !chs.IsRemoved(updatedPerson) {
				t.Fatalf("\t%s\tPerson was updated and removed then, so must be in deleted state", failed)
			}
			t.Logf("\t%s\tPerson must be in deleted state", success)
		}
	}
}

func TestChangeSetAttach(t *testing.T) {
	cleanPerson := Person{id: "b9507aa9-74ec-4046-b9cd-95b007211fe7", name: "Michael", age: 45}
	createdPerson := Person{id: "6449c882-c870-440d-9a03-5d900d282b5f", name: "Robert", age: 32}
	updatedPerson := Person{id: "be006ab5-b868-4276-8bb4-1d078294bbb8", name: "Nicolas", age: 32}
	deletedPerson := Person{id: "9ed80689-7270-4d62-8063-12086ef72cc8", name: "Francis", age: 32}

	chs := NewChangeSet[Person]()

	t.Log("Given the need to test ChangeSet behaviour on attach")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen Person is not tracked", testId)
		{
			chs.Attach(cleanPerson)
			if len(chs.Clean()) == 0 {
				t.Fatalf("\t%s\tPerson was untracked, but wasn't moved to unchanged state", failed)
			}
			t.Logf("\t%s\tPerson must be in unchanged state", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is in created state", testId)
		{
			if err := chs.Add(createdPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on add: %v", failed, err)
			}

			chs.Attach(createdPerson)
			if chs.IsUnchanged(createdPerson) && !chs.IsCreated(createdPerson) {
				t.Fatalf("\t%s\tPerson was in created state, but moved to unchanged afterwards", failed)
			}
			t.Logf("\t%s\tPerson must stay in created state", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is in updated state", testId)
		{
			chs.Attach(updatedPerson)

			updatedPerson.age = 88
			if err := chs.Update(updatedPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on update: %v", failed, err)
			}

			chs.Attach(updatedPerson)
			if chs.IsUnchanged(updatedPerson) && !chs.IsUpdated(updatedPerson) {
				t.Fatalf("\t%s\tPerson was in updated state, but moved to unchanged afterwards", failed)
			}

			t.Logf("\t%s\tPerson must stay in updated state", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is in updated state", testId)
		{
			chs.Attach(deletedPerson)
			if err := chs.Remove(deletedPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on remove: %v", failed, err)
			}

			chs.Attach(deletedPerson)
			if chs.IsUnchanged(deletedPerson) && !chs.IsRemoved(deletedPerson) {
				t.Fatalf("\t%s\tPerson was in deleted state, but moved to unchanged afterwards", failed)
			}

			t.Logf("\t%s\tPerson must stay in deleted state", success)
		}
	}
}

func TestChangeSetFindByKey(t *testing.T) {
	createdPerson := Person{id: "733cc08d-094d-4dae-8d8b-37fbb36bc5e9", name: "Benjamin", age: 12}
	updatedPerson := Person{id: "275450af-5eac-47c6-b1b0-cf7f89676516", name: "Carl", age: 22}
	cleanPerson := Person{id: "77b8bb20-9f3b-4b0d-bafe-a2c2c7cbeb83", name: "Max", age: 45}
	deletedPerson := Person{id: "d7d75061-d5c8-456e-8bf8-faa0e3058123", name: "Wincent", age: 21}
	notTrackedPerson := Person{id: "c5c8e893-e785-433c-8cb1-cb908f58d320", name: "Alan", age: 43}

	chs := NewChangeSet[Person]()

	t.Log("Given the need to test ChangeSet find by id behaviour")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen Person is in created state", testId)
		{
			if err := chs.Add(createdPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on add: %v", failed, err)
			}

			p := chs.FindByKey(createdPerson).Receive()
			if !p.IsPresent() {
				t.Fatalf("\t%s\tPerson wasn't found even though present in ChangeSet", failed)
			}
			t.Logf("\t%s\tPerson must be found", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is in updated state", testId)
		{
			chs.Attach(updatedPerson)

			updatedPerson.age = 25
			if err := chs.Update(updatedPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on update: %v", failed, err)
			}

			p := chs.FindByKey(updatedPerson).Receive()
			if !p.IsPresent() {
				t.Fatalf("\t%s\tPerson wasn't found even though present in ChangeSet", failed)
			}
			t.Logf("\t%s\tPerson must be found", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is in unchanged state", testId)
		{
			chs.Attach(cleanPerson)
			p := chs.FindByKey(cleanPerson).Receive()
			if !p.IsPresent() {
				t.Fatalf("\t%s\tPerson wasn't found even though present in ChangeSet", failed)
			}
			t.Logf("\t%s\tPerson must be found", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is in deleted state", testId)
		{
			chs.Attach(deletedPerson)

			if err := chs.Remove(deletedPerson); err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on remove: %v", failed, err)
			}

			p := chs.FindByKey(deletedPerson).Receive()
			if p.IsPresent() {
				t.Fatalf("\t%s\tPerson was found even though in deleted state", failed)
			}
			t.Logf("\t%s\tPerson must not be found", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person is not tracked", testId)
		{
			p := chs.FindByKey(notTrackedPerson).Receive()
			if p.IsPresent() {
				t.Fatalf("\t%s\tPerson was found even though not tracked", failed)
			}
			t.Logf("\t%s\tPerson must not be found", success)
		}
	}
}

func TestChangeSetFind(t *testing.T) {
	chs := basicChangesetData(func(err error) {
		t.Fatalf("\t%s\tUnexpected error occurred: %v", failed, err)
	})

	t.Log("Given the need to test ChangeSet find behaviour")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen Person Noah is in created state", testId)
		{
			p := chs.Find(func(person Person) bool {
				return person.name == "Noah"
			}).Receive()

			if !p.IsPresent() {
				t.Fatalf("\t%s\tPerson wasn't found even though present in ChangeSet", failed)
			}
			t.Logf("\t%s\tPerson Noah must be found by name", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person Elijah and Lucas is in updated state", testId)
		{
			p := chs.Find(func(person Person) bool {
				return person.age == 35
			}).Receive()

			if !p.IsPresent() {
				t.Fatalf("\t%s\tPerson wasn't found even though present in ChangeSet", failed)
			}
			t.Logf("\t%s\tPerson Elijah must be found by age=35 (was first in the slice, so was added first)", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person Elijah and Lucas is in updated state", testId)
		{
			p := chs.Find(func(person Person) bool {
				return person.age == 35
			}).Receive()

			if !p.IsPresent() {
				t.Fatalf("\t%s\tPerson wasn't found even though present in ChangeSet", failed)
			}

			if p.name != "Elijah" && p.name != "Lucas" {
				t.Fatalf("\t%s\tPerson %s was found, Elijah or Lucas must be found", failed, p.name)
			}
			t.Logf("\t%s\tPerson Elijah or Lucas must be found by age=35", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person Henry is in unchanged state", testId)
		{
			p := chs.Find(func(person Person) bool {
				return person.name == "Henry"
			}).Receive()

			if !p.IsPresent() {
				t.Fatalf("\t%s\tPerson wasn't found even though present in ChangeSet", failed)
			}
			t.Logf("\t%s\tPerson Henry must be found by name", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person Thomas and Joseph are in deleted state", testId)
		{
			p := chs.Find(func(person Person) bool {
				return person.name == "Thomas" || person.name == "Joseph"
			}).Receive()

			if p.IsPresent() {
				t.Fatalf("\t%s\tThomas and Joseph are in deleted state even though were found", failed)
			}
			t.Logf("\t%s\tPerson wasn't found", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Person Charles is not tracked", testId)
		{
			p := chs.Find(func(person Person) bool {
				return person.name == "Charles"
			}).Receive()

			if p.IsPresent() {
				t.Fatalf("\t%s\tPerson Charles is untracked even though was found", failed)
			}
			t.Logf("\t%s\tPerson Charles wasn't found", success)
		}
	}
}

func TestChangeSetFilter(t *testing.T) {
	chs := basicChangesetData(func(err error) {
		t.Fatalf("\t%s\tUnexpected error occurred: %v", failed, err)
	})

	t.Log("Given the need to test ChangeSet find behaviour")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen Henry, Noah and Joseph with age > 30 in ChangeSet", testId)
		{
			persons := chs.Filter(func(person Person) bool {
				return person.age > 30
			})

			if len(persons) != 3 {
				t.Fatalf("\t%s\tHenry, Elijah and Lucas are older than 30, but not all of them were found", failed)
			}
			t.Logf("\t%s\tHenry, Elijah and Lucas were found", success)
		}
	}
}

func TestChangeSetExists(t *testing.T) {
	chs := basicChangesetData(func(err error) {
		t.Fatalf("\t%s\tUnexpected error occurred: %v", failed, err)
	})

	t.Log("Given the need to test ChangeSet exists behaviour")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen Liam exist in ChangeSet in created state", testId)
		{
			liam := Person{id: "9a8c5f3a-2b3b-4841-8b3f-633735efd349", name: "Liam", age: 20}
			if !chs.Exists(liam) {
				t.Fatalf("\t%s\tPerson Liam wasn't found even though has created state", failed)
			}
			t.Logf("\t%s\tLiam existence verified", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Lucas exist in ChangeSet in updated state", testId)
		{
			lucas := Person{id: "5cbeb127-55bd-4514-af69-16aa07b720dc", name: "Lucas", age: 18}
			if !chs.Exists(lucas) {
				t.Fatalf("\t%s\tPerson Lucas wasn't found even though has updated state", failed)
			}
			t.Logf("\t%s\tLucas existence verified", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Henry exist in ChangeSet in unchanged state", testId)
		{
			henry := Person{id: "ebd28805-b914-4f30-aebd-e3caa9f2a5ee", name: "Henry", age: 40}
			if !chs.Exists(henry) {
				t.Fatalf("\t%s\tPerson Henry wasn't found even though has unchanged state", failed)
			}
			t.Logf("\t%s\tHenry existence verified", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen Thomas exist in ChangeSet in deleted state", testId)
		{
			thomas := Person{id: "19280583-d137-4094-b605-3a7375d3e791", name: "Thomas", age: 22}
			if chs.Exists(thomas) {
				t.Fatalf("\t%s\tPerson Thomas wasn found even though has deleted state", failed)
			}
			t.Logf("\t%s\tThomas must not be found", success)
		}
	}
}

func TestChangeSetDelta(t *testing.T) {
	chs := basicChangesetData(func(err error) {
		t.Fatalf("\t%s\tUnexpected error occurred: %v", failed, err)
	})

	t.Log("Given the need to test ChangeSet delta behaviour")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen ChangeSet has 1 clean, 3 created, 2 updated and 2 deleted Persons", testId)
		{
			created, updated, deleted := chs.Delta()
			if len(created) != 3 {
				t.Fatalf("\t%s\tThere were 3 Persons in created state, but got %d in result", failed, len(created))
			}

			if len(updated) != 2 {
				t.Fatalf("\t%s\tThere were 2 Persons in updated state, but got %d in result", failed, len(updated))
			}

			if len(deleted) != 2 {
				t.Fatalf("\t%s\tThere were 2 Persons in deleted state, but got %d in result", failed, len(deleted))
			}

			t.Logf("\t%s\tThere were 3 created, 2 updated and 2 deleted entities found in result", success)
		}
	}
}

func TestChangeSetCleanup(t *testing.T) {
	chs := basicChangesetData(func(err error) {
		t.Fatalf("\t%s\tUnexpected error occurred: %v", failed, err)
	})

	t.Log("Given the need to test ChangeSet cleanup behaviour")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen ChangeSet has 1 clean, 3 created, 2 updated and 2 deleted Persons", testId)
		{
			chs.Cleanup()

			if len(chs.Created()) != 0 {
				t.Fatalf("\t%s\tChangeset was cleanuped, but has %d entries in created state", failed, len(chs.Created()))
			}

			if len(chs.Updated()) != 0 {
				t.Fatalf("\t%s\tChangeset was cleanuped, but has %d entries in updated state", failed, len(chs.Updated()))
			}

			if len(chs.Deleted()) != 0 {
				t.Fatalf("\t%s\tChangeset was cleanuped, but has %d entries in deleted state", failed, len(chs.Deleted()))
			}

			t.Logf("\t%s\tChangeSet is empty", success)
		}
	}
}

func TestChangeSetAll(t *testing.T) {
	chs := basicChangesetData(func(err error) {
		t.Fatalf("\t%s\tUnexpected error occurred: %v", failed, err)
	})

	t.Log("Given the need to test ChangeSet all behaviour")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen ChangeSet has 1 clean, 3 created, 2 updated and 2 deleted Persons", testId)
		{
			persons := chs.All()
			if len(persons) != 6 {
				t.Fatalf("\t%s\tThere were 1 clean, 3 created, 2 updated, so 1 + 3 + 2 = 6, but got %d in result", failed, len(persons))
			}

			t.Logf("\t%s\tChangeSet has 6 entries in total", success)
		}
	}
}

func TestChangeSetDeltaWithMatched(t *testing.T) {
	chs := basicChangesetData(func(err error) {
		t.Fatalf("\t%s\tUnexpected error occurred: %v", failed, err)
	})

	t.Log("Given the need to test ChangeSet with delta matched behaviour")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen ChangeSet has 3 created, 2 updated and 2 deleted Persons and after adjustment state there 2 new entries and 1 deleted in compare slice", testId)
		{
			created, updated, _ := chs.Delta()
			clean := chs.Clean()

			created = append(created, Person{id: "c650b012-5852-4dd5-abb4-e215c96f0891", name: "Walter", age: 87})
			updated[0].age = 46
			clean = slices.Delete(clean, 0, 1)

			all := append(created, updated...)
			all = append(all, clean...)

			deltaCreated, deltaUpdated, deltaDeleted := chs.DeltaWithMatched(all, func(_ Person) bool {
				return true
			})

			if len(deltaCreated) != 1 {
				t.Fatalf("\t%s\tThere was 1 Person entry created, nevertheless %d created entries were determined", failed, len(deltaCreated))
			}

			if len(deltaUpdated) != 1 {
				t.Fatalf("\t%s\tThere was 1 Person entry updated, nevertheless %d updated entries were determined", failed, len(deltaUpdated))
			}

			if len(deltaDeleted) != 1 {
				t.Fatalf("\t%s\tThere was 1 Person entry deleted, nevertheless %d deleted entries were determined", failed, len(deltaDeleted))
			}

			t.Logf("\t%s\tChangeSet has delta with 1 created, updated and deleted entries in total", success)
		}
	}
}

type Pet struct {
	id   string
	name string
}

func (p *Pet) Key() string {
	return p.id
}

func (p *Pet) IsPresent() bool {
	return p.id != ""
}

func (p *Pet) Equal(other *Pet) bool {
	return p.id == other.id &&
		p.name == other.name
}

func (p *Pet) Clone() *Pet {
	return &Pet{
		id:   p.id,
		name: p.name,
	}
}

func TestChangeSetWithPtr(t *testing.T) {
	chs := NewChangeSet[*Pet]()

	chs.Add(&Pet{id: "bc991fdd-9ee6-4726-89f8-2a96cdb21d2a", name: "Milo"})
	chs.Add(&Pet{id: "8e286260-c73d-4514-b941-ff2665f07498", name: "Buddy"})

	rocky := &Pet{id: "95e151db-ef48-4d94-a161-dd259a096379", name: "Rocky"}
	teddy := &Pet{id: "1fd01ce3-eea2-4f87-9923-4ee022ec4a92", name: "Teddy"}
	tucker := &Pet{id: "71756bd7-0db4-4945-8ef9-eb188b426fdb", name: "Tucker"}

	chs.AttachRange(rocky, teddy, tucker)

	teddy.name = "Zoe"
	if err := chs.Update(teddy); err != nil {
		t.Fatalf("\t%s\tUnexpected error occurred on update: %v", failed, err)
	}

	if err := chs.Remove(rocky); err != nil {
		t.Fatalf("\t%s\tUnexpected error occurred on remove: %v", failed, err)
	}

	t.Log("Given the need to test ChangeSet with struct pointer")
	{
		testId := 1
		t.Logf("\tTest %d:\tWhen ChangeSet has 2 created, 1 updated, 1 deleted and 1 clean entries and delta called", testId)
		{
			created, updated, deleted := chs.Delta()
			clean := chs.Clean()
			if len(created) != 2 || len(updated) != 1 || len(deleted) != 1 || len(clean) != 1 {
				t.Fatalf(
					"\t%s\tThere must be 2 created, 1 deleted, 1 updated and 1 clean entries in changeset, got %d created, %d updated, %d deleted, %d clean",
					failed,
					len(created),
					len(updated),
					len(deleted),
					len(clean),
				)
			}
			t.Logf("\t%s\tThere must be 2 created, 1 deleted, 1 updated and 1 clean entries in changeset", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen ChangeSet has Pet with corresponding id", testId)
		{
			pet := chs.FindByKey(tucker).Receive()
			if pet == nil || !pet.IsPresent() {
				t.Fatalf("\t%s\tThere must be entry found, but it wasn't", failed)
			}

			if !chs.IsUnchanged(pet) {
				t.Fatalf("\t%s\tFound entry must be in unchanged state, but it is not", failed)
			}

			t.Logf("\t%s\tCorresponding Pet is found", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen ChangeSet doesn't have Pet with corresponding id - value receiver returns result", testId)
		{
			jasper := &Pet{id: "ac2f9f5f-058d-40a5-82d5-af0456e6ab50", name: "Jasper"}

			pet, err := chs.FindByKey(jasper).IfNotPresent(func() (*Pet, error) { return jasper, nil })
			if err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on find by id: %v", failed, err)
			}

			if pet == nil {
				t.Fatalf("\t%s\tThere must be entry found from receiver, but it wasn't", failed)
			}

			if !pet.Equal(jasper) {
				t.Fatalf("\t%s\tPet must be Jasper, but it is %s", failed, pet.name)
			}

			if !chs.IsUnchanged(pet) {
				t.Fatalf("\t%s\tFound entry must be added to unchanged state, but it is not", failed)
			}

			t.Logf("\t%s\tCorresponding Pet is returned from value receiver", success)
		}

		testId++
		t.Logf("\tTest %d:\tWhen ChangeSet doesn't have Pet with corresponding id - value receiver has no result", testId)
		{
			leo := &Pet{id: "ada2575a-9b68-47e6-8919-fa35843b4e5b", name: "Leo"}
			cleanBeforeReceive := len(chs.Clean())

			pet, err := chs.FindByKey(leo).IfNotPresent(func() (*Pet, error) { return nil, nil })
			if err != nil {
				t.Fatalf("\t%s\tUnexpected error occurred on find by id: %v", failed, err)
			}

			if pet != nil {
				t.Fatalf("\t%s\tThere must be no entry found from receiver, but it was", failed)
			}

			cleanAfterReceive := len(chs.Clean())

			if cleanBeforeReceive != cleanAfterReceive {
				t.Fatalf("\t%s\tSince there was no entry found there must not be any new entries in unchanged state, but something was added", failed)
			}

			t.Logf("\t%s\tValue nil is returned from value receiver", success)
		}
	}
}

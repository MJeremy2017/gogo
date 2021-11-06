package main

import "testing"

func TestSearch(t *testing.T) {
	t.Run("known word", func(t *testing.T) {
		dict := Dictionary{"test": "this is a test"}

		got, _ := dict.Search("test")
		want := "this is a test"
		assertStrings(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		dict := Dictionary{"test": "is a test"}

		_, err := dict.Search("something")
		want := ErrKeyNotFound
		if err == nil {
			t.Fatal("expect an error but not find")
		}
		assertError(t, err, want)
	})

}

func TestAdd(t *testing.T) {
	t.Run("Add new word", func(t *testing.T) {
		dict := Dictionary{}
		key := "test"
		value := "this is a test"
		dict.Add(key, value)
		assertDefinition(t, dict, key, value)
	})

	t.Run("Add existing word", func(t *testing.T) {
		dict := Dictionary{}
		key := "test"
		value := "this is a test"
		dict[key] = value
		err := dict.Add(key, value)

		assertError(t, err, ErrKeyExisted)
		assertDefinition(t, dict, key, value)
	})
	
}

func TestUpdate(t *testing.T) {
	t.Run("Update existing word", func(t *testing.T) {
		dict := Dictionary{"test": "this is a test"}
		key := "test"
		newValue := "this an updated test"

		err := dict.Update(key, newValue)
		assertError(t, err, nil)
		assertDefinition(t, dict, key, newValue)
	})

	t.Run("Update non-existing word", func(t *testing.T) {
		dict := Dictionary{}
		key := "test"
		newValue := "this an updated test"

		err := dict.Update(key, newValue)
		assertError(t, err, ErrKeyNotFound)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete existing word", func(t *testing.T) {
		dict := Dictionary{"test": "this is a test"}
		key := "test"
		err := dict.Delete(key)
		assertError(t, err, nil)
		_, err2 := dict.Search(key)
		assertError(t, err2, ErrKeyNotFound)
	})

	t.Run("Delete non-existing word", func(t *testing.T) {
		dict := Dictionary{}
		key := "test"
		err := dict.Delete(key)
		assertError(t, err, ErrKeyNotFound)
	})

}

func assertDefinition(t testing.TB, dict Dictionary, key string, value string) {
	got, err := dict.Search(key)
	want := value
	if err != nil {
		t.Fatal("Does not expect an error! error:\n", err)
	}

	assertStrings(t, got, want)
}

func assertStrings(t testing.TB, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
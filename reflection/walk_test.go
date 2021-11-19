package reflection

import (
	"testing"
	"reflect"
)

type Person struct {
	Name string
	Profile Profile
}

type Profile struct {
	Age int
	City string
}

func TestWalk(t *testing.T) {
	cases := []struct{
		Name string
		Input interface{}
		Expected []string		
	} {
		{
			"Struct with one field",
			struct{
				Name string
			}{"Alfred"},
			[]string{"Alfred"},
		},
		{
			"Struct with two string field",
			struct{
				Name string
				City string
			}{
				"Alfred", 
				"London",
			},
			[]string{"Alfred", "London"},

		},
		{
			"Struct with non string field",
			struct{
				Name string
				Age int
			}{
				"Alfred", 
				13,
			},
			[]string{"Alfred"},
		},
		{
			"Struct with nested field",
			Person {
				"Alfred",
				Profile{13, "London"},
			},
			[]string{"Alfred", "London"},
		},
		{
			"Pointer to things",
			&Person{
				"Alfred",
				Profile{13, "London"},
			},
			[]string{"Alfred", "London"},
		},
		{
			"Slices of struct",
			[]Profile{
				{24, "Charlie"},
				{13, "Alfred"},
			},
			[]string{"Charlie", "Alfred"},
		},
		{
			"Arrays",
			[]Profile{
				{24, "Charlie"},
				{13, "Alfred"},
			},
			[]string{"Charlie", "Alfred"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got := []string{}
			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, test.Expected) {
				t.Errorf("got %v want %v", got, test.Expected)
			} 

		})
	}

	t.Run("Maps", func(t *testing.T) {
		aMap := map[string]string{
				"Foo": "Charlie",
				"Baz": "Alfred",
			}
		got := []string{}
		walk(aMap, func(input string) {
			got = append(got, input)
		})

		assertContains(t, got, "Charlie")
		assertContains(t, got, "Alfred")
	})

	t.Run("Channels", func(t *testing.T) {
		aChannel := make(chan Profile)

		go func() {
			aChannel <- Profile{32, "A"}
			aChannel <- Profile{12, "B"}
			close(aChannel)
		}()

		var got []string
		want := []string{"A", "B"}
		walk(aChannel, func(input string) {
			got = append(got, input)
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}

	})

	t.Run("Functions", func(t *testing.T) {
		aFunction := func() (Profile, Profile) {
			return Profile{12, "A"}, Profile{23, "B"}
		}

		var got []string
		want := []string{"A", "B"}
		walk(aFunction, func(input string) {
			got = append(got, input)
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}

	})

}

func assertContains(t testing.TB, haystack []string, needle string) {
	t.Helper()
	var contains bool
	for _, x := range haystack {
		if x == needle {
			contains = true
		}
	}

	if !contains {
		t.Errorf("expected %+v to contain %q but did not", haystack, needle)
	}
}

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

}

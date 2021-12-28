package poker

import "fmt"


type Eat interface {
	Eat() string
}

type Person struct {
	Age int
	Children []int
}

func (p *Person) getChildren() []int {
	return p.Children
}

func (p *Person) getAge() int {
	return p.Age
}

func (p *Person) Eat() string {
	return "shurush"
}


func Sample(eat Eat) {
	fmt.Println(eat.Eat())
}

func main2() {
	person := Person{
		Age: 23,
		Children: []int{12, 11, 1},
	}

	Sample(&person)

	children := person.getChildren()
	fmt.Printf("children before %v\n", children)

	children[0]++
	fmt.Printf("children after %v\n", person.Children)

	age := person.getAge()
	fmt.Printf("age before %v\n", age)
	age++
	fmt.Printf("age after %v\n", person.Age)
}
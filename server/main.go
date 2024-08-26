package main

import (
	"fmt"
)

type leha struct {
	Name string
	Age  int
}

func main() {
	fmt.Println("Hello, World!")

	a := []leha{}

	a = append(a, leha{
		Name: "test",
		Age:  1,
	})

	a = append(a, leha{
		Name: "test2",
		Age:  2,
	})

	for i := range a {
		if a[i].Name == "test" {
			a[i].Name = "test3"
			a[i].Age = 3
		}
	}

	for i := range a {
		fmt.Println(a[i])
	}
}

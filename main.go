package main

import "github.com/juliangruber/go-intersect"

import (
	"fmt"
	"os"
)

func main() {
	a := []string{"cumquat", "apple", "orange", "pear"}
	b := []string{"apple", "orange", "banana"}

	i := intersect.Simple(a, b)
	fmt.Println("a length:", len(a))
	fmt.Println("b length:", len(b))

	fmt.Println("intersection:", i)
	fmt.Println("intersection length:", len(i))

	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	arg := os.Args[3]
	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	fmt.Println(arg)

	fmt.Println(len(os.Args))

}

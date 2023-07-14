package data

import (
	"fmt"
)

func ExampleInt2AAA() {
	s := Int2AAA(703)
	fmt.Println(s)
	// Output:AAA
}

func ExampleAAA2Int() {
	s := AAA2Int("AAA")
	fmt.Println(s)
	// Output:703
}

func ExampleContains() {
	fmt.Println(Contains([]string{"a", "B", "c"}, "a"))
	fmt.Println(Contains([]string{"a", "B", "c"}, "b"))
	// Output:true
	//false
}

func ExampleInsert() {
	a := []int{1, 2, 3, 4, 5}
	a = Insert(a, 2, 7, 8, 9)
	fmt.Println(a)
	// Output:[1 2 7 8 9 3 4 5]
}

func ExamplePop() {
	a := []int{1, 2, 3, 4, 5}
	a, b := Pop(a, -3)
	fmt.Println(a, b)
	// Output:[1 2 4 5] 3
}

func ExampleEqual() {
	fmt.Println(Equal([]string{"1"}, []string{"1"}))
	fmt.Println(Equal([]int{1}, []int{2}))
	// Output:true
	//false
}

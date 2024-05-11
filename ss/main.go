package main

import "fmt"

func main() {
	ss([]int{1, 2, 3})
}

func ss(s any) {
	fmt.Println(s.([]any))
}

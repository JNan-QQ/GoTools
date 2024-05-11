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

func ExampleCn2an() {
	n, _ := Cn2an("二十一")
	fmt.Println(n)
	// Output:1000006002400
}

func ExampleAn2cn() {
	s, _ := An2cn(1000006002400)
	fmt.Println(s)
	// Output:一万亿零六百万零二千四百
}

func ExampleRepeatIndex() {
	a := []int{1, 3, 5, 7, 9, 2, 5, 3, 7, 1, 10}
	fmt.Println(RepeatIndex(a))
	// Output: map[1:[0 9] 3:[1 7] 5:[2 6] 7:[3 8]]
}

func ExampleMD5() {
	var s MD5 = "123456"
	fmt.Println(s.Encrypt())
	fmt.Println(s.Encrypt().IsBig())
	fmt.Println(s.Encrypt().IsShort().IsBig())
	fmt.Printf("%T %s", s.Encrypt(), s.Encrypt().ToString())
	// Output:e10adc3949ba59abbe56e057f20f883e
	//E10ADC3949BA59ABBE56E057F20F883E
	//49BA59ABBE56E057
	//data.MD5 e10adc3949ba59abbe56e057f20f883e
}

func ExampleJsonString_Find() {
	j := JsonString(`{"name":"jjs","age":18,"love":[123,45],"info":{"grade":7,"d":[60,70,80,90],"desc":"good"},"s":{"a":1,"b":2}}`)

	find, _ := j.Find("info.d.[0]")
	fmt.Println(find)

	find = j.FindInt("info.grade")
	fmt.Println(find)

	// Output:60
	//7
}

func ExampleIsDigit() {
	fmt.Println(IsDigit("123"))
	fmt.Println(IsDigit("123a"))

	// Output:true
	//false
}

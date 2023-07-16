package pandas

import "fmt"

func ExampleRead() {
	df := Read("22.xls")
	fmt.Println(df)
	df1 := Read("11.xlsx")
	fmt.Println(df1)
	//Output: a
}

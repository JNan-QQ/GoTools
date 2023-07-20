package pandas

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
)

func ExampleRead() {

	df, err := Read("test.xls")
	if err != nil {
		panic(err)
	}
	fmt.Println(df)

	df1, err := Read("test.xlsx")
	if err != nil {
		panic(err)
	}
	fmt.Println(df1)
}

func ExampleDataFrame_FormatCols() {
	df := DataFrame{}
	df.DataFrame = dataframe.LoadRecords(
		[][]string{
			{"A", "B", "C", "D"},
			{"a", "4", "5.1", "true"},
			{"k", "5", "7.0", "true"},
			{"k", "4", "6.0", "true"},
			{"a", "2", "7.1", "false"},
		},
	)

	df.FormatCols(func(elem any) any {
		switch v := elem.(type) {
		case float64:
			return v + 1
		case int:
			return v - 1
		default:
			return 0
		}
	},
		"B", "C",
	)
	fmt.Println(df)
	// Output: [4x4] DataFrame
	//
	//    A        B     C        D
	// 0: a        3     6.100000 true
	// 1: k        4     8.000000 true
	// 2: k        3     7.000000 true
	// 3: a        1     8.100000 false
	//    <string> <int> <float>  <bool>
}

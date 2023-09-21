package pandas

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

func ExampleRead() {

	df := Read("test.xls")
	if df.Error() != nil {
		panic(df.Error())
	}
	fmt.Println(df)

	df1 := Read("test.xlsx")
	if df1.Error() != nil {
		panic(df1.Error())
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
}

func ExampleDataFrame_WriteXLSX() {
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

	err := df.WriteXLSX("test.xlsx")
	if err != nil {
		panic(err)
	}
}

func ExampleDataFrame_SetType() {
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
	fmt.Println(df)

	df.SetType(map[string]series.Type{"D": series.Int, "C": series.String})

	fmt.Println(df)
}

func ExampleDataFrame_SelectCols() {
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

	df = df.SelectCols("D", "B")

	fmt.Println(df)
}

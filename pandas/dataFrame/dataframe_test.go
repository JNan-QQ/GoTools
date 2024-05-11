package dataFrame

import (
	"fmt"
	"gitee.com/jn-qq/go-tools/data"
	"gitee.com/jn-qq/go-tools/pandas/series"
	"testing"
)

func TestLoadMap(t *testing.T) {
	df, err := LoadMap(map[string]any{
		"name":  data.CreateSlice("Join", 100),
		"phone": data.CreateSlice("15963578965", 100),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(df)
}

func TestLoadRecord(t *testing.T) {
	df, err := LoadRecord(
		data.CreateSlice([]string{"Join", "15963578965"}, 100),
		[]string{"name", "phone"},
		[]series.Type{series.String, series.Int},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(df)
}

func TestNew(t *testing.T) {
	df, err := New(
		[]any{data.CreateSlice("Join", 100), data.CreateSlice(15963578965, 100)},
		[]string{"name", "phone"},
	)
	if err != nil {
		return
	}
	fmt.Println(df)
}

func TestDataFrame_Records(t *testing.T) {
	df, err := New(
		[]any{data.CreateSlice("Join", 10), data.CreateSlice(15963578965, 10)},
		[]string{"name", "phone"},
	)
	if err != nil {
		return
	}
	fmt.Println(df.Records(false, false))
	fmt.Println(df.Records(false, true))
	fmt.Println(df.Records(true, true))
}

func TestEqualLength(t *testing.T) {
	s := [][]any{{1, 2, 3}, {"4", "5", "6"}, {7, 8, 9}}
	s1 := map[string]any{
		"1": []int{1, 2, 3},
		"2": []string{"4", "5", "6"},
		"3": []int{7, 8, 9},
	}
	fmt.Println(equalLength(s, 0))
	fmt.Println(equalLength(s1, 0))
	fmt.Println(equalLength(s1, 4))
}

func TestDataFrame_Set(t *testing.T) {
	df, _ := New(
		[]any{data.CreateSlice("Join", 5), data.CreateSlice(15963578965, 5)},
		[]string{"name", "phone"},
	)
	fmt.Println(df)
	_ = df.Set(0, []any{"Andy", 1111111111})
	fmt.Println(df)
	_ = df.Set(2, map[string]any{"phone": 222222222})
	fmt.Println(df)
	_ = df.Set(df.rows, []any{"Andy", 1111111111})
	fmt.Println(df)

}

func TestDataFrame_AddRows(t *testing.T) {
	df, _ := New(
		[]any{data.CreateSlice("Join", 5), data.CreateSlice(15963578965, 5)},
		[]string{"name", "phone"},
	)
	fmt.Println(df)
	_ = df.AddRows([][]any{{"name1", 12345678}, {"name1", 12345678}, {"name1", 12345678}})
	fmt.Println(df)
}

func TestDataFrame_Arrange(t *testing.T) {
	df, _ := New(
		[]any{
			[]string{"伏旭歆", "管原炳", "仰芝凤", "万茵瑾", "左芊筱", "俞淑允", "宗茹淳", "卓虹", "司丽瑾", "岑泳继"},
			[]int{13935531105, 15665203778, 14583084372, 14779318181, 17606363473, 18950385204, 18659058185, 16628908658, 17590257481, 17254554855},
			[]int{35, 36, 42, 13, 20, 20, 14, 20, 30, 36},
		},
		[]string{"name", "phone", "age"},
	)
	fmt.Println(df)
	err := df.Arrange(SortByForward("name"), Order{ColumnName: "age", Reverse: true})
	if err != nil {
		panic(err)
	}
	fmt.Println(df)
}

func TestDataFrame_AddCol(t *testing.T) {
	df, _ := New(
		[]any{data.CreateSlice("Join", 5), data.CreateSlice(15963578965, 5)},
		[]string{"name", "phone"},
	)
	_ = df.AddCol("addr", data.CreateSlice("xxxxx", 5), nil)
	fmt.Println(df)
	_ = df.AddCol("addr", data.CreateSlice("yyyy", 1), "yyy")
	fmt.Println(df)
}

func TestDataFrame_Columns(t *testing.T) {

}

package pandas

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

type DataFrame struct {
	columns []Series
	cols    int
	rows    int
}

// NewDataFrame 创建 DataFrame 数据对象
//
//	values: 待输入数据，可以为 map[string][]T、[][]T、nil、 []Series
//	 当为 map[string][]T 时，以列切片方式添加
//	 当为 [][]T 时，以行切片方式添加，第一个元素是列名切片。如果 columnName != nil ,以列切片添加
//	 当为 []Series 时，直接添加
//	 当为 nil 时，创建空 DataFrame
func NewDataFrame(values any, columnName []string) (*DataFrame, error) {

	dataFrame := DataFrame{columns: make([]Series, 0), cols: 0, rows: 0}

	if values == nil {
		return &dataFrame, nil
	}

	switch vals := values.(type) {
	case []Series:
		rows := 0
		var series []Series
		for _, val := range vals {
			if rows == 0 {
				rows = val.Len()
			} else if rows != val.Len() {
				return nil, fmt.Errorf("输入数据长度不一致！")
			} else {
				series = append(series, *(val.Copy()))
			}
		}
		dataFrame.columns = series
	case map[string]any:
		cols := 0
		for key, value := range vals {
			if err := dataFrame.newSeries(value, key, &cols); err != nil {
				return nil, err
			}
		}
	case []any:
		if columnName != nil {
			if len(columnName) != len(vals) {
				return nil, fmt.Errorf("列名个数不对应！")
			}
			rows := 0
			for i, value := range vals {
				if err := dataFrame.newSeries(value, columnName[i], &rows); err != nil {
					return nil, err
				}
			}
		} else {
			columnsName := vals[0].([]string)
			columns := make([]any, len(columnName))
			for _, value := range vals[1:] {
				v := reflect.ValueOf(value)
				if v.Len() != len(columnsName) {
					return nil, fmt.Errorf("切片长度不一致")
				}
				for i := 0; i < v.Len(); i++ {
					if columns[i] == nil {
						switch v.Type().Elem().String() {
						case String:
							columns[i] = []string{v.Index(i).String()}
						case Bool:
							columns[i] = []bool{v.Index(i).Bool()}
						case Float:
							columns[i] = []float64{v.Index(i).Float()}
						case Int:
							columns[i] = []int{int(v.Index(i).Int())}
						default:
							return nil, fmt.Errorf("未知数据类型%s", v.Type().Elem().String())
						}
					} else {
						switch v.Type().Elem().String() {
						case String:
							columns[i] = append(columns[i].([]string), v.Index(i).String())
						case Bool:
							columns[i] = append(columns[i].([]bool), v.Index(i).Bool())
						case Float:
							columns[i] = append(columns[i].([]float64), v.Index(i).Float())
						case Int:
							columns[i] = append(columns[i].([]int), int(v.Index(i).Int()))
						}
					}
				}
				n := 0
				for i, column := range columns {
					if err := dataFrame.newSeries(column, columnsName[i], &n); err != nil {
						return nil, err
					}
				}
			}
		}
	default:
		return nil, fmt.Errorf("不支持的数据类型！")
	}

	dataFrame.size()

	return &dataFrame, nil
}

// Records 返回字符串切片，第一个切片为列名
//
//	isRow = false 返回列切片
//	isRow = true  返回行切片
func (d *DataFrame) Records(isRow bool) [][]string {
	var res = [][]string{d.ColumnNames()}
	if isRow {
		for i := 0; i < d.RowNum(); i++ {
			var rows []string
			for _, column := range d.columns {
				rows = append(rows, column.elements.index(i).records())
			}
			res = append(res, rows)
		}
	} else {
		for _, column := range d.columns {
			res = append(res, column.Records())
		}
	}

	return res
}

// ColumnNames 返回列名
func (d *DataFrame) ColumnNames() []string {
	var names []string
	for _, column := range d.columns {
		names = append(names, column.Name)
	}
	return names
}

// ColumnType 返回列类型
func (d *DataFrame) ColumnType() []string {
	var _types []string
	for _, column := range d.columns {
		_types = append(_types, column.Type())
	}
	return _types
}

// 创建 Series 对象
func (d *DataFrame) newSeries(values any, name string, length *int) error {

	if reflect.TypeOf(values).Kind() != reflect.Slice {
		return fmt.Errorf("输入数据 map[string][]T 或 [][]T 类型")
	}

	if *length == 0 {
		*length = reflect.ValueOf(values).Len()
	}
	if *length != reflect.ValueOf(values).Len() {
		return fmt.Errorf("输入数据长度不一致！")
	}

	var series *Series
	var err error

	switch value := values.(type) {
	case []string:
		series, err = NewSeries(value, String, name)
	case []int:
		series, err = NewSeries(value, Int, name)
	case []float64:
		series, err = NewSeries(value, Float, name)
	case []bool:
		series, err = NewSeries(value, Bool, name)
	default:
		return fmt.Errorf("错误的切片类型！")
	}

	if err != nil {
		return err
	}

	d.columns = append(d.columns, *series)

	return nil
}

// 自定义输出
func (d *DataFrame) String() string {
	var str string

	maxWith := make([]int, d.ColNum())
	for i, _strings := range d.Records(false)[1:] {
		for _, s := range append(_strings, d.ColumnNames()[i]) {
			with := utf8.RuneCountInString(s)
			if with > maxWith[i] {
				maxWith[i] = with
			}
		}
	}

	for i, _strings := range d.Records(true) {
		if i == 0 {
			str += fmt.Sprintf("%-5s", "序号")
		} else {
			str += fmt.Sprintf("%-5d", i-1)
		}
		for i2, s := range _strings {
			str += fmt.Sprintf("%"+strconv.Itoa(maxWith[i2]+3)+"s", s)
		}
		str += "\n"
		if i == 0 {
			str += strings.Repeat("=", utf8.RuneCountInString(str)) + "\n"
		}
	}

	return str
}

// 更新二维数组大小
func (d *DataFrame) size() {
	d.cols = len(d.columns)
	if d.cols > 0 {
		d.rows = d.columns[0].Len()
	} else {
		d.rows = 0
	}
}

// RowNum 行数
func (d *DataFrame) RowNum() int {
	return d.rows
}

// ColNum 列数
func (d *DataFrame) ColNum() int {
	return d.cols
}

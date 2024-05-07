package dataFrame

import (
	"fmt"
	"gitee.com/jn-qq/go-tools/data"
	"gitee.com/jn-qq/go-tools/pandas/series"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

type DataFrame struct {
	columns []series.Series
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

	dataFrame := DataFrame{columns: make([]series.Series, 0), cols: 0, rows: 0}

	if values == nil {
		return &dataFrame, nil
	}

	switch vals := values.(type) {
	case []series.Series:
		rows := 0
		var _series []series.Series
		for _, val := range vals {
			if rows == 0 {
				rows = val.Len()
			} else if rows != val.Len() {
				return nil, fmt.Errorf("输入数据长度不一致！")
			} else {
				_series = append(_series, *(val.Copy()))
			}
		}
		dataFrame.columns = _series
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
						case series.String:
							columns[i] = []string{v.Index(i).String()}
						case series.Bool:
							columns[i] = []bool{v.Index(i).Bool()}
						case series.Float:
							columns[i] = []float64{v.Index(i).Float()}
						case series.Int:
							columns[i] = []int{int(v.Index(i).Int())}
						default:
							return nil, fmt.Errorf("未知数据类型%s", v.Type().Elem().String())
						}
					} else {
						switch v.Type().Elem().String() {
						case series.String:
							columns[i] = append(columns[i].([]string), v.Index(i).String())
						case series.Bool:
							columns[i] = append(columns[i].([]bool), v.Index(i).Bool())
						case series.Float:
							columns[i] = append(columns[i].([]float64), v.Index(i).Float())
						case series.Int:
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

	dataFrame.Size()

	return &dataFrame, nil
}

// Records 返回字符串切片
//
//	isRow = false 返回列切片 isRow = true  返回行切片
//	hasColName 是否返回列名，第一个元素，或切片
func (d *DataFrame) Records(isRow bool, hasColName bool) [][]string {
	var res [][]string
	if isRow {
		if hasColName {
			res = append(res, d.ColumnNames())
		}
		for i := 0; i < d.rows; i++ {
			var rows []string
			for _, column := range d.columns {
				rows = append(rows, column.Index(i).Records())
			}
			res = append(res, rows)
		}
	} else {
		for _, column := range d.columns {
			values := column.Records()

			if hasColName {
				values = data.Insert(values, 0, column.Name)
			}
			res = append(res, values)
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

	var _series *series.Series
	var err error

	switch value := values.(type) {
	case []string:
		_series, err = series.NewSeries(value, series.String, name)
	case []int:
		_series, err = series.NewSeries(value, series.Int, name)
	case []float64:
		_series, err = series.NewSeries(value, series.Float, name)
	case []bool:
		_series, err = series.NewSeries(value, series.Bool, name)
	default:
		return fmt.Errorf("错误的切片类型！")
	}

	if err != nil {
		return err
	}

	d.columns = append(d.columns, *_series)

	return nil
}

// 自定义输出
func (d *DataFrame) String() string {
	return d.print(20, false)
}

func (d *DataFrame) print(maxWith int, isComplete bool) (str string) {

	if d.cols == 0 || d.rows == 0 {
		return "DataFrame is Empty!"
	}

	str += fmt.Sprintf("DataFrame size：%d x %d\n\n", d.cols, d.rows)

	if maxWith == 0 {
		maxWith = 20
	}

	if d.rows <= 50 {
		isComplete = true
	}

	step := strings.Repeat("|", 2)
	total := 2 + (maxWith+3)*d.cols

	color := [][]int{
		{0, 0, 36},
		{0, 0, 33},
	}

	// 按行遍历
	for i, rows := range d.Records(true, true) {
		// 默认简略输出
		if !isComplete && i >= 17 && i <= d.rows-6 {
			if i == 17 {
				str += fmt.Sprintf("%*s%s", 4, "...", step) +
					strings.Repeat(fmt.Sprintf("%20s%s", "...", step), d.cols) + "\n"
			}
			continue
		}

		// 添加索引列
		if i != 0 {
			str += fmt.Sprintf("%3d%s", i, step)
		} else {
			str += fmt.Sprintf("%2s%s", "序号", step)
		}

		// 遍历列字符
		for j, chars := range rows {
			cl := color[j%2]
			if i == 0 {
				cl = []int{0, 0, 30}
			}
			// 判断实际长度
			if utf8.RuneCountInString(chars)+data.HanCount(chars) <= maxWith {
				str += data.ColorStr(fmt.Sprintf("%*s%s", maxWith-data.HanCount(chars), chars, step), cl[0], cl[1], cl[2])
			} else {
				// 截取指定最大长度的字符串
				n := 0
				nChar := ""
				for _, char := range chars {
					if unicode.Is(unicode.Han, char) {
						n += 2
					} else {
						n += 1
					}
					if n > maxWith-3 {
						nChar += "...."
						break
					} else if n == maxWith-3 {
						nChar += string(char) + "..."
						break
					} else {
						nChar += string(char)
					}
				}
				str += data.ColorStr(fmt.Sprintf("%*s%s", maxWith-data.HanCount(nChar), nChar, step), cl[0], cl[1], cl[2])
			}

		}
		str += "\n"
		if i == 0 {
			str += strings.Repeat("=", total) + "\n"
		}
	}
	str += strings.Repeat("-", total) + "\n"
	str += strings.Repeat(" ", 3) + step
	for _, s := range d.ColumnType() {
		str += fmt.Sprintf("%*s%s", maxWith, s, step)
	}

	return
}

// Size 更新并返回二维数组大小
func (d *DataFrame) Size() (cols, rows int) {
	d.cols = len(d.columns)
	if d.cols > 0 {
		d.rows = d.columns[0].Len()
	} else {
		d.rows = 0
	}
	return d.cols, d.rows
}

// Copy 复制
func (d *DataFrame) Copy() *DataFrame {
	_copy, _ := NewDataFrame(d.columns, nil)
	return _copy
}

package pandas

import (
	"fmt"
	"reflect"
)

type Series struct {
	Name     string
	elements Elements
	t        Type
}

// NewSeries 创建数据列
//
//	values: 数据切片，可选nil、[]string、[]int、[]float64、[]bool
//	dType: 数据类型，可选String、int、float、bool
//	name: 数据列名称
func NewSeries[E int | float64 | string | bool](values []E, dType Type, name string) (Series, error) {
	series := Series{
		Name: name,
		t:    dType,
	}

	switch dType {
	case String:
		series.elements = &stringElements{}
	case Int:
		series.elements = &intElements{}
	case Float:
		series.elements = &floatElements{}
	case Bool:
		series.elements = &boolElements{}
	}

	if values == nil {
		return series, nil
	} else if reflect.TypeOf(values).String() != fmt.Sprintf("[]%s", dType) {
		return Series{}, fmt.Errorf("输入切片与指定数据类型不匹配")
	} else {
		if err := series.elements.append(values); err != nil {
			return Series{}, err
		}
	}

	return series, nil
}

// Elements 返回数据集元素对象切片
func (s Series) Elements() []Element {
	return s.elements.elems()
}

func (s Series) Format() Series {
	return s
}

func (s Series) ChangeType(t Type) Series {
	return s
}

// 自定义输出
func (s Series) String() string {
	return fmt.Sprintf("字段名：%s\n数 据：%v\n类 型：%s\n", s.Name, s.elements.records(), s.t)
}

// Len 返回数据集大小
func (s Series) Len() int {
	return s.elements.len()
}

// Records 将数据集中的元素作为字符串返回
func (s Series) Records() []string {
	return s.elements.records()
}

// Float 将数据集中的元素作为浮点数返回，如果数据转换失败自动设为 math.NaN()
func (s Series) Float() []float64 {
	return s.elements.float()
}

// Int 将数据集中的元素作为整数返回，如果数据转换失败自动设为 math.MinInt()
func (s Series) Int() []int {
	return s.elements.int()
}

// Bool 将数据集中的元素作为浮布尔值返回
func (s Series) Bool() []bool {
	return s.elements.bool()
}

type Operator string

const (
	Equal          Operator = "=="
	NotEqual       Operator = "!="
	LessThan       Operator = "<"
	LessOrEqual    Operator = "<="
	GreaterThan    Operator = ">"
	GreaterOrEqual Operator = ">="
)

func (s Series) Filter() Series {

}

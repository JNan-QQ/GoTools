package pandas

import (
	"fmt"
	"gitee.com/jn-qq/go-tools/data"
	"math"
	"reflect"
	"strings"
)

type Series struct {
	Name     string
	elements Elements
	t        string
	indexes  []int
}

// NewSeries 创建数据列
//
//	values: 数据切片，可选nil、[]string、[]int、[]float64、[]bool
//	dType: 数据类型，可选String、int、float64、bool
//	name: 数据列名称
func NewSeries[E int | float64 | string | bool](values []E, dType string, name string) (*Series, error) {
	series := Series{
		Name:     name,
		t:        dType,
		elements: createElements(dType, 0),
	}

	if values == nil {
		return &series, nil
	} else if reflect.TypeOf(values).String() != fmt.Sprintf("[]%s", dType) {
		return nil, fmt.Errorf("输入切片与指定数据类型不匹配")
	} else {
		if err := series.elements.append(values); err != nil {
			return nil, err
		}
	}
	series.initIndexes()

	return &series, nil
}

// 重置索引
func (s *Series) initIndexes() {
	s.indexes = make([]int, 0)
	for i := 0; i < s.Len(); i++ {
		s.indexes = append(s.indexes, i)
	}
}

// Append 向数据集后添加元素,可以为单个元素或元素切片
func (s *Series) Append(values interface{}) error {
	if err := s.elements.append(values); err != nil {
		return err
	}
	s.initIndexes()
	return nil
}

// Concat 将新数据集加原数据集后,如果类型不同，以原数据集为准
func (s *Series) Concat(x Series) error {
	if x.t != s.t {
		fmt.Printf("两个数据类型不同，正在尝试转换...")
		if err := x.SetType(s.t); err != nil {
			return err
		}
		fmt.Println(" done")
	}

	var nv any

	switch s.t {
	case String:
		nv = x.Records()
	case Int:
		nv = x.Int()
	case Float:
		nv = x.Float()
	case Bool:
		nv = x.Bool()
	}

	if err := s.elements.append(nv); err != nil {
		return err
	}
	s.initIndexes()

	return nil
}

// SubSet 保留对应索引的元素
func (s *Series) SubSet(indexes ...int) *Series {
	var delIndexes []int
	for _, index := range s.indexes {
		if !data.Contains(indexes, index) {
			delIndexes = append(delIndexes, index)
		}
	}
	ns := s.Copy()
	ns.elements.drop(delIndexes...)
	ns.initIndexes()
	return ns
}

// Elements 返回数据集元素对象切片
func (s *Series) Elements() []Element {
	return s.elements.elems()
}

// Format 批量处理数据集
//
//	index: 元素索引
//	elem: 元素对象
func (s *Series) Format(f func(index int, elem Element) Element) {
	for i, element := range s.elements.elems() {
		s.elements.index(i).update(f(i, element))
	}
}

// SetType 改变数据集类型
func (s *Series) SetType(t string) error {
	if s.elements.dType() == t {
		return nil
	}

	newSeries := Series{
		Name:     s.Name,
		elements: createElements(t, 0),
		t:        t,
		indexes:  s.indexes,
	}

	var values any
	switch t {
	case String:
		values = s.elements.records()
	case Int:
		values = s.elements.int()
	case Float:
		values = s.elements.float()
	case Bool:
		values = s.elements.bool()
	default:
		return fmt.Errorf("未知数据类型！")
	}

	if err := newSeries.Append(values); err != nil {
		return err
	}

	*s = newSeries

	return nil
}

// 自定义输出
func (s *Series) String() string {
	return fmt.Sprintf("字段名：%s\n数 据：%v\n索引：%v\n类 型：%s\n", s.Name, s.elements.records(), s.indexes, s.t)
}

// Len 返回数据集大小
func (s *Series) Len() int {
	return s.elements.len()
}

// Type 返回类型
func (s *Series) Type() string {
	return s.t
}

// Copy 复制
func (s *Series) Copy() *Series {
	ns := Series{
		Name:     s.Name,
		elements: createElements(s.t, s.Len()),
		t:        s.t,
		indexes:  s.indexes,
	}
	for i, element := range s.elements.elems() {
		ns.elements.index(i).update(element)
	}

	return &ns
}

// Records 将数据集中的元素作为字符串返回
func (s *Series) Records() []string {
	return s.elements.records()
}

// Float 将数据集中的元素作为浮点数返回，如果数据转换失败自动设为 math.NaN()
func (s *Series) Float() []float64 {
	return s.elements.float()
}

// Int 将数据集中的元素作为整数返回，如果数据转换失败自动设为 math.MinInt()
func (s *Series) Int() []int {
	return s.elements.int()
}

// Bool 将数据集中的元素作为浮布尔值返回
func (s *Series) Bool() []bool {
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
	Contains       Operator = "contains"
	StartsWith     Operator = "starts_with"
	EndsWith       Operator = "ends_with"
	In             Operator = "in"
	NotIn          Operator = "not_in"
)

func (s *Series) Filter(operator Operator, values any) (*Series, error) {

	vT := reflect.TypeOf(values).String()
	if !(vT == fmt.Sprintf("[]%s", s.t) || vT == s.t) {
		return nil, fmt.Errorf("输入数据类型与数据集类型不同")
	}

	ns := s.Copy()

	var elements = s.Elements()
	switch s.t {
	case String:
		if !data.Contains([]Operator{Equal, NotEqual, Contains, StartsWith, EndsWith, In, NotIn}, operator) {
			return nil, fmt.Errorf("string 类型数据无法执行操作%s", operator)
		}

	case Int:
		if !data.Contains([]Operator{Equal, NotEqual, LessThan, LessOrEqual, GreaterThan, GreaterOrEqual, In, NotIn}, operator) {
			return nil, fmt.Errorf("int 类型数据无法执行操作%s", operator)
		}
	case Float:
		if !data.Contains([]Operator{Equal, NotEqual, LessThan, LessOrEqual, GreaterThan, GreaterOrEqual, In, NotIn}, operator) {
			return nil, fmt.Errorf("float64 类型数据无法执行操作%s", operator)
		}
	case Bool:
		if !data.Contains([]Operator{Equal, NotEqual}, operator) {
			return nil, fmt.Errorf("bool 类型数据无法执行操作%s", operator)
		}
	}

	if (operator == In || operator == NotIn) && reflect.TypeOf(values).Kind() != reflect.Slice {
		return nil, fmt.Errorf("in / NotIn 需要输入切片作为参数")
	}

	_, ns.indexes = data.Filter(elements, func(element Element) bool {
		switch operator {
		case Equal:
			return element.value() == values
		case NotEqual:
			return element.value() != values
		case LessThan, LessOrEqual, GreaterThan, GreaterOrEqual:
			if math.IsNaN(element.float()) {
				return false
			}
			if reflect.TypeOf(values).Kind() == reflect.Int {
				switch operator {
				case LessThan:
					return element.int() < values.(int)
				case LessOrEqual:
					return element.int() <= values.(int)
				case GreaterThan:
					return element.int() > values.(int)
				case GreaterOrEqual:
					return element.int() >= values.(int)
				}
			} else {
				switch operator {
				case LessThan:
					return element.float() < values.(float64)
				case LessOrEqual:
					return element.float() <= values.(float64)
				case GreaterThan:
					return element.float() > values.(float64)
				case GreaterOrEqual:
					return element.float() >= values.(float64)
				}
			}
		case Contains:
			return strings.Contains(element.records(), values.(string))
		case StartsWith:
			return strings.HasPrefix(element.records(), values.(string))
		case EndsWith:
			return strings.HasSuffix(element.records(), values.(string))
		case In, NotIn:
			var newValues []any
			v := reflect.ValueOf(values)
			for i := 0; i < v.Len(); i++ {
				newValues = append(newValues, v.Index(i).Interface())
			}
			if operator == In {
				return data.Contains(newValues, element.value())
			} else {
				return !data.Contains(newValues, element.value())
			}
		}
		return false
	})

	var delIndex []int
	for _, i := range s.indexes {
		if !data.Contains(ns.indexes, i) {
			delIndex = append(delIndex, i)
		}
	}
	ns.elements.drop(delIndex...)

	return ns, nil
}

// 生成 Elements 接口的对象
func createElements(dType string, l int) Elements {
	var ne Elements
	switch dType {
	case String:
		s := make(stringElements, l)
		ne = &s
	case Int:
		i := make(intElements, l)
		ne = &i
	case Float:
		f := make(floatElements, l)
		ne = &f
	case Bool:
		b := make(boolElements, l)
		ne = &b
	}
	return ne
}

package series

import (
	"fmt"
	"gitee.com/jn-qq/go-tools/data"
	"math"
	"sort"
	"strconv"
	"strings"
)

// 数据集类型
const (
	String string = "string"
	Int    string = "int"
	Float  string = "float64"
	Bool   string = "bool"
)

type Elements interface {
	// Index 返回指定索引的值
	index(int) Element
	// Elems 返回值切片
	elems() []Element
	// Append 添加值或切片
	append(any) error
	// Insert 在指定索引处插入值
	insert(int, ...any) error
	// Drop 删除定索引的值
	drop(...int)
	// Len 数据集大小
	len() int
	// String 数据集转化为字符串切片
	records() []string
	// Int 数据集转化为整数切片
	int() []int
	// Float 数据集转化为浮点数切片
	float() []float64
	// Bool 数据集转化为布尔值切片
	bool() []bool
	// 数据类型
	dType() string
}

type Element interface {
	// Set 设置值
	Set(any)
	// Records 返回字符串
	Records() string
	// Int 返回整数
	Int() int
	// Float 返回浮点数
	Float() float64
	// Bool 返回布尔值
	Bool() bool
	// Value 任一类型
	Value() any
	// Type 返回数据类型
	dType() string
	// Copy 复制值
	copy() Element
	// 判断是否为NaN
	isNaN() bool
	// 更新
	update(Element)
}

// 字符串数据格式，实现接口 Element
type stringElement string

// 字符串数据类型,实现接口 Elements
type stringElements []stringElement

// 整数数据格式，实现接口 Element
type intElement int

// 整数数据类型，实现接口 Elements
type intElements []intElement

// 浮点数数据格式，实现接口 Element
type floatElement float64

// 字符串数据类型
type floatElements []floatElement

// 布尔值数据格式，实现接口 Element
type boolElement bool

// 字符串数据类型
type boolElements []boolElement

//////////////////////////////////////
//            元素                  //
/////////////////////////////////////

func (s *stringElement) Set(value any) {
	switch val := value.(type) {
	case string:
		if data.Contains([]string{"", "NaN", "nan", "null", "Null"}, val) {
			*s = "NaN"
		} else {
			*s = stringElement(val)
		}
	case int:
		*s = stringElement(strconv.Itoa(val))
	case float64:
		*s = stringElement(strconv.FormatFloat(value.(float64), 'f', 6, 64))
	case bool:
		if val {
			*s = "true"
		} else {
			*s = "false"
		}
	default:
		*s = "NaN"
	}
}
func (i *intElement) Set(value any) {
	switch val := value.(type) {
	case string:
		if newInt, err := strconv.Atoi(val); err == nil {
			*i = intElement(newInt)
		} else {
			*i = math.MinInt
		}
	case int:
		*i = intElement(val)
	case float64:
		*i = intElement(val)
	case bool:
		if val {
			*i = 1
		} else {
			*i = 0
		}
	default:
		*i = math.MinInt
	}
}
func (f *floatElement) Set(value any) {
	switch val := value.(type) {
	case string:
		if newFloat, err := strconv.ParseFloat(val, 10); err != nil {
			*f = floatElement(math.NaN())
		} else {
			*f = floatElement(newFloat)
		}
	case int:
		*f = floatElement(val)
	case float64:
		*f = floatElement(val)
	case bool:
		if val {
			*f = 1
		} else {
			*f = 0
		}
	default:
		*f = floatElement(math.NaN())
	}
}
func (b *boolElement) Set(value any) {
	switch val := value.(type) {
	case string:
		if data.Contains([]string{"false", "0", "F", "f"}, val) {
			*b = false
		} else {
			*b = true
		}
	case int, float64:
		if val == 0 {
			*b = false
		} else {
			*b = true
		}
	case bool:
		*b = boolElement(val)
	default:
		panic("错误的数据类型")
	}
}

//--------------------------------//

func (s *stringElement) Records() string {
	return string(*s)
}
func (i *intElement) Records() string {
	if i.isNaN() {
		return "NaN"
	}
	return strconv.Itoa(int(*i))
}
func (f *floatElement) Records() string {
	if f.isNaN() {
		return "NaN"
	}
	return strconv.FormatFloat(float64(*f), 'f', -1, 64)
}
func (b *boolElement) Records() string {
	if *b {
		return "true"
	} else {
		return "false"
	}
}

//--------------------------------//

func (s *stringElement) Int() int {
	if s.isNaN() {
		return math.MinInt
	}
	if i, err := strconv.Atoi(string(*s)); err != nil {
		fmt.Printf("%s 不能转换为 Int，已置为无限小\n", string(*s))
		return math.MinInt
	} else {
		return i
	}
}
func (i *intElement) Int() int {
	return int(*i)
}
func (f *floatElement) Int() int {
	return int(*f)
}
func (b *boolElement) Int() int {
	if *b {
		return 1
	} else {
		return 0
	}
}

//--------------------------------//

func (s *stringElement) Float() float64 {
	if s.isNaN() {
		return math.NaN()
	}
	if f, err := strconv.ParseFloat(string(*s), 64); err != nil {
		return math.NaN()
	} else {
		return f
	}
}
func (i *intElement) Float() float64 {
	if *i == math.MinInt {
		return math.NaN()
	}
	return float64(*i)
}
func (f *floatElement) Float() float64 {
	return float64(*f)
}
func (b *boolElement) Float() float64 {
	if *b {
		return 1
	} else {
		return 0
	}
}

//--------------------------------//

func (s *stringElement) Bool() bool {
	if data.Contains([]string{"false", "0", "F", "f"}, strings.ToLower(string(*s))) {
		return false
	} else {
		return true
	}
}
func (i *intElement) Bool() bool {
	if *i > 0 {
		return true
	} else {
		return false
	}
}
func (f *floatElement) Bool() bool {
	if *f > 0 && !math.IsNaN(float64(*f)) {
		return true
	} else {
		return false
	}
}
func (b *boolElement) Bool() bool {
	return bool(*b)
}

//--------------------------------//

func (s *stringElement) dType() string {
	return String
}
func (i *intElement) dType() string {
	return Int
}
func (f *floatElement) dType() string {
	return Float
}
func (b *boolElement) dType() string {
	return Bool
}

//--------------------------------//

func (s *stringElement) copy() Element {
	s2 := new(stringElement)
	*s2 = *s
	return s2
}
func (i *intElement) copy() Element {
	i2 := new(intElement)
	*i2 = *i
	return i2
}
func (f *floatElement) copy() Element {
	f2 := new(floatElement)
	*f2 = *f
	return f2
}
func (b *boolElement) copy() Element {
	b2 := new(boolElement)
	*b2 = *b
	return b2
}

//--------------------------------//

func (s *stringElement) isNaN() bool {
	if *s == "NaN" {
		return true
	}
	return false
}
func (i *intElement) isNaN() bool {
	if *i == math.MinInt {
		return true
	}
	return false
}
func (f *floatElement) isNaN() bool {
	if math.IsNaN(float64(*f)) {
		return true
	}
	return false
}
func (b *boolElement) isNaN() bool {
	return false
}

//--------------------------------//

func (s *stringElement) String() string {
	if s.isNaN() {
		return "NaN"
	} else {
		return s.Records()
	}
}
func (i *intElement) String() string {
	if *i == math.MinInt {
		return "NaN"
	} else {
		return i.Records()
	}
}

//--------------------------------//

func (s *stringElement) update(elem Element) {
	s.Set(elem.Records())
}

func (i *intElement) update(elem Element) {
	i.Set(elem.Int())
}

func (f *floatElement) update(elem Element) {
	f.Set(elem.Float())
}

func (b *boolElement) update(elem Element) {
	b.Set(elem.Bool())
}

//--------------------------------//

func (s *stringElement) Value() any {
	return s.Records()
}

func (f *floatElement) Value() any {
	return f.Float()
}

func (b *boolElement) Value() any {
	return b.Bool()
}

func (i *intElement) Value() any {
	return i.Int()
}

///////////////元素组 //////////////////

func (s *stringElements) append(values any) error {
	switch value := values.(type) {
	case string:
		sm := new(stringElement)
		sm.Set(value)
		*s = append(*s, *sm)
	case []string:
		for _, v := range value {
			sm := new(stringElement)
			sm.Set(v)
			*s = append(*s, *sm)
		}
	default:
		return fmt.Errorf("请确认输入的是字符串")
	}

	return nil
}
func (i *intElements) append(values any) error {
	switch value := values.(type) {
	case int:
		im := new(intElement)
		im.Set(value)
		*i = append(*i, *im)
	case []int:
		for _, v := range value {
			im := new(intElement)
			im.Set(v)
			*i = append(*i, *im)
		}
	default:
		return fmt.Errorf("请确认输入的是整数")
	}
	return nil
}
func (f *floatElements) append(values any) error {
	switch value := values.(type) {
	case float64:
		fm := new(floatElement)
		fm.Set(value)
		*f = append(*f, *fm)
	case []float64:
		for _, v := range value {
			fm := new(floatElement)
			fm.Set(v)
			*f = append(*f, *fm)
		}
	default:
		return fmt.Errorf("请确认输入的是浮点数")
	}
	return nil
}
func (b *boolElements) append(values any) error {
	switch value := values.(type) {
	case bool:
		bm := new(boolElement)
		bm.Set(value)
		*b = append(*b, *bm)
	case []bool:
		for _, v := range value {
			bm := new(boolElement)
			bm.Set(v)
			*b = append(*b, *bm)
		}
	default:
		return fmt.Errorf("请确认输入的是布尔值")
	}
	return nil
}

//--------------------------------//

func (s *stringElements) insert(index int, values ...any) error {

	for i, value := range values {
		switch v := value.(type) {
		case string, int, float64, bool:
			sm := new(stringElement)
			sm.Set(v)
			*s = data.Insert(*s, index+i, *sm)
		default:
			return fmt.Errorf("请确认输入的是字符串")
		}
	}
	return nil
}
func (i *intElements) insert(index int, values ...any) error {
	for i2, value := range values {
		switch v := value.(type) {
		case string, int, float64, bool:
			im := new(intElement)
			im.Set(v)
			*i = data.Insert(*i, index+i2, *im)
		default:
			return fmt.Errorf("请确认输入的是整数")
		}
	}
	return nil
}
func (f *floatElements) insert(index int, values ...any) error {
	for i, value := range values {
		switch v := value.(type) {
		case string, int, float64, bool:
			fm := new(floatElement)
			fm.Set(v)
			*f = data.Insert(*f, index+i, *fm)
		default:
			return fmt.Errorf("请确认输入的是浮点数")
		}
	}
	return nil
}
func (b *boolElements) insert(index int, values ...any) error {

	for i, value := range values {
		switch v := value.(type) {
		case string, int, float64, bool:
			bm := new(boolElement)
			bm.Set(v)
			*b = data.Insert(*b, index+i, *bm)
		default:
			return fmt.Errorf("请确认输入的是布尔值")
		}
	}
	return nil
}

//--------------------------------//

func (s *stringElements) drop(indexes ...int) {
	sort.Ints(indexes)
	for _, i := range data.Reverse(indexes) {
		*s, _ = data.Pop(*s, i)
	}
}
func (i *intElements) drop(indexes ...int) {
	sort.Ints(indexes)
	for _, v := range data.Reverse(indexes) {
		*i, _ = data.Pop(*i, v)
	}
}
func (f *floatElements) drop(indexes ...int) {
	sort.Ints(indexes)
	for _, i := range data.Reverse(indexes) {
		*f, _ = data.Pop(*f, i)
	}
}
func (b *boolElements) drop(indexes ...int) {
	sort.Ints(indexes)
	for _, i := range data.Reverse(indexes) {
		*b, _ = data.Pop(*b, i)
	}
}

//--------------------------------//

func (s *stringElements) index(index int) Element {
	newStr := *s
	return &newStr[index]
}
func (i *intElements) index(index int) Element {
	newInt := *i
	return &newInt[index]
}
func (f *floatElements) index(index int) Element {
	newFloat := *f
	return &newFloat[index]
}
func (b *boolElements) index(index int) Element {
	newBool := *b
	return &newBool[index]
}

//--------------------------------//

func (s *stringElements) elems() []Element {
	var newStr []Element
	for _, element := range *s {
		newStr = append(newStr, element.copy())
	}
	return newStr
}
func (i *intElements) elems() []Element {
	var newInt []Element
	for _, element := range *i {
		newInt = append(newInt, element.copy())
	}
	return newInt
}
func (f *floatElements) elems() []Element {
	var newFloat []Element
	for _, element := range *f {
		newFloat = append(newFloat, element.copy())
	}
	return newFloat
}
func (b *boolElements) elems() []Element {
	var newBool []Element
	for _, element := range *b {
		newBool = append(newBool, element.copy())
	}
	return newBool
}

//--------------------------------//

func (s *stringElements) len() int {
	return len(*s)
}
func (i *intElements) len() int {
	return len(*i)
}
func (f *floatElements) len() int {
	return len(*f)
}
func (b *boolElements) len() int {
	return len(*b)
}

//--------------------------------//

func (s *stringElements) records() []string {
	newStr := make([]string, 0)
	for _, element := range *s {
		newStr = append(newStr, element.Records())
	}
	return newStr
}
func (i *intElements) records() []string {
	newStr := make([]string, 0)
	for _, element := range *i {
		newStr = append(newStr, element.Records())
	}
	return newStr
}
func (f *floatElements) records() []string {
	newStr := make([]string, 0)
	for _, element := range *f {
		newStr = append(newStr, element.Records())
	}
	return newStr
}
func (b *boolElements) records() []string {
	newStr := make([]string, 0)
	for _, element := range *b {
		newStr = append(newStr, element.Records())
	}
	return newStr
}

//--------------------------------//

func (s *stringElements) int() []int {
	newInt := make([]int, 0)
	for _, element := range *s {
		newInt = append(newInt, element.Int())
	}
	return newInt
}
func (i *intElements) int() []int {
	newInt := make([]int, 0)
	for _, element := range *i {
		newInt = append(newInt, element.Int())
	}
	return newInt
}
func (f *floatElements) int() []int {
	newInt := make([]int, 0)
	for _, element := range *f {
		newInt = append(newInt, element.Int())
	}
	return newInt
}
func (b *boolElements) int() []int {
	newInt := make([]int, 0)
	for _, element := range *b {
		newInt = append(newInt, element.Int())
	}
	return newInt
}

//--------------------------------//

func (s *stringElements) float() []float64 {
	newFloat := make([]float64, 0)
	for _, element := range *s {
		newFloat = append(newFloat, element.Float())
	}
	return newFloat
}
func (i *intElements) float() []float64 {
	newFloat := make([]float64, 0)
	for _, element := range *i {
		newFloat = append(newFloat, element.Float())
	}
	return newFloat
}
func (f *floatElements) float() []float64 {
	newFloat := make([]float64, 0)
	for _, element := range *f {
		newFloat = append(newFloat, element.Float())
	}
	return newFloat
}
func (b *boolElements) float() []float64 {
	newFloat := make([]float64, 0)
	for _, element := range *b {
		newFloat = append(newFloat, element.Float())
	}
	return newFloat
}

//--------------------------------//

func (s *stringElements) bool() []bool {
	newBool := make([]bool, 0)
	for _, element := range *s {
		newBool = append(newBool, element.Bool())
	}
	return newBool
}
func (i *intElements) bool() []bool {
	newBool := make([]bool, 0)
	for _, element := range *i {
		newBool = append(newBool, element.Bool())
	}
	return newBool
}
func (f *floatElements) bool() []bool {
	newBool := make([]bool, 0)
	for _, element := range *f {
		newBool = append(newBool, element.Bool())
	}
	return newBool
}
func (b *boolElements) bool() []bool {
	newBool := make([]bool, 0)
	for _, element := range *b {
		newBool = append(newBool, element.Bool())
	}
	return newBool
}

//--------------------------------//

//--------------------------------//

func (s *stringElements) dType() string {
	return String
}
func (i *intElements) dType() string {
	return Int
}
func (f *floatElements) dType() string {
	return Float
}
func (b *boolElements) dType() string {
	return Bool
}

package series

import (
	"fmt"
	"testing"
)

func TestNewSeries(t *testing.T) {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number1")
	fmt.Println(&s1)
	s2, _ := NewSeries([]int{1, 2, 3, 4, 6}, Int, "number2")
	fmt.Println(&s2)
	s3, err := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, Float, "number3")
	fmt.Println(&s3, err)
}

func TestSeries_Append(t *testing.T) {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number1")
	fmt.Println(&s1)
	err := s1.Append([]string{"1", "2", "test", "2", "4", "6"})
	fmt.Println(&s1, err)
	err = s1.Append(1)
	fmt.Println(&s1, err)
}

func TestSeries_ChangeType(t *testing.T) {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number1")
	fmt.Println(&s1)
	err := s1.SetType(Int)
	fmt.Println(&s1, err)
}

func TestSeries_Concat(t *testing.T) {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number1")
	s2, _ := NewSeries([]int{1, 2, 3, 4, 6}, Int, "number2")
	fmt.Println(s1, s2)
	err := s1.Concat(*s2)
	fmt.Println(&s1, err)
}

func TestSeries_Format(t *testing.T) {
	s1, _ := NewSeries([]int{1, 2, 3, 4, 6}, Int, "number1")
	s1.Format(func(index int, elem Element) Element {
		e := elem.Int()
		elem.Set(e + index)
		return elem
	})
	fmt.Println(&s1)

	s2, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number2")
	s2.Format(func(index int, elem Element) Element {
		e := elem.Records()
		elem.Set(e + "ioc")
		return elem
	})
	fmt.Println(&s2)
}

func TestSeries_Filter(t *testing.T) {
	s1, _ := NewSeries([]string{"aee", "2", "test", "2", "any", "6"}, String, "number1")
	ns, _ := s1.Filter(Equal, "2")
	fmt.Println(&ns)
	fmt.Println("-----------------------------")
	ns1, _ := s1.Filter(NotEqual, "2")
	fmt.Println(&ns1)
	fmt.Println("-----------------------------")
	ns2, _ := s1.Filter(Contains, "e")
	fmt.Println(&ns2)
	fmt.Println("-----------------------------")
	ns3, _ := s1.Filter(StartsWith, "a")
	fmt.Println(&ns3)
	fmt.Println("-----------------------------")
	ns4, _ := s1.Filter(EndsWith, "e")
	fmt.Println(&ns4)
	fmt.Println("-----------------------------")
}

func TestSeries_Filter2(t *testing.T) {
	s1, _ := NewSeries([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 3, 100, 200, 300, 3, 8, 4}, Int, "number1")
	ns, _ := s1.Filter(Equal, 3)
	fmt.Println(&ns)
	fmt.Println("-----------------------------")
	ns1, _ := s1.Filter(LessThan, 6)
	fmt.Println(&ns1)
	fmt.Println("-----------------------------")
	ns2, _ := s1.Filter(LessOrEqual, 8)
	fmt.Println(&ns2)
	fmt.Println("-----------------------------")
	ns3, _ := s1.Filter(GreaterOrEqual, 10)
	fmt.Println(&ns3)
	fmt.Println("-----------------------------")
	ns4, _ := s1.Filter(In, []int{3, 8, 200})
	fmt.Println(&ns4)
}

func TestSeries_Copy(t *testing.T) {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number1")
	s2 := s1.Copy()
	_ = s2.Append("3")
	fmt.Println(&s1, &s2)
}

func TestSeries_SubSet(t *testing.T) {
	s1, _ := NewSeries([]string{"1", "2", "test", "2", "4", "6"}, String, "number1")
	ns := s1.SubSet(1, 2, 3)
	fmt.Println(&ns)
}

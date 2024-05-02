package pandas

import (
	"fmt"
	"testing"
)

func TestNewSeries(t *testing.T) {
	s1, _ := NewSeries([]string{"1", "2", "3", "2", "4", "6"}, String, "number")
	fmt.Println(s1)
}

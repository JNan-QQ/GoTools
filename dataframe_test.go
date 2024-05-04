package pandas

import (
	"fmt"
	"testing"
)

func TestNewDataFrame(t *testing.T) {
	dataFrame, err := NewDataFrame(map[string]any{
		"filed1": []int{1, 2, 3},
		"filed2": []string{"q", "w", "江南"},
		"filed3": []float64{3.5, 4.5, 5.5},
	}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(dataFrame)
}

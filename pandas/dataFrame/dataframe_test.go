package dataFrame

import (
	"fmt"
	"testing"
)

func TestNewDataFrame(t *testing.T) {
	dataFrame, err := NewDataFrame(map[string]any{
		"filed1": []int{153, 24553, 453},
		"filed2": []string{"你好封大夫", "dddsw44", "江南d第三方阿地方的斯"},
		"filed3": []float64{3.51, 4.54, 5.5},
	}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(dataFrame.String())
}

package main

import (
	"fmt"
	"gitee.com/jn-qq/go-tools/pandas/dataFrame"
)

func main() {
	frame, err := dataFrame.NewDataFrame([]any{
		[]int{153, 24553, 453}, []string{"你好封大夫", "dddsw44", "江南d第三方阿地方的斯"}, []float64{3.51, 4.54, 5.5},
	}, []string{"file1", "file2", "file3"})
	if err != nil {
		panic(err)
	}
	fmt.Println(frame)

}

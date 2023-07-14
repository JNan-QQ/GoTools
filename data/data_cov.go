/*
	data 为数据处理/转化的简单包
*/

package data

import (
	"math"
)

// Int2AAA 将表格列索引（int）转换为列名(string)
func Int2AAA(colIndex int) (colName string) {
	if colIndex <= 26 {
		colName = string(rune(colIndex + 64))
	} else {
		result := colIndex / 26
		mod := math.Mod(float64(colIndex), 26)
		if mod == 0 {
			result -= 1
			mod = 26
		}

		if result > 26 {
			colName = Int2AAA(result) + string(rune(mod+64))
		} else {
			if result != 0 {
				colName += string(rune(result + 64))
			}
			if mod != 0 {
				colName += string(rune(mod + 64))
			}
		}
	}

	return
}

// AAA2Int 将表格列名(string)转换为列索引（int）
func AAA2Int(colName string) (colIndex int) {
	length := len(colName)
	for i, col := range colName {
		colIndex += int(col-64) * int(math.Pow(26, float64(length-1-i)))
	}
	return
}

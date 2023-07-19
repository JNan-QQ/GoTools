/*
	data 为数据处理/转化的简单包
*/

package data

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
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

var numberCn2an = map[string]int{
	"零": 0,
	"〇": 0,
	"一": 1,
	"壹": 1,
	"幺": 1,
	"二": 2,
	"贰": 2,
	"两": 2,
	"三": 3,
	"叁": 3,
	"四": 4,
	"肆": 4,
	"五": 5,
	"伍": 5,
	"六": 6,
	"陆": 6,
	"七": 7,
	"柒": 7,
	"八": 8,
	"捌": 8,
	"九": 9,
	"玖": 9,
}

var unitCn2an = map[string]int{
	"十": 10,
	"拾": 10,
	"百": 100,
	"佰": 100,
	"千": 1000,
	"仟": 1000,
	"万": 10000,
	"亿": 100000000,
}
var numberAn2cn = []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
var unitAn2cn = []string{"", "十", "百", "千", "万", "十万", "百万", "千万", "亿", "十亿", "百亿", "千亿", "万亿", "亿亿"}

// Cn2an 将中文数字转换为阿拉伯数字，不支持小数
func Cn2an(cn string) (int, error) {

	var out int
	var cnList []string

	for _, b := range cn {
		cnList = Insert(cnList, 0, string(b))
	}

	weight := 1
	big := false
	for i, c := range cnList {
		if a, ok := numberCn2an[c]; ok {
			out += a * weight
			weight *= 10
			big = false
		} else if a, ok := unitCn2an[c]; ok {
			// 保权
			if big {
				weight *= a
			} else {
				weight = a
			}
			big = true
		} else {
			return 0, fmt.Errorf("第%d个字符不是中文数字", i)
		}
	}

	return out, nil
}

func An2cn(an int) (string, error) {

	var out string

	str := strconv.Itoa(an)
	for i, a := range str {
		if string(a) == "0" {
			out += "零"
			continue
		}
		if index := len(str) - 1 - i; index < len(unitAn2cn) {
			a_n, _ := strconv.Atoi(string(a))
			out += numberAn2cn[a_n] + unitAn2cn[index]
		} else {
			return "", fmt.Errorf("数字太大了")
		}
	}
	// 去重 “零”
	re := regexp.MustCompile(`零{2,}`)
	out = strings.Trim(re.ReplaceAllString(out, "零"), "零")

	return out, nil
}

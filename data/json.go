package data

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

type JsonString string

// Check 检查是否为json字符串
func (j JsonString) Check() bool {
	return json.Valid([]byte(j))
}

// 格式化关系
func (j JsonString) formatRelation(relation string) (result []map[string]string, err error) {

	for _, s := range strings.Split(relation, ".") {
		isSlice, err := regexp.MatchString(`\[\d+]`, s)
		if err != nil {
			return nil, err
		}

		// 判断是否为切片
		if isSlice {
			result = append(result, map[string]string{
				"name": regexp.MustCompile(`\d+`).FindString(s),
				"type": "array",
			})
		} else {
			result = append(result, map[string]string{
				"name": s,
				"type": "map",
			})
		}
	}

	return
}

// Find 查询方法
func (j JsonString) Find(relational string) (any, error) {
	relation, err := j.formatRelation(relational)
	if err != nil {
		return nil, err
	}
	var crash any
	for _, m := range relation {
		switch m["type"] {
		case "array":
			if crash == nil {
				c := make([]any, 0)
				if err := json.Unmarshal([]byte(j), &c); err != nil {
					return nil, err
				}
				index, err := strconv.Atoi(m["name"])
				if err != nil {
					return nil, err
				}
				crash = c[index]
			} else {
				index, err := strconv.Atoi(m["name"])
				if err != nil {
					return nil, err
				}
				crash = crash.([]any)[index]
			}
		case "map":
			if crash == nil {
				c := make(map[string]any)
				if err := json.Unmarshal([]byte(j), &c); err != nil {
					return nil, err
				}
				crash = c[m["name"]]
			} else {
				crash = crash.(map[string]any)[m["name"]]
			}
		}
	}
	return crash, nil
}

// FindString 返回字符串
func (j JsonString) FindString(relational string) string {
	find, err := j.Find(relational)
	if err != nil {
		return ""
	}
	return find.(string)
}

// FindInt 返回整数
func (j JsonString) FindInt(relational string) int {
	find, err := j.Find(relational)
	if err != nil {
		return 0
	}
	return int(find.(float64))
}

// FindFloat 返回浮点数
func (j JsonString) FindFloat(relational string) float64 {
	find, err := j.Find(relational)
	if err != nil {
		return 0.0
	}
	return find.(float64)
}

// FindBool 返回布尔值
func (j JsonString) FindBool(relational string) bool {
	find, err := j.Find(relational)
	if err != nil {
		return false
	}
	return find.(bool)
}

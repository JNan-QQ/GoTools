/*
	对 gota 简单扩展
*/

package pandas

import (
	"fmt"
	"gitee.com/jn-qq/go-tools/data"
	"github.com/go-gota/gota/dataframe"
	"github.com/shakinm/xlsReader/xls"
	"github.com/xuri/excelize/v2"
	"path/filepath"
	"strings"
)

type DataFrame struct {
	dataframe.DataFrame
}

const (
	XLSX = "xlsx"
	XLS  = "xls"
)

func Read(path string, sheet ...string) (df DataFrame) {
	switch strings.Split(filepath.Base(path), ".")[1] {
	case XLS:
		return readFormXLS(path, sheet...)
	case XLSX:
		return readFromXLSX(path, sheet...)
	default:
		df.Err = fmt.Errorf("不支持的文件类型")
		return
	}
}

func readFromXLSX(path string, sheet ...string) (df DataFrame) {
	// 读取文档
	f, err := excelize.OpenFile(path)
	if err != nil {
		df.Err = err
		return
	}
	defer func(f *excelize.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)

	// 设置默认读取工作部名称
	if sheet == nil {
		sheet = append(sheet, "Sheet1")
	}

	var header []string
	var xlsxData [][]string

	// 遍历sheet
	for _, sn := range f.GetSheetList() {

		// 跳过
		if !data.Contains(sheet, sn) {
			continue
		}

		// 获取表头与行列
		var colNum []int
		var hd []string

		// 行迭代器
		rows, err := f.Rows(sn)
		if err != nil {
			df.Err = err
			return
		}
		for rows.Next() {
			row, err := rows.Columns()
			if err != nil {
				df.Err = err
				return
			}

			if data.IsEmpty(row) {
				continue
			}

			// 获取表头与行列
			if hd == nil {
				for i, r := range row {
					if len(colNum) == 0 {
						if r != "" {
							colNum = append(colNum, i)
						}
					} else if len(colNum) == 1 {
						if r == "" {
							colNum = append(colNum, i)
							break
						}
						if i == len(row)-1 {
							colNum = append(colNum, i)
						}
					}
				}

				hd = append(hd, row[colNum[0]:colNum[1]]...)
				if header == nil {
					header = append(header, hd...)
				} else {
					if !data.Equal(header, hd) {
						df.Err = fmt.Errorf("读取多工作簿发现表头不一致,请确认")
						return
					}
				}
				continue
			}

			// 获取数据体
			xlsxData = append(xlsxData, row[colNum[0]:colNum[1]])

		}
		_ = rows.Close()
	}

	df.DataFrame = dataframe.LoadRecords(append([][]string{header}, xlsxData...))

	return
}

func readFormXLS(path string, sheet ...string) (df DataFrame) {
	workbook, err := xls.OpenFile(path)
	if err != nil {
		df.Err = err
		return
	}
	// 设置默认读取工作部名称
	if sheet == nil {
		sheet = append(sheet, "Sheet1")
	}

	var header []string
	var xlsxData [][]string

	for _, sn := range workbook.GetSheets() {
		// 跳过
		if !data.Contains(sheet, sn.GetName()) {
			continue
		}

		// 获取表头与行列
		var colNum []int
		var hd []string
		for _, rw := range sn.GetRows() {
			var row []string
			for _, cellData := range rw.GetCols() {
				row = append(row, cellData.GetString())
			}

			// 空行跳过
			if data.IsEmpty(row) {
				continue
			}

			// 获取表头与行列
			if hd == nil {
				for i, r := range row {
					if len(colNum) == 0 {
						if r != "" {
							colNum = append(colNum, i)
						}
					} else if len(colNum) == 1 {
						if r == "" {
							colNum = append(colNum, i)
							break
						}
						if i == len(row)-1 {
							colNum = append(colNum, i)
						}
					}
				}

				hd = append(hd, row[colNum[0]:colNum[1]]...)
				if header == nil {
					header = append(header, hd...)
				} else {
					if !data.Equal(header, hd) {
						df.Err = fmt.Errorf("读取多工作簿发现表头不一致,请确认")
						return
					}
				}
				continue
			}

			// 获取数据体
			xlsxData = append(xlsxData, row[colNum[0]:colNum[1]])

		}

	}

	df.DataFrame = dataframe.LoadRecords(append([][]string{header}, xlsxData...))

	return

}

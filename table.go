package p2mfmt

import (
	"fmt"
	"strings"
)

// 表格
type Table struct {
	Cells [][]tableCell
}

// 儲存格
type tableCell struct {
	Text string

	ColSpan int // colspan 的值
	RowSpan int // rowspan 的值

	IsRowMerged bool // 是否被縱向合併
	IsColMerged bool // 是否被橫向合併

	IsHead bool
}

// 輸入 puki wiki 的 table row
// 將之轉為 mediawiki 格式
func (t *Table) ParseRow(pukiWikiText string) {
	cells := []tableCell{}

	texts := strings.Split(pukiWikiText, "|")

	for i, v := range texts {
		// puki wiki 的 table 格式為
		//
		// |~ 標題               |~ 標題  |~ 標題  |
		// | 儲存格(rowspan = 2) | 儲存格 | 儲存格 |
		// |~|                     儲存格 | 儲存格 |
		// | 儲存格(colspan=2)          |>| 儲存格 |
		//
		if i == 0 || i == len(texts)-1 {
			continue
		}

		cell := tableCell{ColSpan: 1, RowSpan: 1}

		if v == "~" {
			// 處理 rowspan
			cell.IsRowMerged = true
			for index := len(t.Cells) - 1; index >= 0; index-- {
				if !t.Cells[index][i-1].IsRowMerged {
					t.Cells[index][i-1].RowSpan++
					break
				}
			}

		} else if v == ">" {
			// 處理 colspan
			cell.IsColMerged = true

			for index := i - 2; index >= 0; index-- {
				if !cells[index].IsColMerged {
					cells[index].ColSpan++
					break
				}
			}

		} else {
			if len(v) > 0 && v[0:1] == "~" {
				cell.IsHead = true
				cell.Text = ConvertLine(v[1:])
			} else {
				cell.Text = ConvertLine(v)
			}

		}
		cells = append(cells, cell)
	}
	t.Cells = append(t.Cells, cells)
}

func (t Table) String() string {
	s := `{| class="wikitable" style="margin: 0 auto;" ` + "\n"

	for _, row := range t.Cells {
		// 先檢測是否整列都是表格標題

		s += "|-\n"

		allIsHead := true
		for _, cell := range row {
			if cell.IsRowMerged || cell.IsColMerged {
				continue
			}
			if !cell.IsHead {
				allIsHead = false
				break
			}
		}

		for index, cell := range row {
			if cell.IsRowMerged || cell.IsColMerged {
				continue
			}

			// 處理 attributes
			attrs := ""
			if cell.ColSpan > 1 {
				attrs += fmt.Sprintf(`colspan = "%d" `, cell.ColSpan)
			}
			if cell.RowSpan > 1 {
				attrs += fmt.Sprintf(`rowspan = "%d"`, cell.RowSpan)
			}

			// 處理表格標題
			if allIsHead {
				s += "! "
				if attrs != "" {
					s += attrs + " | "
				}
				s += cell.Text

				if index == len(row)-1 {
					// 此列最後一個儲存格
					s += " \n"
				} else {
					s += " !"
				}

			} else if cell.IsHead {
				s += `! scope ` + attrs + ` | ` + cell.Text
				s += "\n"
			} else {
				s += "| "
				if attrs != "" {
					s += attrs + " | "
				}
				s += cell.Text

				if index == len(row)-1 {
					// 此列最後一個儲存格
					s += " \n"
				} else {
					s += " |"
				}
			}
		}

	}

	s += "|}"

	return s
}

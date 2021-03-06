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
	Style map[string]string // 儲存格樣式

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
			var text = "" //儲存格內文
			
			if len(v) > 0 && v[0:1] == "~" {
				cell.IsHead = true
				text = v[1:]
			} else {
				text = v
			}
			
			
			// 處理儲存格樣式
			styles, text := ConvertFirstOfLine(text)
			text = ConvertLine(text)
			
			cell.Text = text
			cell.Style = styles
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

		wroteCellOfRow := -1 // 此列第幾個被繪製的儲存格

		for index, cell := range row {
			if cell.IsRowMerged || cell.IsColMerged {
				continue
			}

			wroteCellOfRow++

			// 處理 attributes
			attrs := ""
			if cell.ColSpan > 1 {
				attrs += fmt.Sprintf(`colspan = "%d" `, cell.ColSpan)
			}
			if cell.RowSpan > 1 {
				attrs += fmt.Sprintf(`rowspan = "%d" `, cell.RowSpan)
			}
			if cell.Style != nil {
				s := ""
				for k, v:= range(cell.Style){
					s += k + ": " + v + "; "
				}
				if s != ""{
					attrs += fmt.Sprintf(`style = "%s" `, s)	
				}
				
			}
			
			

			// 處理表格標題
			if allIsHead {
				if wroteCellOfRow == 0 {
					s += "! "
				} else {
					s += "!! "
				}

				if attrs != "" {
					s += attrs + " | "
				}
				s += cell.Text

				if index == len(row)-1 {
					// 此列最後一個儲存格
					s += " "
				} else {
					s += " "
				}

			} else if cell.IsHead {
				s += `! scope ` + attrs + ` | ` + cell.Text
				s += "\n"
			} else {
				if wroteCellOfRow == 0 {
					s += "| "
				} else {
					s += "|| "
				}

				if attrs != "" {
					s += attrs + " | "
				}
				s += cell.Text

				if index == len(row)-1 {
					// 此列最後一個儲存格
					s += " "
				} else {
					s += " "
				}
			}
		}
		s += "\n"

	}

	s += "|}"

	return s
}

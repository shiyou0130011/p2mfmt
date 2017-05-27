package p2mfmt

import (
	"bytes"
	"regexp"
	"strings"
)


//
func Convert(pukiText string) (string, []string) {
	splitedWiki := strings.Split(pukiText, "\n")

	var b bytes.Buffer
	var table Table

	categories := []string{}

	for _, v := range splitedWiki {
		v = regexp.MustCompile(`#navi\(.*\)`).ReplaceAllStringFunc(
			v,
			func(s string) string {
				s = strings.NewReplacer("#navi(", "", ")", "").Replace(s)
				categories = append(categories, s)

				return ""
			},
		)
		
		v = regexp.MustCompile(`#comment\(.*\)`).ReplaceAllStringFunc(
			v,
			func(s string) string {
				return "<!-- コメントの挿入 -->\n<!-- " + s + " -->"
			},
		)

		if v != "" && v[0:1] == "|" && v[len(v)-1:] == "|" {
			table.ParseRow(v)
		} else {
			if table.Cells != nil {
				b.WriteString(table.String() + "\n")
				table.Cells = nil
			}
			b.WriteString(ConvertLine(v) + "\n")
		}
	}

	result := string(b.Bytes())

	if strings.Contains(result, "<ref>") {
		result += "\n== 備註 ==\n<references/>\n\n"
	}

	return result, categories
}

// 轉換任意時候均可轉換的內容
// 例如換行符號、引用/註解、圖片……等
func ConvertLine(wikiText string) string {
	if wikiText == "" {
		return ""
	}

	
	// 轉換標題
	if wikiText[0:1] == "*" {
		noneTitleIndexes := regexp.MustCompile(`[^\* ]`).FindStringIndex(wikiText)
		if noneTitleIndexes != nil {
			title := wikiText[noneTitleIndexes[0]:]
			title = regexp.MustCompile(`\[#[a-zA-Z0-9]*\]`).ReplaceAllLiteralString(title, "")

			h := wikiText[0:noneTitleIndexes[0]]
			h = strings.Replace(h, "*", "=", -1)
			wikiText = h + title + h
		}
	}
	// 轉換引用

	// 轉換清單
	if wikiText[0:1] == "-" {
		noneListIndexes := regexp.MustCompile(`[^\- ]`).FindStringIndex(wikiText)
		if noneListIndexes != nil {
			l := strings.Trim(wikiText[0:noneListIndexes[0]], " \t")
			if(len(l) != 0){
				l = strings.Repeat(":", len(l) - 1) + "*"
			}
			
			wikiText = l + " " + wikiText[noneListIndexes[0]:]
		}
	}

	//
	// 轉換超連結
	// 一般超連結可不用轉換，但是連結文字和頁面標題不同時兩者格式不同
	// puki wiki 格式是 [[連結文字>頁面標題]]
	// media wiki 則是 [[頁面標題|連結文字]]
	//
	wikiText = regexp.MustCompile(`\[\[[^\>]*>[^\>]*\]\]`).ReplaceAllStringFunc(
		wikiText,
		func(s string) string {
			s = strings.NewReplacer("[[", "", "]]", "").Replace(s)
			vars := strings.Split(s, ">")

			pageTitle := vars[1]
			linkText := vars[0]

			if strings.Contains(pageTitle, "://") {
				// 此為超連結
				return `[` + pageTitle + " " + linkText + `]`
			}

			return `[[` + pageTitle + "|" + linkText + `]]`
		},
	)

	// 轉換圖片
	imgReg := regexp.MustCompile(`#imgr\([\.A-Za-z0-9\/,%]*\)`)
	wikiText = imgReg.ReplaceAllStringFunc(wikiText, func(s string) string {
		s = strings.NewReplacer("#imgr(", "", ")", "").Replace(s)

		imgVars := strings.Split(s, ",")

		img := Image{}

		for i, v := range imgVars {
			if i == 0 {
				img.Url = v
			} else if regexp.MustCompile("[0-9]").MatchString(v) {
				img.Size = v
			} else if v == "around" {
				img.Format = "thumb"
			} else if v == "center" || v == "right" || v == "left" || v == "none" {
				img.Align = v
			}
		}

		return img.String()
	})

	// 轉換刪除線
	delIndex := -1
	wikiText = regexp.MustCompile(`%%`).ReplaceAllStringFunc(
		wikiText,
		func(s string) string {
			delIndex++
			if delIndex%2 == 0 {
				return "<del>"
			} else {
				return "</del>"
			}

		},
	)

	// 轉換 換行符號
	wikiText = strings.Replace(wikiText, "&br;", "<br/>", -1)
	wikiText = strings.Replace(wikiText, "~\n", "<br/>\n", -1)
	if wikiText[len(wikiText):] == "~" {
		wikiText = wikiText[0:len(wikiText)-1] + "\n\n"
	}

	// 轉換 引用/ 註解 （<ref>）
	wikiText = strings.Replace(wikiText, "((", "<ref>", -1)
	wikiText = strings.Replace(wikiText, "))", "</ref>", -1)

	// 轉換粗斜體
	wikiText = strings.Replace(wikiText, "''", "<b>", -1)
	wikiText = strings.Replace(wikiText, "'''", "<i>", -1)
	wikiText = strings.Replace(wikiText, "<b>", "'''", -1)
	wikiText = strings.Replace(wikiText, "<i>", "''", -1)
	// TODO 處理顏色、背景、置左/中/右
	
	
	
	return wikiText
}

package p2mfmt

import (
	"bytes"
	"regexp"
	"strings"
)

//
func Convert(pukiText string) (string, []string) {
	splitedWiki := strings.Split(
		strings.Replace(pukiText, "\r\n", "\n", -1),
		"\n",
	)

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

		if v != "" && v[0:1] == "|" && (v[len(v)-1:] == "|" || v[len(v)-2:] == "|h") {
			if v[len(v)-2:] == "|h" {
				v = v[0:len(v)-1]
			}
			
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
	var styles map[string]string
	var listPrefix string

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
			listPrefix = strings.Trim(wikiText[0:noneListIndexes[0]], " \t")
			if len(listPrefix) != 0 {
				listPrefix = strings.Repeat(":", len(listPrefix)-1) + "*"
			}
			
			wikiText = wikiText[noneListIndexes[0]:]
			
			styles, wikiText = ConvertFirstOfLine(wikiText)
		}
	}else {
		styles, wikiText = ConvertFirstOfLine(wikiText)
	}

	//
	// 轉換超連結
	// 一般超連結可不用轉換，但是連結文字和頁面標題不同時兩者格式不同
	// puki wiki 格式是 [[連結文字>頁面標題]]
	// media wiki 則是 [[頁面標題|連結文字]]
	//
	wikiText = regexp.MustCompile(`\[\[[^\>\[\]]*>[^\>\[\]]*\]\]`).ReplaceAllStringFunc(
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

	// 轉換回復的時間格式
	// pukiwiki 在進行回覆的時候，格式如下
	//
	// - 我的回覆 -- &new{2009-01-26 (月) 01:02:59};
	//
	// 前面已經轉換 list 了，這邊是將 &new{2009-01-26 (月) 01:02:59}; 轉換成
	//
	// <time>2009-01-26 (月) 01:02:59</time>
	//

	wikiText = regexp.MustCompile("&new{.*};").ReplaceAllStringFunc(
		wikiText,
		func(s string) string {
			return "<time>" + s[5:len(s)-2] + "</time>"
		},
	)

	// 轉換註解
	if strings.Contains(wikiText, "//") {
		s := strings.Split(wikiText, "//")
		hasComment := false

		for i, v := range s {
			if len(v) > 0 && v[len(v)-1:] == ":" {
				s[i] += "//"
			} else if i != len(s)-1 {
				s[i] += "<!-- "
				hasComment = true
				break
			}

		}
		wikiText = strings.Join(s, "")
		if hasComment {
			wikiText += " -->"
		}
	}
	
	// 防止簽名
	// ~~~ 在 mediawiki 是著名是誰寫的用的。
	// 但在 pukiwiki，這沒有用途。
	// 避免在轉換時，也將 ~~~ 轉成 mediawiki 簽名。
	wikiText = strings.Replace(wikiText, "~~~", "<nowiki>~~~</nowiki>", -1)
	
	
	// 處理樣式
	if styles != nil{
		_, setTextAlign := styles["text-align"]
		
		styleText := ""
		for i, v:= range styles {
			styleText += i + ": " + v + "; "
		}
		
		if setTextAlign{
			wikiText = `<div style="` + styleText + `">` + wikiText + `</div>`
		}else if styleText != ""{
			
			wikiText = `<span style="` + styleText + `">` + wikiText + `</span>`
		}
	}
	
	if listPrefix != ""{
		wikiText = listPrefix + " " + wikiText
	}
	
	return wikiText
}

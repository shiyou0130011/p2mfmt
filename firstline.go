package p2mfmt

import(
	"strings"
	"regexp"
)
// Pukiwiki 可以每行或每個儲存格的最開始加上下面幾項，為整行或整個儲存格設置樣式
// 
// LEFT:
// 
// CENTER:
// 
// RIGHT:
// 
// BGCOLOR(色):
// 
// COLOR(色):
// 
// SIZE(サイズ):
//
// 此 function 是將之轉換成 CSS style map
//
func ConvertFirstOfLine(pukiText string) (map[string]string, string) {
	pukiText = strings.Trim(pukiText, "\t \r")

	style := map[string]string{}

	for {
		indexes := regexp.MustCompile(`(LEFT|CENTER|RIGHT|BGCOLOR|COLOR\([#0-9a-zA-Z,]+\)|SIZE\([0-9]+\)):`).FindStringIndex(pukiText)

		if indexes == nil || indexes[0] != 0 {
			break
		}

		s := pukiText[indexes[0]:indexes[1]]
		pukiText = pukiText[indexes[1]:]

		switch s {
		case "LEFT:":
			style["text-align"] = "left"
		case "CENTER:":
			style["text-align"] = "center"
		case "RIGHT:":
			style["text-align"] = "right"
		default:
			v := s[strings.Index(s, "(")+1 : len(s)-2]

			styleName := s[0:strings.Index(s, "(")]
			switch styleName {
			case "COLOR":
				colors := append(strings.Split(v, ","), "", "")
				if colors[0] != "" {
					style["color"] = colors[0]
				}
				if colors[1] != "" {
					style["background-color"] = colors[1]
				}

			case "SIZE":
				style["font-size"] = v + "px"
			}

		}
	}

	return style, pukiText
}
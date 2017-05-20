package p2mfmt

import "strings"

// 圖片
type Image struct {
	Url    string
	Size   string
	Align  string
	Format string
}

func (i Image) String() string {
	s := "[[File:"

	urlSplit := strings.Split(i.Url, "/")
	s += urlSplit[len(urlSplit)-1][0:]

	if i.Size != "" {
		s += "|" + i.Size
	}

	if i.Format != "" {
		s += "|" + i.Format
	}
	if i.Align != "" {
		s += "|" + i.Align
	}
	s += "]]"

	return s
}

package createform

import (
	"strings"

	"github.com/gomarkdown/markdown"
)

func CreateHeaderHtml(md string, title string, amber bool) string {
	html := getHtmlBody(md)
	file := "<!DOCTYPE html>\n" + "<head>\n<meta charset=\"UTF-8\">\n" + styles() + "\n</head>\n" + banner(title, amber) + html
	return file
}

func CreateCoverHtml(md string, title string) string {
	html := getHtmlBody(md)
	file := "<!DOCTYPE html>\n" + "<head>\n<meta charset=\"UTF-8\">\n" + styles() + "\n</head>\n" + cover(strings.ToUpper(title)) + html

	return file
}

func getHtmlBody(md string) string {
	md =  "<div class=\"content\">\n\n" + md + "</div>"
	html := markdown.ToHTML([]byte(md), nil, nil)

	return string(html)
}

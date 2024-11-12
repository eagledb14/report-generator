package createform

import (
	"github.com/gomarkdown/markdown"
)

func CreateHtml(md string, title string, amber bool) string {
	md =  "<div class=\"content\">\n\n" + md + "</div>"
	html := markdown.ToHTML([]byte(md), nil, nil)

	file := "<!DOCTYPE html>" + "<head>\n<meta charset=\"UTF-8\">\n" + styles() + "\n</head>\n" + banner(title, amber) + string(html)
	return file
}



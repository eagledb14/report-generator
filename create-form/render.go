package createform

import (
	"github.com/gomarkdown/markdown"
)

func CreateHtml(md string, title string) string {
	md =  "<div class=\"content\">\n\n" + md + "</div>"
	html := markdown.ToHTML([]byte(md), nil, nil)

	file := "<!DOCTYPE html>" + "<head>" + banner(title, true) + "</head>" + string(html)
	return file
}



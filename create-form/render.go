package createform

import (
	"github.com/gomarkdown/markdown"
)

func CreateHtml(md string, title string, amber bool) string {
	md =  "<div class=\"content\">\n\n" + md + "</div>"
	html := markdown.ToHTML([]byte(md), nil, nil)

	file := "<!DOCTYPE html>" + "<head>" + banner(title, amber) + "</head>" + string(html)
	return file
}



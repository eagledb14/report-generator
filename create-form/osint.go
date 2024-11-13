package createform

import "github.com/eagledb14/form-scanner/templates"

// "github.com/eagledb14/form-scanner/templates"
// "github.com/eagledb14/form-scanner/types"

type Osint struct {

}

func (o *Osint) CreateMarkdown() string {
	data := struct {

	} {

	}

	const page = `

`

	return templates.ExecuteText("osintmd", page, data)
}

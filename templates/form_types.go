package templates

//here I'll put the different templates of the generic forms for openport

import (
	"github.com/eagledb14/form-scanner/alerts"
	"github.com/eagledb14/form-scanner/types"
)


func form(formType types.Form, events []*alerts.Event) string {
	// summary := ""
	// body := ""
	// make a match on which type is passed int
	data := struct {
		Summary string
		Body string
	} {
		Summary: "this is a placeholder summary",
		Body: "this is a placeholder body",
	}

	const page = `
	<article>
		<form hx-post="" hx-target="body">
			<fieldset>
                    <label>
                        Form Number
                        <input name="formNumber"/>
                    </label>

					<label>
						Threat Type
						<input name="threat" value="T1133 External Remote Services"/>
					</label>

					<label>
						Summary Paragraph
						<textarea name="summary">{{.Summary}}</textarea>
					</label>

					<label>
						Body Paragraph
						<textarea name="body">{{.Body}}</textarea>
					</label>

					<label>
						Additional References
						<textarea name="reference"></textarea>
					</label>

					<label>TLP Alert</label>
					<label>
						<input type="radio" value="amber" name="tlp" checked/>
						Amber
					</label>
					<label>
						<input type="radio" value="green" name="tlp"/>
						Green
					</label>

					<hr>
					<div class="grid">
						<input type="submit">
						<input type="reset">
					</div>
			</fieldset>
		</form>
	</article>
	`

	return Execute("form", page, data)
}

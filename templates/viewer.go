package templates

import (
	"github.com/eagledb14/form-scanner/types"
)

func MarkdownViewer(state *types.State) string {
	data := struct {
		Markdown string
		Name     string
	}{
		Markdown: state.Markdown,
		Name: state.Name,
	}

	const page = `
    <script>
		function autoResize(textarea) {
		  textarea.style.height = 'auto';
		  textarea.style.height = textarea.scrollHeight + 'px';
		  textarea.style.resize = 'none';
		  textarea.style.overflow = 'hidden';
		}

		function download() {
		  var link = document.createElement('a');
		  link.download = "{{.Name}}";
		  link.href = '/create';
		  link.target= '_blank';
		  link.click();
		  link.remove();
		}

		autoResize(document.getElementById('markdown'));
    </script>
    <h1>Page Preview</h1>
    <article>
		<form hx-post="/preview" hx-swap="none" hx-on::after-request="download()">
		<fieldset>
		  <textarea id="markdown" name="markdown">{{.Markdown}}</textarea>
		  
		  <div class="grid">
			<input type="submit" value="Submit" download/>
			<input type="reset"/>
		  </div>
		</fieldset>
      </form>
    </article>
  `

	return ExecuteText("markdownViewer", page, data)
}

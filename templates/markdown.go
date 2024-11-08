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
		  textarea.style.resize = 'vertical';
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

		async function imageToBase64(event) {
			const items = await navigator.clipboard.read()
			if (items.length <= 0) {
				return
			}
			item = items[0]

			if (item.types.includes("image/png") || item.types.includes("image/jpeg")) {
				const blob = await item.getType(item.types[1]); // Get the image as a Blob
				const reader = new FileReader();

				reader.onload = function (e) {
					const base64Image = e.target.result; // Base64 encoded image string
					console.log(base64Image)

					// Create Markdown image syntax
					const markdownImage = "![Pasted Image](" + base64Image + ")";

					// Assuming you have a textarea to insert the markdown text
					const textarea = document.getElementById("markdown");

					// Insert the markdown at the current cursor position
					const cursorPos = textarea.selectionStart;
					const textBefore = textarea.value.substring(0, cursorPos);
					const textAfter = textarea.value.substring(cursorPos);
					textarea.value = textBefore + markdownImage + textAfter;

					// Set cursor position to the end of the newly added markdown
					textarea.selectionStart = cursorPos + markdownImage.length;
					textarea.selectionEnd = cursorPos + markdownImage.length;
				};

				reader.readAsDataURL(blob); // Convert the image Blob to Base64
			}
		}

		document.getElementById("markdown").addEventListener("paste", (event) => {
			imageToBase64(event)
		})
    </script>
    <h1>Markdown Preview</h1>
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

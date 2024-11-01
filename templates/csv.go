package templates

func Csv() string {

	const page = `
	<script>
		function download() {
		  var link = document.createElement('a');
		  link.download = "{{.Name}}";
		  link.href = '/csv/create';
		  link.target= '_blank';
		  link.click();
		  link.remove();
		}
	</script>
	<h1>CSV</h1>
	<article>
	    <form hx-post="/csv" hx-swap="none" hx-indicator="#load" hx-on::after-request="download()">
		<fieldset>
		    <label>
			    Organization Name
			    <input name="orgName"/>
		    </label>
		    <label>
			    IP Addresses
			    <input name="ipAddress" />
		    </label>

		    <div id="load" class="htmx-indicator center" aria-busy="true">Loading...</div>
		    <div class="grid">
			    <input type="submit" value="Submit"/>
			    <input type="reset"/>
		    </div>
		</fieldset>
	    </form>
	</article>
    `

	return ExecuteText("csv", page, nil)
}

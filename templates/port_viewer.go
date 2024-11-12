package templates

func PortViewer() string {
	data := struct {
	}{}

	const page = `
        <h1>Port Viewer</h1>
		<article>
			<form hx-post="/portview" hx-target="body" hx-indicator="#load">
				<fieldset>
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

	return Execute("portViewer", page, data)
}


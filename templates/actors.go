package templates

import (

)

func Actors() string {
    data := struct {

    } {

    }

    const page = `
        <h1>Actor</h1>
		<article>
			<form hx-post="/actor" hx-target="body" hx-push-url="/preview">
				<fieldset>
					<label>
						Primary Name
						<input name="name" />
					</label>
					<label>
						Alias
						<input name="alias" />
					</label>
					<label>
						First Seen Activity [ DD MM, YYYY ]
						<input name="date" />
					</label>
					<label>
						Country of Origin
						<input name="country" />
					</label>
					<label>
						Motivation
						<input name="motivation" />
					</label>
					<label>
						Targeting
						<input name="target" />
					</label>
					<label>
						Malware Name
						<input name="malware" />
					</label>
					<label>
						Third Party Reporting
						<input name="report" />
					</label>

					<hr>
					<label>Assessment Confidence</label>
					<label>
						<input type="radio" name="confidence" value="High" checked/>
						High
					</label>
					<label>
						<input type="radio" name="confidence" value="Medium" />
						Medium
					</label>
					<label>
						<input type="radio" name="confidence" value="Low" />
						Low
					</label>
					<hr>

					<label>
						Exploits
						<textarea name="exploits"></textarea>
					</label>
					<label>
						Exploits
						<textarea name="exploits"></textarea>
					</label>
					<label>
						Attack Chain Summary
						<textarea name="summary"></textarea>
					</label>
					<label>
						Capabilities
						<textarea name="capabilities"></textarea>
					</label>
					<label>
						Detection Names
						<textarea name="detection"></textarea>
					</label>
					<label>
						TTPS 
						<textarea name="ttps"></textarea>
					</label>
					<label>
						Detection Names
						<textarea name="detection"></textarea>
					</label>
					<label>
						Infrastructure
						<textarea name="infra"></textarea>
					</label>

					<hr>
					<div class="grid">
						<input type="submit" value="Submit">
						<input type="reset">
					</div>
				</fieldset>
			</form>
		</article>
        `

    return Execute("actor", page, data)
}

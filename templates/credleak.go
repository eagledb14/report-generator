package templates

import (

)

func CredLeak() string {
    data := struct {

    } {

    }

    const page = `
        <h1>Cred Leak</h1>
		<article>
			<form hx-post="/credleak" hx-target="body">
				<fieldset>
					<div class="grid">
						<label>
							Organization Name
							<input name="orgName"/>
						</label>
						<label>
							Form Number
							<input name="formNumber"/>
						</label>
					</div>
					<label>
						Victim Org Name
						<input name="victimOrg"/>
					</label>

					<label>
						<input type="checkbox" name="leaks"/>
						Multiple Leaks Found
					</label>

					<label>
						<input type="checkbox" name="creds"/>
						Multiple Creadentials
					</label>
					
					<hr>
					<label>For password add either:</label>
					<label>
						<input type="radio" value="The passwords have been obfuscated to show only the first two letters." name="password"/>
						The passwords have been obfuscated to show only the first two letters.  
					</label>
					<label>
						<input type="radio" value="The passwords are not included due to being posted in plain text." name="password"/>
						The passwords are not included due to being posted in plain text.
					</label>
					<label>
						<input type="radio" value="The passwords are not included due to the threat actor not disclosing them." name="password"/>
						The passwords are not included due to the threat actor not disclosing them.
					</label>
					<label>
						<input type="radio" value="No Passwords were included." name="password" checked/>
						No Passwords were included.
					</label>
					<hr>

					<label>
						IP Address
						<input name="ipAddress" />
					</label>
					
					<label>
						Insert Username: Password
						<textarea name="userPass"></textarea>
					</label>

					<label>
						Additional Information
						<textarea name="addInfo"></textarea>
					</label>

					<label>
						Additional References
						<textarea name="reference"></textarea>
					</label>

					<hr>
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

    return Execute("credleak", page, data)
}

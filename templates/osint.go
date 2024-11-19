package templates

import (

)

func Osint() string {
    data := struct {

    } {

    }


    const page = `
<h1>Osint</h1>
<article>
    <form hx-post="/osint" hx-target="body" hx-push-url="preview" hx-indicator="#load">
	<fieldset>
	    <label>
		Organization Name
		<input name="orgName"/>
	    </label>
	    <label>
		Organization Url
		<input name="url"/>
	    </label>


	    <label>
		Number of Vulnerable Urls 
		<input name="vulnerableUrls"  />
	    </label>

	    <hr>

	    <label>
		In Scope IP Addresses
		<input name="inScope">
	    </label>
	    <label>
		Out of Scope IP Addresses
		<input name="outScope">
	    </label>

	    <hr>
	    <label>Asset Severity</label>
	    <div class="grid">
		<label>
		    <input type="radio" value="LOW" name="assetSeverity" checked>
		    LOW
		</label>
		<label>
		    <input type="radio" value="MEDIUM" name="assetSeverity">
		    MEDIUM
		</label>
		<label>
		    <input type="radio" value="HIGH" name="assetSeverity">
		    HIGH
		</label>
		<label>
		    <input type="radio" value="CRITICAL" name="assetSeverity">
		    CRITICAL
		</label>
	    </div>

	    <hr>
	    <label>Email And Password Severity</label>
	    <div class="grid">
		<label>
		    <input type="radio" value="LOW" name="accountSeverity" checked>
		    LOW
		</label>
		<label>
		    <input type="radio" value="MEDIUM" name="accountSeverity">
		    MEDIUM
		</label>
		<label>
		    <input type="radio" value="HIGH" name="accountSeverity">
		    HIGH
		</label>
		<label>
		    <input type="radio" value="CRITICAL" name="accountSeverity">
		    CRITICAL
		</label>
	    </div>

	    <hr>
	    <label>Website Severity</label>
	    <div class="grid">
		<label>
		    <input type="radio" value="LOW" name="websiteSeverity" checked>
		    LOW
		</label>
		<label>
		    <input type="radio" value="MEDIUM" name="websiteSeverity">
		    MEDIUM
		</label>
		<label>
		    <input type="radio" value="HIGH" name="websiteSeverity">
		    HIGH
		</label>
		<label>
		    <input type="radio" value="CRITICAL" name="websiteSeverity">
		    CRITICAL
		</label>
	    </div>

	    <hr>
	    <label>
		    Recorded Future Credentials
		    <textarea name="recordedFutureCreds"></textarea>
	    </label>
	    <label>
		    Additional Credentials
		    <textarea name="otherCreds"></textarea>
	    </label>

	    <div id="load" class="htmx-indicator center" aria-busy="true">Loading...</div>
	    <div class="grid">
		    <input type="submit" value="Submit">
		    <input type="reset">
	    </div>
	</fieldset>
    </form>
</article>
`

    return Execute("osint", page, data)
}

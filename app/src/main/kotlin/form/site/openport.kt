package form.site

import form.site.buildFormRadio
import form.site.buildFormCheck
import form.site.buildFormText
import form.create.OpenPortBuilder
import form.create.createDocument

import form.events.Feed
import form.events.Event
import form.Sessions

fun createOpenPortPage(
    e: List<Event>, 
    path: String = "/openport", 
    summary: (List<Event>) -> String = ::openPortSummary,
    body: (List<Event>) -> String = ::openPortBody
): String {
    return """
    <form class="mt-4 flex flex-col justify-start" hx-post="${path}" hx-target="#toast" hx-swap="innerHTML">
        ${buildFormText("Form Number", "formNumber")}

        ${buildFormText("Threat Type (MITRE ATT&CK T-Code)", "threatType", default="T1133 External Remote Services")}

        ${buildFormArea("Summary Paragraph", "summary", default=summary(e))}
        ${buildFormArea("Body Paragraph", "body", default=body(e))}

        ${buildFormArea("Additional References", "additionalReferences")}

        ${buildFormRadio(listOf("Amber", "Green"), "changeTlp", title="TLP Alert" )}
        ${buildFormCheck("Create CSV", "ips")}

        <div class="flex justify-evenly">
            ${buildButton("Create", "submit")}
            <div class="min-w-2"></div>
            ${buildButton("Reset", "reset")}
        </div>
    </form>
    <div id="toast"></div>
    """
}

fun openPortSummary(e: List<Event>): String {
    var cves = false
    for (i in e) {
        if (i.cves.isNotEmpty()) {
            cves = true
        }
    }

    if (cves) {
        return "The North Carolina National Guard Cyber Security Response Force (NCNG CSRF) received an alert indicating the ${e.first().name} domain is publicly exposed to the internet via several risky open ports and CVEs of concern."
    } else {
        return "The North Carolina National Guard Cyber Security Response Force (NCNG CSRF) received an alert indicating the ${e.first().name} domain is publicly exposed to the internet via several risky open ports."
    }
}

fun openPortBody(e: List<Event>): String {
    val ips = e.map{it.ip}.joinToString(", ")
    return "A threat actor may have an easier pathway to conducting a cyber attack or cyber espionage against your organization based on your current configuration. We encourage ${e.first().name} to review the infrastructure at the following IP addresses: ${ips} and evaluate the risk of leaving them in their current state. We also encourage ${e.first().name} to search for indicators of unauthorized access because threat actors exploit this configuration often for initial access."
}

fun endOfLifeSummary(e: List<Event>): String {
    var cves = false
    for (i in e) {
        if (i.cves.isNotEmpty()) {
            cves = true
        }
    }

    if (cves) {
        return """The North Carolina National Guard Cyber Security Response Force (NCNG CSRF) received an alert indicating the ${e.first().name} domain is publicly exposing end of life infrastructure via several risky open ports and CVEs of concern."""
    } else {
        return """The North Carolina National Guard Cyber Security Response Force (NCNG CSRF) received an alert indicating the ${e.first().name} domain is publicly exposing end of life infrastructure via several risky open ports of concern."""
    }
}

fun loginPageSummary(e: List<Event>): String {
    return """The North Carolina National Guard Cyber Security Response Force (NCNG CSRF) received an alert indicating the ${e.first().name} domain is publicly exposing risky login pages to the internet."""
}

fun loginPageBody(e: List<Event>): String {
    val ips = e.map{it.ip}.joinToString(", ")
    return """A threat actor may have an easier pathway to conducting a cyber attack or cyber espionage against your organization based on your current configuration, through repeated login attempts against possible weak user login credentials. We encourage ${e.first().name} to review the infrastructure at the following IP addresses: ${ips} and evaluate the risk of leaving them in their current state. We also encourage ${e.first().name} to search for indicators of unauthorized access because threat actors exploit this configuration often for initial access."""
}

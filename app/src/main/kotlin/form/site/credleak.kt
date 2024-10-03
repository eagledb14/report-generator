package form.site

import io.ktor.server.netty.*
import io.ktor.server.routing.*
import io.ktor.server.application.*
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.engine.*
import io.ktor.server.request.receiveParameters

import form.create.CredLeakFileBuilder
import form.create.createDocument
import form.events.Feed
import form.events.Event

fun Route.credLeak() {
    route("/credleak") {
        get {
            call.respondText(buildContent(createCredLeakPage(Event.buildEmpty()), Page.CredLeak), ContentType.Text.Html)
        }

        post {
            val params = call.receiveParameters().entries().associate { it.key.toString() to it.value.toString().removePrefix("[").removeSuffix("]")} as HashMap<String, String>
            val doc = CredLeakFileBuilder().buildFromMap(params)

            call.respondText("Created form: ${doc.getAlertName()}")
            createDocument(doc, params["ips"] == "on")
        }
    }
}

fun createCredLeakPage(e: Event): String {
    return """
    <form class="mt-4 flex flex-col justify-start" hx-post="/credleak" hx-target="#toast" hx-swap="innerHTML">
        <div class="flex">
            ${buildFormText("Name of organization", "orgName", default=e.name)}
            ${buildFormText("Form Number", "formNumber")}
        </div>
        ${buildFormText("Name of Victim Org", "victimOrg", default=e.name)}

        ${buildFormCheck("Multiple Leaks Found", "leakOrLeaks")}
        ${buildFormCheck("Multiple Credentials Leaked", "credOrCreds")}

        ${buildFormRadio(listOf("The passwords have been obfuscated to show only the first two letters.", "The passwords are not included due to being posted in plain text.", "The passwords are not included due to the threat actor not disclosing them.", "No Passwords were included"), "passwordOption", title="For passwords add either", checked=3)}

        ${buildFormText("Ip Address", "ipAddress", default=e.ip)}

        ${buildFormArea("Insert Username: Password", "userPass")}
        ${buildFormArea("Additional Information", "addInfo")}
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


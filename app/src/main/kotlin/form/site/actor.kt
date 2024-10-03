package form.site

import io.ktor.server.netty.*
import io.ktor.server.routing.*
import io.ktor.server.application.*
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.engine.*
import io.ktor.server.request.receiveParameters

import form.site.buildFormRadio
import form.site.buildFormCheck
import form.site.buildFormText
import form.site.buildFormArea
import form.create.ActorProfileBuilder
import form.create.createDocument

import form.events.Feed
import form.events.Event
import form.Sessions


fun Route.actor() {
    route("/actor") {
        get {
            call.respondText(buildContent(createActorPage(), Page.Actor), ContentType.Text.Html)
        }
        
        post {
            // val cookie = call.request.cookies[]
            val params = call.receiveParameters().entries().associate { it.key.toString() to it.value.toString().removePrefix("[").removeSuffix("]")} as HashMap<String, String>

            val doc = ActorProfileBuilder().buildFromMap(params)
            println(params)
            call.respondText("Created form: ")
            createDocument(doc, false)
        }

    }
}

fun createActorPage(
    
): String {
    return """
    <form class="mt-4 flex flex-col justify-start" hx-post="/actor" hx-target="#toast">
        <div id="stuff"></div>
        ${buildFormText("Primary Name", "name")}
        ${buildFormText("Alias", "alias")}
        ${buildFormText("First seen activity [ DD MM, YYYY]", "date")}
        ${buildFormText("Country of Origin", "country")}
        ${buildFormText("Motivation", "motivation")}
        ${buildFormText("Targeting", "target")}
        ${buildFormText("Malware Name", "malware")}
        ${buildFormText("Third Party Reporting", "reporter")}

        ${buildFormRadio(listOf("High", "Medium", "Low"), "confidence", checked=0, title="Assessment Confidence")}

        ${buildFormArea("Exploits", "exploits")}
        ${buildFormArea("Attack Chain Summary", "summary")}
        ${buildFormArea("Capabilities", "capabilities")}
        ${buildFormArea("Detection names", "detection")}
        ${buildFormArea("TTPS", "ttps")}
        ${buildFormArea("Infrastructure", "infra")}


        <div class="flex justify-evenly">
            ${buildButton("Create", "submit")}
            <div class="min-w-2"></div>
            ${buildButton("Reset", "reset")}
        </div>
    </form>
    <div id="toast"></div>
    """
}

// Primary Name:
// Alias:
// First seen activity: [ DD MM, YYYY]
// Country of Origin:
// Motivation:

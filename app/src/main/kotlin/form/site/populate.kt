package form.site

import io.ktor.server.netty.*
import io.ktor.server.routing.*
import io.ktor.server.application.*
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.engine.*
import io.ktor.server.request.receiveParameters

import kotlinx.coroutines.runBlocking
import kotlinx.coroutines.launch

import form.events.Feed
import form.events.Event
import form.events.getIpsFromCitr
import form.Sessions
import form.create.OpenPortBuilder
import form.create.createDocument

fun Route.populate() {
    route("/populate") {
        get {
            val cookie = call.request.cookies["session"]
            val event = Sessions.getStaticEvent(cookie)
            if (event.first() == Event.buildEmpty()) {
                call.respondText(buildContent(createDownloadPage(event), Page.OpenPort), ContentType.Text.Html)
            }
            call.respondText(buildContent(displayDownloadEvent(event), Page.OpenPort), ContentType.Text.Html)
        }

        put {
            val cookie = call.request.cookies["session"]
            Sessions.setStaticEvent(cookie, listOf(Event.buildEmpty()))
            call.respondText(buildContent(createDownloadPage(listOf(Event.buildEmpty())), Page.OpenPort), ContentType.Text.Html)
        }

        post {
            val params = call.receiveParameters()

            val ips = params["ipAddress"]!!.split(",").map { it.trim() }
            val org = params["orgName"]!!
            val newEvent = loadEvents(ips, org)

            val cookie = call.request.cookies["session"]
            Sessions.setStaticEvent(cookie, newEvent)
            
            call.respondText(displayDownloadEvent(newEvent))
        }

        post("/openport") {
            val params = call.receiveParameters().entries().associate { it.key.toString() to it.value.toString().removePrefix("[").removeSuffix("]") } as HashMap<String, String>

            val cookie = call.request.cookies["session"]
            val events = Sessions.getStaticEvent(cookie)

            val doc = OpenPortBuilder().buildFromMap(params, events)

            call.respondText("Created form: ${doc.getAlertName()}")
            createDocument(doc, params["ips"] == "on")
        }

        get("/openport") {
            val cookie = call.request.cookies["session"]
            val events = Sessions.getStaticEvent(cookie)
            call.respondText(createFormTabs(events, Form.OpenPort, "/populate"), ContentType.Text.Html)
        }

        get("/endoflife") {
            val cookie = call.request.cookies["session"]
            val events = Sessions.getStaticEvent(cookie)
            call.respondText(createFormTabs(events, Form.EndOfLife, "/populate"), ContentType.Text.Html)
        }

        get("/loginpage") {
            val cookie = call.request.cookies["session"]
            val events = Sessions.getStaticEvent(cookie)
            call.respondText(createFormTabs(events, Form.LoginPage, "/populate"), ContentType.Text.Html)
        }
    }
}

private fun createDownloadPage(e: List<Event>): String {
    return """
    <form class="mt-4 flex flex-col justify-start" hx-post="/populate" hx-target="#ring" hx-indicator="#load">
        ${buildFormText("Name of organization", "orgName", default=e.first().name)}
        ${buildFormText("Ip Address", "ipAddress", default=e.first().ip)}

        <div id="load" class="htmx-indicator text-center animate-bounce">loading...</div>
        <div class="flex justify-evenly">
            ${buildButton("Populate", "submit")}
            <div class="min-w-2"></div>
            ${buildButton("Reset", "reset")}
        </div>
    </form>
    """ 
}

private fun displayDownloadEvent(e: List<Event>, form: Form = Form.OpenPort): String {
    return """
    <div id="events">
        <h2 class="text-center my-4 text-xl font-bold">${e.first().name}</h2>

        <div class="my-2">IPs:</div>
        <div class="mx-2"> ${e.map{it.ip}.joinToString("<br> ")}</div>
        <div class="my-2">Ports:</div>
        <div class="mx-2"> ${e.flatMap{it.ports}.joinToString("<br> ")}</div>
        <br>

        <div>
            ${e.flatMap{it.cves}.map{
                "${it.first}: Priority ${it.second}"}
            .joinToString("<br> ")}
        </div>

        <button class="w-full mt-4 px-4 py-2 bg-blue-800 text-white rounded-md" hx-put="/populate" hx-target="#content">Return</button>
        ${createFormTabs(e, form, "/populate")}
    </div>
    """
}

//if there is a slash in the ip address, then get all the sub ips and create events with them
private fun loadEvents(ips: List<String>, org: String): List<Event> {
    var events = mutableSetOf<Event>()
    for (ip in ips) {
        if (ip.contains("/")) {
            getIpsFromCitr(ip.trim()).map {
                events.add(Event(it.trim(), org))
            }
        } else {
            events.add(Event(ip.trim(), org))
        }
    }

    runBlocking {
        events.forEach {
            launch {
                it.load()
            }
        }
    }

    return events.filter { it.ports.isNotEmpty() }.sortedBy { it.ip }
}

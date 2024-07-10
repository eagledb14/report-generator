package form.site

import io.ktor.http.*
import io.ktor.http.cookies
import io.ktor.server.application.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.request.receiveParameters
import io.ktor.server.response.*
import io.ktor.server.routing.*

import form.events.Feed
import form.events.Event
import form.Sessions

import form.site.displayEvent
import form.site.createFormTabs
import form.site.createOpenPortPage

import form.create.OpenPortBuilder
import form.create.createDocument

fun Route.events() {
    route("/events") {
        get {
            val cookie = call.request.cookies["session"]
            val feed = Sessions.getFeed(cookie)
            if (!feed.isEmpty()) {
                val event = feed.event()
                val index = feed.index
                call.respondText(buildContent(displayEvent(listOf(event), index), Page.Events), ContentType.Text.Html)
            }
            call.respondText(buildContent(createEventsDownloadPage(), Page.Events), ContentType.Text.Html)
        }

        post("/download") {
            val cookie = call.request.cookies["session"]
            val feed = Sessions.getFeed(cookie)
            feed.download()
            val event = feed.event()

            val index = feed.index
            call.respondText(displayEvent(listOf(event), index))
        }

        get("/next") {
            val cookie = call.request.cookies["session"]
            val feed = Sessions.getFeed(cookie)
            val event = feed.next()
            val index = feed.index

            call.respondText(displayEvent(listOf(event), index))
        }

        get("/prev") {
            val cookie = call.request.cookies["session"]
            val feed = Sessions.getFeed(cookie)
            val event = feed.prev()
            val index = feed.index

            call.respondText(displayEvent(listOf(event), index))
        }

        get("/openport") {
            val cookie = call.request.cookies["session"]
            val events = listOf(Sessions.getFeed(cookie).event())
            call.respondText(createFormTabs(events, Form.OpenPort, "/events"), ContentType.Text.Html)
        }

        post("/openport") {
            println("received the event")
            val params = call.receiveParameters().entries().associate { it.key.toString() to it.value.toString().removePrefix("[").removeSuffix("]") } as HashMap<String, String>

            val cookie = call.request.cookies["session"]
            val event = Sessions.getFeed(cookie).event()

            val doc = OpenPortBuilder().buildFromMap(params, listOf(event))
            call.respondText("Created form: ${doc.getAlertName()}")
            createDocument(doc, params["ips"] == "on")
        }

        get("/endoflife") {
            val cookie = call.request.cookies["session"]
            val events = listOf(Sessions.getFeed(cookie).event())
            call.respondText(createFormTabs(events, Form.EndOfLife, "/events"), ContentType.Text.Html)
        }

        get("/loginpage") {
            val cookie = call.request.cookies["session"]
            val events = listOf(Sessions.getFeed(cookie).event())
            call.respondText(createFormTabs(events, Form.LoginPage, "/events"), ContentType.Text.Html)
        }
    }
}

fun createEventsDownloadPage(): String {
    return """
    <h2 class="my-4 text-xl font-bold text-center">Start Viewing Events</h2>
    <div id="load" class="htmx-indicator text-center animate-bounce">loading...</div>
    <button hx-post="/events/download" hx-target="#ring" hx-indicator="#load" class="w-full mt-4 px-4 py-2 bg-blue-800 text-white rounded-md">
        Download
    </button>
    """
}

fun displayEvent(e: List<Event>, index: Int): String {
    return """
    <div id="events">

        <div>${index + 1}</div>
        <h2 class="text-center my-4 text-xl font-bold">${e.first().name}</h2>

        <div class="my-2">IPs:</div>
        <div class="mx-2"> ${e.map{it.ip}.joinToString("<br> ")}</div>
        <div class="my-2">Ports:</div>
        <div class="mx-2"> ${e.flatMap{it.ports}.joinToString("<br> ")}</div>
        <br>

        <div class="my-2">Trigger: ${e.first().trigger}</div>
        <div class="my-2">Time Stamp: ${e.first().timestamp}</div>
        <div class="my-2">Description: ${e.first().desc}</div>
        <div>
            ${e.flatMap{it.cves}.map{
                "${it.first}: Priority ${it.second}"}
            .joinToString("<br> ")}
        </div>

        <div class="my-4 text-center">
            <a href="${e.first().alert_link}" target="_blank" class="w-full mt-4 px-4 py-2 bg-blue-800 text-white rounded-md">Alert Link</a>
            <a href="${e.first().host_link}" target="_blank" class="w-full mt-4 px-4 py-2 bg-blue-800 text-white rounded-md">Shodan Link</a>
        </div>

        ${navButtons()}

        ${createFormTabs(e)}
    </div>
    """
}

fun displayUnknownEvent(): String {
    return """
    <div id="events">
        <div> Oops, you ran out of events</div>
        <div>Click next or previous to load more, or wait for Shodan to load more</div>
        ${navButtons()}
    </div>
    """
}

fun navButtons(): String {
    return """
        <div id="load" class="htmx-indicator text-center animate-bounce">loading...</div>
        <div class="flex justify-evenly">
            <button hx-get="/events/prev" hx-target="#events" hx-swap="outerHTML transition:true" hx-indicator="#load" class="w-full mt-4 px-4 py-2 bg-blue-800 text-white rounded-md">
                Previous
            </button>
            <div class="min-w-2"></div>
            <button hx-get="/events/next" hx-target="#events" hx-swap="outerHTML transition:true" hx-indicator="#load" class="w-full mt-4 px-4 py-2 bg-blue-800 text-white rounded-md">
                Next
            </button>
        </div>
    """
}

fun createFormTabs(e: List<Event>, form: Form = Form.OpenPort, endpoint: String = "/events"): String {

    return when(form) {
        Form.OpenPort, Form.CredLeak -> """
            <div id="event-form" class="mt-10 overflow-y-auto">
                <div id="tabs" class="flex justify-start">
                    <button hx-get="${endpoint}/openport" hx-target="#event-form" hx-swap="outerHTML" class="ring-2 ring-inset ring-blue-800 rounded-md px-4 py-2">Open Port</button>
                    <button hx-get="${endpoint}/endoflife" hx-target="#event-form" hx-swap="outerHTML" class="bg-blue-800 px-4 py-2 rounded-md text-white">End of Life</button>
                    <button hx-get="${endpoint}/loginpage" hx-target="#event-form" hx-swap="outerHTML" class="bg-blue-800 px-4 py-2 rounded-md text-white">Login Page</button>
                </div>
                ${createOpenPortPage(e, "${endpoint}/openport")}
            </div>
            """
        Form.EndOfLife -> """
            <div id="event-form" class="mt-10 overflow-y-auto">
                <div id="tabs" class="flex justify-start">
                    <button hx-get="${endpoint}/openport" hx-target="#event-form" hx-swap="outerHTML" class="bg-blue-800 px-4 py-2 rounded-md text-white">Open Port</button>
                    <button hx-get="${endpoint}/endoflife" hx-target="#event-form" hx-swap="outerHTML" class="ring-2 ring-inset ring-blue-800 rounded-md px-4 py-2">End of Life</button>
                    <button hx-get="${endpoint}/loginpage" hx-target="#event-form" hx-swap="outerHTML" class="bg-blue-800 px-4 py-2 rounded-md text-white">Login Page</button>
                </div>
                ${createOpenPortPage(e, "${endpoint}/openport", summary=::endOfLifeSummary)}
            </div>
            """
            Form.LoginPage -> """
                <div id="event-form" class="mt-10 overflow-y-auto">
                    <div id="tabs" class="flex justify-start">
                        <button hx-get="${endpoint}/openport" hx-target="#event-form" hx-swap="outerHTML" class="bg-blue-800 px-4 py-2 rounded-md text-white">Open Port</button>
                        <button hx-get="${endpoint}/endoflife" hx-target="#event-form" hx-swap="outerHTML" class="bg-blue-800 px-4 py-2 rounded-md text-white">End of Life</button>
                        <button hx-get="${endpoint}/loginpage" hx-target="#event-form" hx-swap="outerHTML" class="ring-2 ring-inset ring-blue-800 rounded-md px-4 py-2">Login Page</button>
                    </div>
                    ${createOpenPortPage(e, "${endpoint}/openport", summary=::loginPageSummary, body=::loginPageBody)}
                </div>
            """
    }
}

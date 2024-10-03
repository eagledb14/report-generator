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
import form.create.OpenPortBuilder
import form.create.createDocument

import form.events.Feed
import form.events.Event
import form.Sessions

fun Route.filter() {
    route("/filter") {
        get {
            val cookie = call.request.cookies["session"]
            val filters = Sessions.getFeed(cookie).filters

            call.respondText(buildContent(createFilterPage(filters), Page.Filter), ContentType.Text.Html)
        }

        post("/on/{filter}") {
            val filter = call.parameters["filter"]!!
            val cookie = call.request.cookies["session"]
            Sessions.getFeed(cookie).filters[filter] = true

            call.respondText(createFilterOn(filter, filter), ContentType.Text.Html)
        }

        post("/off/{filter}") {
            val filter = call.parameters["filter"]!!
            val cookie = call.request.cookies["session"]
            Sessions.getFeed(cookie).filters[filter] = false

            call.respondText(createFilterOff(filter, filter), ContentType.Text.Html)
        }
    }
}

fun createFilterPage(filters: Map<String, Boolean>): String {
    val builder = StringBuilder()
    builder.append("""<div class="mt-4 flex flex-col justify-start">""")

    for ((key, value) in filters) {
        if (value == true) {
            builder.append(createFilterOn(key, key))
        } else {
            builder.append(createFilterOff(key, key))
        }

    }

    builder.append("</div>")
    return builder.toString()
}

fun createFilterOn(filter: String, param: String): String {
    return """
    <div id ="${param.replace(" ", "_")}" class="flex items-center">
        <button class="rounded bg-blue-800 text-white mt-1 ml-10 p-2" hx-post="/filter/off/$filter" hx-target="#${param.replace(" ", "_")}">On</button>
        <label class="block text-sm font-medium text-black ml-5">$filter</label><br>
    </div>
    """
}

fun createFilterOff(filter: String, param: String): String {
    return """
    <div id ="${param.replace(" ", "_")}" class="flex items-center">
        <button class="rounded bg-white text-black mt-1 ml-10 p-2" hx-post="/filter/on/$filter" hx-target="#${param.replace(" ", "_")}">Off</button>
        <label class="block text-sm font-medium text-black ml-5">$filter</label><br>
    </div>
    """
}

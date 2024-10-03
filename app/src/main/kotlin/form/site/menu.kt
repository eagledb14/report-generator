package form.site

import io.ktor.server.netty.*
import io.ktor.server.routing.*
import io.ktor.server.application.*
import io.ktor.server.http.content.*
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.engine.*
import io.ktor.server.websocket.*
import io.ktor.http.cookies

import form.events.Event
import form.Sessions

import java.io.File

fun launch(port: Int) {
    embeddedServer(Netty, port = port, module = Application::start).start(wait = true)
}

fun Application.start() {
    install(WebSockets)

    routing {
        //remove if putting on website
        webSocket("/conn") {
            try {
                for (frame in incoming) {
                    // println("Message")
                }
            } catch (e: Exception) {
                System.exit(0)
            } finally {
                System.exit(0)
            }
        }

        get("/") {
            val cookie = call.request.cookies["session"]
            
            if (cookie == null) {
                val newCookie = Sessions.createHash()
                call.response.cookies.append("session", "$newCookie")
            } else if (!Sessions.hashExists(cookie)) {
                Sessions.addHash(cookie)
            }

            call.respondText(buildPage(createCredLeakPage(Event.buildEmpty()), Page.CredLeak), ContentType.Text.Html)
        }


        credLeak()

        events()
        populate()

        filter()
        actor()
    }
}

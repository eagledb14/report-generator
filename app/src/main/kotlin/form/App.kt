package form

import java.io.File
import java.awt.Desktop

import form.site.launch

import kotlinx.coroutines.runBlocking
import kotlinx.coroutines.launch
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock

import com.github.kittinunf.fuel.Fuel
import com.github.kittinunf.fuel.httpGet
import com.github.kittinunf.result.Result.*
import com.github.kittinunf.fuel.core.Headers

import java.net.ServerSocket
import java.net.URI
import form.events.*

fun main() {
    loadProperties()

    //you want to set the port to a static port if you are hosting it from a web server
    val port = if (System.getProperty("STATIC_PORT") == "true") {
        8084
    } else {
        val socket = ServerSocket(0)
        val port = socket.localPort
        socket.close()
        println(port)
        port
    }

    val dev = System.getProperty("DEV") == "true"
    if (dev) {
        launch(port)
    } else {
        runBlocking {
            launch {
                launch(port)
            }

            try {
                Desktop.getDesktop().browse(URI("http://localhost:${port}"))
            } catch(e: Exception) {}
        }
    }
}


fun loadProperties() {
    val envFile = File(".${File.separator}reference${File.separator}key.env")
    if (!envFile.exists()) {
        return
    }

    val lines = envFile.readLines()
    for (line in lines) {
        val parts = line.split("=")
        if (parts.size == 2) {
            System.setProperty(parts[0], parts[1])
        }
    }
}

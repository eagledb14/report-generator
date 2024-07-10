package form

import io.ktor.http.*
import io.ktor.server.engine.embeddedServer
import io.ktor.server.netty.Netty
import io.ktor.http.cookies
import java.util.concurrent.ConcurrentHashMap

import form.events.*

object Sessions {
    val feeds = ConcurrentHashMap<String, Feed>()
    private val staticEvent = ConcurrentHashMap<String, List<Event>>()

    fun createHash(): String {
        val newHash = (1..50).map { ('a'..'z').random() }.joinToString("")
        if (feeds.containsKey(newHash)) {
            return createHash()
        }

        feeds.put(newHash, Feed())
        return newHash
    }

    fun addHash(hash: String) {
        feeds.put(hash, Feed())
    }

    fun hashExists(hash: String?): Boolean {
        return feeds.containsKey(hash)
    }

    fun removeHash(hash: String?) {
        feeds.remove(hash)
    }

    fun getFeed(hash: String?): Feed {
        if (feeds.containsKey(hash)) {
            return feeds[hash]!!
        }

        addHash(hash!!)
        return getFeed(hash)
    }

    fun getStaticEvent(hash: String?): List<Event> {
        if (feeds.containsKey(hash)) {
            if (staticEvent[hash] == null) {
                return listOf(Event.buildEmpty())
            }
            return staticEvent[hash]!!
        }
        return listOf(Event.buildEmpty())
    }

    fun setStaticEvent(hash: String?, e: List<Event>) {
        staticEvent[hash!!] = e
    }

}





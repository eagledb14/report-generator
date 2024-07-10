package form.events

import com.github.kittinunf.fuel.Fuel
import com.github.kittinunf.fuel.httpGet
import com.github.kittinunf.result.Result.*
import java.net.URL
import java.nio.file.Paths
import com.rometools.rome.io.SyndFeedInput
import com.rometools.rome.feed.synd.SyndFeed
import com.rometools.rome.io.XmlReader
import com.rometools.rome.feed.synd.SyndEntry
import com.google.gson.Gson
import javax.xml.parsers.SAXParserFactory
import kotlinx.coroutines.*
import form.create.getKey

class Feed(): Iterator<Event> {
    private val key: String = getKey()
    private val events = ArrayList<Event>()
    var index = 0
        private set

    val filters: MutableMap<String, Boolean> = mutableMapOf(
        "end of life" to true,
        "industrial control system" to true, 
        "internet scanner" to false,
        "iot" to true,
        "malware" to true,
        "new service" to true,
        "open database" to true,
        "ssl expired" to false,
        "uncommon" to false,
        "uncommon plus" to false,
        "vulnerable" to true,
        "vulnerable unverified" to false,
    )

    override fun hasNext(): Boolean {
        if (filters.values.all { it == false }) {
            return false
        }

        return index < events.size
    }

    fun event(): Event {
        if (filters.values.all { it == false }) {
            return Event.buildEmpty()
        }

        waitForNext(index)

        return events.get(index).load()
    }

    override fun next(): Event {

        if (filters.values.all { it == false }) {
            return Event.buildEmpty()
        }

        if (index == events.size) {
            return event()
        }

        index = minOf(index + 1, events.size)
        val event = events.get(index)
        if (filters.getOrDefault(event.trigger, false) == true) {
            waitForNext(index)
            return event
        }
            
        return next()
    }

    private fun waitForNext(next: Int) {
        while (events.get(next).loaded == false) {
            Thread.sleep(1000)
        }
    }

    fun prev(): Event {

        if (filters.values.all { it == false }) {
            return Event.buildEmpty()
        }

        if (index == 0) {
            return event()
        }

        index = maxOf(index - 1, 0)
        val event = events.get(index)
        if (filters.getOrDefault(event.trigger, false) == true) {
            return event.load()
        }

        return prev()
    }

    //returns download progress / 100
    fun progress(): Int {
        return events.count { it.loaded == true }
        
    }

    fun isEmpty(): Boolean {
        return events.isNullOrEmpty()
    }

    fun resetFilters() {
        for ((key, _) in filters) {
            filters[key] = true
        }
    }

    //Idk how to fix the deprecation issues here, but it works
    fun download() {
        index = 0
        events.clear()

        // val url = Paths.get("https://monitor.shodan.io/events.rss?key=${key}").toUri().toURL() //Idk why this doesn't work
        val url = URL("https://monitor.shodan.io/events.rss?key=${key}")
        SyndFeedInput().build(XmlReader(url)).entries
            .forEach { events.add(Event(it)) }
        events.reverse()

        GlobalScope.launch {
            events.forEach { 
                if (it.ip != "Asset") {
                    it.load() 
                }
            }
        }
    }
}

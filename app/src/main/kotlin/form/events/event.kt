package form.events

import com.github.kittinunf.fuel.Fuel
import com.github.kittinunf.fuel.httpGet
import com.github.kittinunf.result.Result.*
import com.github.kittinunf.fuel.core.Headers
import com.rometools.rome.io.SyndFeedInput
import com.rometools.rome.feed.synd.SyndFeed
import com.rometools.rome.io.XmlReader
import com.rometools.rome.feed.synd.SyndEntry
import com.google.gson.Gson
import javax.xml.parsers.SAXParserFactory
import kotlinx.coroutines.*
import form.create.getKey

data class Event (
    val ip: String,
    val trigger: String,
    val alert_link: String,
    val host_link: String,
    val alert_id: String,
    val desc: String,
    val timestamp: String,
) {
    public var ports = listOf<String>()
    var name = ""
    var id = ""
        private set
    val cves: ArrayList<Triple<String, Int, String>> = ArrayList()

    var loaded = false
        private set

    public constructor(entry: SyndEntry): this (
        entry.title.split(" ")[0],
        entry.title.split(" ").last().replace("`", "").split("_").joinToString(" "),
        entry.link,
        "https://www.shodan.io/host/${entry.title.split(" ")[0]}",
        entry.link.split("/").last(),
        entry.description.value + " on port: " + entry.title.split(" ")[3],
        entry.publishedDate.toString(),
    ) 

    public constructor(ip: String, org: String): this(
        ip, 
        "",
        "", 
        ip.split(", ").map{ "https://www.shodan.io/host/${it.trim()}" }.joinToString(", "),
        "", 
        "", 
        ""
    ) {
        this.name = org
    }

    private constructor(): this("", "", "", "", "", "", "") {}

    companion object {
        fun buildEmpty(): Event {
            return Event()
        }
    }

    public fun load(): Event {
        if (loaded == true) {
            return this
        }

        runBlocking {
            launch {
                if (alert_link != "") {
                    id = getId(alert_link)
                    name = getName(id)
                }
            }
            launch {
                val banner = getHost()

                if (banner.data.isEmpty()) {
                    return@launch
                }

                val cvesFromHost = HashMap<String, Vuln>()
                for (product in banner.data) {
                    product.vulns?.forEach {
                        cvesFromHost[it.key] = it.value
                    }
                }

                for ((key, value) in cvesFromHost) {
                    cves.add(getPriority(key, value))
                }
                cves.sortBy { it.second }

                ports = banner.data
                    .map { 
                        if (it.product == null) {
                            "${it.port}"
                        } else {
                            "${it.port} (${it.product})" 
                        }
                    }
                
            }
        }
        loaded = true

        return this
    }

    fun cveList(): String {
        if (cves.isEmpty()) {
            return ""
        }
        return cves.map { it.first }.joinToString(", ")
    }

    fun cveParagraph(): String {
        var defaultCveList = ""

        ports.forEach {
            defaultCveList += "${it.replace("(", "").replace(")", "")} \n"
        }

        cves.forEach {
            defaultCveList += "${it.first}\u3000Priority-${it.second} \u3000${it.third}\n"
        }
        return defaultCveList
    }

    private fun getName(id: String, retries: Int = 5): String {
        val url = "https://api.shodan.io/shodan/alert/$id/info?key=${getKey()}"
        val(_,resp, res) = Fuel.get(url).header(headers).responseString()

        when (res) {
            is Success -> {
                val data = res.get()
                val alert = Gson().fromJson(data, Alert::class.java)
                return alert.name
            }
            is Failure -> {
                if (retries == 0) {
                    println("name ${resp.statusCode}")
                    return "ERROR: ${resp.statusCode}, $url"
                }
                Thread.sleep(1000)
                return getName(id, retries - 1)
            }
        }
    }

    private fun getId(alert_link: String, retries: Int = 5): String {
        val url = "${alert_link}?key=${getKey()}"
        val(_,resp, res) = Fuel.get(url).header(headers).responseString()

        when (res) {
            is Success -> {
                val data = res.get().toString().split("\n")
                for (line in data) {
                    if (line.contains("let data = ")) {
                        return line.split("\"")[3]
                    }
                }

                return ""
            }
            is Failure -> {
                if (retries == 0) {
                    println("id ${resp.statusCode}")
                    return "ERROR: ${resp.statusCode}, ${url}"
                }
                Thread.sleep(1000)
                return getId(alert_link, retries - 1)
            }
        }
    }

    private fun getHost(retries: Int = 5): Banner {
        val url = "https://api.shodan.io/shodan/host/${ip}?key=${getKey()}"
        val (_,resp, res) = Fuel.get(url).header(headers).responseString()

        when (res) {
            is Success -> {
                val data = res.get()
                val banner = Gson().fromJson(data, Banner::class.java)
                return banner
            }
            is Failure -> {
                if (retries == 0) {
                    println("Host ${resp.statusCode}")
                    return Banner(listOf())
                }
                Thread.sleep(1000)
                getHost(retries - 1)
            }
        }
        return Banner(listOf())
    }
}

data class Alert(
    val name: String,
)

val headers = mapOf(
    Headers.USER_AGENT to """Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36""",
)

data class Banner(
    val data: List<Product>,
)

data class Product(
    val product: String?,
    val port: String?,
    val vulns: HashMap<String, Vuln>?,
)

data class Vuln(
    val cvss: Float?,
    val cvss_v2: Float?,
    val epss: Float?,
    val kev: Boolean?,
    val summary: String,
)

data class Net(
    val matches: List<Match>?,
)

data class Match(
    val ip_str: String,
)

fun getIpsFromCitr(citr: String, retries: Int = 5): List<String> {
    val url = "https://api.shodan.io/shodan/host/search?key=${getKey()}&query=net:${citr}"
    val (_,resp, res) = Fuel.get(url).responseString()

    when(res) {
        is Success -> {
            val data = res.get()
            val matches = Gson().fromJson(data, Net::class.java)
            return matches.matches!!.map { it.ip_str }
        }
        is Failure -> {
            if (retries == 0) {
                println("Host ${resp.statusCode}")
                return listOf()
            }
            Thread.sleep(1000)
            return getIpsFromCitr(citr, retries - 1)
        }
    }
}

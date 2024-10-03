package form.create

import com.github.kittinunf.fuel.Fuel
import com.github.kittinunf.fuel.*
import kotlinx.coroutines.runBlocking
import kotlinx.coroutines.launch
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock
import com.github.kittinunf.result.Result.*
import com.opencsv.CSVWriter
import java.io.FileWriter
import java.io.Writer

import java.awt.image.BufferedImage
import java.io.ByteArrayInputStream
import java.io.File
import java.util.Base64
import javax.imageio.ImageIO

import com.google.gson.Gson

fun getHosts(ips: String): List<Host> {
    if (ips == "") {
        return listOf()
    }

    val key = getKey()
    val ip_list = getIpList(ips)

    val url = "https://api.shodan.io/shodan/host/search?key=${key}&query=net:$ip_list"

    val (_,_, result) = Fuel.get(url).responseString()
    when (result) {
        is Success -> {
            val data = result.get()
            val ranges = Gson().fromJson(data, Ranges::class.java)
            return ranges.matches.sortedWith(compareBy<Host> { it.ip_str }.thenBy { it.port })
        }
        is Failure -> {
            return listOf()
        }
    }
}

fun getIpList(ips: String): String {
    return ips.split("[\n,\\s]+".toRegex()).joinToString(",") { it.trim() }
}

fun getKey(): String {
    return System.getProperty("SHODAN_API_KEY")
}

fun writeToCsv(orgName: String, alertId: String, hosts: List<Host>) {
    val newFolder = File(".${File.separator}$orgName")

    if (!newFolder.exists()) {
        newFolder.mkdirs()
    }

    val newCsv = File("${newFolder}${File.separator}${orgName}-${alertId}.csv")
    var writer = FileWriter(newCsv)
    val csvWriter = CSVWriter(writer)
    
    val header = arrayOf("asn", "ip", "port", "timestamp", "domains", "data", "hostnames", "isp", "org", "os", "country", "country code", "region code", "city", "product")
    csvWriter.writeNext(header)

    hosts.forEach{ host ->
        val row = arrayOf(
            host.asn, 
            host.ip_str, 
            host.port.toString(),
            host.timestamp, 
            host.domains.joinToString(", "),
            host.data,
            host.hostnames.joinToString(", "),
            host.isp,
            host.org, 
            host.os,
            host.location.country_name, 
            host.location.country_code, 
            host.location.region_code,
            host.location.city, 
            host.product,
        )
        csvWriter.writeNext(row)
    }

    csvWriter.close()
}

data class Ranges(
    val matches: List<Host>,
    val total: Int
)

data class Host(
    val asn: String,
    val timestamp: String,
    val domains: List<String>,
    val hostnames: List<String>,
    val product: String,
    val location: Location,
    val org: String,
    val data: String,
    val isp: String,
    val os: String,
    val transport: String,
    val port: Int,
    val ip_str: String,
    val screenshot: Screenshot?
)

data class Location(
    val city: String,
    val country_name: String,
    val country_code: String,
    val region_code: String
)

data class Screenshot(
    val data: String
)

fun base64ToImage(base64String: String, orgName: String, alertId: String, address: String) {
    var str = base64String.replace(" ", "").replace("\n", "")
    val imageBytes = Base64.getDecoder().decode(str)
    File(".${File.separator}${orgName}${File.separator}${orgName}-${alertId}-${address}.jpeg").writeBytes(imageBytes)
}


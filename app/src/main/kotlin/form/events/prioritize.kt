package form.events

fun getPriority(cve: String, vuln: Vuln): Triple<String, Int, String> {
    val cvss = if (vuln.cvss_v2 == null) {
        vuln.cvss
    } else {
        vuln.cvss_v2
    }
    val rank = getRank(cvss, vuln.epss, vuln.kev)
    return Triple(cve, rank, vuln.summary.replace(">", ")").replace("<", "(").replace("\n", " "))
}

private fun getRank(cvss: Float?, epss: Float?, kev: Boolean?): Int {
    val cvss_score = 6.0
    val epss_score = 0.2

    if (kev != null && kev) {
        return 0
    } else if (cvss == null || epss == null) {
        return 5
    } else if (cvss >= cvss_score) {
        if (epss >= epss_score) {
            return 1
        }
        return 2
    } else {
        if (epss >= epss_score) {
            return 3
        }
        return 4
    }
}


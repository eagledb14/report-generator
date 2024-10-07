package form.create

import java.io.File
import java.time.format.DateTimeFormatter
import java.time.LocalDateTime
import net.lingala.zip4j.ZipFile
import form.events.Event

interface CreateDoc {
    fun createDoc()
    fun writeToFile(tempDoc: String)
    fun getOrgname(): String
    fun getAlertId(): String
    fun getIps(): String
    fun getAlertName(): String
}

class CreateOpenPortDoc(
    val formNumber: String,
    val threatType: String,
    val changeTlp: Boolean,
    val additionalReferences: String,
    val summaryInsert: String,
    val bodyInsert: String,
    val events: List<Event>,
) : CreateDoc {
    val docname = ".${File.separator}open_port_edited.docx"
    var orgName = events.first().name.replace(Regex("[\\/:*?\"<>|]"), "-")
    var ipAddress = events.map{ it.ip }.joinToString(", ")

    override fun createDoc() {
        val tempDoc = ".${File.separator}temp"
        this.writeToFile(tempDoc)
        zipDocument(tempDoc, orgName, getAlertId(formNumber))
    }

    override fun writeToFile(tempDoc: String) {
        ZipFile(".${File.separator}reference${File.separator}${this.docname}").extractAll(tempDoc)

        val doc = File(".${File.separator}temp${File.separator}word${File.separator}document.xml")

        var docxFile = doc.readText()
        if (changeTlp) {
            docxFile = docxFile.replace("tlp_insert", "Recipients may share TLP:GREEN information with peers and partner organizations within their community, but not via publicly accessible channels. Unless otherwise specified, TLP:GREEN information may not be shared outside of the cybersecurity or cyber defense community.")
            writeToHeaderFooter()
        } else {
            docxFile = docxFile.replace("tlp_insert", "Recipients may share TLP:AMBER information with members of their own organization and its clients on a need-to-know basis to protect their organization and its clients and prevent further harm.")
        }

        docxFile = docxFile.replace("org_name_insert", orgName)
        docxFile = docxFile.replace("alert_id", getAlertId(formNumber))

        docxFile = docxFile.replace("threat_type", threatType)

        docxFile = docxFile.replace("summary_insert", summaryInsert)
        docxFile = docxFile.replace("body_insert", bodyInsert)

        //inf_impact_insert
        if (ipAddress != "") {
            val infString = docxFile.indexOf("inf_impact_insert")
            val start = docxFile.lastIndexOf("<w:p ", infString)
            val end = docxFile.indexOf("</w:p>", infString) + 6
            if (start != -1 && end != -1 && infString != -1) {
                val infSub = docxFile.substring(start, end)
                docxFile = docxFile.replace(infSub, splitList(ipAddress))
            }
        } else {
            docxFile = docxFile.replace("inf_impact_insert", "") 
        }

        val source = events.map { it.host_link }.joinToString(" ")
        if (source != "") {
            val sourceString = docxFile.indexOf("source_insert")
            val start = docxFile.lastIndexOf("<w:p ", sourceString)
            val end = docxFile.indexOf("</w:p>", sourceString) + 6
            if (start != -1 && end != -1 && sourceString != -1) {
                val sourceSub = docxFile.substring(start, end)
                docxFile = docxFile.replace(sourceSub, splitList(source))
            }
        } else {
            docxFile = docxFile.replace("source_insert", "") 
        }
        
        //port_or_cve_list_insert
        if (events.size > 0) {
            val refString = docxFile.indexOf("port_or_cve_list_insert")
            val start = docxFile.lastIndexOf("<w:tr ", refString)
            val end = docxFile.indexOf("</w:tr>", refString) + 7
            if (start != -1 && end != -1 && refString != -1) {
                val refSub = docxFile.substring(start, end)
                docxFile = docxFile.replace(refSub, this.splitTable())
            }
        } else {
            docxFile = docxFile.replace("port_or_cve_list_insert", "")
        }

        //reference_insert
        if (additionalReferences != "") {
            docxFile = docxFile.replace("reference_insert", additionalReferences)
        } else {
            val refString = docxFile.indexOf("reference_insert")
            val start = docxFile.lastIndexOf("<w:tr ", refString)
            val end = docxFile.indexOf("</w:tr>", refString) + 7
            if (start != -1 && end != -1 && refString != -1) {
                val refsub = docxFile.substring(start, end)
                docxFile = docxFile.replace(refsub, "")
            }
        }

        //Deletes the references if there are no cves
        if (!events.any { it.cves.isNotEmpty() }) {
            val refString = docxFile.indexOf("CVE Priority Key")
            val start = docxFile.lastIndexOf("<w:tr ", refString)

            val endString = docxFile.indexOf("Priority-5")
            val end = docxFile.indexOf("</w:tr>", endString) + 7
            if (start != -1 && end != -1 && refString != -1 && endString != -1) {
                val refsub = docxFile.substring(start, end)
                docxFile = docxFile.replace(refsub, "")
            }
        }

        doc.writeText(docxFile)
    }

    override fun getOrgname(): String {
        return this.orgName
    }

    override fun getAlertId(): String {
        return getAlertId(this.formNumber)
    }

    override fun getIps(): String {
        return events.map{ it.ip }.joinToString(",")
    }

    override fun getAlertName(): String {
        if (changeTlp) {
            return "TLP GREEN: ${orgName} EP-${ipAddress}: ${getAlertId()}"
        } else {
            return "TLP AMBER: ${orgName} EP-${ipAddress}: ${getAlertId()}"
        }
    }

    fun splitTable(): String {
        return events.map {
            var title = "Ports"
            if (it.cves.size > 0) {
                title = "Ports and CVEs"
            }
            """<w:tr w:rsidR="007438A5" w:rsidRPr="00CE4213" w14:paraId="1CC9F227" w14:textId="77777777" w:rsidTr="009E01CE"><w:trPr><w:trHeight w:val="294"/></w:trPr><w:tc><w:tcPr><w:tcW w:w="1710" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="E6E6E6"/></w:tcPr><w:p w14:paraId="373B8676" w14:textId="3DF70593" w:rsidR="007438A5" w:rsidRPr="00CE4213" w:rsidRDefault="00B068CA" w:rsidP="00BB2555"><w:pPr><w:rPr><w:rFonts w:asciiTheme="minorHAnsi" w:hAnsiTheme="minorHAnsi" w:cstheme="minorHAnsi"/><w:b/><w:bCs/><w:sz w:val="24"/></w:rPr></w:pPr><w:r w:rsidRPr="00385E00"><w:rPr><w:rFonts w:asciiTheme="minorHAnsi" w:hAnsiTheme="minorHAnsi" w:cstheme="minorHAnsi"/><w:b/><w:bCs/><w:sz w:val="24"/></w:rPr><w:t xml:space="preserve">List of </w:t></w:r><w:r w:rsidR="00DE2088" w:rsidRPr="00385E00"><w:rPr><w:rFonts w:asciiTheme="minorHAnsi" w:hAnsiTheme="minorHAnsi" w:cstheme="minorHAnsi"/><w:b/><w:bCs/><w:sz w:val="24"/></w:rPr><w:t>${title}</w:t></w:r><w:r w:rsidR="00BB2555"><w:rPr><w:rFonts w:asciiTheme="minorHAnsi" w:hAnsiTheme="minorHAnsi" w:cstheme="minorHAnsi"/><w:b/><w:bCs/><w:sz w:val="24"/></w:rPr><w:t xml:space="preserve"> on</w:t></w:r><w:r w:rsidRPr="00385E00"><w:rPr><w:rFonts w:asciiTheme="minorHAnsi" w:hAnsiTheme="minorHAnsi" w:cstheme="minorHAnsi"/><w:b/><w:bCs/><w:sz w:val="24"/></w:rPr><w:t xml:space="preserve"> </w:t></w:r><w:r w:rsidRPr="0075021F"><w:rPr><w:rFonts w:asciiTheme="minorHAnsi" w:hAnsiTheme="minorHAnsi" w:cstheme="minorHAnsi"/><w:b/><w:bCs/></w:rPr><w:t>${it.ip}</w:t></w:r></w:p></w:tc><w:tc><w:tcPr><w:tcW w:w="9450" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="auto"/></w:tcPr>${splitHeadedList(it.cveParagraph())}</w:tc></w:tr>"""
            
        }.joinToString("")
    }
}

class CreateCredLeakFileDoc(
    val orgName: String,
    val formNumber: String,
    val leakOrLeaks: Boolean,
    val victimOrg: String,
    val credOrCreds: Boolean,
    val passwordOption: String,
    val userPass: String,
    val addInfo: String,
    val changeTlp: Boolean,
    val additionalReferences: String,
    val ipAddress: String,
) : CreateDoc {
    val docname = ".${File.separator}credential_leak_edited.docx"

    override fun createDoc() {
        val tempDoc = ".${File.separator}temp"
        this.writeToFile(tempDoc)
        zipDocument(tempDoc, orgName, getAlertId(formNumber))
    }

    override fun writeToFile(tempDoc: String) {
        ZipFile(".${File.separator}reference${File.separator}${this.docname}").extractAll(tempDoc)

        val doc = File(".${File.separator}temp${File.separator}word${File.separator}document.xml")

        var docxFile = doc.readText()
        if (changeTlp) {
            docxFile = docxFile.replace("tlp_insert", "Recipients may share TLP:GREEN information with peers and partner organizations within their community, but not via publicly accessible channels. Unless otherwise specified, TLP:GREEN information may not be shared outside of the cybersecurity or cyber defense community.")
            writeToHeaderFooter()
        } else {
            docxFile = docxFile.replace("tlp_insert", "Recipients may share TLP:AMBER information with members of their own organization and its clients on a need-to-know basis to protect their organization and its clients and prevent further harm.")
        }

        docxFile = docxFile.replace("org_name_insert", orgName)
        docxFile = docxFile.replace("alert_id", getAlertId(formNumber))

        if (leakOrLeaks) {
            docxFile = docxFile.replace("leak_or_leaks_insert", "credential leaks")
        } else {
            docxFile = docxFile.replace("leak_or_leaks_insert", "a credential leak")
        }

        if (victimOrg == "") {
            docxFile = docxFile.replace("victim_org_insert", orgName)
        } else {
            docxFile = docxFile.replace("victim_org_insert", victimOrg)
        }

        if (credOrCreds) {
            docxFile = docxFile.replace("cred_or_creds", "these credentials")
            docxFile = docxFile.replace("cred_is_are", "any credentials are")
            docxFile = docxFile.replace("inc_or_incs", "any of these credentials as a security incident")
        } else {
            docxFile = docxFile.replace("cred_or_creds", "this credential")
            docxFile = docxFile.replace("cred_is_are", "this credential is")
            docxFile = docxFile.replace("inc_or_incs", "this credential as a security incident")
        }

        docxFile = docxFile.replace("password_option_insert", passwordOption)

        if (userPass != "") {
            docxFile = docxFile.replace("""<w:p w14:paraId="63BAFEA2" w14:textId="73770817" w:rsidR="001F1CB8" w:rsidRPr="008853E4" w:rsidRDefault="001F1CB8" w:rsidP="00EE6CD1"><w:pPr><w:pStyle w:val="xxxmsonormal"/><w:shd w:val="clear" w:color="auto" w:fill="FFFFFF"/><w:spacing w:before="0" w:beforeAutospacing="0" w:after="0" w:afterAutospacing="0"/><w:rPr><w:rFonts w:ascii="Calibri" w:hAnsi="Calibri" w:cs="Calibri"/><w:color w:val="000000"/><w:sz w:val="20"/><w:szCs w:val="20"/><w:bdr w:val="none" w:sz="0" w:space="0" w:color="auto" w:frame="1"/><w:shd w:val="clear" w:color="auto" w:fill="FFFFFF"/></w:rPr></w:pPr><w:proofErr w:type="spellStart"/><w:r w:rsidRPr="008853E4"><w:rPr><w:rFonts w:ascii="Calibri" w:hAnsi="Calibri" w:cs="Calibri"/><w:color w:val="000000"/><w:sz w:val="20"/><w:szCs w:val="20"/><w:bdr w:val="none" w:sz="0" w:space="0" w:color="auto" w:frame="1"/><w:shd w:val="clear" w:color="auto" w:fill="FFFFFF"/></w:rPr><w:t>user_pass_insert</w:t></w:r><w:proofErr w:type="spellEnd"/></w:p>""", splitParagraph(userPass))
        } else {
            docxFile = docxFile.replace("user_pass_insert", "")
        }

        docxFile = docxFile.replace("add_info_insert", addInfo)

        //reference_insert
        if (additionalReferences != "") {
            docxFile = docxFile.replace("reference_insert", additionalReferences)
        } else {
            val refstring = docxFile.indexOf("reference_insert")
            val start = docxFile.lastIndexOf("<w:tr ", refstring)
            val end = docxFile.indexOf("</w:tr>", refstring) + 7
            val refsub = docxFile.substring(start, end)
            docxFile = docxFile.replace(refsub, "")
        }

        doc.writeText(docxFile)
    }

    override fun getOrgname(): String {
        return this.orgName
    }

    override fun getAlertId(): String {
        return getAlertId(this.formNumber)
    }

    override fun getIps(): String {
        return this.ipAddress
    }

    override fun getAlertName(): String {
        if (changeTlp) {
            return "TLP GREEN: ${orgName} CL-${ipAddress}: ${getAlertId()}"
        } else {
            return "TLP AMBER: ${orgName} CL-${ipAddress}: ${getAlertId()}"
        }
    }
}

class CreateActorProfile(
    val name: String, 
    val alias: String, 
    val date: String, 
    val country: String, 
    val motivation: String, 
    val target: String, 
    val malware: String, 
    val reporter: String,
    val confidence: String, 
    val exploits: String, 
    val summary: String,
    val capabilities: String, 
    val detection: String, 
    val ttps: String, 
    val infra: String, 
) : CreateDoc {
    val docname = ".${File.separator}actor_profile.docx"

    override fun createDoc() {
        val tempDoc = ".${File.separator}temp"
        this.writeToFile(tempDoc)
        zipDocument(tempDoc, name.replace(Regex("[\\/:*?\"<>|]"), "-"), "")

    }

    override fun writeToFile(tempDoc: String) {
        ZipFile(".${File.separator}reference${File.separator}${this.docname}").extractAll(tempDoc)

        val doc = File(".${File.separator}temp${File.separator}word${File.separator}document.xml")

        var docxFile = doc.readText()
        docxFile = docxFile.replace("{name}", name)
        docxFile = docxFile.replace("{alias}", alias)
        docxFile = docxFile.replace("{date}", date)
        docxFile = docxFile.replace("{country}", country)
        docxFile = docxFile.replace("{motivation}", motivation)
        docxFile = docxFile.replace("{confidence}", confidence)

        val details = """This malware is unique and only seen with ${country}s' nation state attacks. According to ${reporter}, the ${name} malware family is the name given to malware developed and controlled by an intelligence directorate supporting the nation state ${country}.""""
        docxFile = docxFile.replace("{details}", details)
        docxFile = docxFile.replace("{cve}", exploits)
        docxFile = docxFile.replace("{summary}", summary)

        if (capabilities != "") {
            val refString = docxFile.indexOf("{capabilities}")
            var start = docxFile.lastIndexOf("<w:p>", refString)

            val end = docxFile.indexOf("</w:p>", refString) + 6
            if (start != -1 && end != -1 && refString != -1) {
                val refsub = docxFile.substring(start, end)
                docxFile = docxFile.replace(refsub, splitBulletPoint(capabilities))
            }
        } else {
            docxFile = docxFile.replace("{capabilities}", "")
        }

        if (detection != "") {
            val refString = docxFile.indexOf("{detection}")
            var start = docxFile.lastIndexOf("<w:p>", refString)

            val end = docxFile.indexOf("</w:p>", refString) + 6
            if (start != -1 && end != -1 && refString != -1) {
                val refsub = docxFile.substring(start, end)
                docxFile = docxFile.replace(refsub, splitBulletPoint(detection))
            }
        } else {
            docxFile = docxFile.replace("{detection}", "")
        }

        if (ttps != "") {
            val refString = docxFile.indexOf("{ttps}")
            var start = docxFile.lastIndexOf("<w:p>", refString)

            val end = docxFile.indexOf("</w:p>", refString) + 6
            if (start != -1 && end != -1 && refString != -1) {
                val refsub = docxFile.substring(start, end)
                docxFile = docxFile.replace(refsub, splitBulletPoint(ttps))
            }
        } else {
            docxFile = docxFile.replace("{ttps}", "")
        }

        if (infra != "") {
            val refString = docxFile.indexOf("{infra}")
            var start = docxFile.lastIndexOf("<w:p>", refString)

            val end = docxFile.indexOf("</w:p>", refString) + 6
            if (start != -1 && end != -1 && refString != -1) {
                val refsub = docxFile.substring(start, end)
                docxFile = docxFile.replace(refsub, splitBulletPoint(infra))
            }
        } else {
            docxFile = docxFile.replace("{infra}", "")
        }

        doc.writeText(docxFile)
    }

    override fun getOrgname(): String {
        return ""
    }

    override fun getAlertId(): String {

        return ""
    }

    override fun getIps(): String {

        return ""
    }

    override fun getAlertName(): String {

        return ""
    }
}

fun splitList(list: String): String {
    return list
    .split("[,\\s]".toRegex())
    .filter { it != "" }
    .map{ "<w:p w14:paraId=\"63BAFEA2\" w14:textId=\"73770817\" w:rsidR=\"001F1CB8\" w:rsidRPr=\"008853E4\" w:rsidRDefault=\"001F1CB8\" w:rsidP=\"00EE6CD1\"><w:pPr><w:pStyle w:val=\"xxxmsonormal\"/><w:shd w:val=\"clear\" w:color=\"auto\" w:fill=\"FFFFFF\"/><w:spacing w:before=\"0\" w:beforeAutospacing=\"0\" w:after=\"0\" w:afterAutospacing=\"0\"/><w:rPr><w:rFonts w:ascii=\"Calibri\" w:hAnsi=\"Calibri\" w:cs=\"Calibri\"/><w:color w:val=\"000000\"/><w:sz w:val=\"20\"/><w:szCs w:val=\"20\"/><w:bdr w:val=\"none\" w:sz=\"0\" w:space=\"0\" w:color=\"auto\" w:frame=\"1\"/><w:shd w:val=\"clear\" w:color=\"auto\" w:fill=\"FFFFFF\"/></w:rPr></w:pPr><w:r w:rsidRPr=\"008853E4\"><w:rPr><w:rFonts w:ascii=\"Calibri\" w:hAnsi=\"Calibri\" w:cs=\"Calibri\"/><w:color w:val=\"000000\"/><w:sz w:val=\"20\"/><w:szCs w:val=\"20\"/><w:bdr w:val=\"none\" w:sz=\"0\" w:space=\"0\" w:color=\"auto\" w:frame=\"1\"/><w:shd w:val=\"clear\" w:color=\"auto\" w:fill=\"FFFFFF\"/></w:rPr><w:t>${it}</w:t></w:r></w:p>"}
    .joinToString("")
}

fun splitHeadedList(list: String): String {
    var builder = StringBuilder()

    for (line in list.split("\n")) {
        if (line == "") {
            continue
        }
        
        val parts = line.split(" ", limit=2)

        if (parts.size > 1) {
            val head = parts[0]
            val body = parts[1]
            builder.append("""<w:p><w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">${head} </w:t></w:r><w:r><w:t xml:space="preserve">${body}</w:t></w:r></w:p>""")
        } else {
            val head = parts[0]
            builder.append("""<w:p><w:t xml:space="preserve">${head}</w:t></w:p>""")
        }
        builder.append("<w:p></w:p>")
    }
    
    return builder.toString()
}

fun splitBulletPoint(list: String): String {
    return list
    .split("\n")
    .map{ """<w:p w14:paraId="485619EC" w14:textId="04A28E1E" w:rsidR="001961C3" w:rsidRPr="00196A91" w:rsidRDefault="00196A91" w:rsidP="001961C3"><w:pPr><w:pStyle w:val="ListParagraph"/><w:numPr><w:ilvl w:val="0"/><w:numId w:val="1"/></w:numPr><w:jc w:val="both"/><w:rPr><w:rFonts w:ascii="Arial" w:hAnsi="Arial" w:cs="Arial"/><w:color w:val="000000"/></w:rPr></w:pPr><w:r w:rsidRPr="00196A91"><w:rPr><w:rFonts w:ascii="Arial" w:hAnsi="Arial" w:cs="Arial"/><w:color w:val="000000"/></w:rPr><w:t>${it}</w:t></w:r></w:p>"""}
    .joinToString("")
}

fun splitParagraph(list: String): String {
    return list
    .split("\n")
    .map{ "<w:p w14:paraId=\"63BAFEA2\" w14:textId=\"73770817\" w:rsidR=\"001F1CB8\" w:rsidRPr=\"008853E4\" w:rsidRDefault=\"001F1CB8\" w:rsidP=\"00EE6CD1\"><w:pPr><w:pStyle w:val=\"xxxmsonormal\"/><w:shd w:val=\"clear\" w:color=\"auto\" w:fill=\"FFFFFF\"/><w:spacing w:before=\"0\" w:beforeAutospacing=\"0\" w:after=\"0\" w:afterAutospacing=\"0\"/><w:rPr><w:rFonts w:ascii=\"Calibri\" w:hAnsi=\"Calibri\" w:cs=\"Calibri\"/><w:color w:val=\"000000\"/><w:sz w:val=\"20\"/><w:szCs w:val=\"20\"/><w:bdr w:val=\"none\" w:sz=\"0\" w:space=\"0\" w:color=\"auto\" w:frame=\"1\"/><w:shd w:val=\"clear\" w:color=\"auto\" w:fill=\"FFFFFF\"/></w:rPr></w:pPr><w:r w:rsidRPr=\"008853E4\"><w:rPr><w:rFonts w:ascii=\"Calibri\" w:hAnsi=\"Calibri\" w:cs=\"Calibri\"/><w:color w:val=\"000000\"/><w:sz w:val=\"20\"/><w:szCs w:val=\"20\"/><w:bdr w:val=\"none\" w:sz=\"0\" w:space=\"0\" w:color=\"auto\" w:frame=\"1\"/><w:shd w:val=\"clear\" w:color=\"auto\" w:fill=\"FFFFFF\"/></w:rPr><w:t>${it}</w:t></w:r></w:p>"}
    .joinToString("")
}

fun getAlertId(form_number: String): String {
    val formatter = DateTimeFormatter.ofPattern("yyyyMMdd")
    val current = LocalDateTime.now().format(formatter)
    if (form_number.length == 1) {
        return "${current}0${form_number}"
    }
    return "${current}${form_number}"
}

fun zipDocument(tempDoc: String, orgName: String, alertId: String) {
    val tempDir = "${tempDoc}${File.separator}"
    val newFolder = File(".${File.separator}${orgName.trim()}")
    val newReport = "${newFolder}${File.separator}${orgName.trim()}-${alertId}.docx"

    if (!newFolder.exists()) {
        if (!newFolder.mkdirs()) {
            throw Exception("Unable to create folder: ${orgName}")
        }
    }

    val newZip = ZipFile(newReport)
    File(tempDir).listFiles().forEach {
        if (it.isDirectory()) {
            newZip.addFolder(it)
        } else {
            newZip.addFile(it)
        }
    }

    File(tempDoc).deleteRecursively()
}

fun writeToHeaderFooter() {
    //#33FF00
    // <w:color w:val="FFC000"/><w:sz w:val="24"/><w:szCs w:val="16"/></w:rPr><w:t>TLP: AMBER</w:t>
    for (i in 0..3) {
        val header = File(".${File.separator}temp${File.separator}word${File.separator}header${i}.xml")
        if (header.exists()) {
            var headerFile = header.readText()

            headerFile = headerFile.replace("<w:color w:val=\"FFC000\"/><w:sz w:val=\"16\"/><w:szCs w:val=\"16\"/></w:rPr><w:t>TLP: AMBER</w:t>", "<w:color w:val=\"33FF00\"/><w:sz w:val=\"16\"/><w:szCs w:val=\"16\"/></w:rPr><w:t>TLP: GREEN</w:t>")
            header.writeText(headerFile)
        }

        val footer = File(".${File.separator}temp${File.separator}word${File.separator}footer${i}.xml")
        if (footer.exists()) {
            var footerFile = footer.readText()

            footerFile =  footerFile.replace("<w:color w:val=\"FFC000\"/><w:sz w:val=\"16\"/><w:szCs w:val=\"16\"/></w:rPr><w:t>TLP: AMBER</w:t>", "<w:color w:val=\"33FF00\"/><w:sz w:val=\"16\"/><w:szCs w:val=\"16\"/></w:rPr><w:t>TLP: GREEN</w:t>")
            footer.writeText(footerFile)
        }
    }
}

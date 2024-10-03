package form.create

import form.create.CreateCredLeakFileDoc
import form.create.CreateOpenPortDoc
import form.events.Event


class OpenPortBuilder() {

    fun buildFromMap(map: HashMap<String, String>, events: List<Event> = listOf()): CreateDoc {
        return CreateOpenPortDoc(
            map["formNumber"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["threatType"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["changeTlp"] == "Green",
            map["additionalReferences"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["summary"]?: "",
            map["body"]?: "",
            events,
        )
    }
}

class CredLeakFileBuilder() {

    fun buildFromMap(map: HashMap<String, String>): CreateDoc {
        return CreateCredLeakFileDoc(
        map["orgName"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
        map["formNumber"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "", 
        map["leakOrLeaks"] == "on", 
        map["victimOrg"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "", 
        map["credOrCreds"] == "on", 
        map["passwordOption"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "", 
        map["userPass"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "", 
        map["addInfo"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "", 
        map["changeTlp"] == "Green", 
        map["additionalReferences"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "", 
        map["ipAddress"]?.trim()?.replace(">", ")")?.replace("<", "(")?: ""
        )
    }
}

class ActorProfileBuilder() {
    fun buildFromMap(map: HashMap<String, String>): CreateDoc {
        return CreateActorProfile(
            map["name"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["alias"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["date"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["country"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["motivation"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["target"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["malware"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["reporter"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["confidence"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["exploits"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["summary"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["capabilities"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["detection"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["ttps"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
            map["infra"]?.trim()?.replace(">", ")")?.replace("<", "(")?: "",
        )
    }
}

fun createDocument(doc: CreateDoc, ips: Boolean) {

    if (ips) {
        val hosts = getHosts(doc.getIps())
        writeToCsv(doc.getOrgname(), doc.getAlertId(), hosts)
    }

    doc.createDoc()
}

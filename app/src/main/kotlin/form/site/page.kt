package form.site

enum class Page {
    CredLeak,
    OpenPort,
    Filter,
    Events,
    Actor,
}

enum class Form {
    OpenPort,
    CredLeak,
    EndOfLife,
    LoginPage,
}

//builds the entire webpage, should only use on first loaf of the app
fun buildPage(content: String, page: Page): String {
    
    return """
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>JCTF Form Generator</title>

    <!-- Include HTMX library -->
    <script src="https://unpkg.com/htmx.org@1.7.0/dist/htmx.js"></script>
    <script src="https://unpkg.com/htmx.org/dist/ext/ws.js"></script>
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        document.addEventListener('htmx:wsError', function(event) {
            window.location.href = "/"
        });
    </script>
</head>
<body hx-boost="true">
    <!-- Allows the server to know when to quit -->
    <div hx-ext="ws" ws-connect="/conn"></div>

    <div id="content">
        ${buildContent(content, page)}
    </div>
</body>
</html>
    """
}

//buils the content part of the webpage
fun buildContent(content: String, page: Page): String {
    return """
        ${buildHeader(page)}
        ${buildTitle(page)}
        <div class="flex justify-center">
            <div id="ring" class="w-1/2 ring-2 ring-inset ring-yellow-500 p-5 mx-5 mt-2 rounded mb-20">
                $content
            </div>
        </div
    """
}

fun buildTitle(page: Page): String {
    val title = when (page) {
        // Page.Form -> "Generate Form"
        Page.CredLeak -> "CredLeak"
        Page.OpenPort -> "OpenPort"
        Page.Events -> "Events"
        Page.Filter -> "Filters"
        Page.Actor -> "Actor Profile"
        // Page.Prio -> "Prioritizer"
        // Page.Populate -> "Populate"
    }
    return """
        <h1 class="my-4 text-3xl font-bold text-center">$title</h1>
    """
}

//creates the tabs header at the top of the pages
fun buildHeader(page: Page): String {

    return when(page) {
        Page.CredLeak ->
    """
    <div class="py-4 bg-blue-800 rounded-b-lg flex justify-between">
        <div>
            <button class="rounded bg-white p-2 mx-2 ring-2 ring-yellow-500 ring-inset" hx-get="/credleak" hx-target="#content" >Cred Leak</button>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/populate" hx-target="#content" >Open Port</button>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/actor" hx-target="#content" >Actor</button>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/events" hx-target="#content" >Events</button>
        </div>
    </div>
    """
    Page.OpenPort ->
    """
    <div class="py-4 bg-blue-800 rounded-b-lg flex justify-between">
        <div>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/credleak" hx-target="#content" >Cred Leak</button>
            <button class="rounded bg-white p-2 mx-2 ring-2 ring-yellow-500 ring-inset" hx-get="/populate" hx-target="#content" >Open Port</button>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/actor" hx-target="#content" >Actor</button>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/events" hx-target="#content" >Events</button>
        </div>
    </div>
    """
        Page.Events ->
    """
    <div class="py-4 bg-blue-800 rounded-b-lg flex justify-between">
        <div>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/credleak" hx-target="#content" >Cred Leak</button>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/populate" hx-target="#content" >Open Port</button>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/actor" hx-target="#content" >Actor</button>
            <button class="rounded bg-white p-2 mx-2 ring-2 ring-yellow-500 ring-inset" hx-get="/events" hx-target="#content" >Events</button>
        </div>
        <div>
            <button class="rounded text-white hover:text-black hover:bg-white p-2 mx-2" hx-post="/events/download" hx-indicator="#load" hx-target="#ring">Update Events</button>
            <button class="rounded text-white hover:text-black hover:bg-white p-2 mx-2" hx-get="/filter" hx-indicator="#load" hx-target="#content">Filters</button>
        </div>
    </div>
    """
        Page.Filter -> """
    <div class="py-4 bg-blue-800 rounded-b-lg flex justify-between">
        <div>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/credleak" hx-target="#content" >Cred Leak</button>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/populate" hx-target="#content" >Open Port</button>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/actor" hx-target="#content" >Actor</button>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/events" hx-target="#content" >Events</button>
        </div>
        <div>
            <button class="rounded text-white bg-white p-2 mx-2 ring-2 ring-yellow-500 ring-inset" hx-get="/filter" hx-indicator="#load" hx-target="#content">Filters</button>
        </div>
    </div>
    """
    Page.Actor -> """
    <div class="py-4 bg-blue-800 rounded-b-lg flex justify-between">
        <div>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/credleak" hx-target="#content" >Cred Leak</button>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/populate" hx-target="#content" >Open Port</button>
            <button class="rounded bg-white p-2 mx-2 ring-2 ring-yellow-500 ring-inset" hx-get="/actor" hx-target="#content" >Actor</button>
            <button class="rounded hover:bg-white hover:text-black text-white p-2 mx-2" hx-get="/events" hx-target="#content" >Events</button>
        </div>
    </div>
    """
    }
}

//wraps form parts in similar spacing
private fun buildFormEntry(input: String): String {
    return """
        <div class="flex items-center justify-start mx-4 my-2">
            $input
        </div>
    """
}

//creates form textfield
fun buildFormText(name: String, param: String, default: String = ""): String {
    val entry = """
    <label class="block text-sm font-medium text-gray-700 mr-5">$name:</label>
    <input type="text" name="$param" value="$default" class="mt-1 p-2 border border-gray-300 rounded-md w-full">
    """
    return buildFormEntry(entry)
}

//creates textarea in form
fun buildFormArea(name: String, param: String, default: String = ""): String {
    val entry = """
    <label class="block text-sm font-medium text-gray-700 mr-5">$name:</label>
    <textarea name="$param" rows=8 class="mt-1 p-2 border border-gray-300 rounded-md w-full">$default</textarea>
    """
    return buildFormEntry(entry)
}

//creates radio buttons
fun buildFormRadio(names: List<String>, param: String, checked: Int = 0, title: String = ""): String {
    val entryBuilder = StringBuilder()

    entryBuilder.append("""<hr><div class="flex flex-col justify-start mx-4 my-2">""")
    if (title != "") {
        entryBuilder.append("""<label class="block text-sm font-medium text-gray-700 mr-5">$title:</label>""")
    }


    for ((i, name) in names.withIndex()) {
        val check = when(i) {
            checked -> "checked"
            else -> ""
        }

        entryBuilder.append("""
        <div class="flex items-center">
            <input type="radio" id="$name" value="$name" name="$param" class="mt-1 ml-10 p-2" $check>
            <label for="$name" class="block text-sm font-medium text-gray-700 ml-5">$name</label>
        </div>
        <br>
        """
        )
    }

    entryBuilder.append("</div><hr>")
    return entryBuilder.toString()
}

//creates check mark buttons
fun buildFormCheck(name: String, param: String): String {
    val entry = """
    <input type="checkbox" id="$param" name="$param" class="mt-1 ml-10 p-2">
    <label for="$param" class="block text-sm font-medium text-gray-700 ml-5">$name</label><br>
    """

    return buildFormEntry(entry)
}

fun buildButton(name: String, type: String): String {
    return """
       <button type="$type" class="w-full mt-4 px-4 py-2 bg-blue-800 text-white rounded-md">$name</button>
    """
}

fun progress(percent: Int): String {
    return """
    <div id="progress">
        <div id="load" class="text-center animate-bounce">loading...${percent}%</div>

        <div>
            <div class="bg-gray-200 h-7 rounded-full relative">
                <div class="bg-blue-800 absolute h-full rounded-full" style="width: ${percent}%"></div>
            </div>
        <div>
    </div>
    """
}

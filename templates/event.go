package templates

import (

)

func Event() string {
    data := struct {

    } {

    }

    const page = `
        <h1>Event</h1>
        `

    return Execute("event", page, data)
}

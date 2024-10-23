package templates

import (

)

func OpenPort() string {
    data := struct {

    } {

    }

    const page = `
        <h1>Open Port</h1>
        `

    return Execute("openport", page, data)
}

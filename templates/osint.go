package templates

import (

)

func Osint() string {
    data := struct {

    } {

    }


    const page = `

`

    return Execute("osint", page, data)
}

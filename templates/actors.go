package templates

import (

)

func Actors() string {
    data := struct {

    } {

    }

    const page = `
        <h1>Actor</h1>
        `

    return Execute("actor", page, data)
}

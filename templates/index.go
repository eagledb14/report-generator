package templates

import (
)

func Index() string {
    data := struct {

    } {

    }

    const page = `
        <div>
            index
        </div>
        `

    return ExecuteText("index", page, data)
}

package types

type Page int

const (
    Actor Page = iota
    CredLeak
    Events
    OpenPort
)

type Form int

const (
    Open Form = iota
    EOL 
    Login
)

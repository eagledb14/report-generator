package types

type ReportType int

const (
    Header ReportType = iota
    Cover
)

type Form int

const (
    Open Form = iota
    EOL 
    Login
)

var FormName = map[Form]string{
	Open: "Open Port",
	EOL: "End of Life",
	Login: "Login Page",
}

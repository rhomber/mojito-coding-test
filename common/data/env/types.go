package env

type Env string

const (
	Local Env = "LOCAL"
	Test  Env = "TEST"
	Prod  Env = "PROD"
)

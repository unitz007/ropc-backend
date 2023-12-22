package kernel

type Logger interface {
	Info(v string)
	Error(v string)
	Warn(v string)
	Fatal(v string)
}

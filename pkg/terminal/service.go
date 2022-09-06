package terminal

type Service interface {
	Start() (err error)
	Stop()
}

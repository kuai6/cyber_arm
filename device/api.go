package device

type Device interface {
	Start() error
	Stop()
}

package consumer

type Consumer interface {
	GetName() string
	GetType() string
	Start() error
	Stop() error
}

package domain

type Domain interface {
	Init() func() error
}

package server

type Storer interface {
	GetAll() []Response
	Add(Response) error
	RemoveAll() error
}

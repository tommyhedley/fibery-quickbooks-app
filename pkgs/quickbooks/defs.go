package quickbooks

type Entity interface {
	Read() (Entity, error)
}

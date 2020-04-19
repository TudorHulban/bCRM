package interfaces

// IAccount Interface provides decoupling of persistence solution for account related operations.
type IStore interface {
	Open() error
	Close() error
}

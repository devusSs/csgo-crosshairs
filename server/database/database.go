package database

type Service interface {
	TestConnection() error
	CloseConnection() error
}

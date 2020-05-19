package models

type SagaInstance struct {
	ID           string
	SagaType     string
	SagaData     []byte
	CurrentState string
}

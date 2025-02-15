package model

type Transaction struct {
	IsSender               bool
	CounterpartyEmployeeID int64
	CounterpartyUsername   string
	Amount                 int64
}

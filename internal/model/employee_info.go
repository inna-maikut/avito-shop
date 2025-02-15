package model

type EmployeeInfo struct {
	Coins                int64
	Inventory            []Inventory
	ReceivedTransactions []Transaction
	SentTransactions     []Transaction
}

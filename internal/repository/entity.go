package repository

type Employee struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Balance  int64  `db:"balance"`
}

type Merch struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Price int64  `db:"price"`
}

type EmployeeTransaction struct {
	ID                     int64  `db:"id"`
	IsSender               bool   `db:"is_sender"`
	CounterpartyEmployeeID int64  `db:"counterparty_employee_id"`
	CounterpartyUsername   string `db:"counterparty_username"`
	Amount                 int64  `db:"amount"`
}

type InventoryWithMerchName struct {
	EmployeeID int64  `db:"employee_id"`
	MerchID    int64  `db:"merch_id"`
	Quantity   int64  `db:"quantity"`
	MerchName  string `db:"merch_name"`
}

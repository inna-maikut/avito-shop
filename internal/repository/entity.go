package repository

type Employee struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Balance  int    `db:"balance"`
}

type Merch struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Price int    `db:"price"`
}

type Transaction struct {
	SenderID   int64 `json:"sender_id"`
	ReceiverID int64 `json:"receiver_id"`
	Amount     int   `json:"amount"`
}

package Model

type User struct {
	ID         int64  `db:"user_id" json:"user_id"`
	Username   string `db:"user_name" json:"user_name"`
	Firstname  string `db:"first_name" json:"first_name"`
	Lastname   string `db:"last_name" json:"last_name"`
	Email      string `db:"email" json:"email"`
	UserStatus string `db:"user_status" json:"user_status"`
	Department string `db:"department" json:"department"`
}

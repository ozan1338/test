package domain

type Job struct {
	ID          int    `db:"id"`
	Id_User     int    `db:"id_user"`
	City        string `db:"city"`
	Full_Time   bool   `db:"full_time"`
	Description string `db:"description"`
}
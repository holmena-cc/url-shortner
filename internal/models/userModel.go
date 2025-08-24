package models

import "time"

type User struct {
    UserID       int       `db:"user_id" json:"userId"`
    Email        string    `db:"email" json:"email"`
    PasswordHash string    `db:"password_hash" json:"-"`
    Name         string    `db:"name" json:"name"`
    RegisterDate time.Time `db:"register_date" json:"registerDate"`
}
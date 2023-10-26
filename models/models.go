package models

import (
	uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Sn_no    uint64    `gorm:"autoIncrement:true"`
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserName string    `gorm:"NOT NULL"`
	Password string    `gorm:"NOT NULL"`
	Email    string    `gorm:"NOT NULL;UNIQUE"`
	Is_Admin bool      `gorm:"DEFAULT=false"`
}

type Serial struct{

	Ser int

}

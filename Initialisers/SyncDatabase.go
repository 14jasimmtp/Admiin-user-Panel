package initialisers

import "main.go/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})

	// DB.Exec("ALTER TABLE users ADD COLUMN sn_no bigserial ;")
}

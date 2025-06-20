package main

import (
	// _ "log"
	// _ "simpleauth_server/src/database"
	"simpleauth_server/src/server"
)

func main() {
	// var dberr error = database.Connect("postgres://postgres:1qaz!QAZ1qaz@localhost/postgres")
	// if dberr != nil {
	// 	log.Fatalf("Database connection error: %v", dberr)
	// } else {
	// 	log.Println("Database Connect OK")
	// }

	var cfg *server.Config = server.ConfigLoad()

	//TODO: Блок проверки параметров конфигурации
	//TODO: Блок проверки переменных окружения SECRET_KEY и DB_CONNECTION
	server.Start(cfg)
}

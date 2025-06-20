package handlers

import (
	"simpleauth_server/src/database/models"

	"simpleauth_server/src/database"

	"github.com/gofiber/fiber/v2"
)

// получение списка всех продуктов
// select current_database(), version()
func GetDatabaseInfo(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT current_database(), version()")
	if err != nil {
		return c.Status(500).SendString("Ошибка выполнения запроса к базе данных")
	}
	defer rows.Close()

	var dbinfo []models.Dbinfo
	for rows.Next() {
		var dbrow models.Dbinfo
		err := rows.Scan(&dbrow.Current_database, &dbrow.Version)
		if err != nil {
			return c.Status(500).SendString("Ошибка сканирования данных")
		}
		dbinfo = append(dbinfo, dbrow)
	}

	return c.JSON(dbinfo)
}

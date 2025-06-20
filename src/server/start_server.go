package server

import (
	"log"
	"os"
	"simpleauth_server/src/database"
	"simpleauth_server/src/database/handlers"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

const (
	ContextKeyUser = "user"
	DefSecretKey   = "very-secret-key"
)

var (
	authStorage *AuthStorage = &AuthStorage{map[string]User{}}
	authHandler *AuthHandler = &AuthHandler{storage: authStorage}
	userHandler *UserHandler = &UserHandler{storage: authStorage}
	useDatabase bool         = false
)

func Start(cfg *Config) {

	app := fiber.New()

	// Инициализация ключа шифрования
	Init_SecretKey()
	// Соединени с базой данных
	Connect_Database()
	// Создание пользователей из файла конфигурации (user1, user2 и user3)
	Create_Users(cfg)

	// Группа обработчиков, которые доступны неавторизованным пользователям
	publicGroup := app.Group("")
	publicGroup.Post("/register", authHandler.Register)
	publicGroup.Post("/login", authHandler.Login)

	// Группа обработчиков, которые требуют авторизации
	authorizedGroup := app.Group("")
	authorizedGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: jwtSecretKey,
		},
		ContextKey: ContextKeyUser,
	}))
	authorizedGroup.Get("/profile", userHandler.Profile)
	if useDatabase {
		authorizedGroup.Get("/dbinfo", handlers.GetDatabaseInfo)
	}

	log.Fatal(app.Listen(cfg.Address))
}

func Connect_Database() {

	db_connection := os.Getenv("DB_CONNECTION")
	if db_connection == "" {
		log.Println("Environment variable DB_CONNECTION is not set. Database handlers not used.")
		return
	}
	//database.Connect("postgres://postgres:password@host/database")
	//database.Connect("postgres://postgres:1qaz!QAZ1qaz@localhost/postgres")
	var dberr error = database.Connect(db_connection)
	if dberr != nil {
		log.Fatalf("Database connection error: %v", dberr)
	} else {
		useDatabase = true
		log.Println("Connect to database.")
	}

}

func Init_SecretKey() {
	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Printf("Environment variable JWT_SECRET is not set. Using default key \"%s\".", DefSecretKey)
		jwtSecretKey = []byte(DefSecretKey)
	} else {
		jwtSecretKey = []byte(jwt_secret)
		log.Println("JWT key applied.")
	}
}

func Create_Users(cfg *Config) {
	if cfg.User1.User_Email != "" && cfg.User1.User_Password != "" {
		authHandler.storage.users[cfg.User1.User_Email] = User{
			Email:    cfg.User1.User_Email,
			Name:     cfg.User1.User_Name,
			password: cfg.User1.User_Password,
			payload:  cfg.User1.User_Payload,
		}
		log.Printf("User 1 created: %s", cfg.User1.User_Email)
	} else {
		log.Println("User 1 not created.")
	}

	_, user2_exist := authHandler.storage.users[cfg.User2.User_Email]
	if user2_exist {
		log.Printf("User %s duplicated", cfg.User2.User_Email)
	}
	if cfg.User2.User_Email != "" && cfg.User2.User_Password != "" && !user2_exist {
		authHandler.storage.users[cfg.User2.User_Email] = User{
			Email:    cfg.User2.User_Email,
			Name:     cfg.User2.User_Name,
			password: cfg.User2.User_Password,
			payload:  cfg.User2.User_Payload,
		}
		log.Printf("User 2 created: %s", cfg.User2.User_Email)
	} else {
		log.Println("User 2 not created.")
	}

	_, user3_exist := authHandler.storage.users[cfg.User3.User_Email]
	if user3_exist {
		log.Printf("User %s duplicated", cfg.User3.User_Email)
	}
	if cfg.User3.User_Email != "" && cfg.User3.User_Password != "" && !user3_exist {
		authHandler.storage.users[cfg.User3.User_Email] = User{
			Email:    cfg.User3.User_Email,
			Name:     cfg.User3.User_Name,
			password: cfg.User3.User_Password,
			payload:  cfg.User3.User_Payload,
		}
		log.Printf("User 3 created: %s", cfg.User3.User_Email)
	} else {
		log.Println("User 3 not created.")
	}

}

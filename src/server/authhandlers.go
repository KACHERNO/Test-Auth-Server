package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"maps"
	// "time"

	// jwtware "github.com/gofiber/contrib/jwt"
	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
)

type (
	// Обработчик HTTP-запросов на регистрацию и аутентификацию пользователей
	AuthHandler struct {
		storage *AuthStorage
	}

	// Хранилище зарегистрированных пользователей
	// Данные хранятся в оперативной памяти
	AuthStorage struct {
		users map[string]User
	}

	// Структура данных с информацией о пользователе
	User struct {
		Email    string
		Name     string
		password string
		payload  string
	}
)

// Структура HTTP-запроса на регистрацию пользователя
type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Структура HTTP-запроса на вход в аккаунт
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Структура HTTP-ответа на вход в аккаунт
// В ответе содержится JWT-токен авторизованного пользователя
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

// Обработчик HTTP-запросов, которые связаны с пользователем
type UserHandler struct {
	storage *AuthStorage
}

// Структура HTTP-ответа с информацией о пользователе
type ProfileResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

var (
	errBadCredentials = errors.New("email or password is incorrect")
)

// Секретный ключ для JWT-токена необходимо хранить в безопасном месте
var jwtSecretKey []byte

// var jwtSecretKey = []byte("very-secret-key")

// Обработчик HTTP-запросов на регистрацию пользователя
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	regReq := RegisterRequest{}
	if err := c.BodyParser(&regReq); err != nil {
		log.Printf("Register failed: %v.", err)
		return fmt.Errorf("body parser: %w", err)
	}

	// Проверяем, что пользователь с таким email еще не зарегистрирован
	if _, exists := h.storage.users[regReq.Email]; exists {
		log.Printf("Email %s already registered...", h.storage.users[regReq.Email].Email)
		return errors.New("the user already exists")
	}

	// Сохраняем в память нового зарегистрированного пользователя
	h.storage.users[regReq.Email] = User{
		Email:    regReq.Email,
		Name:     regReq.Name,
		password: regReq.Password,
	}

	// log.Printf("Registered user %s (%s).", h.storage.users[regReq.Email].Name, h.storage.users[regReq.Email].Email)
	return c.SendStatus(fiber.StatusCreated)
}

// Обработчик HTTP-запросов на вход в аккаунт
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	regReq := LoginRequest{}
	if err := c.BodyParser(&regReq); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	// Ищем пользователя в памяти приложения по электронной почте
	user, exists := h.storage.users[regReq.Email]
	// Если пользователь не найден, возвращаем ошибку
	if !exists {
		return errBadCredentials
	}
	// Если пользователь найден, но у него другой пароль, возвращаем ошибку
	if user.password != regReq.Password {
		return errBadCredentials
	}

	// Генерируем JWT-токен для пользователя,
	// который он будет использовать в будущих HTTP-запросах

	// str := `{"https://hasura.io/jwt/claims":{"x-hasura-allowed-roles":["Admins","Users"],"x-hasura-default-role":"Users","x-hasura-user-id":"1"}}`
	const test_payload = `{"testkey":"testvalue"}`
	// user_payload := h.storage.users[regReq.Email].payload
	user_payload := user.payload
	var token_payload map[string]interface{}

	if user_payload != "" {
		if err := json.Unmarshal([]byte(user_payload), &token_payload); err != nil {
			log.Printf("Error payload %v", err)
			json.Unmarshal([]byte(test_payload), &token_payload)
		}
	} else {
		json.Unmarshal([]byte(test_payload), &token_payload)
	}
	// log.Printf("user_payload = %s", user_payload)
	// log.Printf("token_payload = %s", token_payload)

	payload_add := jwt.MapClaims{
		"sub": user.Email,
		// 	"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	maps.Copy(token_payload, payload_add)

	// Генерируем полезные данные, которые будут храниться в токене
	payload := jwt.MapClaims(token_payload)

	// Создаем новый JWT-токен и подписываем его по алгоритму HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		log.Printf("JWT token signing %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(LoginResponse{AccessToken: t})
}

// Обработчик HTTP-запросов на получение информации о пользователе
func (h *UserHandler) Profile(c *fiber.Ctx) error {
	jwtPayload, ok := jwtPayloadFromRequest(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userInfo, ok := h.storage.users[jwtPayload["sub"].(string)]
	if !ok {
		return errors.New("user not found")
	}

	return c.JSON(ProfileResponse{
		Email: userInfo.Email,
		Name:  userInfo.Name,
	})
}

func jwtPayloadFromRequest(c *fiber.Ctx) (jwt.MapClaims, bool) {
	jwtToken, ok := c.Context().Value(ContextKeyUser).(*jwt.Token)
	if !ok {
		// logrus.WithFields(logrus.Fields{
		// 	"jwt_token_context_value": c.Context().Value(ContextKeyUser),
		// }).Error("wrong type of JWT token in context")
		log.Printf("Wrong type of JWT token in context: %s", c.Context().Value(ContextKeyUser))
		return nil, false
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		// logrus.WithFields(logrus.Fields{
		// 	"jwt_token_claims": jwtToken.Claims,
		// }).Error("wrong type of JWT token claims")
		log.Printf("Wrong type of JWT token claims: %s", jwtToken.Claims)
		return nil, false
	}

	return payload, true
}

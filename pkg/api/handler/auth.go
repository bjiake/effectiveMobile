package handler

import (
	"effectiveMobile/pkg/db"
	"effectiveMobile/pkg/domain/people"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

var jwtKey = []byte("secret_key")

func createToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString(jwtKey)
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
}

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags User
// @Accept  json
// @Produce  json
// @Param registration body people.Registration true "User registration info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {string} string "Email already exist"
// @Failure 500 {object} map[string]string
// @Router /registration [post]
func (h *Handler) Registration(c *gin.Context) {
	var acc people.Registration
	if err := c.BindJSON(&acc); err != nil {
		c.JSON(400, gin.H{"error bind Registration people": err.Error()})
		log.Error("error bind Registration people %v", err.Error())
		return
	}
	result, err := h.service.Registration(c.Request.Context(), acc)
	if err != nil {
		switch err.Error() {
		case db.ErrDuplicate.Error():
			c.JSON(409, "passportNumber already exist")
			log.Error("Register email failed %v", err.Error())
			break
		default:
			c.JSON(500, gin.H{"error Registration people": err.Error()})
			log.Error("error service Registration people %v", err.Error())
		}
		return
	}

	c.JSON(201, result)
	log.Info("Success registration: %v", result)
}

// @Summary Login a user
// @Description Login with email and password
// @Tags User
// @Accept  json
// @Produce  json
// @Param login body people.Registration true "User login info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	var acc people.Registration
	if err := c.BindJSON(&acc); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id, err := h.service.Login(c.Request.Context(), acc)
	if err != nil {
		switch err.Error() {
		case db.ErrNotExist.Error():
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Account does not exist"})
			log.Error("Account does not exist")
			break
		default:
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("error service Login error %v", err.Error())
		}
		return
	}

	token, err := createToken(strconv.FormatInt(id, 10))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error creating token"})
		return
	}
	sameSite := http.SameSiteNoneMode // Важно установить SameSite=None с Secure=true

	//c.SetCookie("token", token, 3600*72, "/", "", true, true)
	//c.SetSameSite(sameSite)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		Path:     "/",
		Domain:   "",
		SameSite: sameSite,
		Secure:   true,
		HttpOnly: true,
	})

	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
	log.Infof("Success login: %v", id)
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := parseToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userId", claims["id"])

		c.Next()
	}
}

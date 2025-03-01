package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/bigelle/online-shop/backend/internal/database"
	"github.com/bigelle/online-shop/backend/internal/models"
	"github.com/bigelle/online-shop/backend/internal/schemas"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var l schemas.Login
	if err := ctx.Bind(&l); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusBadRequest,
				Description: http.StatusText(http.StatusBadRequest),
			},
		)
		return
	}

	_, err := database.FindUser(h.DB, l.Email)
	if err != gorm.ErrRecordNotFound {
		ctx.JSON(
			http.StatusBadRequest,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusBadRequest,
				Description: "user with this username already exists",
			},
		)
		return
	}

	hashed, err := hashPassword(l.Password)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusBadRequest,
				Description: http.StatusText(http.StatusInternalServerError),
			},
		)
		return
	}
	l.Password = hashed

	if err := database.AddUser(h.DB, l); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusBadRequest,
				Description: http.StatusText(http.StatusInternalServerError),
			},
		)
		return
	}

	ctx.JSON(
		http.StatusAccepted,
		schemas.Response{
			Ok:   true,
			Code: http.StatusAccepted,
		},
	)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var l schemas.Login
	if err := ctx.Bind(&l); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusBadRequest,
				Description: http.StatusText(http.StatusBadRequest),
			},
		)
		return
	}

	usr, err := database.FindUser(h.DB, l.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(
				http.StatusNotFound,
				schemas.Response{
					Ok:          false,
					Code:        http.StatusNotFound,
					Description: "user with this email not found",
				},
			)
			return
		}
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusInternalServerError,
				Description: http.StatusText(http.StatusInternalServerError),
			},
		)
		return
	}

	if !checkPassword(l.Password, usr.HashedPassword) {
		ctx.JSON(
			http.StatusBadRequest,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusBadRequest,
				Description: "wrong password", // FIXME change the logic the way that you can't tell if its a non-existing username or wrong password
			},
		)
		return
	}

	h.addAuthCookies(ctx, *usr)

	ctx.JSON(
		http.StatusAccepted,
		schemas.Response{
			Ok:   true,
			Code: http.StatusAccepted,
		},
	)
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	usr, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusUnauthorized,
				Description: "unauthorized",
			},
		)
	}
	u := usr.(models.User)

	h.expireAuthCookies(ctx, u)

	ctx.JSON(
		http.StatusAccepted,
		schemas.Response{
			Ok:   true,
			Code: http.StatusAccepted,
		},
	)
}

func (h *AuthHandler) updateUser(usr models.User) error {
	err := h.DB.Model(&models.User{}).
		Where("email = ?", usr.Email).
		Select("*").
		Updates(map[string]interface{}{
			"session_token": usr.SessionToken,
			"csrf_token":    usr.CsrfToken,
		}).Error
	return err
}

func (h *AuthHandler) addAuthCookies(ctx *gin.Context, usr models.User) {
	st, err := generateToken(32)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusInternalServerError,
				Description: http.StatusText(http.StatusInternalServerError),
			},
		)
		return
	}
	csrf, err := generateToken(32)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusInternalServerError,
				Description: http.StatusText(http.StatusInternalServerError),
			},
		)
		return
	}
	ctx.SetCookie(
		"session_token",
		st,
		int(24*time.Hour),
		"/",
		"",
		false,
		true,
	)
	ctx.SetCookie(
		"csrf_token",
		csrf,
		int(24*time.Hour),
		"/",
		"",
		false,
		false,
	)
	usr.CsrfToken = csrf
	usr.SessionToken = st
	if err := h.updateUser(usr); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
}

func (h *AuthHandler) expireAuthCookies(ctx *gin.Context, usr models.User) {
	ctx.SetCookie(
		"session_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
	ctx.SetCookie(
		"csrf_token",
		"",
		-1,
		"/",
		"",
		false,
		false,
	)
	usr.SessionToken = ""
	usr.CsrfToken = ""
	if err := h.updateUser(usr); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
}

func hashPassword(pass string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	return string(b), err
}

func checkPassword(pass, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}

func generateToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

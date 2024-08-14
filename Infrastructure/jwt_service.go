package infrastructure

import (
    "os"
    "time"
    "errors"

    "github.com/dgrijalva/jwt-go"
    "github.com/Hailemari/clean_architecture_task_manager/Domain"
)

func GenerateToken(user *domain.User) (string, error) {
    claims := jwt.MapClaims{
        "id":       user.ID,
        "username": user.Username,
        "role":     user.Role,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
}
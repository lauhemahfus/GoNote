package utils

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v4"
)

type JWTUtil struct {
    secretKey string
}

func NewJWTUtil(secretKey string) *JWTUtil {
    return &JWTUtil{
        secretKey: secretKey,
    }
}

func (j *JWTUtil) GenerateToken(userID int) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
    })
    
    tokenString, err := token.SignedString([]byte(j.secretKey))
    if err != nil {
        return "", err
    }
    
    return tokenString, nil
}

func (j *JWTUtil) ValidateToken(tokenString string) (int, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(j.secretKey), nil
    })
    
    if err != nil {
        return 0, err
    }
    
    if !token.Valid {
        return 0, errors.New("invalid token")
    }
    
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return 0, errors.New("invalid token claims")
    }
    
    userID, ok := claims["user_id"].(float64)
    if !ok {
        return 0, errors.New("invalid user_id in token")
    }
    
    return int(userID), nil
}
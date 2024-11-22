package helper

import (
	"ekak_kabupaten_madiun/model/web"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var jwtIssuer = os.Getenv("JWT_ISSUER")
var jwtExpiration = os.Getenv("JWT_EXPIRATION")

func CreateNewJWT(userId int, email string, nip string, roles []string) string {
	exp := 24 * time.Hour
	if jwtExpiration != "" {
		if duration, err := time.ParseDuration(jwtExpiration + "h"); err == nil {
			exp = duration
		}
	}

	claims := jwt.MapClaims{
		"iss":     jwtIssuer,
		"user_id": userId,
		"email":   email,
		"nip":     nip,
		"roles":   roles,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		fmt.Println(err)
	}

	return signedToken
}

func ValidateJWT(tokenString string) web.JWTClaim {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		fmt.Println(err)
		return web.JWTClaim{}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var roles []string
		if rolesInterface, exists := claims["roles"]; exists {
			if rolesArray, ok := rolesInterface.([]interface{}); ok {
				for _, role := range rolesArray {
					if roleStr, ok := role.(string); ok {
						roles = append(roles, roleStr)
					}
				}
			}
		}

		return web.JWTClaim{
			Issuer: claims["iss"].(string),
			UserId: int(claims["user_id"].(float64)),
			Email:  claims["email"].(string),
			Nip:    claims["nip"].(string),
			Roles:  roles,
			Iat:    int64(claims["iat"].(float64)),
			Exp:    int64(claims["exp"].(float64)),
		}
	}

	return web.JWTClaim{}
}

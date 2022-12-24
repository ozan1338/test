package jwt

import (
	resError "test/util/errors_response"
	"time"
)

//go:generate mockgen -destination=../../mocks/pkg/jwt/mockJwt.go -package=jwt test/pkg/jwt Maker
type Maker interface {
	//Create Token creates a new token for a specific username and duration
	CreateToken(id int, duration time.Duration) (string, *Payload, resError.RespError)

	//Verify Token checks if the token is valid or not
	VerifyToken(token string) (*Payload, resError.RespError)
}
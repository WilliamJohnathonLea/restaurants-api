package types

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

const UserAdmin string = "admin"
const UserRestaurant string = "restaurant"
const UserCustomer string = "customer"

type UserClaims struct {
	jwt.RegisteredClaims
	UserType string `json:"type"`
}

func (uc UserClaims) Validate() error {
	switch uc.UserType {
	case UserAdmin, UserCustomer, UserRestaurant:
		return nil
	default:
		return errors.New("user type is invalid")
	}
}

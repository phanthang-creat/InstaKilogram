package model

import (
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        uint32 `orm:"auto" json:"id" form:"id" query:"id"`
	Username  string `orm:"size(50)" json:"username" form:"username" query:"username" validate:"required"`
	Password  string `orm:"size(255)" json:"password" form:"password" query:"password" validate:"required,gte=6"`
	Email     string `orm:"size(50)" json:"email" form:"email" query:"email" validate:"required,email"`
	Phone     string `orm:"size(20)" json:"phone" form:"phone" query:"phone" validate:"required"`
	Avatar    string `orm:"size(255)" json:"avatar" form:"avatar" query:"avatar"`
	Role      int    `json:"role" form:"role" query:"role"`
	Status    int    `json:"status" form:"status" query:"status"`
	CreatedAt int64  `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt int64  `orm:"auto_now;type(datetime)" json:"updated_at"`
	Age       int    `json:"age" form:"age" query:"age"`
	Name      string `json:"name" form:"name" query:"name"`
}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func Sanitize(data string) string {
	data = html.EscapeString(strings.TrimSpace(data))
	return data
}

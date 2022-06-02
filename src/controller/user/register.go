package user

import (

	// "os/user"
	"ytb/src/model"
	myVld "ytb/src/validate"

	"github.com/astaxie/beego/orm"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	//get user info from request body x-www-form-urlencoded
	c.Request().ParseForm()
	o := orm.NewOrm()

	validate := validator.New()

	username := c.Request().FormValue("username")
	password := c.Request().FormValue("password")
	email := c.Request().FormValue("email")
	phone := c.Request().FormValue("phone")

	username = model.Sanitize(username)
	password = model.Sanitize(password)
	email = model.Sanitize(email)
	phone = model.Sanitize(phone)

	//hash password
	hashedPassword, e := model.Hash(password)
	if e != nil {
		return c.JSON(400, "failed to hash password")
	}

	user := model.User{
		Username: username,
		Password: password,
		Email:    email,
		Phone:    phone,
		Avatar:   "default.jpg",
	}
	err := validate.Struct(user)

	user.Password = hashedPassword

	if err != nil {
		// log.Println("err", err)
		return c.JSON(400, "password must be at least 6 characters")
	}

	//validate phone number
	if !myVld.IsPhoneNumber(phone) {
		return c.JSON(400, "phone number is not valid")
	}

	if o.QueryTable("user").Filter("username", username).Exist() {
		return c.JSON(400, "username already exists")
	}

	if o.QueryTable("user").Filter("email", email).Exist() {
		return c.JSON(400, "email already exists")
	}

	//insert user
	_, err = o.Insert(&user)
	if err != nil {
		return c.JSON(400, "failed to insert user")
	}

	return c.JSON(200, "user added successfully")
}

package user

import (
	// "fmt"

	api "ytb/src/controller/api/user"
	myJwt "ytb/src/controller/jwt"
	"ytb/src/model"

	"github.com/astaxie/beego/orm"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserLogin struct {
	Username string `json:"username" validate:"required" form:"username" query:"username"`
	Password string `json:"password" validate:"required" form:"password" query:"password"`
}

func Login(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	// c.Request().ParseForm()
	validate := validator.New()
	o := orm.NewOrm()

	var res = model.LoginRespone{}
	var td = model.TokenDetails{}

	var user UserLogin
	if err := c.Bind(&user); err != nil {
		return err
	}

	if err := validate.Struct(user); err != nil {
		return c.JSON(401, "invalid username or password")
	}

	user.Username = model.Sanitize(user.Username)
	user.Password = model.Sanitize(user.Password)

	err := validate.Struct(user)
	if err != nil {
		return c.JSON(400, err)
	}

	//find user by username
	var maps []orm.Params
	num, err := o.Raw("SELECT * FROM user WHERE username = ?", user.Username).Values(&maps)

	var hashedPassword string

	if err == nil && num > 0 {
		hashedPassword = maps[0]["password"].(string)
	} else {
		return c.JSON(401, "user not found")
	}

	//check password
	err = model.CheckPasswordHash(hashedPassword, user.Password)
	if err != nil {
		return c.JSON(401, "password is not correct")
	}

	td, _ = myJwt.CreateJWT(maps[0]["id"].(string), user.Username, maps[0]["role"].(string))

	u, err := api.GetUserById(maps[0]["id"].(string))
	if err != nil {
		return c.JSON(400, err)
	}

	res.Token = td.AccessToken
	res.RefreshToken = td.RefreshToken
	res.Message = "login success"
	res.User.ID = u[0]["id"].(string)
	res.User.Username = u[0]["username"].(string)
	res.User.Email = u[0]["email"].(string)
	res.User.Avatar = u[0]["avatar"].(string)
	res.User.Age = u[0]["age"].(string)
	res.User.Phone = u[0]["phone"].(string)
	res.User.Name = u[0]["name"].(string)

	if err != nil {

		return c.JSON(401, err)
	}

	return c.JSON(200, res)
}

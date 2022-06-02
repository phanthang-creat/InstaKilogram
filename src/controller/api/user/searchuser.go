package api

import (
	// "github.com/astaxie/beego/orm"

	"github.com/astaxie/beego/orm"
	"github.com/labstack/echo/v4"
)

func APISearchUser(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	var q = c.Param("q")
	// c.Request().ParseForm()
	o := orm.NewOrm()
	var maps []orm.Params

	//select * from user where username like 'q%'
	num, err := o.Raw("SELECT username, avatar, id, name FROM user WHERE username LIKE ?", q+"%").Values(&maps)
	if err != nil {
		return c.JSON(400, err)
	}
	if num == 0 {
		return c.JSON(200, "user not found")
	}
	return c.JSON(200, maps)
}

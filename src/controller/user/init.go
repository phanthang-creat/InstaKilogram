package user

import (
	"ytb/src/model"

	"github.com/astaxie/beego/orm"
	// "github.com/labstack/echo/v4"
)

func init() {

	orm.RegisterModel(new(model.User))
	orm.RegisterModel(new(model.Following))
	orm.RegisterModel(new(model.Post))
	orm.RegisterModel(new(model.PostComment))
	orm.RegisterModel(new(model.PostLike))
	orm.RegisterModel(new(model.PostImage))
}

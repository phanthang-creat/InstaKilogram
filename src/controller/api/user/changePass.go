package api

import (
	"log"
	"ytb/src/model"

	"github.com/astaxie/beego/orm"
	"github.com/labstack/echo/v4"
)

type FormPassword struct {
	OldPassword string `json:"oldPassword" form:"oldPassword"`
	NewPassword string `json:"newPassword" form:"newPassword"`
}

func ChangePass(e echo.Context) error {
	uid := e.Get("userID").(string)
	if e.Get("isOwner") == false {
		return e.JSON(403, "forbidden")
	}
	var form FormPassword
	if err := e.Bind(&form); err != nil {
		return e.JSON(400, err)
	}

	log.Println("form", form)

	if form.OldPassword == "" || form.NewPassword == "" {
		return e.JSON(400, "oldPassword or newPassword is empty")
	}

	if form.OldPassword == form.NewPassword {
		return e.JSON(400, "oldPassword and newPassword is same")
	}

	form.NewPassword = model.Sanitize(form.NewPassword)

	hashedNewPassword, _ := model.Hash(form.NewPassword)
	// hashedOldPassword, err := model.Hash(form.OldPassword)

	// if err != nil {
	// 	return e.JSON(400, err)
	// }

	//fill hashed password in database
	o := orm.NewOrm()

	//get password from database by id
	var maps []orm.Params
	num, err := o.Raw("SELECT * FROM user WHERE id = ?", uid).Values(&maps)

	var hashedPassword string

	if err == nil && num > 0 {
		hashedPassword = maps[0]["password"].(string)
	} else {
		return e.JSON(401, "user not found")
	}

	//check password
	err = model.CheckPasswordHash(hashedPassword, form.OldPassword)
	if err != nil {
		return e.JSON(401, "password is not correct")
	}
	_, err = o.Raw("UPDATE user SET password = ? WHERE username = ?", hashedNewPassword, e.Get("username")).Exec()
	if err != nil {
		return e.JSON(400, err)
	}

	return e.JSON(200, "success")
}

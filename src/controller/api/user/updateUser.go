package api

import (
	// "errors"
	// "log"

	// "github.com/astaxie/beego/orm"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
	"ytb/src/model"

	"github.com/astaxie/beego/orm"
	"github.com/labstack/echo/v4"
)

type FormUpdateUser struct {
	Email  string `json:"email" form:"email"`
	Phone  string `json:"phone" form:"phone"`
	Age    string `json:"age" form:"age"`
	Name   string `json:"name" form:"name"`
	Avatar string `json:"avatar" form:"avatar"`
	// Username string `json:"username" form:"username"`
	// Id string `json:"id" form:"id"`
}

func GetMD5Hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

func APIUpdateUserById(e echo.Context) error {
	uid := e.Param("id")                              //uid is user id
	o := orm.NewOrm()                                 //new orm
	var err error                                     //error
	var user model.User                               //user
	var newAvt string                                 //new avatar
	o.QueryTable("user").Filter("id", uid).One(&user) //get user from database by id
	oldAvatar := user.Avatar                          //old avatar

	if e.Get("isOwner") == true {
		var form FormUpdateUser
		if err = e.Bind(&form); err != nil {
			return e.JSON(400, err)
		}
		if file, err := e.FormFile("avatar"); err == nil {
			src, err := file.Open()
			if err != nil {
				return e.JSON(400, err)
			}
			defer src.Close()

			//hash the file name
			typeImage := strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, "."))-1]
			fileNameHash := GetMD5Hash(file.Filename + time.Now().String())
			newAvt = fileNameHash + "." + typeImage

			//destination file
			dst, err := os.Create("public/image/avt/" + newAvt)
			if err != nil {
				return e.JSON(400, err)
			}
			defer dst.Close()

			//copy file
			if _, err = io.Copy(dst, src); err != nil {
				return e.JSON(400, err)
			}
			if (oldAvatar != "default.png") && (oldAvatar != "") {
				cmd := "rm public/image/avt/" + oldAvatar
				exec.Command("sh", "-c", cmd).Run()
			}
		} else {
			newAvt = oldAvatar
		}
		form.Avatar = newAvt

		//sanitize
		form.Email = model.Sanitize(form.Email)
		form.Phone = model.Sanitize(form.Phone)
		form.Age = model.Sanitize(form.Age)
		form.Name = model.Sanitize(form.Name)
		form.Avatar = model.Sanitize(form.Avatar)
		//update
		_, err = o.Raw("UPDATE user SET email = ?, phone = ?, age = ?, name = ?, avatar = ? WHERE id = ?", form.Email, form.Phone, form.Age, form.Name, newAvt, uid).Exec()
		if err != nil {
			return e.JSON(400, err)
		}
		//get old avatar of user and delete
		return e.JSON(200, form)
	} else {
		return e.JSON(400, "you are not owner")
	}
}

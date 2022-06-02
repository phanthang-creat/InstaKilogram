package api

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/labstack/echo/v4"
)

func FollowHandler(e echo.Context) error {
	uid := e.Get("userID").(string)
	fid := e.Param("id")
	if uid == fid {
		return e.JSON(400, errors.New("you can't follow yourself"))
	}
	o := orm.NewOrm()
	if o.QueryTable("user").Filter("id", fid).Exist() {
		if o.QueryTable("following").Filter("uid", uid).Filter("following_id", fid).Exist() {
			return e.JSON(400, "already followed")
		}
		_, err := o.Raw("INSERT INTO following (uid, following_id) VALUES (?, ?)", uid, fid).Exec()
		if err != nil {
			return e.JSON(400, err)
		} else {
			return e.JSON(200, "success")
		}
	} else {
		return e.JSON(400, "user not found")
	}
}

func UnFollowHandler(e echo.Context) error {
	uid := e.Get("userID").(string)
	fid := e.Param("id")

	o := orm.NewOrm()
	if o.QueryTable("user").Filter("id", fid).Exist() {
		if o.QueryTable("following").Filter("uid", uid).Filter("following_id", fid).Exist() {
			_, err := o.Raw("DELETE FROM following WHERE following_id = ? AND uid = ? ", fid, uid).Exec()
			if err != nil {
				return e.JSON(400, err)
			} else {
				return e.JSON(200, "success")
			}
		} else {
			return e.JSON(400, "not followed")
		}
	} else {
		return e.JSON(400, "user not found")
	}
}

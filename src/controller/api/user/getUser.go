package api

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/labstack/echo/v4"
)

func GetUserById(uid string) ([]orm.Params, error) {
	//get user from database by id
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.Raw("SELECT * FROM user WHERE id = ?", uid).Values(&maps)
	if err != nil {
		log.Println("Get", err)
		return maps, err
	} else if num == 0 {
		return maps, errors.New("user not found")
	}
	delete(maps[0], "password")
	delete(maps[0], "role")
	delete(maps[0], "created_at")
	delete(maps[0], "updated_at")
	delete(maps[0], "status")
	return maps, nil
}

func APIGetUserById(e echo.Context) error {
	e.Response().Header().Set("Access-Control-Allow-Origin", "*")
	id := e.Param("id")
	uid := e.Get("userID").(string)
	maps, err := GetUserById(id)
	if err != nil {
		return e.JSON(400, err)
	}
	if e.Get("isOwner") == true {
		maps[0]["isOwner"] = true
	} else {
		//check if user is following
		o := orm.NewOrm()
		var maps2 []orm.Params
		num, err := o.Raw("SELECT * FROM following WHERE uid = ? AND following_id = ?", uid, id).Values(&maps2)
		if err != nil {
			return e.JSON(400, err)
		}
		if num == 0 {
			maps[0]["isFollowing"] = false
		} else {
			maps[0]["isFollowing"] = true
		}
	}
	return e.JSON(200, maps)
}

func APIGetFollowUserById(e echo.Context) error {
	e.Response().Header().Set("Access-Control-Allow-Origin", "*")

	idMaps, err := GetFlowerByUserId(e)
	if err != nil {
		return e.JSON(400, err)
	}
	var key string
	if e.Get("column") == "following_id" {
		key = "uid"
	} else {
		key = "following_id"
	}

	e.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	e.Response().WriteHeader(http.StatusOK)
	type Response struct {
		ID       string `json:"id"`
		Age      string `json:"age"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}
	var res []Response
	for _, maps := range idMaps {
		uid := maps[key].(string)
		user, err := GetUserById(uid)
		if err != nil {
			return e.JSON(400, err)
		}
		res = append(res, Response{
			ID:       user[0]["id"].(string),
			Age:      user[0]["age"].(string),
			Email:    user[0]["email"].(string),
			Name:     user[0]["name"].(string),
			Phone:    user[0]["phone"].(string),
			Username: user[0]["username"].(string),
			Avatar:   user[0]["avatar"].(string),
		})
	}
	return e.JSON(200, res)
}

func GetFlowerByUserId(e echo.Context) ([]orm.Params, error) {
	// uid := e.Get("userID").(string)
	uid := e.Param("id")
	//get user from database by id
	o := orm.NewOrm()
	var maps []orm.Params

	req := strings.Split(e.Request().URL.Path, "/")[3]

	var column string
	if req == "following" {
		column = "uid"
	} else if req == "follower" {
		column = "following_id"
	} else {
		return maps, errors.New("invalid request")
	}
	e.Set("column", column)
	num, err := o.Raw("SELECT * FROM following WHERE "+column+" = ?", uid).Values(&maps)
	if err != nil {
		return maps, err
	} else if num == 0 {
		return maps, errors.New("user not found")
	}
	return maps, nil
}

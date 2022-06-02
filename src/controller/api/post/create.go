package post

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"ytb/src/model"

	"github.com/astaxie/beego/orm"
	"github.com/labstack/echo/v4"
)

func GetMD5Hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

func CreatePost(c echo.Context) error {

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	o := orm.NewOrm()
	var post model.Post
	if err := c.Bind(&post); err != nil {
		return c.JSON(400, err)
	}

	qImage, err := strconv.Atoi(c.FormValue("quantityImage"))
	if err != nil {
		return c.JSON(400, err)
	}

	//get last id auto_increment of table post
	var lastId int
	o.Raw("SELECT MAX(id) FROM post").QueryRow(&lastId)
	//post.Id = lastId + 1
	post.Id = uint32(lastId + 1)
	post.CreatedAt = time.Now().String()

	log.Println(post, lastId)

	log.Println(lastId)

	for i := 0; i < qImage; i++ {
		str := "image" + strconv.Itoa(i)
		if file, err := c.FormFile(str); err == nil {
			log.Printf(file.Filename)
			src, err := file.Open()
			if err != nil {
				return c.JSON(400, err)
			}
			defer src.Close()

			//hash the file name
			typeImage := strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, "."))-1]
			fileNameHash := GetMD5Hash(file.Filename + time.Now().String())
			newImage := fileNameHash + "." + typeImage

			log.Println(newImage)

			// destination file
			dst, err := os.Create("public/image/post/" + newImage)
			if err != nil {
				return c.JSON(400, err)
			}
			defer dst.Close()

			// //copy file
			if _, err = io.Copy(dst, src); err != nil {
				return c.JSON(400, err)
			}

			//insert image to table image
			image := model.PostImage{
				Image:  newImage,
				PostId: uint32(lastId + 1),
			}
			o.Insert(&image)
		}
	}

	if _, err := o.Insert(&post); err != nil {
		return c.JSON(400, err)
	} else {
		return c.JSON(200, "success")
	}
}

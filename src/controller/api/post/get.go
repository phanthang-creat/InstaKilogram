package post

import (
	"log"
	"strconv"
	"ytb/src/model"

	"github.com/astaxie/beego/orm"
	"github.com/labstack/echo/v4"
)

type ResponsePost struct {
	Id      uint32   `json:"id"`
	Uid     uint32   `json:"uid"`
	Caption string   `json:"caption"`
	Image   []string `json:"image"`
	Created string   `json:"created"`
}

func GetPostByIdPost(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	postId, err := strconv.Atoi(c.Param("id"))
	log.Println(postId)
	if err != nil {
		return c.JSON(400, err)
	}
	o := orm.NewOrm()
	var post model.Post
	o.QueryTable("post").Filter("id", postId).One(&post)
	log.Println(post)
	var postImage []model.PostImage
	o.QueryTable("post_image").Filter("post_id", post.Id).All(&postImage)
	// var postComment []model.PostComment
	// o.QueryTable("post_comment").Filter("post_id", post.Id).All(&postComment)
	// var postLike []model.PostLike
	// o.QueryTable("post_like").Filter("post_id", post.Id).All(&postLike)
	var responsePost ResponsePost
	responsePost.Id = post.Id
	responsePost.Uid = post.Uid
	responsePost.Caption = post.Caption
	responsePost.Created = post.CreatedAt
	for _, image := range postImage {
		responsePost.Image = append(responsePost.Image, image.Image)
	}
	return c.JSON(200, responsePost)
}

func GetPostByUid(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	uid, err := strconv.Atoi(c.Param("uid"))
	log.Println(uid)
	if err != nil {
		return c.JSON(400, err)
	}
	o := orm.NewOrm()
	var post []model.Post
	o.QueryTable("post").Filter("uid", uid).All(&post)
	log.Println(post)
	var responsePost []ResponsePost
	for _, post := range post {
		var postImage []model.PostImage
		o.QueryTable("post_image").Filter("post_id", post.Id).All(&postImage)
		var responsePostResponse ResponsePost
		responsePostResponse.Id = post.Id
		responsePostResponse.Uid = post.Uid
		responsePostResponse.Caption = post.Caption
		responsePostResponse.Created = post.CreatedAt
		for _, image := range postImage {
			responsePostResponse.Image = append(responsePostResponse.Image, image.Image)
		}
		responsePost = append(responsePost, responsePostResponse)
	}
	return c.JSON(200, responsePost)
}

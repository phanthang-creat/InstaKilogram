package post

import (
	"strconv"
	api "ytb/src/controller/api/user"
	"ytb/src/model"

	"github.com/astaxie/beego/orm"
	"github.com/labstack/echo/v4"
)

type CommentResponse struct {
	Comment model.PostComment `json:"comment"`
	User    []orm.Params      `json:"user"`
}

func CreateComment(c echo.Context) error {
	var comment model.PostComment
	if err := c.Bind(&comment); err != nil {
		return c.JSON(400, err)
	}

	postId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(400, err)
	}
	if uint32(postId) != comment.PostId {
		return c.JSON(400, "postId not match")
	}

	uid := c.Get("userID").(string)

	uidInt, err := strconv.Atoi(uid)

	if err != nil {
		return c.JSON(400, err)
	}

	comment.Uid = uint32(uidInt)

	o := orm.NewOrm()
	if _, err := o.Insert(&comment); err != nil {
		return c.JSON(400, err)
	}
	return c.JSON(200, "success")
}

func GetCommentById(c echo.Context) error {
	postId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(400, err)
	}

	o := orm.NewOrm()
	var comments []model.PostComment
	o.QueryTable("post_comment").Filter("post_id", uint32(postId)).All(&comments)

	//get user name and uid of comment
	var commentResponse []CommentResponse
	for i := 0; i < len(comments); i++ {
		var user []orm.Params
		user, _ = api.GetUserById(strconv.Itoa(int(comments[i].Uid)))
		delete(user[0], "age")
		delete(user[0], "email")
		delete(user[0], "name")
		delete(user[0], "phone")
		commentResponse = append(commentResponse, CommentResponse{
			Comment: comments[i],
			User:    user,
		})
	}

	return c.JSON(200, commentResponse)
}

func DeleteCommentById(c echo.Context) error {

	//get post id and comment id from form data
	postId := c.FormValue("postId")
	commentId := c.FormValue("commentId")

	o := orm.NewOrm()
	var comment model.PostComment
	o.QueryTable("post_comment").Filter("post_id", postId).Filter("id", commentId).One(&comment)

	//if not found comment
	if comment.Id == 0 {
		return c.JSON(400, "comment not found")
	}

	//if not owner of comment
	uid := c.Get("userID").(string)
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return c.JSON(400, err)
	}
	if comment.Uid != uint32(uidInt) {
		return c.JSON(400, "you are not the owner of this comment")
	}

	if _, err := o.Delete(&comment); err != nil {
		return c.JSON(400, err)
	}

	return c.JSON(200, "success")
}

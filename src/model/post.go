package model

type Post struct {
	Id        uint32 `orm:"auto" json:"id" form:"id" query:"id"`
	Uid       uint32 `json:"uid" form:"uid" query:"uid" validate:"required"`
	Caption   string `json:"caption" form:"caption" query:"caption" validate:"required"`
	Status    uint8  `json:"status" form:"status" query:"status"`
	CreatedAt string `json:"created_at" form:"created_at" query:"created_at"`
}

type PostComment struct {
	Id        uint32 `orm:"auto" json:"id" form:"id" query:"id"`
	Uid       uint32 `json:"uid" form:"uid" query:"uid" validate:"required"`
	PostId    uint32 `json:"post_id" form:"post_id" query:"post_id" validate:"required"`
	Comment   string `json:"comment" form:"comment" query:"comment" validate:"required"`
	Status    uint8  `json:"status" form:"status" query:"status"`
	CreatedAt string `json:"created_at" form:"created_at" query:"created_at"`
}

type PostLike struct {
	Id     uint32 `orm:"auto" json:"id" form:"id" query:"id"`
	Uid    uint32 `json:"uid" form:"uid" query:"uid" validate:"required"`
	PostId uint32 `json:"post_id" form:"post_id" query:"post_id" validate:"required"`
}

type PostImage struct {
	Id     uint32 `orm:"auto" json:"id" form:"id" query:"id"`
	PostId uint32 `json:"post_id" form:"post_id" query:"post_id" validate:"required"`
	Image  string `json:"image" form:"image" query:"image" validate:"required"`
}

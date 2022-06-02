package model

type Following struct {
	Id          uint32 `orm:"auto" json:"id" form:"id" query:"id"`
	Uid         uint32 `json:"uid" form:"uid" query:"uid" validate:"required"`
	FollowingId uint32 `json:"following_id" form:"following_id" query:"following_id" validate:"required"`
	Status      uint8  `json:"status" form:"status" query:"status"`
}

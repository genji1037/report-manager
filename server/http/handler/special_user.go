package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"report-manager/db"
	"report-manager/model"
	"report-manager/server/http/respond"
)

func CreateSpecialUser(c *gin.Context) {
	var req CreateSpecialUserReq
	if err := c.ShouldBind(&req); err != nil {
		respond.BadRequest(c, 400, err.Error())
		return
	}
	var u db.ExchangeSpecialUser
	u.UID = req.UID
	u.Email = req.Email
	u.Remark = req.Remark
	u.Role = db.ExchangeSpecialUserRoleFina
	if err := u.GetByUID(); err != gorm.ErrRecordNotFound {
		if err == nil {
			respond.BizFailed(c, respond.SpecialUserExists(req.UID))
			return
		}
		respond.InternalError(c, err)
		return
	}

	if err := u.Create(); err != nil {
		respond.InternalError(c, err)
		return
	}
	respond.Success(c, req.SpecialUser)
}

func DeleteSpecialUser(c *gin.Context) {
	var req DeleteSpecialUserReq
	if err := c.ShouldBind(&req); err != nil {
		respond.BadRequest(c, 400, err.Error())
		return
	}
	err := db.ExchangeSpecialUser{}.DeleteByUID(req.UID)
	if err != nil {
		respond.InternalError(c, err)
		return
	}
	respond.Success(c, nil)
}

func UpdateSpecialUser(c *gin.Context) {
	var req UpdateSpecialUserReq
	if err := c.ShouldBind(&req); err != nil {
		respond.BadRequest(c, 400, err.Error())
		return
	}
	u := db.ExchangeSpecialUser{
		UID:    req.UID,
		Email:  req.Email,
		Remark: req.Remark,
	}
	if err := u.UpdateByUID(); err != nil {
		respond.InternalError(c, err)
		return
	}
	respond.Success(c, req.SpecialUser)
}

func ListSpecialUsers(c *gin.Context) {
	us, err := db.ExchangeSpecialUser{}.List(db.ExchangeSpecialUserRoleFina)
	if err != nil {
		respond.InternalError(c, err)
		return
	}
	respond.Success(c, ListSpecialUsersResp{us})
}

func GetSpecialUserReport(c *gin.Context) {
	uid := c.Query("uid")
	date := c.Query("date")
	u := db.ExchangeSpecialUser{UID: uid}
	if err := u.GetByUID(); err != nil {

	}

	rs, err := db.ExchangeSpecialUserReport{}.QueryByDatUID(date, uid)
	if err != nil {
		respond.InternalError(c, err)
		return
	}
	respond.Success(c, GetSpecialUserReportResp{
		SpecialUser: SpecialUser{
			UID:    u.UID,
			Email:  u.Email,
			Remark: u.Remark,
		},
		Reports: rs,
	})
}

type CreateSpecialUserReq struct {
	SpecialUser
}

type DeleteSpecialUserReq struct {
	UID string `form:"uid" json:"uid" binding:"uuid"`
}

type UpdateSpecialUserReq struct {
	SpecialUser
}

type ListSpecialUsersResp struct {
	SpecialUsers []db.SpecialUser `json:"special_users"`
}

type SpecialUser struct {
	UID    string `form:"uid" json:"uid" binding:"uuid"`
	Email  string `form:"email" json:"email"`
	Remark string `form:"remark" json:"remark"`
}

type GetSpecialUserReportResp struct {
	SpecialUser
	Reports []model.ExchangeSpecialUserReport `json:"reports"`
}

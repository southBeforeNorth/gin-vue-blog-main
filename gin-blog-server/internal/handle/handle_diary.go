package handle

import (
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Diary struct{}

// CreateDiaryRequest 创建日记请求
type CreateDiaryRequest struct {
	ID      int      `json:"id"`
	Content string   `json:"content"`
	Status  int      `json:"status" `
	Imgs    []string `json:"imgs"`
	AddTime int64    `json:"add_time"`
}

// 日记查询请求
type QueryDiary struct {
	PageQuery
	Content      string `form:"content"`
	Status       int    `form:"status"`
	IsDelete     *bool  `form:"is_delete"`
	AddTimeStart int64  `form:"add_time_start"`
	AddTimeEnd   int64  `form:"add_time_end"`
}

func (*Diary) SaveOrUpdate(c *gin.Context) {
	var req CreateDiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	auth, userErr := CurrentUserAuth(c)
	if userErr != nil {
		ReturnError(c, g.ErrRequest, userErr)
		return
	}
	diary := model.Diary{
		Model: model.Model{
			ID: req.ID,
		},
		Content: req.Content,
		Status:  req.Status,
		Imgs:    req.Imgs,
		UserId:  auth.UserInfoId,
		AddTime: req.AddTime,
	}

	err := model.SaveOrUpdateDiary(db, &diary)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, diary)
}

// 物理删除
func (*Diary) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, err := model.DeleteDiary(GetDB(c), ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// 软删除
func (*Diary) UpdateSoftDelete(c *gin.Context) {
	var req SoftDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, err := model.UpdateDiarySoftDelete(GetDB(c), req.Ids, req.IsDelete)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// 管理后台查询日志列表
func (*Diary) GetListByAdmin(c *gin.Context) {
	var query QueryDiary
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	db := GetDB(c)
	list, total, err := model.GetDiaryList(db, query.Page, query.Size, query.Content, query.IsDelete, query.Status, query.AddTimeStart, query.AddTimeEnd)
	if err != nil || list == nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, PageResult[model.Diary]{
		Size:  query.Size,
		Page:  query.Page,
		Total: total,
		List:  list,
	})
}

// 前台查询日记，包括列表展示
func (*Diary) GetListByFront(c *gin.Context) {
	slog.Debug("GetListByFront start")
	var query QueryDiary
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	slog.Debug(fmt.Sprintf("GetListByFront page:%d, size:%d", query.Page, query.Size))
	db := GetDB(c)
	//前台固定查询非删除，状态为展示 1
	isDelete := false
	list, total, err := model.GetDiaryList(db, query.Page, query.Size, query.Content, &isDelete, 1, query.AddTimeStart, query.AddTimeEnd)
	if err != nil || list == nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, PageResult[model.Diary]{
		Size:  query.Size,
		Page:  query.Page,
		Total: total,
		List:  list,
	})
}

// 获取日记详细信息
func (*Diary) GetDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	article, err := model.GetDiary(GetDB(c), id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, article)
}

package model

import (
	"log/slog"

	"gorm.io/gorm"
)

// Diary 日记模型
type Diary struct {
	Model
	Content  string    `json:"content"`
	Status   int       `json:"status" `
	IsDelete bool      `json:"is_delete"`
	Imgs     []string  `json:"imgs" gorm:"serializer:json"`
	AddTime  int64     `json:"add_time" `
	UserId   int       `json:"-"` // user_auth_id
	User     *UserAuth `gorm:"foreignkey:UserId" json:"user"`
}

// 日记的详细信息
func GetDiary(db *gorm.DB, id int) (data *Diary, err error) {
	result := db.
		Where(Diary{Model: Model{ID: id}}).
		First(&data)
	return data, result.Error
}

func GetDiaryList(db *gorm.DB, page, size int, content string, isDelete *bool, status int, createStart, createEnd int64) (list []Diary, total int64, err error) {
	slog.Debug("GetDiaryList start")
	db = db.Model(Diary{})
	// 开启调试模式
	//db = db.Debug().Model(Diary{})

	if content != "" {
		db = db.Where("content LIKE ?", "%"+content+"%")
	}
	if isDelete != nil {
		db = db.Where("is_delete", isDelete)
	}
	if status != 0 {
		db = db.Where("status", status)
	}
	if createStart != 0 {
		db = db.Where("add_time >= ?", createStart)
	}
	if createEnd != 0 {
		db = db.Where("add_time <= ?", createEnd)
	}

	db = db.Group("id") // 去重

	result := db.Count(&total).
		Scopes(Paginate(page, size)).
		Order("add_time desc").
		Find(&list)
	return list, total, result.Error
}

// 物理删除文章
func DeleteDiary(db *gorm.DB, ids []int) (int64, error) {
	result := db.Where("id IN ?", ids).Delete(&Diary{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// 软删除日记（修改）
func UpdateDiarySoftDelete(db *gorm.DB, ids []int, isDelete bool) (int64, error) {
	result := db.Model(Diary{}).
		Where("id IN ?", ids).
		Update("is_delete", isDelete)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// 新增/编辑日记, 同时根据 分类名称, 标签名称 维护关联表
func SaveOrUpdateDiary(db *gorm.DB, diary *Diary) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var result *gorm.DB
		// 先 添加/更新 日志, 获取到其 ID
		if diary.ID == 0 {
			result = db.Create(&diary)
		} else {
			result = db.Model(&diary).Where("id", diary.ID).Updates(diary)
		}
		if result.Error != nil {
			return result.Error
		}
		return result.Error
	})
}

package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type SiteVisits struct {
	Model
	UserId    int       `json:"user_id"`
	UUID      string    `json:"uuid"`
	IP        string    `json:"ip"`
	IPSource  string    `json:"ip_source"`
	Browser   string    `json:"browser"`
	OS        string    `json:"os"`
	Device    string    `json:"device"`
	PageURL   string    `json:"page_url"`
	VisitTime time.Time `json:"visit_time"`
}

// 保存用户访问记录
func SaveSiteVisit(db *gorm.DB, siteVisit *SiteVisits) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var result *gorm.DB
		if siteVisit.ID == 0 {
			siteVisit.VisitTime = time.Now()
			result = db.Create(&siteVisit)
		} else {
			return errors.New("has exist")
		}
		if result.Error != nil {
			return result.Error
		}
		return result.Error
	})
}

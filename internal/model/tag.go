package model

import (
	"github.com/jinzhu/gorm"
)

type Tag struct {
	*Model
	Name  string `json:"name" form:"name" binding:"max=100"`
	State uint8  `json:"state" form:"state,default=1" binding:"oneof=0 1"`
}

func (t Tag) TableName() string {
	return "blog_tag"
}

// Count 符合的标签数量
func (t Tag) Count(db *gorm.DB) (int, error) {
	var count int
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// List 分页查询到到标签列表
func (t Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
	var tags []*Tag
	var err error
	if pageOffset >= 0 && pageSize > 0 {
		// Offset：偏移量，用于指定开始返回记录之前要跳过的记录数。
		// Limit：限制检索的记录数。
		//global.DBEngine.Offset()
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	if t.Name != "" {
		// Where：设置筛选条件，接受 map，struct 或 string 作为条件。
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	if err = db.Where("is_del = ?", 0).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

// Create 创建标签
func (t Tag) Create(db *gorm.DB) error {
	return db.Create(&t).Error
}

// Update 更新标签
//func (t Tag) Update(db *gorm.DB) error {
//	return db.Model(&Tag{}).Where("id = ? AND is_del = ?", t.ID, 0).Update(t).Error
//}
func (t Tag) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(t).Where("id = ? AND is_del = ?", t.ID, 0).Updates(values).Error; err != nil {
		return err
	}

	return nil
}

// Delete 删除标签
func (t Tag) Delete(db *gorm.DB) error {
	return db.Where("id = ? AND is_del = ?", t.Model.ID, 0).Delete(&t).Error
}

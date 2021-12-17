package model

import (
	"github.com/small-ek/antgo/crypto/hash"
	"github.com/small-ek/antgo/crypto/uuid"
	"gorm.io/gorm"
	"time"
)

//管理员
type Admin struct {
	Id             int            `json:"id" form:"id" uri:"id" gorm:"primaryKey;autoIncrement;comment:'标识'" `
	Name           string         `json:"name" form:"name" gorm:"comment:'名称'"`
	ProfilePicture string         `json:"profile_picture" form:"profile_picture" gorm:"default:Null;comment:'头像'"`
	Username       string         `json:"username" form:"username" comment:"用户名"`
	Salt           string         `json:"salt,omitempty" form:"salt" gorm:"comment:'盐'" `
	Password       string         `json:"password,omitempty" form:"password" gorm:"comment:'密码'"`
	Status         string         `json:"status" form:"status" gorm:"default:true;comment:'状态:启用=true,禁用=false'"`
	Super          string         `json:"super" form:"super" gorm:"default:true;comment:'超级管理员:是=true,否=false'"`
	RoleId         int            `json:"role_id" form:"role_id" gorm:"comment:'角色标识'"`
	CreatedAt      time.Time      `json:"created_at" form:"created_at" gorm:"comment:'创建时间'"`
	UpdatedAt      time.Time      `json:"updated_at" form:"updated_at" gorm:"comment:'修改时间'"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" form:"deleted_at" gorm:"comment:'删除时间'"`
}

//TableName
func (Admin) TableName() string {
	return "s_admin"

}

//BeforeCreate 在创建之前
func (m *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	m.Salt = uuid.Create().String()
	m.Password = hash.Sha256(m.Salt + m.Password)
	return
}

//BeforeUpdate 在跟新之前
func (m *Admin) BeforeUpdate(tx *gorm.DB) (err error) {
	if m.Password != "" {
		m.Salt = uuid.Create().String()
		m.Password = hash.Sha256(m.Salt + m.Password)
	}
	return
}

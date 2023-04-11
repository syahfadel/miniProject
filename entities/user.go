package entities

import (
	"miniProject/helpers"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	UserName  string   `gorm:"not null" json:"username"`
	Email     string   `gorm:"not null;uniqueIndex" json:"email" form:"email" valid:"required,email"`
	Password  string   `gorm:"not null" json:"password" form:"password" valid:"required,length(6|50)"`
	Following []Follow `gorm:"foreignKey:FollowingID"`
	Followed  []Follow `gorm:"foreignKey:FollowedID"`
	Posts     []Post
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helpers.HashPass(u.Password)
	err = nil
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

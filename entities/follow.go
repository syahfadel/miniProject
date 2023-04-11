package entities

import "time"

type Follow struct {
	FollowingID uint
	FollowedID  uint
	Following   User `gorm:"foreignKey:FollowingID"`
	Followed    User `gorm:"foreignKey:FollowedID"`
	CreatedAt   time.Time
}

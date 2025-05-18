package domain

import "time"

type Friend struct {
	ID           uint    `gorm:"primaryKey"`
	UserID       uint    `gorm:"index;not null"` // кто добавил
	FriendUserID uint    `gorm:"index;not null"` // кого добавил
	Nickname     *string `gorm:"size:64"`
	CreatedAt    time.Time
}

package file

import "time"

type File struct {
	ID         uint      `gorm:"primaryKey"`
	Name       string    `gorm:"size:255;not null"`
	Extension  string    `gorm:"size:255;not null"`
	Path       string    `gorm:"size:255;not null"`
	CreatorID  uint      `gorm:"not null"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	LastUpdate time.Time `gorm:"default:null"`
}

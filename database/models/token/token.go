package token

type UserToken struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"not null;unique"`
	Token  string `gorm:"size:255;not null;unique"`
}

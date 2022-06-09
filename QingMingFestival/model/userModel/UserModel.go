package userModel

type UserInfo struct {
	Id     int64  `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Name   string `gorm:"size:20;not null;unique" json:"name"`
	Avatar string `gorm:"size:50" json:"avatar"`
	Email  string `gorm:"size:50" json:"email"`
	Role   string `gorm:"size:25" json:"role"`
}

type Role struct {
	Id    int64  `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Role  string `gorm:"size:20" json:"role"`
	Token string `gorm:"size:50" json:"token"`
}

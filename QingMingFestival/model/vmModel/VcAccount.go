package vmModel

// VcAccount vc集群账户密码表
// VcAccount todo: 字段非空没有做好，密码加密还没做
type VcAccount struct {
	Id          int64  `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Name        string `gorm:"size:20;not null;unique" json:"name,omitempty"`
	Host        string `gorm:"size:25;not null;unique" json:"host,omitempty"`
	User        string `gorm:"size:50;not null" json:"user,omitempty"`
	Password    string `gorm:"size:25;not null" json:"password,omitempty"`
	Tags        string `gorm:"size:100" json:"tags,omitempty"`
	CreatedTime int64  `gorm:"autoUpdateTime" json:"created_time,omitempty"`
	UpdatedTime int64  `gorm:"autoUpdateTime" json:"updated_time,omitempty"`
}

package vmModel

import "time"

type CloneVmJob struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Name        string    `gorm:"size:60;not null" json:"name,omitempty"`
	VmName      string    `gorm:"size:30;not null" json:"vm_name,omitempty"`
	VcName      string    `gorm:"size:30;not null" json:"vc_name,omitempty"`
	JobData     string    `gorm:"size:1000;not null" json:"job_data,omitempty"`
	Status      string    `gorm:"size:20;not null" json:"status,omitempty"`
	Message     string    `gorm:"size:100" json:"message,omitempty"`
	CreatedTime time.Time `gorm:"autoCreateTime" json:"created_time,omitempty"`
	UpdatedTime time.Time `gorm:"autoUpdateTime" json:"updated_time,omitempty"`
}

package vmOperate

import (
	"QingMingFestival/model/vmModel"
	"QingMingFestival/tools"
	"QingMingFestival/tools/Variable"
	"QingMingFestival/types/vmTypes"
	"gorm.io/gorm"
	"time"
)

type VmOperate struct {
	*gorm.DB
}

func NewVmOperate() *VmOperate {
	return &VmOperate{tools.DbEngine}
}

// InsertCluster 添加vc集群账户
func (vo *VmOperate) InsertCluster(vc *vmModel.VcAccount) error {
	res := vo.Create(vc)
	return res.Error
}

// DeleteCluster 删除vc集群账户
func (vo *VmOperate) DeleteCluster(vc *vmModel.VcAccount) error {
	res := vo.Delete(vc)
	return res.Error
}

// QueryClusters 查询vc集群账户列表
func (vo *VmOperate) QueryClusters() ([]vmModel.VcAccount, error) {
	var vcAccounts []vmModel.VcAccount
	res := vo.Find(&vcAccounts)
	if res.Error != nil {
		return nil, res.Error
	}
	return vcAccounts, nil
}

// QueryClusterById QueryCluster 查询单个vc集群账户信息
func (vo *VmOperate) QueryClusterById(id int64) (vmModel.VcAccount, error) {
	var vcAccount vmModel.VcAccount
	res := vo.Where("id = ?", id).First(&vcAccount)
	if res.Error != nil {
		return vcAccount, res.Error
	}
	return vcAccount, nil
}

func (vo *VmOperate) QueryClusterByName(name string) (vmModel.VcAccount, error) {
	var vcAccount vmModel.VcAccount
	res := vo.Where("name = ?", name).First(&vcAccount)
	if res.Error != nil {
		return vcAccount, res.Error
	}
	return vcAccount, nil
}

func (vo *VmOperate) InsertCloneVmJob(cvj []vmModel.CloneVmJob) error {
	res := vo.CreateInBatches(cvj, len(cvj))
	return res.Error
}

func (vo *VmOperate) UpdateCloneVmJob(cvj vmModel.CloneVmJob) {
	vo.Where("id = ?", cvj.Id).Updates(cvj)
	//cvj.UpdatedTime = time.Now().UTC()
	Variable.CloneVmMessageChan <- vmTypes.ResultCloneVmMessage{
		Id:         cvj.Id,
		Status:     cvj.Status,
		Message:    cvj.Message,
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
}

func (vo *VmOperate) QueryCloneVmJobs() ([]vmModel.CloneVmJob, error) {
	var cvjs []vmModel.CloneVmJob
	res := vo.Find(&cvjs)
	if res.Error != nil {
		return nil, res.Error
	}
	return cvjs, nil
}

func (vo *VmOperate) QueryCloneVmJobsById(ids []int) ([]vmModel.CloneVmJob, error) {
	var cvjs []vmModel.CloneVmJob
	res := vo.Find(&cvjs, ids)
	if res.Error != nil {
		return nil, res.Error
	}
	return cvjs, nil
}

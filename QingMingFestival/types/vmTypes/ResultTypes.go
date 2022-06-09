package vmTypes

import (
	"github.com/vmware/govmomi/vim25/types"
)

type ResultMgmtVcCluster struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	Host     string   `json:"host"`
	User     string   `json:"user"`
	Password string   `json:"password"`
	Tags     []string `json:"tags,omitempty"`
}

type RequestMgmtVcClusterByObj struct {
	Data ResultMgmtVcCluster `json:"data"`
}

type RequestClusterById struct {
	Data struct {
		Id int64 `json:"id"`
	} `json:"data"`
}

type ResultHostInfo struct {
	Name    string `json:"name"`
	Cluster string `json:"cluster"`
	//HostSelf        types.ManagedObjectReference    `json:"host_self"`
	Vms        int             `json:"vms"`
	CPU        HostCpu         `json:"cpu"`
	Memory     HostMemory      `json:"memory"`
	Datastore  []HostDataStore `json:"datastore"`
	Hardware   Hardware        `json:"hardware"`
	Status     HostStatus      `json:"status"`
	Datacenter string          `json:"datacenter"`
}

type VmGuestNetInfo struct {
	Name      string `json:"name"`
	Ip        string `json:"ip"`
	Connected bool   `json:"connected"`
}

type VmStorageInfo struct {
	Name       string `json:"name"`
	TotalSpace int64  `json:"total_space"`
	UsageSpace int64  `json:"usage_space"`
	//Committed   int64  `json:"committed"`
	//Uncommitted int64  `json:"uncommitted"`
	//Unshared    int64  `json:"unshared"`
}

type VmCpuInfo struct {
	NumCPU          int32 `json:"num_cpu"`
	CpuUsagePercent int32 `json:"cpu_usage_percent"`
}

type VmMemoryInfo struct {
	MemoryGB               int32 `json:"memory_gb"`
	HostMemoryUsagePercent int32 `json:"host_memory_usage_percent"`
	VmMemoryUsagePercent   int32 `json:"vm_memory_usage_percent"`
}

type ResultVmInfo struct {
	Name          string                `json:"name"`
	Host          string                `json:"host"`
	Folder        string                `json:"folder"`
	Status        VirtualMachineStatus  `json:"status"`
	CPU           VmCpuInfo             `json:"cpu"`
	Memory        VmMemoryInfo          `json:"memory"`
	Template      bool                  `json:"template"`
	GuestId       string                `json:"guest_id"`
	GuestFullName string                `json:"guest_full_name"`
	Storage       VmStorageInfo         `json:"storage"`
	Disk          []types.GuestDiskInfo `json:"disk"`
	Net           []VmGuestNetInfo      `json:"net"`
}

type ResultCloneVmMetaData struct {
	Vms       []string        `json:"vms"`
	Templates []CloneTemplate `json:"templates"`
	Hosts     []CloneHost     `json:"hosts"`
	Folders   []string        `json:"folders"`
}

type CloneTemplate struct {
	Name    string      `json:"name"`
	CpuNum  int32       `json:"cpu_num"`
	Memory  int32       `json:"memory"`
	Storage []int64     `json:"storage"`
	Network []VmNetwork `json:"network"`
}

type CloneHost struct {
	Name               string          `json:"name"`
	CpuNum             int16           `json:"cpu_num"`
	CPUUsagePercent    int8            `json:"cpu_usage_percent"`
	Memory             int64           `json:"memory"`
	MemoryUsagePercent int8            `json:"memory_usage_percent"`
	Storages           []HostDataStore `json:"storages"`
	Networks           []VmNetwork     `json:"networks"`
}

type ResultCloneVmJob struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	VcName     string `json:"vc_name"`
	VmName     string `json:"vm_name"`
	Template   string `json:"template"`
	Host       string `json:"host"`
	Datastore  string `json:"datastore"`
	Folder     string `json:"folder"`
	CPU        int32  `json:"cpu"`
	Memory     int64  `json:"memory"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type RequestCloneVmJob struct {
	Ids []int `json:"ids"`
}

type ResultCloneVmMessage struct {
	Id         int64  `json:"id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	UpdateTime string `json:"update_time"`
}

// DataCenterInfo 套娃info
type DataCenterInfo struct {
	Name          string                    `json:"name"`
	OverallStatus types.ManagedEntityStatus `json:"overall_status"`
	Cluster       []ClusterInfo             `json:"cluster"`
}

type ClusterInfo struct {
	Name           string                    `json:"name"`
	OverallStatus  types.ManagedEntityStatus `json:"overall_status"`
	NumCpuCores    int64                     `json:"num_cpu_cores"`
	TotalMemory    int64                     `json:"total_memory"`
	NumHosts       int64                     `json:"num_hosts"`
	NumActiveHosts int64                     `json:"num_active_hosts"`
	Host           []HostInfo                `json:"host"`
}

type HostInfo struct {
	Name          string                    `json:"name"`
	HostSelf      string                    `json:"host_self"`
	OverallStatus types.ManagedEntityStatus `json:"overall_status"`
	Vms           int                       `json:"vms"`
	CPU           HostCpu                   `json:"cpu"`
	Memory        HostMemory                `json:"memory"`
	Hardware      Hardware                  `json:"hardware"`
	Runtime       Runtime                   `json:"runtime"`
	DataStore     []HostDataStore           `json:"datastore"`
}

type VmFoldersInfo struct {
	Folder string               `json:"folder"`
	Vms    []VirtualMachineInfo `json:"vms"`
	Child  []VmFoldersInfo      `json:"child"`
}

type VirtualMachineInfo struct {
	Name          string                `json:"name"`
	Host          string                `json:"host"`
	ParentFolder  string                `json:"parent_folder"`
	Status        VirtualMachineStatus  `json:"status"`
	GuestFullName string                `json:"guest_full_name"`
	HostName      string                `json:"host_name"`
	CPU           VirtualMachineCpu     `json:"cpu"`
	Memory        VirtualMachineMemory  `json:"memory"`
	Storage       string                `json:"storage"`
	Disk          string                `json:"disk"`
	Net           string                `json:"net"`
	Runtime       VirtualMachineRuntime `json:"runtime"`
}

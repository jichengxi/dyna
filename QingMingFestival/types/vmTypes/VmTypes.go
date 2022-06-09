package vmTypes

import (
	"github.com/vmware/govmomi/vim25/types"
	"time"
)

type HostSummary struct {
	Name     string                       `json:"name"`
	Host     types.ManagedObjectReference `json:"host"`
	HostSelf types.ManagedObjectReference `json:"host_self"`
	Parent   types.ManagedObjectReference `json:"parent"`
	//OverallStatus types.ManagedEntityStatus      `json:"overall_status"`
	CPU       HostCpu                        `json:"cpu"`
	Memory    HostMemory                     `json:"memory"`
	Vm        []types.ManagedObjectReference `json:"vm"`
	Datastore []types.ManagedObjectReference `json:"datastore"`
	Network   []types.ManagedObjectReference `json:"network"`
	Hardware  Hardware                       `json:"hardware"`
	//Runtime       Runtime                        `json:"runtime"`
	Status HostStatus `json:"status"`
}

type Hardware struct {
	Vendor        string `json:"vendor"`
	Model         string `json:"model"`
	CpuModel      string `json:"cpu_model"`
	NumCpuPkgs    int16  `json:"num_cpu_pkgs"`
	NumCpuCores   int16  `json:"num_cpu_cores"`
	NumCpuThreads int16  `json:"num_cpu_threads"`
	EsxiFullName  string `json:"esxi_full_name"`
	Version       string `json:"version"`
}

type Runtime struct {
	ConnectionState types.HostSystemConnectionState `json:"connection_state"`
	PowerState      types.HostSystemPowerState      `json:"power_state"`
}

type HostStatus struct {
	OverallStatus   types.ManagedEntityStatus       `json:"overall_status"`
	ConnectionState types.HostSystemConnectionState `json:"connection_state"`
	PowerState      types.HostSystemPowerState      `json:"power_state"`
}

type HostCpu struct {
	TotalCPU        int64 `json:"total_cpu"`
	UsedCPU         int64 `json:"used_cpu"`
	FreeCPU         int64 `json:"free_cpu"`
	CPUUsagePercent int8  `json:"cpu_usage_percent"`
}

type HostMemory struct {
	TotalMemory        int64 `json:"total_memory"`
	UsedMemory         int64 `json:"used_memory"`
	FreeMemory         int64 `json:"free_memory"`
	MemoryUsagePercent int8  `json:"memory_usage_percent"`
}

type HostDataStore struct {
	Name                 string `json:"name"`
	TotalSpace           int64  `json:"total_space"`
	FreeSpace            int64  `json:"free_space"`
	DatastoreFreePercent int8   `json:"datastore_free_percent"`
}

type NetworkSummary struct {
	Name        string                       `json:"name"`
	NetworkSelf types.ManagedObjectReference `json:"network_self"`
}

type DataStoreSummary struct {
	Name          string                       `json:"name"`
	Datastore     types.ManagedObjectReference `json:"datastore"`
	DatastoreSelf types.ManagedObjectReference `json:"datastore_self"`
	Capacity      int64                        `json:"capacity"`
	FreeSpace     int64                        `json:"free_space"`
	Test          string                       `json:"test"`
}

type ClusterResourceSummary struct {
	Name           string                         `json:"name"`
	ClusterSelf    types.ManagedObjectReference   `json:"cluster_self"`
	Parent         types.ManagedObjectReference   `json:"parent"`
	OverallStatus  types.ManagedEntityStatus      `json:"overall_status"`
	Host           []types.ManagedObjectReference `json:"host"`
	NumCpuCores    int64                          `json:"num_cpu_cores"`
	TotalMemory    int64                          `json:"total_memory"`
	NumHosts       int64                          `json:"num_hosts"`
	NumActiveHosts int64                          `json:"num_active_hosts"`
}

type FolderSummary struct {
	Name          string                         `json:"name"`
	FolderSelf    types.ManagedObjectReference   `json:"folder_self"`
	Parent        types.ManagedObjectReference   `json:"parent"`
	OverallStatus types.ManagedEntityStatus      `json:"overall_status"`
	ChildEntity   []types.ManagedObjectReference `json:"child_entity"`
}

type DatacenterSummary struct {
	Name           string                       `json:"name"`
	DatacenterSelf types.ManagedObjectReference `json:"datacenter_self"`
	Parent         types.ManagedObjectReference `json:"parent"`
	OverallStatus  types.ManagedEntityStatus    `json:"overall_status"`
	HostFolder     types.ManagedObjectReference `json:"host_folder"`
	VmFolder       types.ManagedObjectReference `json:"vm_folder"`
}

type VirtualMachineSummary struct {
	Name               string                                 `json:"name"`
	Template           bool                                   `json:"template"`
	VirtualMachineSelf types.ManagedObjectReference           `json:"virtual_machine_self"`
	Parent             types.ManagedObjectReference           `json:"parent"`
	Host               types.ManagedObjectReference           `json:"host"`
	Status             VirtualMachineStatus                   `json:"status"`
	GuestFullName      string                                 `json:"guest_full_name"`
	GuestId            string                                 `json:"guest_id"`
	HostName           string                                 `json:"host_name"`
	CPU                VirtualMachineCpu                      `json:"cpu"`
	Memory             VirtualMachineMemory                   `json:"memory"`
	Storage            []types.VirtualMachineUsageOnDatastore `json:"storage"`
	Disk               []types.GuestDiskInfo                  `json:"disk"`
	Net                []types.GuestNicInfo                   `json:"net"`
	Nic                []types.ManagedObjectReference         `json:"nic"`
	Runtime            VirtualMachineRuntime                  `json:"runtime"`
}

type VirtualMachineStatus struct {
	OverallStatus        types.ManagedEntityStatus           `json:"overall_status"`
	PowerState           types.VirtualMachinePowerState      `json:"power_state"`
	ConnectionState      types.VirtualMachineConnectionState `json:"connection_state"`
	GuestState           string                              `json:"guest_state"`
	GuestHeartbeatStatus types.ManagedEntityStatus           `json:"guest_heartbeat_status"`
}

type VirtualMachineCpu struct {
	NumCPU      int32 `json:"num_cpu"`
	MaxCpuUsage int32 `json:"max_cpu_usage"`
	CpuUsage    int32 `json:"cpu_usage"`
}

type VirtualMachineMemory struct {
	MemoryMB         int32 `json:"memory_mb"`
	MaxMemoryUsage   int32 `json:"max_memory_usage"`
	HostMemoryUsage  int32 `json:"host_memory_usage"`
	GuestMemoryUsage int32 `json:"guest_memory_usage"`
}

type VirtualMachineRuntime struct {
	BootTime      *time.Time `json:"boot_time"`
	UptimeSeconds int32      `json:"uptime_seconds"`
}

type DatacenterByVm struct {
	Vms      []VirtualMachineSummary
	Folders  string
	Children []DatacenterByVm
}

type VirtualMachineTemplateSummary struct {
	Name               string                                 `json:"name"`
	VirtualMachineSelf types.ManagedObjectReference           `json:"virtual_machine_self"`
	Parent             types.ManagedObjectReference           `json:"parent"`
	Host               types.ManagedObjectReference           `json:"host"`
	GuestFullName      string                                 `json:"guest_full_name"`
	GuestId            string                                 `json:"guest_id"`
	HostName           string                                 `json:"host_name"`
	CPU                VirtualMachineCpu                      `json:"cpu"`
	Memory             VirtualMachineMemory                   `json:"memory"`
	Storage            []types.VirtualMachineUsageOnDatastore `json:"storage"`
	Disk               []types.GuestDiskInfo                  `json:"disk"`
	Net                []types.GuestNicInfo                   `json:"net"`
	Nic                []types.ManagedObjectReference         `json:"nic"`
}

type VirtualSwitchType string

const (
	VirtualSwitchVssType VirtualSwitchType = "vss"
	VirtualSwitchVdsType VirtualSwitchType = "vds"
)

type CloneVmData struct {
	VmName    string      `json:"vm_name"`
	VcName    string      `json:"vc_name"`
	Template  string      `json:"template"`
	CpuNum    int32       `json:"cpu_num"`
	MemoryMB  int64       `json:"memory_mb"`
	Folder    string      `json:"folder"`
	Host      string      `json:"host"`
	Datastore string      `json:"datastore"`
	Disks     []VmDisk    `json:"disks"`
	Networks  []VmNetwork `json:"networks"`
	Ip        IpConfig    `json:"ip"`
	PowerOn   bool        `json:"power_on"`
}

type IpConfig struct {
	IPAddr  string   `json:"ip_addr"`
	NetMask string   `json:"net_mask"`
	Gateway string   `json:"gateway"`
	Dns     []string `json:"dns"`
}

type VmDisk struct {
	DiskCapacityGB int64 `json:"disk_capacity_gb"`
	DiskThin       bool  `json:"disk_thin"`
}

type VmNetwork struct {
	NetworkName string `json:"network_name"`
	// todo: 这里有问题,想办法限制
	NetworkType VirtualSwitchType `json:"network_type"`
}

type MigrateVMData struct {
	VmName     string `json:"vm_name"`
	TargetHost string `json:"target_host"`
}

type RelocateVMData struct {
	VmName          string `json:"vm_name"`
	TargetHost      string `json:"target_host"`
	TargetDatastore string `json:"target_datastore"`
}

const (
	CloneVmInitFailedStatus = "InitFailed"
	CloneVmReadyStatus      = "Ready"
	CloneVmPreStartStatus   = "PreStart"
	CloneVmStartStatus      = "Start"
	CloneVmFailedStatus     = "Failed"
	CloneVmSuccessStatus    = "Success"
)

type CloneVmMessage struct {
}

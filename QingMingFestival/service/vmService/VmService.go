package vmService

import (
	"QingMingFestival/model/vmModel"
	"QingMingFestival/ormOperate/vmOperate"
	"QingMingFestival/tools/middleware"
	"QingMingFestival/tools/vmUtils"
	"QingMingFestival/types/vmTypes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"strconv"
	"strings"
)

type VmService struct {
}

func NewVmService() *VmService {
	return &VmService{}
}

func (vs *VmService) GetClusters() ([]vmTypes.ResultMgmtVcCluster, error) {
	clusters, err := vmOperate.NewVmOperate().QueryClusters()
	if err != nil {
		return nil, err
	}
	var resultClusters []vmTypes.ResultMgmtVcCluster
	for _, cluster := range clusters {
		resultClusters = append(resultClusters, vmTypes.ResultMgmtVcCluster{
			Id:       cluster.Id,
			Name:     cluster.Name,
			Host:     cluster.Host,
			User:     cluster.User,
			Password: cluster.Password,
			Tags:     strings.Split(cluster.Tags, ","),
		})
	}
	return resultClusters, nil
}

func (vs *VmService) AddCluster(vc vmTypes.ResultMgmtVcCluster) ([]vmTypes.ResultMgmtVcCluster, error) {
	tags := strings.Join(vc.Tags, ",")
	vcAccount := vmModel.VcAccount{
		Name:     vc.Name,
		Host:     vc.Host,
		User:     vc.User,
		Password: vc.Password,
		Tags:     tags,
	}
	vmOp := vmOperate.NewVmOperate()
	err := vmOp.InsertCluster(&vcAccount)
	if err != nil {
		return nil, err
	}
	clusters, err := vmOp.QueryClusters()
	if err != nil {
		return nil, err
	}
	var resultClusters []vmTypes.ResultMgmtVcCluster
	for _, cluster := range clusters {
		clusterTags := strings.Split(cluster.Tags, ",")
		resultClusters = append(resultClusters, vmTypes.ResultMgmtVcCluster{
			Id:       cluster.Id,
			Name:     cluster.Name,
			Host:     cluster.Host,
			User:     cluster.User,
			Password: cluster.Password,
			Tags:     clusterTags,
		})
	}
	return resultClusters, nil

}

func (vs *VmService) DelCluster(id int64) error {
	vcAccount := vmModel.VcAccount{Id: id}
	vmOp := vmOperate.NewVmOperate()
	err := vmOp.DeleteCluster(&vcAccount)
	if err != nil {
		return err
	}
	return nil
}

func (vs *VmService) GetCluster(id int64) (vmModel.VcAccount, error) {
	return vmOperate.NewVmOperate().QueryClusterById(id)
}

func (vs *VmService) LoginCluster(id int64) error {
	oldClient := vmOperate.GetVmClient()
	if oldClient != nil {
		oldClient.LogoutVc()
	}
	vcAccount, err := vmOperate.NewVmOperate().QueryClusterById(id)
	if err != nil {
		middleware.Logf.Error("登录ID查询失败, ", err.Error())
		return err
	}
	middleware.Logf.Info("登录ID查询成功, ", vcAccount.Name)
	_, err = vmOperate.NewVmClient(vcAccount.Host, vcAccount.User, vcAccount.Password)
	if err != nil {
		middleware.Logf.Error("登录失败, ", err.Error())
		return err
	}
	middleware.Logf.Info("登录成功, ", vcAccount.Name)
	return nil
}

func (vs *VmService) GetAllHosts() ([]vmTypes.ResultHostInfo, error) {
	client := vmOperate.GetVmClient()
	if client == nil {
		return nil, errors.New("未登录")
	}
	folderMap, err := client.QueryAllFolder()
	if err != nil {
		middleware.Logf.Error("加载文件夹对象失败.")
		return nil, err
	}
	dataCenterMap, err := client.QueryAllDataCenterByHost()
	if err != nil {
		middleware.Logf.Error("加载数据中心对象失败.")
		return nil, err
	}
	clusterMap, err := client.QueryAllClusterComputeResource()
	if err != nil {
		middleware.Logf.Error("加载集群资源对象失败.")
		return nil, err
	}
	hostsMap, err := client.QueryAllHosts()
	if err != nil {
		middleware.Logf.Error("加载主机对象失败.")
		return nil, err
	}
	dataStoresMap, err := client.QueryAllDatastore()
	if err != nil {
		middleware.Logf.Error("加载数据存储对象失败.")
		return nil, err
	}

	var resultHostsInfo []vmTypes.ResultHostInfo
	for _, hostMap := range hostsMap {
		cluster := clusterMap[hostMap.Parent.Value].Name
		clusterParent := clusterMap[hostMap.Parent.Value].Parent.Value
		datacenter := dataCenterMap[folderMap[clusterParent].Parent.Value].Name
		var dataStores []vmTypes.HostDataStore
		for _, ds := range hostMap.Datastore {
			dsInfo := dataStoresMap[ds.Value]
			dataStores = append(dataStores, vmTypes.HostDataStore{
				Name:       dsInfo.Name,
				TotalSpace: dsInfo.Capacity / 1024 / 1024 / 1024,
				FreeSpace:  dsInfo.FreeSpace / 1024 / 1024 / 1024,
			})
		}

		resultHostsInfo = append(resultHostsInfo, vmTypes.ResultHostInfo{
			Name:    hostMap.Name,
			Cluster: cluster,
			//HostSelf:      hostMap.HostSelf,
			Vms: len(hostMap.Vm),
			CPU: vmTypes.HostCpu{
				TotalCPU: hostMap.CPU.TotalCPU,
				UsedCPU:  hostMap.CPU.UsedCPU,
				FreeCPU:  hostMap.CPU.FreeCPU,
			},
			Memory: vmTypes.HostMemory{
				TotalMemory: hostMap.Memory.TotalMemory / 1024 / 1024 / 1024,
				UsedMemory:  hostMap.Memory.UsedMemory / 1024 / 1024 / 1024,
				FreeMemory:  hostMap.Memory.FreeMemory / 1024 / 1024 / 1024,
			},
			Datastore: dataStores,
			//Hardware:        hostMap.Hardware,
			Status:     hostMap.Status,
			Datacenter: datacenter,
		})
	}

	//var dcIs []vmTypes.DataCenterInfo
	// 第一层循环摘到dc
	//for dcK, dcV := range dataCenterMap {
	//	var dcI vmTypes.DataCenterInfo
	//	dcI.Name = folderMap[dcK].Name
	//	dcI.OverallStatus = dcV.OverallStatus
	//
	//	// 第二层循环cluster
	//	clusterObjList := folderMap[dcV.HostFolder.Value].ChildEntity
	//	for _, clusterObj := range clusterObjList {
	//		cluster := clusterMap[clusterObj.Value]
	//		var cI vmTypes.ClusterInfo
	//		cI.Name = cluster.Name
	//		cI.OverallStatus = cluster.OverallStatus
	//		cI.NumCpuCores = cluster.NumCpuCores
	//		cI.TotalMemory = cluster.TotalMemory
	//		cI.NumHosts = cluster.NumHosts
	//		cI.NumActiveHosts = cluster.NumActiveHosts
	//
	//		// 第三层循环host
	//		hostObjList := cluster.Host
	//		for _, hostObj := range hostObjList {
	//			host := hostMap[hostObj.Value]
	//			var hI vmTypes.HostInfo
	//			hI.Name = host.Name
	//			hI.HostSelf = host.HostSelf.Value
	//			hI.OverallStatus = host.OverallStatus
	//			hI.Vms = len(host.Vm)
	//			hI.Hardware = host.Hardware
	//			hI.Runtime = host.Runtime
	//			hI.CPU = host.CPU
	//			hI.Memory = host.Memory
	//
	//			// 第四层循环datastore
	//			datastoreObjList := host.Datastore
	//			for _, datastoreObj := range datastoreObjList {
	//				datastore := dataStoreMap[datastoreObj.Value]
	//				var dsI vmTypes.DataStoreInfo
	//				dsI.Name = datastore.Name
	//				dsI.Capacity = datastore.Capacity
	//				dsI.FreeSpace = datastore.FreeSpace
	//				hI.DataStore = append(hI.DataStore, dsI)
	//			}
	//
	//			cI.Host = append(cI.Host, hI)
	//		}
	//
	//		dcI.Cluster = append(dcI.Cluster, cI)
	//	}
	//	dcIs = append(dcIs, dcI)
	//}

	//var hostsInfo []vmTypes.HostInfo
	//for _, hostObj := range hostsObjList {
	//	hostsInfo = append(hostsInfo, vmTypes.HostInfo{
	//		Name:          hostObj.Name,
	//		OverallStatus: hostObj.OverallStatus,
	//		Vms:           len(hostObj.Vm),
	//		CPU:           hostObj.CPU,
	//		Memory:        hostObj.Memory,
	//		Hardware:      hostObj.Hardware,
	//		Runtime:       hostObj.Runtime,
	//	})
	//}
	return resultHostsInfo, nil
}

func (vs *VmService) GetHost(host string) (*vmTypes.ResultHostInfo, error) {
	client := vmOperate.GetVmClient()
	if client == nil {
		return nil, errors.New("未登录")
	}
	folderMap, err := client.QueryAllFolder()
	if err != nil {
		middleware.Logf.Error("加载文件夹对象失败.")
		return nil, err
	}
	dataCenterMap, err := client.QueryAllDataCenterByHost()
	if err != nil {
		middleware.Logf.Error("加载数据中心对象失败.")
		return nil, err
	}
	clusterMap, err := client.QueryAllClusterComputeResource()
	if err != nil {
		middleware.Logf.Error("加载集群资源对象失败.")
		return nil, err
	}
	hostsMap, err := client.QueryAllHosts()
	if err != nil {
		middleware.Logf.Error("加载主机对象失败.")
		return nil, err
	}
	dataStoresMap, err := client.QueryAllDatastore()
	if err != nil {
		middleware.Logf.Error("加载数据存储对象失败.")
		return nil, err
	}

	//var resultHostInfo *vmTypes.ResultHostInfo
	for _, hostMap := range hostsMap {
		if hostMap.Name == host {
			cluster := clusterMap[hostMap.Parent.Value].Name
			clusterParent := clusterMap[hostMap.Parent.Value].Parent.Value
			datacenter := dataCenterMap[folderMap[clusterParent].Parent.Value].Name
			var dataStores []vmTypes.HostDataStore
			for _, ds := range hostMap.Datastore {
				dsInfo := dataStoresMap[ds.Value]
				dataRes := vmUtils.PercentCalc(float64(dsInfo.FreeSpace) / float64(dsInfo.Capacity))
				datastoreFreePercent := int8(dataRes * 100)
				dataStores = append(dataStores, vmTypes.HostDataStore{
					Name:                 dsInfo.Name,
					TotalSpace:           dsInfo.Capacity / 1024 / 1024 / 1024,
					FreeSpace:            dsInfo.FreeSpace / 1024 / 1024 / 1024,
					DatastoreFreePercent: datastoreFreePercent,
				})
			}
			cpuRes := vmUtils.PercentCalc(float64(hostMap.CPU.UsedCPU) / float64(hostMap.CPU.TotalCPU))
			cpuUsagePercent := int8(cpuRes * 100)
			memoryRes := vmUtils.PercentCalc(float64(hostMap.Memory.UsedMemory) / float64(hostMap.Memory.TotalMemory))
			memoryUsagePercent := int8(memoryRes * 100)
			return &vmTypes.ResultHostInfo{
				Name:    hostMap.Name,
				Cluster: cluster,
				//HostSelf:      hostMap.HostSelf,
				Vms: len(hostMap.Vm),
				CPU: vmTypes.HostCpu{
					TotalCPU:        hostMap.CPU.TotalCPU,
					UsedCPU:         hostMap.CPU.UsedCPU,
					FreeCPU:         hostMap.CPU.FreeCPU,
					CPUUsagePercent: cpuUsagePercent,
				},
				Memory: vmTypes.HostMemory{
					TotalMemory:        hostMap.Memory.TotalMemory / 1024 / 1024 / 1024,
					UsedMemory:         hostMap.Memory.UsedMemory / 1024 / 1024 / 1024,
					FreeMemory:         hostMap.Memory.FreeMemory / 1024 / 1024 / 1024,
					MemoryUsagePercent: memoryUsagePercent,
				},
				Datastore:  dataStores,
				Hardware:   hostMap.Hardware,
				Status:     hostMap.Status,
				Datacenter: datacenter,
			}, nil
		}
	}
	return nil, errors.New("未找到主机: " + host)
}

func (vs *VmService) GetAllVms() ([]vmTypes.ResultVmInfo, error) {
	client := vmOperate.GetVmClient()
	if client == nil {
		middleware.Logf.Error("vc集群未登录.")
		return nil, errors.New("未登录")
	}

	folderMap, err := client.QueryAllFolder()
	if err != nil {
		middleware.Logf.Errorf("获取文件夹失败, error: %s", err.Error())
		return nil, err
	}

	datastoreMap, err := client.QueryAllDatastore()
	if err != nil {
		middleware.Logf.Errorf("获取存储失败, error: %s", err.Error())
		return nil, err
	}

	hostsMap, err := client.QueryAllHosts()
	if err != nil {
		middleware.Logf.Errorf("获取主机失败, error: %s", err.Error())
		return nil, err
	}

	vmsMap, _, err := client.QueryAllVirtualMachine()
	if err != nil {
		middleware.Logf.Errorf("获取虚拟机失败, error: %s", err.Error())
		return nil, err
	}

	var resultVmsInfo []vmTypes.ResultVmInfo
	for _, vmMap := range vmsMap {
		hostObj := hostsMap[vmMap.Host.Value]
		host := hostObj.Name
		hostMemoryMB := hostObj.Memory.TotalMemory / 1024 / 1024
		folder := folderMap[vmMap.Parent.Value].Name
		//var storagesInfo []vmTypes.VmStorageInfo
		//for _, storage := range vmMap.Storage {
		//	storageName := datastoreMap[storage.Datastore.Value].Name
		//	storagesInfo = append(storagesInfo, vmTypes.VmStorageInfo{
		//		Name:        storageName,
		//		Committed:   storage.Committed / 1024 / 1024 / 1024,
		//		Uncommitted: storage.Uncommitted / 1024 / 1024 / 1024,
		//		Unshared:    storage.Unshared / 1024 / 1024 / 1024,
		//	})
		//}
		storageObj := vmMap.Storage[0]
		storageName := datastoreMap[storageObj.Datastore.Value].Name

		var guestNetInfo []vmTypes.VmGuestNetInfo
		for _, net := range vmMap.Net {
			guestNetInfo = append(guestNetInfo, vmTypes.VmGuestNetInfo{
				Name:      net.Network,
				Ip:        net.IpAddress[0],
				Connected: net.Connected,
			})
		}

		var cpuUsagePercent int32
		if vmMap.CPU.MaxCpuUsage != 0 {
			resCPU := vmUtils.PercentCalc(float64(vmMap.CPU.CpuUsage) / float64(vmMap.CPU.MaxCpuUsage))
			cpuUsagePercent = int32(resCPU * 100)
		}
		var vmMemoryUsagePercent int32
		if vmMap.Memory.MaxMemoryUsage != 0 {
			resVmMem := vmUtils.PercentCalc(float64(vmMap.Memory.GuestMemoryUsage) / float64(vmMap.Memory.MemoryMB))
			vmMemoryUsagePercent = int32(resVmMem * 100)
		}

		var HostMemoryUsagePercent int32
		resHostMem := vmUtils.PercentCalc(float64(vmMap.Memory.HostMemoryUsage) / float64(hostMemoryMB))
		HostMemoryUsagePercent = int32(resHostMem * 100)

		resultVmsInfo = append(resultVmsInfo, vmTypes.ResultVmInfo{
			Name:   vmMap.Name,
			Host:   host,
			Folder: folder,
			Status: vmMap.Status,
			CPU: vmTypes.VmCpuInfo{
				NumCPU:          vmMap.CPU.NumCPU,
				CpuUsagePercent: cpuUsagePercent,
			},
			Memory: vmTypes.VmMemoryInfo{
				MemoryGB:               vmMap.Memory.MemoryMB / 1024,
				VmMemoryUsagePercent:   vmMemoryUsagePercent,
				HostMemoryUsagePercent: HostMemoryUsagePercent,
			},
			Template:      vmMap.Template,
			GuestId:       vmMap.GuestId,
			GuestFullName: vmMap.GuestFullName,
			Storage: vmTypes.VmStorageInfo{
				Name:       storageName,
				TotalSpace: (storageObj.Committed + storageObj.Uncommitted) / 1024 / 1024 / 1024,
				UsageSpace: storageObj.Committed / 1024 / 1024 / 1024,
			},
			//Disk: vmMap.Disk,
			//Net:  guestNetInfo,
		})
	}

	// 层叠式api方案
	//var dcVms []vmTypes.DatacenterByVm
	//for _, dc := range datacenterMap {
	//	var dcVm vmTypes.DatacenterByVm
	//	vmFoldersObjList := folderMap[dc.VmFolder.Value].ChildEntity
	//	fmt.Println("vmFoldersObjList: ", vmFoldersObjList)
	//	dcVm.Folders = folderMap[dc.VmFolder.Value].Name
	//	fmt.Println("dcVm.Folders: ", dcVm.Folders)
	//	vmUtils.ForAllVms(&dcVm, vmFoldersObjList, vmsMap, folderMap)
	//	dcVms = append(dcVms, dcVm)
	//}

	return resultVmsInfo, nil
}

func (vs *VmService) GetVm(vm string) (*vmTypes.ResultVmInfo, error) {
	client := vmOperate.GetVmClient()
	if client == nil {
		middleware.Logf.Error("vc集群未登录.")
		return nil, errors.New("未登录")
	}

	folderMap, err := client.QueryAllFolder()
	if err != nil {
		middleware.Logf.Errorf("获取文件夹失败, error: %s", err.Error())
		return nil, err
	}

	datastoreMap, err := client.QueryAllDatastore()
	if err != nil {
		middleware.Logf.Errorf("获取存储失败, error: %s", err.Error())
		return nil, err
	}

	hostsMap, err := client.QueryAllHosts()
	if err != nil {
		middleware.Logf.Errorf("获取主机失败, error: %s", err.Error())
		return nil, err
	}

	vmsMap, _, err := client.QueryAllVirtualMachine()
	if err != nil {
		middleware.Logf.Errorf("获取虚拟机失败, error: %s", err.Error())
		return nil, err
	}

	for _, vmMap := range vmsMap {
		if vmMap.Name == vm {
			hostObj := hostsMap[vmMap.Host.Value]
			host := hostObj.Name
			hostMemoryMB := hostObj.Memory.TotalMemory / 1024 / 1024
			folder := folderMap[vmMap.Parent.Value].Name
			storageObj := vmMap.Storage[0]
			storageName := datastoreMap[storageObj.Datastore.Value].Name

			var guestNetInfo []vmTypes.VmGuestNetInfo
			for _, net := range vmMap.Net {
				guestNetInfo = append(guestNetInfo, vmTypes.VmGuestNetInfo{
					Name:      net.Network,
					Ip:        net.IpAddress[0],
					Connected: net.Connected,
				})
			}

			var cpuUsagePercent int32
			if vmMap.CPU.MaxCpuUsage != 0 {
				resCPU := vmUtils.PercentCalc(float64(vmMap.CPU.CpuUsage) / float64(vmMap.CPU.MaxCpuUsage))
				cpuUsagePercent = int32(resCPU * 100)
			}
			var vmMemoryUsagePercent int32
			if vmMap.Memory.MaxMemoryUsage != 0 {
				resVmMem := vmUtils.PercentCalc(float64(vmMap.Memory.GuestMemoryUsage) / float64(vmMap.Memory.MemoryMB))
				vmMemoryUsagePercent = int32(resVmMem * 100)
			}

			var HostMemoryUsagePercent int32
			resHostMem := vmUtils.PercentCalc(float64(vmMap.Memory.HostMemoryUsage) / float64(hostMemoryMB))
			HostMemoryUsagePercent = int32(resHostMem * 100)

			return &vmTypes.ResultVmInfo{
				Name:   vmMap.Name,
				Host:   host,
				Folder: folder,
				Status: vmMap.Status,
				CPU: vmTypes.VmCpuInfo{
					NumCPU:          vmMap.CPU.NumCPU,
					CpuUsagePercent: cpuUsagePercent,
				},
				Memory: vmTypes.VmMemoryInfo{
					MemoryGB:               vmMap.Memory.MemoryMB / 1024,
					VmMemoryUsagePercent:   vmMemoryUsagePercent,
					HostMemoryUsagePercent: HostMemoryUsagePercent,
				},
				Template:      vmMap.Template,
				GuestId:       vmMap.GuestId,
				GuestFullName: vmMap.GuestFullName,
				Storage: vmTypes.VmStorageInfo{
					Name:       storageName,
					TotalSpace: (storageObj.Committed + storageObj.Uncommitted) / 1024 / 1024 / 1024,
					UsageSpace: storageObj.Committed / 1024 / 1024 / 1024,
				},
				Disk: vmMap.Disk,
				Net:  guestNetInfo,
			}, nil
		}
	}
	return nil, errors.New("未找到虚拟机: " + vm)
}

func (vs *VmService) GetCloneVmMetaData() (*vmTypes.ResultCloneVmMetaData, error) {
	client := vmOperate.GetVmClient()
	if client == nil {
		middleware.Logf.Error("vc集群未登录.")
		return nil, errors.New("未登录")
	}
	finder := find.NewFinder(client.Client.Client)
	vmsMap, templatesMap, err := client.QueryAllVirtualMachine()
	if err != nil {
		middleware.Logf.Errorf("获取虚拟机与模版失败, error: %s", err.Error())
		return nil, err
	}

	hostsMap, err := client.QueryAllHosts()
	if err != nil {
		middleware.Logf.Errorf("获取宿主机失败, error: %s", err.Error())
		return nil, err
	}

	dataStoresMap, err := client.QueryAllDatastore()
	if err != nil {
		middleware.Logf.Error("获取数据存储对象失败.")
		return nil, err
	}

	folders, err := finder.FolderList(client.Ctx, "*")
	if err != nil {
		middleware.Logf.Error("查找文件夹失败.")
		return nil, err
	}

	netsMap, err := client.QueryAllNetwork()
	if err != nil {
		middleware.Logf.Error("获取网络失败.")
		return nil, err
	}

	var resultCloneVmMetadata vmTypes.ResultCloneVmMetaData

	for _, hostMap := range hostsMap {
		var dataStores []vmTypes.HostDataStore
		for _, ds := range hostMap.Datastore {
			dsInfo := dataStoresMap[ds.Value]
			dataRes := vmUtils.PercentCalc(float64(dsInfo.FreeSpace) / float64(dsInfo.Capacity))
			datastoreFreePercent := int8(dataRes * 100)
			dataStores = append(dataStores, vmTypes.HostDataStore{
				Name:                 dsInfo.Name,
				TotalSpace:           dsInfo.Capacity / 1024 / 1024 / 1024,
				FreeSpace:            dsInfo.FreeSpace / 1024 / 1024 / 1024,
				DatastoreFreePercent: datastoreFreePercent,
			})
		}
		var networks []vmTypes.VmNetwork
		for _, net := range hostMap.Network {
			if net.Type == "Network" {
				networks = append(networks, vmTypes.VmNetwork{
					NetworkType: "vss",
					NetworkName: netsMap[net.Value].Name,
				})
			}
			if net.Type == "DistributedVirtualPortgroup" {
				networks = append(networks, vmTypes.VmNetwork{
					NetworkType: "vds",
					NetworkName: netsMap[net.Value].Name,
				})
			}
		}

		cpuRes := vmUtils.PercentCalc(float64(hostMap.CPU.UsedCPU) / float64(hostMap.CPU.TotalCPU))
		cpuUsagePercent := int8(cpuRes * 100)
		memoryRes := vmUtils.PercentCalc(float64(hostMap.Memory.UsedMemory) / float64(hostMap.Memory.TotalMemory))
		memoryUsagePercent := int8(memoryRes * 100)
		resultCloneVmMetadata.Hosts = append(resultCloneVmMetadata.Hosts, vmTypes.CloneHost{
			Name:               hostMap.Name,
			CpuNum:             hostMap.Hardware.NumCpuThreads,
			CPUUsagePercent:    cpuUsagePercent,
			Memory:             hostMap.Memory.TotalMemory / 1024 / 1024 / 1024,
			MemoryUsagePercent: memoryUsagePercent,
			Storages:           dataStores,
			Networks:           networks,
		})
	}

	for _, folder := range folders {
		if strings.HasPrefix(folder.Reference().Value, "group-v") {
			resultCloneVmMetadata.Folders = append(resultCloneVmMetadata.Folders, folder.Name())
		}
	}

	for _, vmObj := range vmsMap {
		resultCloneVmMetadata.Vms = append(resultCloneVmMetadata.Vms, vmObj.Name)
	}

	for _, templateObj := range templatesMap {
		var storages []int64
		for _, storage := range templateObj.Storage {
			storages = append(storages, (storage.Committed+storage.Uncommitted)/1024/1024/1024)
		}
		var networks []vmTypes.VmNetwork
		for _, nic := range templateObj.Nic {
			if nic.Type == "Network" {
				networks = append(networks, vmTypes.VmNetwork{
					NetworkType: "vss",
					NetworkName: netsMap[nic.Value].Name,
				})
			}
			if nic.Type == "DistributedVirtualPortgroup" {
				networks = append(networks, vmTypes.VmNetwork{
					NetworkType: "vds",
					NetworkName: netsMap[nic.Value].Name,
				})
			}
		}
		resultCloneVmMetadata.Templates = append(resultCloneVmMetadata.Templates, vmTypes.CloneTemplate{
			Name:    templateObj.Name,
			CpuNum:  templateObj.CPU.NumCPU,
			Memory:  templateObj.Memory.MemoryMB / 1024,
			Storage: storages,
			Network: networks,
		})
	}

	return &resultCloneVmMetadata, nil
}

func (vs *VmService) CloneVm(ids []int) error {
	vmOp := vmOperate.NewVmOperate()
	jobs, err := vmOp.QueryCloneVmJobsById(ids)
	if err != nil {
		return err
	}
	vcPool := make(map[string]*vmOperate.VmClient)
	for _, job := range jobs {
		_, ok := vcPool[job.VcName]
		if !ok {
			vc, err := vmOp.QueryClusterByName(job.VcName)
			if err != nil {
				return err
			}
			client, err := vmOperate.NewVmClient(vc.Host, vc.User, vc.Password)
			if err != nil {
				middleware.Logf.Error("登录失败, ", err.Error())
				return err
			}
			vcPool[job.VcName] = client
		}
	}
	for _, job := range jobs {
		go func(job vmModel.CloneVmJob) {
			vmUtils.CloneVm(vcPool[job.VcName], vmOp, job)
		}(job)
	}

	return nil
}

func (vs *VmService) GetCloneVmJobs() ([]vmTypes.ResultCloneVmJob, error) {
	var resultCloneVmJob []vmTypes.ResultCloneVmJob
	jobs, err := vmOperate.NewVmOperate().QueryCloneVmJobs()
	if err != nil {
		return nil, err
	}
	for _, job := range jobs {
		var jobData vmTypes.CloneVmData
		err = json.Unmarshal([]byte(job.JobData), &jobData)
		if err != nil {
			return nil, err
		}
		resultCloneVmJob = append(resultCloneVmJob, vmTypes.ResultCloneVmJob{
			Id:         job.Id,
			Name:       job.Name,
			VcName:     job.VcName,
			VmName:     job.VmName,
			Template:   jobData.Template,
			Host:       jobData.Host,
			Datastore:  jobData.Datastore,
			Folder:     jobData.Folder,
			CPU:        jobData.CpuNum,
			Memory:     jobData.MemoryMB,
			Status:     job.Status,
			Message:    job.Message,
			CreateTime: job.CreatedTime.Format("2006-01-02 15:04:05"),
			UpdateTime: job.UpdatedTime.Format("2006-01-02 15:04:05"),
		})

	}
	return resultCloneVmJob, nil
}

func (vs *VmService) CreateCloneVmJob(cloneVmData []vmTypes.CloneVmData) ([]vmTypes.ResultCloneVmJob, error) {
	vmOp := vmOperate.NewVmOperate()
	var cloneVmJob []vmModel.CloneVmJob
	for _, cloneData := range cloneVmData {
		var message string
		var status string
		jobData, err := json.Marshal(cloneData)
		if err != nil {
			status = vmTypes.CloneVmInitFailedStatus
			message = err.Error()
		}
		status = vmTypes.CloneVmReadyStatus
		cloneVmJob = append(cloneVmJob, vmModel.CloneVmJob{
			Name:    cloneData.VcName + "-" + cloneData.VmName,
			VmName:  cloneData.VmName,
			VcName:  cloneData.VcName,
			JobData: string(jobData),
			Status:  status,
			Message: message,
		})
	}
	err := vmOp.InsertCloneVmJob(cloneVmJob)
	if err != nil {
		return nil, err
	}
	jobs, err := vmOp.QueryCloneVmJobs()
	if err != nil {
		return nil, err
	}
	var resultCloneVmJob []vmTypes.ResultCloneVmJob
	for _, job := range jobs {
		var jobData vmTypes.CloneVmData
		err = json.Unmarshal([]byte(job.JobData), &jobData)
		if err != nil {
			return nil, err
		}
		resultCloneVmJob = append(resultCloneVmJob, vmTypes.ResultCloneVmJob{
			Id:         job.Id,
			Name:       job.Name,
			VcName:     job.VcName,
			VmName:     job.VmName,
			Template:   jobData.Template,
			Host:       jobData.Host,
			Datastore:  jobData.Datastore,
			Folder:     jobData.Folder,
			CPU:        jobData.CpuNum,
			Memory:     jobData.MemoryMB,
			Status:     job.Status,
			Message:    job.Message,
			CreateTime: job.CreatedTime.Format("2006-01-02 15:04:05"),
			UpdateTime: job.UpdatedTime.Format("2006-01-02 15:04:05"),
		})

	}
	return resultCloneVmJob, nil
}

func (vs *VmService) ParseCloneVmXlsx(rows [][]string) ([]vmTypes.CloneVmData, error) {
	xlsxTitle := rows[0]

	var xlsxDataMaps []map[string]string
	for _, row := range rows[1:] {
		xlsxDataMap := make(map[string]string)
		for i, value := range row {
			xlsxDataMap[xlsxTitle[i]] = value
		}
		xlsxDataMaps = append(xlsxDataMaps, xlsxDataMap)

	}
	//fmt.Println("xlsxDataMaps", xlsxDataMaps)
	var cloneVmDatas []vmTypes.CloneVmData
	for _, xlsxDataMap := range xlsxDataMaps {
		cpuNum, err := strconv.ParseInt(xlsxDataMap["cpu_num"], 10, 32)
		if err != nil {
			return nil, err
		}
		memoryMB, err := strconv.ParseInt(xlsxDataMap["memory_mb"], 10, 64)
		if err != nil {
			return nil, err
		}
		var vmDisks []vmTypes.VmDisk
		disks := strings.Split(xlsxDataMap["disk_gb"], ",")
		for _, disk := range disks {
			diskGB, err := strconv.ParseInt(disk, 10, 64)
			if err != nil {
				return nil, err
			}
			vmDisks = append(vmDisks, vmTypes.VmDisk{
				DiskCapacityGB: diskGB,
				DiskThin:       xlsxDataMap["disk_type"] == "1",
			})
		}
		var vmNetworks []vmTypes.VmNetwork
		for i := 0; i < 3; i++ {
			if network, ok := xlsxDataMap["network"+strconv.Itoa(i)]; ok {
				netInfo := strings.Split(network, ",")
				vmNetworks = append(vmNetworks, vmTypes.VmNetwork{
					NetworkName: netInfo[0],
					NetworkType: vmTypes.VirtualSwitchType(netInfo[1]),
				})

			}
		}

		cloneVmDatas = append(cloneVmDatas, vmTypes.CloneVmData{
			VcName:    xlsxDataMap["vc_name"],
			VmName:    xlsxDataMap["vm_name"],
			Template:  xlsxDataMap["template"],
			CpuNum:    int32(cpuNum),
			MemoryMB:  memoryMB,
			Folder:    xlsxDataMap["folder"],
			Host:      xlsxDataMap["host"],
			Datastore: xlsxDataMap["datastore"],
			Disks:     vmDisks,
			Networks:  vmNetworks,
			Ip: vmTypes.IpConfig{
				IPAddr:  xlsxDataMap["ipaddr"],
				NetMask: xlsxDataMap["netmask"],
				Gateway: xlsxDataMap["gateway"],
				Dns:     strings.Split(xlsxDataMap["dns"], ","),
			},
			PowerOn: xlsxDataMap["power_on"] == "1",
		})
	}
	return cloneVmDatas, nil
}

func (vs *VmService) MigrateVM(migrateVMData vmTypes.MigrateVMData) error {
	client := vmOperate.GetVmClient()
	if client == nil {
		return errors.New("还未登录vc, 先登录")
	}
	finder := find.NewFinder(client.Client.Client)
	vms, err := finder.VirtualMachineList(client.Ctx, "*")
	if err != nil {
		return err
	}
	var vmObj *object.VirtualMachine
	for _, vm := range vms {
		if vm.Name() == migrateVMData.VmName {
			vmObj = vm
		}
	}
	fmt.Println("迁移虚拟机名字: ", vmObj.Name())

	// 目标资源池, 这里用默认的就可以
	fmt.Println("获取默认资源池")
	pool, err := finder.DefaultResourcePool(client.Ctx)
	if err != nil {
		return errors.New("获取默认资源池失败")
	}

	// 查询目标主机
	fmt.Println("获取主机", migrateVMData.TargetHost)
	if migrateVMData.TargetHost == "" {
		return errors.New("目标主机不允许为空")
	}
	var host *object.HostSystem
	hosts, err := finder.HostSystemList(client.Ctx, "*")
	if err != nil {
		return errors.New("获取主机列表失败")
	}
	isHost := false
	for _, h := range hosts {
		if h.Name() == migrateVMData.TargetHost {
			host = h
			isHost = true
			break
		}
	}
	if !isHost {
		return errors.New("未找到目标主机")
	}

	state, err := vmObj.PowerState(client.Ctx)
	if err != nil {
		return err
	}
	fmt.Println("迁移虚拟机")
	migrate, err := vmObj.Migrate(client.Ctx, pool, host, types.VirtualMachineMovePriorityDefaultPriority, state)
	if err != nil {
		return err
	}
	result, err := migrate.WaitForResult(client.Ctx)
	if err != nil {
		return err
	}
	fmt.Println(result.State)
	return nil
}

func (vs *VmService) RelocateVM(relocateVMData vmTypes.RelocateVMData) error {
	client := vmOperate.GetVmClient()
	if client == nil {
		return errors.New("还未登录vc, 先登录")
	}
	finder := find.NewFinder(client.Client.Client)
	vms, err := finder.VirtualMachineList(client.Ctx, "*")
	if err != nil {
		return err
	}
	var vmObj *object.VirtualMachine
	for _, vm := range vms {
		if vm.Name() == relocateVMData.VmName {
			vmObj = vm
		}
	}
	middleware.Logf.Infof("迁移虚拟机名字: %s", vmObj.Name())

	// 目标资源池, 这里用默认的就可以
	middleware.Logf.Info("使用默认资源池")
	defaultPool, err := finder.DefaultResourcePool(client.Ctx)
	if err != nil {
		return errors.New("获取默认资源池失败")
	}
	pool := defaultPool.Reference()

	// 查询目标主机
	middleware.Logf.Infof("获取主机 %s", relocateVMData.TargetHost)
	if relocateVMData.TargetHost == "" {
		return errors.New("目标主机不允许为空")
	}
	var host types.ManagedObjectReference
	hosts, err := finder.HostSystemList(client.Ctx, "*")
	if err != nil {
		return errors.New("获取主机列表失败")
	}
	isHost := false
	for _, h := range hosts {
		if h.Name() == relocateVMData.TargetHost {
			host = h.Reference()
			middleware.Logf.Infof("获取主机成功 %s", host.Value)
			isHost = true
			break
		}
	}
	if !isHost {
		return errors.New("未找到目标主机")
	}

	middleware.Logf.Infof("获取存储 %s", relocateVMData.TargetDatastore)
	var datastore types.ManagedObjectReference
	dataStores, err := finder.DatastoreList(client.Ctx, "*")
	if err != nil {
		return errors.New("获取存储列表失败")
	}
	isDatastore := false
	for _, ds := range dataStores {
		if ds.Name() == relocateVMData.TargetDatastore {
			datastore = ds.Reference()
			middleware.Logf.Infof("获取存储成功 %s", datastore.Value)
			isDatastore = true
			break
		}
	}
	if !isDatastore {
		return errors.New("未找到指定存储")
	}

	cm := types.VirtualMachineRelocateSpec{
		Pool:      &pool,
		Host:      &host,
		Datastore: &datastore,
	}
	middleware.Logf.Infof("虚拟机 %s 准备迁移配置", vmObj.Name())
	relocateTask, err := vmObj.Relocate(client.Ctx, cm, types.VirtualMachineMovePriorityDefaultPriority)
	if err != nil {
		middleware.Logf.Errorf("虚拟机 %s 迁移配置失败", vmObj.Name())
		return err
	}
	middleware.Logf.Infof("虚拟机 %s 准备迁移", vmObj.Name())
	result, err := relocateTask.WaitForResult(client.Ctx)
	if err != nil {
		middleware.Logf.Errorf("虚拟机 %s 迁移任务失败", vmObj.Name())
		return err
	}
	middleware.Logf.Infof("虚拟机 %s 迁移任务结束: %s", vmObj.Name(), result.State)
	return nil
}

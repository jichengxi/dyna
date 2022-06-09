package vmUtils

import (
	"QingMingFestival/model/vmModel"
	"QingMingFestival/ormOperate/vmOperate"
	"QingMingFestival/tools/middleware"
	"QingMingFestival/types/vmTypes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"strconv"
	"strings"
)

func ForAllVms(a *vmTypes.DatacenterByVm, vmFoldersObjList []types.ManagedObjectReference, vmMap map[string]vmTypes.VirtualMachineSummary, folderMap map[string]vmTypes.FolderSummary) {
	for _, vmFoldersObj := range vmFoldersObjList {
		if vmFoldersObj.Type == "VirtualMachine" {
			a.Vms = append(a.Vms, vmMap[vmFoldersObj.Value])
		}
		if vmFoldersObj.Type == "Folder" {
			var b vmTypes.DatacenterByVm
			b.Folders = folderMap[vmFoldersObj.Value].Name
			folderVms := folderMap[vmFoldersObj.Value].ChildEntity

			ForAllVms(&b, folderVms, vmMap, folderMap)
			a.Children = append(a.Children, b)
		}
	}
}

func AddDisk(ctx context.Context, vm *object.VirtualMachine, vmDisk vmTypes.VmDisk, ds *types.ManagedObjectReference) error {
	middleware.Logf.Infof("查询 %s 虚拟机设备列表.", vm.Name())
	devices, err := vm.Device(ctx)
	if err != nil {
		middleware.Logf.Errorf("查找虚拟机设备列表失败 %s, %v", vm.Name(), err)
		return err
	}
	middleware.Logf.Infof("查询 %s 虚拟机scsi磁盘控制器.", vm.Name())
	// 这里要看你的磁盘类型，如果你有 nvme，就选择 nvme；否则就选 scsi。当然还有 ide，但是还有人用么
	controller, err := devices.FindDiskController("scsi")
	if err != nil {
		middleware.Logf.Errorf("查找 %s 虚拟机scsi磁盘控制器失败 , %v", vm.Name(), err)
		return err
	}
	device := types.VirtualDisk{
		CapacityInKB: vmDisk.DiskCapacityGB * 1024 * 1024,
		VirtualDevice: types.VirtualDevice{
			Backing: &types.VirtualDiskFlatVer2BackingInfo{
				DiskMode:        string(types.VirtualDiskModePersistent),
				ThinProvisioned: types.NewBool(vmDisk.DiskThin),
				VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
					Datastore: ds,
				},
			},
		},
	}

	devices.AssignController(&device, controller)
	deviceSpec := &types.VirtualDeviceConfigSpec{
		Operation:     types.VirtualDeviceConfigSpecOperationAdd,
		FileOperation: types.VirtualDeviceConfigSpecFileOperationCreate,
		Device:        &device,
	}
	deviceConfigSpec := types.VirtualMachineConfigSpec{}
	deviceConfigSpec.DeviceChange = append(deviceConfigSpec.DeviceChange, deviceSpec)
	taskDisk, err := vm.Reconfigure(ctx, deviceConfigSpec)
	if err != nil {
		middleware.Logf.Error("重置虚拟机磁盘配置失败: ", err)
		return err
	}
	result, err := taskDisk.WaitForResult(ctx)
	if err != nil {
		middleware.Logf.Error("重置虚拟机磁盘配置任务失败: ", err)
		return err
	}
	middleware.Logf.Infof("重置虚拟机 %s 磁盘配置任务结束,状态: %s.", vm.Name(), result.State)
	return nil
}

func SetNetwork(ctx context.Context, edit string, dev *types.VirtualVmxnet3, net object.NetworkReference) (*types.VirtualDeviceConfigSpec, error) {
	netBackInfo, err := net.EthernetCardBackingInfo(ctx)
	if err != nil {
		return nil, err
	}
	if edit == "add" {
		dev.Backing = netBackInfo
		dev.DeviceInfo = &types.Description{}
		dev.WakeOnLanEnabled = types.NewBool(true)
		dev.Connectable = &types.VirtualDeviceConnectInfo{
			StartConnected:    true,
			AllowGuestControl: true,
		}

		deviceSpec := &types.VirtualDeviceConfigSpec{
			Operation: types.VirtualDeviceConfigSpecOperationAdd,
			Device:    dev,
		}
		return deviceSpec, nil
	}
	if edit == "edit" {
		dev.Backing = netBackInfo
		dev.WakeOnLanEnabled = types.NewBool(true)
		dev.Connectable = &types.VirtualDeviceConnectInfo{
			StartConnected:    true,
			AllowGuestControl: true,
		}

		deviceSpec := &types.VirtualDeviceConfigSpec{
			Operation: types.VirtualDeviceConfigSpecOperationEdit,
			Device:    dev,
		}

		return deviceSpec, nil
	}
	return nil, errors.New("edit必须填写, add或edit")

}

func SetIpAddr(cloneVmData vmTypes.CloneVmData) types.CustomizationSpec {
	// 设置IP
	ipAddr := cloneVmData.Ip
	cam := types.CustomizationAdapterMapping{
		Adapter: types.CustomizationIPSettings{
			Ip:            &types.CustomizationFixedIp{IpAddress: ipAddr.IPAddr},
			SubnetMask:    ipAddr.NetMask,
			Gateway:       []string{ipAddr.Gateway},
			DnsDomain:     "local",
			DnsServerList: ipAddr.Dns,
		},
	}
	customSpec := types.CustomizationSpec{
		NicSettingMap: []types.CustomizationAdapterMapping{cam},
		Identity: &types.CustomizationLinuxPrep{
			HostName: &types.CustomizationFixedName{Name: cloneVmData.VmName},
			TimeZone: "Asia/Shanghai1",
			Domain:   "local",
		},
	}
	return customSpec
}

func PowerOn(ctx context.Context, vm *object.VirtualMachine) error {
	taskPowerOn, err := vm.PowerOn(ctx)
	if err != nil {
		middleware.Logf.Errorf("虚拟机 %s 开机配置失败: %s", vm.Name(), err)
		return err
	}
	result, err := taskPowerOn.WaitForResult(ctx)
	if err != nil {
		middleware.Logf.Errorf("虚拟机 %s 开机任务失败: %s", vm.Name(), err)
		return err
	}
	middleware.Logf.Infof("虚拟机 %s 开机.状态: %s", vm.Name(), result.State)
	return nil
}

func PercentCalc(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func CloneVm(client *vmOperate.VmClient, vmOp *vmOperate.VmOperate, cloneVmJob vmModel.CloneVmJob) {
	vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
		Id:      cloneVmJob.Id,
		Status:  vmTypes.CloneVmPreStartStatus,
		Message: "",
	})
	var cloneVmData vmTypes.CloneVmData
	err := json.Unmarshal([]byte(cloneVmJob.JobData), &cloneVmData)
	if err != nil {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmFailedStatus,
			Message: fmt.Sprintf("%s 解析克隆参数失败, error: %s", cloneVmJob.Name, err.Error()),
		})
		middleware.Logf.Errorf("%s 解析克隆参数失败, error: %s", cloneVmJob.Name, err.Error())
		return
	}

	finder := find.NewFinder(client.Client.Client)
	vms, err := finder.VirtualMachineList(client.Ctx, "*")
	if err != nil {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmFailedStatus,
			Message: fmt.Sprintf("%s 获取虚拟机列表失败, error: %s", cloneVmJob.Name, err.Error()),
		})
		middleware.Logf.Errorf("%s 获取虚拟机列表失败, error: %s", cloneVmJob.Name, err.Error())
		return
	}

	var vmTemplate types.ManagedObjectReference
	hasTemplate := false
	for _, vm := range vms {
		if vm.Name() == cloneVmData.VmName {
			vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
				Id:      cloneVmJob.Id,
				Status:  vmTypes.CloneVmFailedStatus,
				Message: fmt.Sprintf("名称为: %s 的虚拟机已存在", cloneVmData.VmName),
			})
			middleware.Logf.Errorf("名称为: %s 的虚拟机已存在", cloneVmData.VmName)
			return
		}
		if vm.Name() == cloneVmData.Template {
			hasTemplate, err = vm.IsTemplate(client.Ctx)
			if err != nil {
				vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
					Id:      cloneVmJob.Id,
					Status:  vmTypes.CloneVmFailedStatus,
					Message: fmt.Sprintf("%s 任务虚拟机模版有误, error: %s", cloneVmJob.Name, err.Error()),
				})
				middleware.Logf.Errorf("%s 任务虚拟机模版有误, error: %s", cloneVmJob.Name, err.Error())
				return
			} else {
				vmTemplate = vm.Reference()
			}
		}
	}
	if !hasTemplate {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmFailedStatus,
			Message: fmt.Sprintf("%s 模版不存在，克隆失败", cloneVmJob.Name),
		})
		middleware.Logf.Errorf("%s 模版不存在，克隆失败", cloneVmJob.Name)
		return
	}

	var configSpecs []types.BaseVirtualDeviceConfigSpec

	// 查询文件夹，如果有值，则使用传入的文件夹，如果没有或者传入文件夹未获取到，则使用默认文件夹
	var folderRef types.ManagedObjectReference
	var newFolder *object.Folder
	hasFolder := false
	if cloneVmData.Folder != "" {
		folders, err := finder.FolderList(client.Ctx, "*")
		if err != nil {
			vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
				Id:      cloneVmJob.Id,
				Status:  vmTypes.CloneVmFailedStatus,
				Message: fmt.Sprintf("%s 获取文件夹列表失败", cloneVmJob.Name),
			})
			middleware.Logf.Errorf("%s 获取文件夹列表失败", cloneVmJob.Name)
			return
		}
		for _, folder := range folders {
			if folder.Name() == cloneVmData.Folder && strings.HasPrefix(folder.Reference().Value, "group-v") {
				folderRef = folder.Reference()
				newFolder = folder
				hasFolder = true
				break
			}
		}
	}
	if !hasFolder || cloneVmData.Folder == "" {
		folder, err := finder.DefaultFolder(client.Ctx)
		if err != nil {
			vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
				Id:      cloneVmJob.Id,
				Status:  vmTypes.CloneVmFailedStatus,
				Message: fmt.Sprintf("%s 获取默认文件夹失败", cloneVmJob.Name),
			})
			middleware.Logf.Errorf("%s 获取默认文件夹失败", cloneVmJob.Name)
			return
		}
		folderRef = folder.Reference()
		newFolder = folder
	}

	// 查询资源池
	var poolRef types.ManagedObjectReference
	pool, err := finder.DefaultResourcePool(client.Ctx)
	if err != nil {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmFailedStatus,
			Message: fmt.Sprintf("%s 获取默认资源池失败", cloneVmJob.Name),
		})
		middleware.Logf.Errorf("%s 获取默认资源池失败", cloneVmJob.Name)
		return
	}
	poolRef = pool.Reference()

	var hostRef types.ManagedObjectReference
	if cloneVmData.Host == "" {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmFailedStatus,
			Message: fmt.Sprintf("%s 主机不允许为空", cloneVmJob.Name),
		})
		middleware.Logf.Errorf("%s 主机不允许为空", cloneVmJob.Name)
		return
	}
	hosts, err := finder.HostSystemList(client.Ctx, "*")
	if err != nil {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmFailedStatus,
			Message: fmt.Sprintf("%s 获取主机列表失败", cloneVmJob.Name),
		})
		middleware.Logf.Errorf("%s 获取主机列表失败", cloneVmJob.Name)
		return
	}
	isHost := false
	for _, h := range hosts {
		if h.Name() == cloneVmData.Host {
			hostRef = h.Reference()
			isHost = true
			break
		}
	}
	if !isHost {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmFailedStatus,
			Message: fmt.Sprintf("%s 未找到指定主机", cloneVmJob.Name),
		})
		middleware.Logf.Errorf("%s 未找到指定主机", cloneVmJob.Name)
		return
	}

	var datastoreRef types.ManagedObjectReference
	dataStores, err := finder.DatastoreList(client.Ctx, "*")
	if err != nil {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmFailedStatus,
			Message: fmt.Sprintf("%s 获取存储列表失败", cloneVmJob.Name),
		})
		middleware.Logf.Errorf("%s 获取存储列表失败", cloneVmJob.Name)
		return
	}
	isDatastore := false
	for _, ds := range dataStores {
		if ds.Name() == cloneVmData.Datastore {
			datastoreRef = ds.Reference()
			isDatastore = true
			break
		}
	}
	if !isDatastore {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmFailedStatus,
			Message: fmt.Sprintf("%s 未找到指定存储", cloneVmJob.Name),
		})
		middleware.Logf.Errorf("%s 未找到指定存储", cloneVmJob.Name)
		return
	}

	// 设置cpu内存
	vmConf := &types.VirtualMachineConfigSpec{
		NumCPUs:  cloneVmData.CpuNum,
		MemoryMB: cloneVmData.MemoryMB,
	}
	relocateSpec := types.VirtualMachineRelocateSpec{
		DeviceChange: configSpecs,
		Folder:       &folderRef,
		Pool:         &poolRef,
		Host:         &hostRef,
		Datastore:    &datastoreRef,
	}

	vmConfig := types.VirtualMachineCloneSpec{
		PowerOn:  false,
		Template: false,
		Config:   vmConf,
		Location: relocateSpec,
	}

	t := object.NewVirtualMachine(client.Client.Client, vmTemplate)
	fmt.Println(t.InventoryPath)

	task, err := t.Clone(client.Ctx, newFolder, cloneVmData.VmName, vmConfig)
	if err != nil {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmFailedStatus,
			Message: fmt.Sprintf("%s 配置克隆任务失败， error: %s", cloneVmJob.Name, err.Error()),
		})
		middleware.Logf.Errorf("%s 配置克隆任务失败， error: %s", cloneVmJob.Name, err.Error())
		return
	}

	middleware.Logf.Infof("%s 克隆任务开始.", cloneVmData.VmName)
	vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
		Id:      cloneVmJob.Id,
		Status:  vmTypes.CloneVmStartStatus,
		Message: fmt.Sprintf("%s 克隆任务开始.", cloneVmJob.Name),
	})
	middleware.Logf.Infof("%s 克隆任务开始.", cloneVmJob.Name)
	result, err := task.WaitForResult(client.Ctx, nil)
	if err != nil {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmFailedStatus,
			Message: fmt.Sprintf("%s 克隆任务失败， error: %s", cloneVmJob.Name, err.Error()),
		})
		middleware.Logf.Errorf("%s 克隆任务失败， error: %s", cloneVmJob.Name, err.Error())
		return
	}
	//vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
	//	Id:      cloneVmJob.Id,
	//	Status:  vmTypes.CloneVmSuccessStatus,
	//	Message: fmt.Sprintf("%s 克隆任务结束, 状态: %s", cloneVmJob.Name, result.State),
	//})
	middleware.Logf.Infof("%s 克隆任务结束, 状态: %s.", cloneVmData.VmName, result.State)

	vmObj := object.NewVirtualMachine(client.Client.Client, result.Result.(types.ManagedObjectReference))

	// 添加磁盘
	if result.State == "success" && len(cloneVmData.Disks) != 0 {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmStartStatus,
			Message: fmt.Sprintf("%s 开始添加磁盘.", cloneVmJob.Name),
		})
		middleware.Logf.Infof("虚拟机 %s 开始添加磁盘.", cloneVmData.VmName)
		for _, disk := range cloneVmData.Disks {
			err = AddDisk(client.Ctx, vmObj, disk, &datastoreRef)
			if err != nil {
				vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
					Id:      cloneVmJob.Id,
					Status:  vmTypes.CloneVmFailedStatus,
					Message: fmt.Sprintf("%s 添加磁盘失败， error: %s", cloneVmJob.Name, err.Error()),
				})
				middleware.Logf.Errorf("%s 添加磁盘失败， error: %s", cloneVmJob.Name, err.Error())
				return
			}
		}
		//vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
		//	Id:      cloneVmJob.Id,
		//	Status:  vmTypes.CloneVmSuccessStatus,
		//	Message: fmt.Sprintf("%s 添加磁盘任务结束.", cloneVmJob.Name),
		//})
		middleware.Logf.Infof("%s 添加磁盘任务结束.", cloneVmData.VmName)
	}

	// 添加网络
	if result.State == "success" && len(cloneVmData.Networks) != 0 {
		vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
			Id:      cloneVmJob.Id,
			Status:  vmTypes.CloneVmStartStatus,
			Message: fmt.Sprintf("%s 添加网络开始.", cloneVmJob.Name),
		})
		middleware.Logf.Info(cloneVmData.VmName + "添加网络开始.")

		netsList, err := finder.NetworkList(client.Ctx, "*")
		if err != nil {
			vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
				Id:      cloneVmJob.Id,
				Status:  vmTypes.CloneVmFailedStatus,
				Message: fmt.Sprintf("%s 获取网络列表失败, error: %s", cloneVmJob.Name, err.Error()),
			})
			return
		}

		for index, network := range cloneVmData.Networks {
			if index == 0 {
				var deviceConfigSpec types.VirtualMachineConfigSpec
				var dev *types.VirtualVmxnet3
				devices, err := vmObj.Device(client.Ctx)
				if err != nil {
					vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
						Id:      cloneVmJob.Id,
						Status:  vmTypes.CloneVmFailedStatus,
						Message: fmt.Sprintf("%s 查找虚拟机设备列表失败, error: %s", cloneVmJob.Name, err.Error()),
					})
					middleware.Logf.Errorf("查找虚拟机设备列表失败 %s, error: %s", vmObj.Name(), err.Error())
					return
				}
				hasDevice := false
				for _, device := range devices {
					switch device.(type) {
					case *types.VirtualVmxnet3:
						dev = device.(*types.VirtualVmxnet3)
						hasDevice = true
					}
					if hasDevice {
						break
					}
				}
				if !hasDevice {
					vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
						Id:      cloneVmJob.Id,
						Status:  vmTypes.CloneVmFailedStatus,
						Message: fmt.Sprintf("%s 未找到网络适配器, error: %s", cloneVmJob.Name, err.Error()),
					})
					middleware.Logf.Errorf("%s 未找到网络适配器, error: %s", vmObj.Name(), err.Error())
					return
				}
				for _, net := range netsList {
					pathArr := strings.Split(net.GetInventoryPath(), "/")
					netName := pathArr[len(pathArr)-1]
					if network.NetworkType == "vss" && net.Reference().Type == "Network" && network.NetworkName == netName {
						middleware.Logf.Infof("虚拟机 %s 更改vss网络: %s", cloneVmData.VmName, netName)
						deviceSpec, err := SetNetwork(client.Ctx, "edit", dev, net)
						if err != nil {
							vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
								Id:      cloneVmJob.Id,
								Status:  vmTypes.CloneVmFailedStatus,
								Message: fmt.Sprintf("%s 配置网络适配器失败, error: %s", cloneVmJob.Name, err.Error()),
							})
							middleware.Logf.Errorf("%s 配置网络适配器失败, error: %s", vmObj.Name(), err.Error())
							return
						}
						deviceConfigSpec.DeviceChange = append(deviceConfigSpec.DeviceChange, deviceSpec)
						break
					}
					if network.NetworkType == "vds" && net.Reference().Type == "DistributedVirtualPortgroup" && network.NetworkName == netName {
						middleware.Logf.Infof("虚拟机 %s 更改vds网络: %s", cloneVmData.VmName, netName)
						deviceSpec, err := SetNetwork(client.Ctx, "edit", dev, net)
						if err != nil {
							vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
								Id:      cloneVmJob.Id,
								Status:  vmTypes.CloneVmFailedStatus,
								Message: fmt.Sprintf("%s 配置网络适配器失败, error: %s", cloneVmJob.Name, err.Error()),
							})
							middleware.Logf.Errorf("%s 配置网络适配器失败, error: %s", vmObj.Name(), err.Error())
							return
						}
						deviceConfigSpec.DeviceChange = append(deviceConfigSpec.DeviceChange, deviceSpec)
						break
					}
				}
				taskNet, err := vmObj.Reconfigure(client.Ctx, deviceConfigSpec)
				if err != nil {
					vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
						Id:      cloneVmJob.Id,
						Status:  vmTypes.CloneVmFailedStatus,
						Message: fmt.Sprintf("%s 重置虚拟机网络配置失败, error: %s", cloneVmJob.Name, err.Error()),
					})
					middleware.Logf.Errorf("%s 重置虚拟机网络配置失败, error: %s", vmObj.Name(), err.Error())
					return
				}

				result, err = taskNet.WaitForResult(client.Ctx)
				if err != nil {
					vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
						Id:      cloneVmJob.Id,
						Status:  vmTypes.CloneVmFailedStatus,
						Message: fmt.Sprintf("%s 重置虚拟机网络配置任务失败, error: %s", cloneVmJob.Name, err.Error()),
					})
					middleware.Logf.Errorf("%s 重置虚拟机网络配置任务失败, error: %s", vmObj.Name(), err.Error())
					return
				}
				//vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
				//	Id:      cloneVmJob.Id,
				//	Status:  vmTypes.CloneVmSuccessStatus,
				//	Message: fmt.Sprintf("%s 重置虚拟机网络配置任务结束.", cloneVmJob.Name),
				//})
				middleware.Logf.Infof("重置虚拟机 %s 网络配置任务结束: %s", cloneVmData.VmName, result.State)

				// 配置IP
				if len(cloneVmData.Ip.IPAddr) != 0 && len(cloneVmData.Ip.NetMask) != 0 && len(cloneVmData.Ip.Gateway) != 0 && len(cloneVmData.Ip.Dns) != 0 {
					vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
						Id:      cloneVmJob.Id,
						Status:  vmTypes.CloneVmStartStatus,
						Message: fmt.Sprintf("%s 虚拟机配置IP任务开始.", cloneVmJob.Name),
					})
					middleware.Logf.Infof("%s 虚拟机配置IP任务开始.", vmObj.Name())
					customSpec := SetIpAddr(cloneVmData)
					taskIp, err := vmObj.Customize(client.Ctx, customSpec)
					if err != nil {
						vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
							Id:      cloneVmJob.Id,
							Status:  vmTypes.CloneVmFailedStatus,
							Message: fmt.Sprintf("%s 虚拟机IP配置失败, error: %s", cloneVmJob.Name, err.Error()),
						})
						middleware.Logf.Errorf("%s 虚拟机IP配置失败, error: %s", vmObj.Name(), err.Error())
						return
					}
					result, err = taskIp.WaitForResult(client.Ctx)
					if err != nil {
						vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
							Id:      cloneVmJob.Id,
							Status:  vmTypes.CloneVmFailedStatus,
							Message: fmt.Sprintf("%s 虚拟机IP配置任务失败, error: %s", cloneVmJob.Name, err.Error()),
						})
						middleware.Logf.Errorf("%s 虚拟机IP配置任务失败, error: %s", vmObj.Name(), err.Error())
						return
					}
					//vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
					//	Id:      cloneVmJob.Id,
					//	Status:  vmTypes.CloneVmSuccessStatus,
					//	Message: fmt.Sprintf("%s 虚拟机配置IP任务结束, 状态: %s", cloneVmJob.Name, result.State),
					//})
					middleware.Logf.Infof("%s 虚拟机配置IP任务结束, 状态: %s", vmObj.Name(), result.State)
				}
			} else {
				var deviceConfigSpec types.VirtualMachineConfigSpec
				var dev types.VirtualVmxnet3
				for _, net := range netsList {
					pathArr := strings.Split(net.GetInventoryPath(), "/")
					netName := pathArr[len(pathArr)-1]
					if network.NetworkType == "vss" && net.Reference().Type == "Network" && network.NetworkName == netName {
						middleware.Logf.Infof("虚拟机 %s 新增vss网络: %s", cloneVmData.VmName, netName)
						deviceSpec, err := SetNetwork(client.Ctx, "add", &dev, net)
						if err != nil {
							vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
								Id:      cloneVmJob.Id,
								Status:  vmTypes.CloneVmFailedStatus,
								Message: fmt.Sprintf("%s 配置网络适配器失败, error: %s", cloneVmJob.Name, err.Error()),
							})
							middleware.Logf.Errorf("%s 配置网络适配器失败, error: %s", vmObj.Name(), err.Error())
							return
						}
						deviceConfigSpec.DeviceChange = append(deviceConfigSpec.DeviceChange, deviceSpec)
						break
					}
					if network.NetworkType == "vds" && net.Reference().Type == "DistributedVirtualPortgroup" && network.NetworkName == netName {
						middleware.Logf.Infof("虚拟机 %s 新增vds网络: %s", cloneVmData.VmName, netName)
						deviceSpec, err := SetNetwork(client.Ctx, "add", &dev, net)
						if err != nil {
							vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
								Id:      cloneVmJob.Id,
								Status:  vmTypes.CloneVmFailedStatus,
								Message: fmt.Sprintf("%s 配置网络适配器失败, error: %s", cloneVmJob.Name, err.Error()),
							})
							middleware.Logf.Errorf("%s 配置网络适配器失败, error: %s", vmObj.Name(), err.Error())
							return
						}
						deviceConfigSpec.DeviceChange = append(deviceConfigSpec.DeviceChange, deviceSpec)
						break
					}
				}
				taskNet, err := vmObj.Reconfigure(client.Ctx, deviceConfigSpec)
				if err != nil {
					vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
						Id:      cloneVmJob.Id,
						Status:  vmTypes.CloneVmFailedStatus,
						Message: fmt.Sprintf("%s 重置虚拟机网络配置失败, error: %s", cloneVmJob.Name, err.Error()),
					})
					middleware.Logf.Errorf("%s 重置虚拟机网络配置失败, error: %s", vmObj.Name(), err.Error())
					return
				}
				result, err = taskNet.WaitForResult(client.Ctx)
				if err != nil {
					vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
						Id:      cloneVmJob.Id,
						Status:  vmTypes.CloneVmFailedStatus,
						Message: fmt.Sprintf("%s 重置虚拟机网络配置任务失败, error: %s", cloneVmJob.Name, err.Error()),
					})
					middleware.Logf.Errorf("%s 重置虚拟机网络配置任务失败, error: %s", vmObj.Name(), err.Error())
					return
				}
				//vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
				//	Id:      cloneVmJob.Id,
				//	Status:  vmTypes.CloneVmSuccessStatus,
				//	Message: fmt.Sprintf("%s 重置虚拟机网络配置任务结束.", cloneVmJob.Name),
				//})
				middleware.Logf.Infof("重置虚拟机 %s 网络配置任务结束: %s", cloneVmData.VmName, result.State)
			}

		}
	}

	// 是否需要开机
	if cloneVmData.PowerOn {
		err = PowerOn(client.Ctx, vmObj)
		if err != nil {
			vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
				Id:      cloneVmJob.Id,
				Status:  vmTypes.CloneVmFailedStatus,
				Message: fmt.Sprintf("%s 虚拟机开机失败, error: %s", cloneVmJob.Name, err.Error()),
			})
			middleware.Logf.Errorf("%s 虚拟机开机失败, error: %s", vmObj.Name(), err.Error())
			return
		}
	}
	vmOp.UpdateCloneVmJob(vmModel.CloneVmJob{
		Id:      cloneVmJob.Id,
		Status:  vmTypes.CloneVmSuccessStatus,
		Message: fmt.Sprintf("%s 克隆虚拟机任务结束.", cloneVmJob.Name),
	})
	middleware.Logf.Infof("%s 克隆虚拟机任务结束: %s", cloneVmData.VmName, result.State)
	return
}

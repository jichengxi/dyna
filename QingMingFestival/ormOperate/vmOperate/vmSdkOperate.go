package vmOperate

import (
	"QingMingFestival/types/vmTypes"
	"context"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"net/url"
)

var _client *VmClient = nil

type VmClient struct {
	Ctx context.Context
	*govmomi.Client
}

func GetVmClient() *VmClient {
	return _client
}

func NewVmClient(host, user, password string) (*VmClient, error) {
	//ctx, cancel := context.WithCancel(context.Background())
	ctx := context.Background()
	//defer cancel()
	userInfo := url.UserPassword(user, password)
	u := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/sdk",
		User:   userInfo,
	}
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		return nil, err
	}
	vmClient := &VmClient{
		Ctx:    ctx,
		Client: client,
	}
	_client = vmClient
	return vmClient, nil
}

func (vc *VmClient) LogoutVc() {
	vc.Logout(vc.Ctx)
}

func (vc *VmClient) getBase(tp string) (*view.ContainerView, error) {
	m := view.NewManager(vc.Client.Client)
	v, err := m.CreateContainerView(vc.Ctx, vc.Client.ServiceContent.RootFolder, []string{tp}, true)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (vc *VmClient) QueryAllHosts() (map[string]vmTypes.HostSummary, error) {
	v, err := vc.getBase("HostSystem")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vc.Ctx)
	var hostSystem []mo.HostSystem
	err = v.Retrieve(vc.Ctx, []string{"HostSystem"}, []string{"parent", "overallStatus", "summary", "datastore", "vm", "network"}, &hostSystem)
	if err != nil {
		return nil, err
	}
	hostsMap := make(map[string]vmTypes.HostSummary)
	for _, hs := range hostSystem {
		totalCPU := int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores)
		freeCPU := totalCPU - int64(hs.Summary.QuickStats.OverallCpuUsage)
		freeMemory := hs.Summary.Hardware.MemorySize - (int64(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)
		hostsMap[hs.Self.Value] = vmTypes.HostSummary{
			Name:     hs.Summary.Config.Name,
			Host:     *hs.Summary.Host,
			HostSelf: hs.Self,
			Parent:   *hs.Parent,
			CPU:      vmTypes.HostCpu{TotalCPU: totalCPU, UsedCPU: int64(hs.Summary.QuickStats.OverallCpuUsage), FreeCPU: freeCPU},
			Memory: vmTypes.HostMemory{TotalMemory: int64(units.ByteSize(hs.Summary.Hardware.MemorySize)),
				UsedMemory: int64(units.ByteSize(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024),
				FreeMemory: freeMemory},
			Vm:        hs.Vm,
			Datastore: hs.Datastore,
			Network:   hs.Network,
			Hardware: vmTypes.Hardware{
				Vendor:        hs.Summary.Hardware.Vendor,
				Model:         hs.Summary.Hardware.Model,
				CpuModel:      hs.Summary.Hardware.CpuModel,
				NumCpuPkgs:    hs.Summary.Hardware.NumCpuPkgs,
				NumCpuCores:   hs.Summary.Hardware.NumCpuCores,
				NumCpuThreads: hs.Summary.Hardware.NumCpuThreads,
				EsxiFullName:  hs.Summary.Config.Product.FullName,
				Version:       hs.Summary.Config.Product.Version,
			},
			Status: vmTypes.HostStatus{
				OverallStatus:   hs.OverallStatus,
				ConnectionState: hs.Summary.Runtime.ConnectionState,
				PowerState:      hs.Summary.Runtime.PowerState,
			},
		}
	}
	return hostsMap, err
}

func (vc *VmClient) QueryAllDatastore() (map[string]vmTypes.DataStoreSummary, error) {
	v, err := vc.getBase("Datastore")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vc.Ctx)
	var dataStore []mo.Datastore
	err = v.Retrieve(vc.Ctx, []string{"Datastore"}, []string{"summary"}, &dataStore)
	if err != nil {
		return nil, err
	}
	dataStoresMap := make(map[string]vmTypes.DataStoreSummary)
	for _, ds := range dataStore {
		dataStoresMap[ds.Self.Value] = vmTypes.DataStoreSummary{
			Name:          ds.Summary.Name,
			Datastore:     *ds.Summary.Datastore,
			DatastoreSelf: ds.Self,
			Capacity:      int64(units.ByteSize(ds.Summary.Capacity)),
			FreeSpace:     int64(units.ByteSize(ds.Summary.FreeSpace)),
			Test:          ds.Name,
		}
	}
	return dataStoresMap, nil
}

func (vc *VmClient) QueryAllClusterComputeResource() (map[string]vmTypes.ClusterResourceSummary, error) {
	v, err := vc.getBase("ClusterComputeResource")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vc.Ctx)
	var cCRs []mo.ClusterComputeResource
	err = v.Retrieve(vc.Ctx, []string{"ClusterComputeResource"}, []string{"name", "parent", "overallStatus", "summary", "host"}, &cCRs)
	if err != nil {
		return nil, err
	}
	cRSs := make(map[string]vmTypes.ClusterResourceSummary)
	for _, ccr := range cCRs {
		cRS := ccr.Summary.GetComputeResourceSummary()
		cRSs[ccr.Self.Value] = vmTypes.ClusterResourceSummary{
			Name:           ccr.Name,
			ClusterSelf:    ccr.Self,
			Parent:         *ccr.Parent,
			OverallStatus:  ccr.OverallStatus,
			Host:           ccr.Host,
			NumCpuCores:    int64(cRS.NumCpuCores),
			TotalMemory:    cRS.TotalMemory,
			NumHosts:       int64(cRS.NumHosts),
			NumActiveHosts: int64(cRS.NumEffectiveHosts),
		}
	}
	return cRSs, nil
}

func (vc *VmClient) QueryAllFolder() (map[string]vmTypes.FolderSummary, error) {
	v, err := vc.getBase("Folder")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vc.Ctx)
	var fds []mo.Folder
	err = v.Retrieve(vc.Ctx, []string{"Folder"}, []string{"name", "parent", "overallStatus", "childEntity"}, &fds)
	if err != nil {
		return nil, err
	}
	crsMap := make(map[string]vmTypes.FolderSummary)
	for _, fd := range fds {
		crsMap[fd.Self.Value] = vmTypes.FolderSummary{
			Name:          fd.Name,
			FolderSelf:    fd.Self,
			Parent:        *fd.Parent,
			OverallStatus: fd.OverallStatus,
			ChildEntity:   fd.ChildEntity,
		}
	}
	return crsMap, nil
}

func (vc *VmClient) QueryAllDataCenterByHost() (map[string]vmTypes.DatacenterSummary, error) {
	v, err := vc.getBase("Datacenter")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vc.Ctx)
	var dcs []mo.Datacenter
	err = v.Retrieve(vc.Ctx, []string{"Datacenter"}, []string{"name", "parent", "hostFolder"}, &dcs) //
	if err != nil {
		return nil, err
	}
	dcsMap := make(map[string]vmTypes.DatacenterSummary)
	for _, dc := range dcs {
		dcsMap[dc.Self.Value] = vmTypes.DatacenterSummary{
			Name:           dc.Name,
			DatacenterSelf: dc.Self,
			Parent:         *dc.Parent,
			OverallStatus:  dc.OverallStatus,
			HostFolder:     dc.HostFolder,
		}
	}
	return dcsMap, nil
}

func (vc *VmClient) QueryAllVirtualMachine() (map[string]vmTypes.VirtualMachineSummary, map[string]vmTypes.VirtualMachineTemplateSummary, error) {
	v, err := vc.getBase("VirtualMachine")
	if err != nil {
		return nil, nil, err
	}
	defer v.Destroy(vc.Ctx)
	var vms []mo.VirtualMachine
	// "name", "parent", "config", "runtime", "guest", "summary", "guestHeartbeatStatus"
	err = v.Retrieve(vc.Ctx, []string{"VirtualMachine"}, []string{"name", "parent", "config", "runtime", "guest", "summary", "guestHeartbeatStatus", "storage", "network"}, &vms)
	if err != nil {
		return nil, nil, err
	}
	vmsMap := make(map[string]vmTypes.VirtualMachineSummary)
	vmsTempMap := make(map[string]vmTypes.VirtualMachineTemplateSummary)
	for _, vm := range vms {
		vmsMap[vm.Self.Value] = vmTypes.VirtualMachineSummary{
			Name:               vm.Name,
			Template:           vm.Summary.Config.Template,
			VirtualMachineSelf: vm.Self,
			Parent:             *vm.Parent,
			Host:               *vm.Runtime.Host,
			Status: vmTypes.VirtualMachineStatus{
				OverallStatus:        vm.Summary.OverallStatus,
				PowerState:           vm.Summary.Runtime.PowerState,
				ConnectionState:      vm.Summary.Runtime.ConnectionState,
				GuestState:           vm.Guest.GuestState,
				GuestHeartbeatStatus: vm.GuestHeartbeatStatus,
			},
			GuestFullName: vm.Config.GuestFullName,
			GuestId:       vm.Config.GuestId,
			HostName:      vm.Guest.HostName,
			CPU: vmTypes.VirtualMachineCpu{
				NumCPU:      vm.Summary.Config.NumCpu,
				MaxCpuUsage: vm.Summary.Runtime.MaxCpuUsage,
				CpuUsage:    vm.Summary.QuickStats.OverallCpuUsage,
			},
			Memory: vmTypes.VirtualMachineMemory{
				MemoryMB:         vm.Summary.Config.MemorySizeMB,
				MaxMemoryUsage:   vm.Summary.Runtime.MaxMemoryUsage,
				HostMemoryUsage:  vm.Summary.QuickStats.HostMemoryUsage,
				GuestMemoryUsage: vm.Summary.QuickStats.GuestMemoryUsage,
			},
			Storage: vm.Storage.PerDatastoreUsage,
			Disk:    vm.Guest.Disk,
			Net:     vm.Guest.Net,
			Nic:     vm.Network,
			Runtime: vmTypes.VirtualMachineRuntime{
				BootTime:      vm.Summary.Runtime.BootTime,
				UptimeSeconds: vm.Summary.QuickStats.UptimeSeconds,
			},
		}
		if vm.Summary.Config.Template {
			vmsTempMap[vm.Self.Value] = vmTypes.VirtualMachineTemplateSummary{
				Name:               vm.Name,
				VirtualMachineSelf: vm.Self,
				Parent:             *vm.Parent,
				Host:               *vm.Runtime.Host,
				GuestFullName:      vm.Config.GuestFullName,
				GuestId:            vm.Config.GuestId,
				HostName:           vm.Guest.HostName,
				CPU: vmTypes.VirtualMachineCpu{
					NumCPU:      vm.Summary.Config.NumCpu,
					MaxCpuUsage: vm.Summary.Runtime.MaxCpuUsage,
					CpuUsage:    vm.Summary.QuickStats.OverallCpuUsage,
				},
				Memory: vmTypes.VirtualMachineMemory{
					MemoryMB:         vm.Summary.Config.MemorySizeMB,
					MaxMemoryUsage:   vm.Summary.Runtime.MaxMemoryUsage,
					HostMemoryUsage:  vm.Summary.QuickStats.HostMemoryUsage,
					GuestMemoryUsage: vm.Summary.QuickStats.GuestMemoryUsage,
				},
				Storage: vm.Storage.PerDatastoreUsage,
				Disk:    vm.Guest.Disk,
				Net:     vm.Guest.Net,
				Nic:     vm.Network,
			}
		}
	}

	return vmsMap, vmsTempMap, nil
}

func (vc *VmClient) QueryAllDataCenterByVm() (map[string]vmTypes.DatacenterSummary, error) {
	v, err := vc.getBase("Datacenter")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vc.Ctx)
	var dcs []mo.Datacenter
	err = v.Retrieve(vc.Ctx, []string{"Datacenter"}, []string{"name", "parent", "overallStatus", "vmFolder"}, &dcs)
	if err != nil {
		return nil, err
	}
	dcsMap := make(map[string]vmTypes.DatacenterSummary)
	for _, dc := range dcs {
		dcsMap[dc.Self.Value] = vmTypes.DatacenterSummary{
			Name:           dc.Name,
			DatacenterSelf: dc.Self,
			Parent:         *dc.Parent,
			OverallStatus:  dc.OverallStatus,
			VmFolder:       dc.VmFolder,
		}
	}
	return dcsMap, nil
}

func (vc *VmClient) QueryAllNetwork() (map[string]vmTypes.NetworkSummary, error) {
	v, err := vc.getBase("Network")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vc.Ctx)
	var nets []mo.Network
	err = v.Retrieve(vc.Ctx, []string{"Network"}, []string{"name"}, &nets)
	if err != nil {
		return nil, err
	}
	netsMap := make(map[string]vmTypes.NetworkSummary)
	for _, net := range nets {
		netsMap[net.Self.Value] = vmTypes.NetworkSummary{
			Name:        net.Name,
			NetworkSelf: net.Self,
		}
	}
	return netsMap, nil
}

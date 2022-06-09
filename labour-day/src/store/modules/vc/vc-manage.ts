import { defineStore } from 'pinia'
import {addVcClustersApi, getVcClustersApi, delVcClustersApi} from "@/api/management/vmware";
import {getVmHostsApi, getVmsApi, loginVcClusterApi} from "@/api/vmware/dashboard"
import {h} from 'vue'
import {NButton} from 'naive-ui'
import { IconCheckBlackCircle } from "@/components/AppIcons"
import {getCloneVmMetadataApi} from "@/api/vmware/clone-vm";
interface vcCluster {
    id: number
    name: string
    ip: string
    user: string
    password: string
    tags: []
}

export interface vmHost {
    name: string
    cluster: string
    vms: number
    cpu: {
        total_cpu: number
        used_cpu: number
        free_cpu: number
        cpu_usage_percent: number
    }
    memory: {
        total_memory: number
        used_memory: number
        free_memory: number
        memory_usage_percent: number
    }
    datastore: {
        name: string
        total_space: number
        free_space: number
        datastore_free_percent: number
    }[]
    hardware: {
        vendor: string
        model: string
        cpu_model: string
        num_cpu_pkgs: number
        num_cpu_cores: number
        num_cpu_threads: number
        esxi_full_name: string
        version: string
    }
    status: {
        overall_status: string
        connection_state: string
        power_state: string
    }
    datacenter: string
}

interface vm {
    name: string
    host: string
    folder: string
    status: {
        overall_status: string
        power_state: string
        connection_state: string
        guest_state: string
        guest_heartbeat_status: string
    }
    cpu: {
        num_cpu: number
        cpu_usage_percent: number
    }
    memory: {
        memory_gb: number
        host_memory_usage_percent: number
        vm_memory_usage_percent: number
    }
    storage: {
        name: string
        total_space: number
        usage_space: number
    }
    template: boolean
    guest_id: string
    guest_full_name: string
    disk?: []
    net?: []
}


export const useVcStore = defineStore('vc-clusters', {
    state() {
        return {
            vcClustersList: <vcCluster[]>[],
            dropDownVcList: <{label: string, key: number}[]>[],
            currentVcCluster: <{label: string, key: number, isLogin: boolean}>{},
            tableLoading: false,

            vmHostsTableData: [],
            vmHostsTableColumns: <{}[]>[],

            vmsTableData: <vm[]>[],
            vmsTableColumns: <{}[]>[],

            vmTemplatesTableData: <vm[]>[],
            vmTemplatesTableColumns: <{}[]>[],
        }
    },
    actions: {
        changeTableLoading() {
            this.tableLoading = !this.tableLoading
        },

        // 获取所有vc集群
        async getVcClusters() {
            try {
                const res = await getVcClustersApi()
                if (res.data.code === 0) {
                    this.vcClustersList = res.data.data
                    this.dropDownVcList = []
                    this.vcClustersList.map((vc) => {
                        const vcPush = {label: vc.name, key: vc.id}
                        this.dropDownVcList.push(vcPush)
                    })
                    console.log("this.vcClustersList",this.vcClustersList)
                    return
                }
                console.warn(res.data.message)
                return
            } catch (error) {
                console.error(error)
                return
            }
        },

        // 添加vc集群
        async addVcCluster(data: {}) {
            try {
                const res = await addVcClustersApi(data)
                if (res.data.code === 0) {
                    this.vcClustersList = res.data.data
                    return
                }
                console.warn(res.data.message)
                return
            } catch (error) {
                console.error(error)
                return
            }
        },

        // 删除vc集群
        async delVcCluster(data: {}) {
            try {
                console.log("data1", data)
                const res = await delVcClustersApi(data)
                if (res.data.code === 0) {
                    return true
                }
                console.warn(res.data.message)
                return false
            } catch (error) {
                console.error(error)
                return false
            }
        },

        // 登录vc集群
        async loginVcCluster(data: {}) {
            try {
                const res = await loginVcClusterApi(data)
                if (res.data.code === 0) {
                    this.currentVcCluster.isLogin = true
                    return true
                }
                console.warn(res.data.message)
                return false
            } catch (error) {
                console.error(error)
                return false
            }
        },

        // 获取所有宿主机
        async getVmHosts() {
            try {
                const res = await getVmHostsApi()
                if (res.data.code === 0) {
                    this.vmHostsTableData = res.data.data
                    return
                }
                console.warn(res.data.message)
                return
            } catch (error) {
                console.error(error)
                return
            }
        },

        // 获取所有虚拟机
        async getVms() {
            try {
                const res = await getVmsApi()
                if (res.data.code === 0) {
                    console.log("res.data.data", res.data.data)
                    // this.vmsTableData = res.data.data
                    const resData = res.data.data
                    resData.map((resVm: vm) => {
                        if (resVm.template) {
                            this.vmTemplatesTableData.push(resVm)
                        } else {
                            this.vmsTableData.push(resVm)
                        }
                    })
                    return
                }
                console.warn(res.data.message)
                return
            } catch (error) {
                console.error(error)
                return
            }
        },

        async getCloneVmMetadata() {
            try {
                const res = await getCloneVmMetadataApi()
                if (res.data.code === 0) {
                    console.log("res.data.data", res.data.data)
                    return res.data.data
                }
                console.warn(res.data.message)
                return {}
            } catch (error) {
                console.error(error)
                return {}
            }
        },

        createVmHostsTableColumns(rows: vmHost[], routeFunc: (host: string) => void) {
            let maxDsCount = 0
            let dsColumns = []
            rows.map((row: vmHost) =>{
                if (row.datastore.length > maxDsCount) {
                    maxDsCount = row.datastore.length
                }
            })
            for (let i = 0; i < maxDsCount; i++) {
                const pushData = {
                    title: "Datastore-" + (i+1).toString(),
                    Key: "datastore_index",
                    align: "center",
                    children: [
                        {
                            title: "name",
                            key: "datastore_name",
                            align: "center",
                            render (row: vmHost) {
                                return h('span', [row.datastore[i].name])
                            }
                        },
                        {
                            title: "Total(GB)",
                            key: "total_space",
                            align: "center",
                            render (row: vmHost) {
                                return h('span', [row.datastore[i].total_space])
                            }
                        },
                        {
                            title: "Free(GB)",
                            key: "free_space",
                            align: "center",
                            render (row: vmHost) {
                                return h('span', [row.datastore[i].free_space])
                            }
                        }
                    ]
                }
                dsColumns.push(pushData)
            }
            this.vmHostsTableColumns = [
                {
                    title: 'Name',
                    key: 'name',
                    align: "center",
                    render (row: vmHost) {
                        return h(
                            NButton,
                            {
                                text: true,
                                size: 'small',
                                onClick: () => routeFunc(row.name),
                            },
                            {default: () => row.name}

                        )
                    }
                },
                {
                    title: 'CPU',
                    key: 'cpu_info',
                    align: "center",
                    children: [
                        {
                            title: 'Total',
                            key: 'total_cpu',
                            align: "center",
                            render (row: vmHost) {
                                return h('span', [row.cpu.total_cpu])
                            }
                        },
                        {
                            title: 'Used',
                            key: 'used_cpu',
                            align: "center",
                            render (row: vmHost) {
                                return h('span', [row.cpu.used_cpu])
                            }
                        },
                        {
                            title: 'Free',
                            key: 'free_cpu',
                            align: "center",
                            render (row: vmHost) {
                                return h('span', [row.cpu.free_cpu])
                            }
                        }
                    ]
                },
                {
                    title: 'Memory',
                    key: 'memory_info',
                    align: "center",
                    children: [
                        {
                            title: 'Total(GB)',
                            key: 'total_memory',
                            align: "center",
                            render (row: vmHost) {
                                return h('span', [row.memory.total_memory])
                            }
                        },
                        {
                            title: 'Used(GB)',
                            key: 'used_memory',
                            align: "center",
                            render (row: vmHost) {
                                return h('span', [row.memory.used_memory])
                            }
                        },
                        {
                            title: 'Free(GB)',
                            key: 'free_memory',
                            align: "center",
                            render (row: vmHost) {
                                return h('span', [row.memory.free_memory])
                            }
                        }
                    ]
                },
                {
                    title: 'DataStore',
                    key: 'datastore',
                    align: "center",
                    children: dsColumns
                },
                {
                    title: 'Vms',
                    key: 'vms',
                    align: "center",
                },
                {
                    title: 'Status',
                    key: 'status',
                    align: "center",
                    children: [
                        {
                            title: 'Overall',
                            key: 'overall_status',
                            align: "center",
                            render (row: vmHost) {
                                return h('span', [row.status.overall_status])
                            }
                        },
                        {
                            title: 'Connection',
                            key: 'connection_state',
                            align: "center",
                            render (row: vmHost) {
                                return h('span', [row.status.connection_state])
                            }
                        },
                        {
                            title: 'Power',
                            key: 'power_state',
                            align: "center",
                            render (row: vmHost) {
                                return h('span', [row.status.power_state])
                            }
                        },
                    ]
                },

            ]
        },

        createVmsTableColumns() {
            this.vmsTableColumns = [
                // {
                //   type: 'expand',
                //   expandable: (rowData) => rowData.name !== 'Jim Green',
                //   renderExpand: (rowData) => {
                //     return `${rowData.name} is a good guy.`
                //   }
                // },
                {
                    title: 'Name',
                    key: 'name',
                    align: "left",
                    width: 150,
                    ellipsis: {
                        tooltip: true
                    }
                },
                {
                    title: 'CPU',
                    key: 'cpu',
                    align: "center",
                    children: [
                        {
                            title: 'Num',
                            key: 'num_cpu',
                            align: "center",
                            render (row: vm) {
                                return h('span', [row.cpu.num_cpu])
                            }
                        },
                        {
                            title: 'Usage(%)',
                            key: 'cpu_usage_percent',
                            align: "center",
                            render (row: vm) {
                                return h('span', [row.cpu.cpu_usage_percent])
                            }
                        },
                    ]
                },
                {
                    title: 'Memory',
                    key: 'memory',
                    align: "center",
                    children: [
                        {
                            title: 'Total(GB)',
                            key: 'memory_gb',
                            align: "center",
                            render (row: vm) {
                                return h('span', [row.memory.memory_gb])
                            }
                        },
                        {
                            title: 'Host Usage(%)',
                            key: 'host_memory_usage_percent',
                            align: "center",
                            width: 120,
                            render (row: vm) {
                                return h('span', [row.memory.host_memory_usage_percent])
                            }
                        },
                        {
                            title: 'Vm Usage(%)',
                            key: 'vm_memory_usage_percent',
                            align: "center",
                            width: 120,
                            render (row: vm) {
                                return h('span', [row.memory.vm_memory_usage_percent])
                            }
                        },
                    ]
                },
                {
                    title: 'Storage',
                    key: 'storage',
                    align: "center",
                    children: [
                        // {
                        //   title: 'Name',
                        //   key: 'storage_name',
                        //   align: "center",
                        //   width: 100,
                        //   ellipsis: {
                        //     tooltip: true
                        //   },
                        //   render (row) {
                        //     return h('span', [row.storage.name])
                        //   }
                        // },
                        {
                            title: 'Total',
                            key: 'total_space',
                            align: "center",
                            render (row: vm) {
                                return h('span', [row.storage.total_space])
                            }
                        },
                        {
                            title: 'Usage',
                            key: 'usage_space',
                            align: "center",
                            render (row: vm) {
                                return h('span', [row.storage.usage_space])
                            }
                        },
                    ]
                },
                {
                    title: 'Status',
                    key: 'status',
                    align: "center",
                    children: [
                        {
                            title: 'Overall',
                            key: 'overall_status',
                            align: "center",
                            render (row: vm) {
                                return h(IconCheckBlackCircle, {color: row.status.overall_status, width: "30", height:"30"})
                            }
                        },
                        {
                            title: 'Power',
                            key: 'power_state',
                            align: "center",
                            render (row: vm) {
                                // IconEmoticonOutline
                                if (row.status.power_state === "poweredOn") {
                                    return h(IconCheckBlackCircle, {color: "green", width: "30", height:"30"})
                                } else {
                                    return h(IconCheckBlackCircle, {color: "red", width: "30", height:"30"})
                                }
                            }
                        },
                        {
                            title: 'Connection',
                            key: 'connection_state',
                            align: "center",
                            render (row: vm) {
                                if (row.status.connection_state === "connected") {
                                    return h(IconCheckBlackCircle, {color: "green", width: "30", height:"30"})
                                } else {
                                    return h(IconCheckBlackCircle, {color: "red", width: "30", height:"30"})
                                }
                            }
                        },
                        {
                            title: 'Guest',
                            key: 'guest_state',
                            align: "center",
                            render (row: vm) {
                                if (row.status.guest_state === "running") {
                                    return h(IconCheckBlackCircle, {color: "green", width: "30", height:"30"})
                                } else {
                                    return h(IconCheckBlackCircle, {color: "red", width: "30", height:"30"})
                                }
                                // return h('span', [row.status.guest_state])
                            }
                        },
                        {
                            title: 'Heartbeat',
                            key: 'guest_heartbeat_status',
                            align: "center",
                            render (row: vm) {
                                // if (row.status.guest_heartbeat_status === "running") {
                                //   return h(IconCheckBlackCircle, {color: "green", width: "30", height:"30"})
                                // } else {
                                //   return h(IconCheckBlackCircle, {color: "red", width: "30", height:"30"})
                                // }
                                return h(IconCheckBlackCircle, {color: row.status.guest_heartbeat_status, width: "30", height:"30"})
                                // return h('span', [row.status.guest_heartbeat_status])
                            }
                        },
                    ],
                },
                {
                    title: 'Host',
                    key: 'host',
                    align: "center",
                    width: 100,
                },
                {
                    title: 'Folder',
                    key: 'folder',
                    align: "center",
                },
                {
                    title: 'System',
                    key: 'guest_full_name',
                    align: "center",
                    width: 100,
                    ellipsis: {
                        tooltip: true
                    }
                },
                {
                    title: 'Action',
                    key: 'actions',
                    align: "center",
                    render (row: vm) {
                        return h(
                            NButton,
                            {
                                size: 'small',
                            },
                            { default: () => 'Send Email' }
                        )
                    }
                }
            ]
        },

        createVmTemplateTableColumns(){
            this.vmTemplatesTableColumns = [
                {
                    title: 'Name',
                    key: 'name',
                    align: "left",
                    // width: 150,
                    // ellipsis: {
                    //     tooltip: true
                    // }
                },
                {
                    title: 'Host',
                    key: 'host',
                    align: "center",
                },
                {
                    title: 'Folder',
                    key: 'folder',
                    align: "center",
                },
                {
                    title: 'Cpu Num',
                    key: 'cpu_num',
                    align: "center",
                    render (row: vm) {
                        return h('span', [row.cpu.num_cpu])
                    }
                },
                {
                    title: 'Memory(GB)',
                    key: 'memory_gb',
                    align: "center",
                    render (row: vm) {
                        return h('span', [row.memory.memory_gb])
                    }
                },
                {
                    title: 'System',
                    key: 'guest_full_name',
                    align: "center"
                },
                {
                    title: 'Storage',
                    key: 'total_space',
                    align: "center",
                    render (row: vm) {
                        return h('span', [row.storage.total_space])
                    }
                }
            ]
        }
    }
})



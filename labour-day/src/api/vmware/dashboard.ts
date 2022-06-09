import {request} from "@/utils/axios";

export const loginVcClusterApi = (data: any) => request({url: '/vmware/cluster/login', method: 'post', data: data})

export const getVmHostsApi = () => request({url: '/vmware/cluster/hosts', method: 'get'})

export const getVmHostApi = (params: {}) => request({url: '/vmware/cluster/host', method:'get', params: params})

export const getVmsApi = () => request({url: '/vmware/cluster/vms', method: 'get'})


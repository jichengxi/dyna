import { request } from "@/utils/axios"

export const getVcClustersApi = () => request({url: '/vmware/clusters', method: 'get'})
export const addVcClustersApi = (data: any) => request({url: '/vmware/clusters/add', method: 'post', data: data})
export const delVcClustersApi = (data: any) => request({url: '/vmware/clusters/del', method: 'post', data: data})

export const login = (data: {}) => request({url: '/auth/login', method: 'post', data: data})

import {request} from "@/utils/axios"

export const getCloneVmMetadataApi = () => request({url: '/vmware/clone-vm/metadata', method: 'get'})

export const getCloneVmJobsApi = () => request({url: '/vmware/clone-vm/jobs', method: 'get'})

export const postCloneVmJobsApi = (data: {}) => request({url: '/vmware/clone-vm/jobs', method: 'post', data: [data]})

export const getCloneVmTemplateApi = () => request({url: '/vmware/clone-vm/template', method: 'get', headers: { 'Content-Type': 'application/x-download' }, responseType: 'blob'})

export const postCloneVmRunApi = (data: []) => request({url: '/vmware/clone-vm/run', method: 'post', data: {ids: data}})


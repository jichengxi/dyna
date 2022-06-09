import {getVmHostApi} from '@/api/vmware/dashboard'

export const getVmHost = async (host: string) => {
    const res = await getVmHostApi({"host": host})
    try {
        if (res.data.code === 0) {
            return res.data.data
        }
        console.warn(res.data.message)
        return {}
    } catch (error) {
        console.error(error)
        return {}
    }
}



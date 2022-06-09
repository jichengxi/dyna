import {getCloneVmJobsApi, getCloneVmTemplateApi, postCloneVmJobsApi, postCloneVmRunApi} from "@/api/vmware/clone-vm";

export const getCloneVmJobs = async () => {
    const res = await getCloneVmJobsApi()
    try {
        if (res.data.code === 0) {
            console.log("res.data.data", res.data.data)
            return res.data.data
        }
        console.warn(res.data.message)
        return []
    } catch (error) {
        console.error(error)
        return []
    }
}

export const postCloneVmJobs = async (data: {}) => {
    const res = await postCloneVmJobsApi(data)
    try {
        if (res.data.code === 0) {
            console.log("res.data.data", res.data.data)
            return res.data.data
        }
        console.warn(res.data.message)
        return []
    } catch (error) {
        console.error(error)
        return []
    }
}

export const getCloneVmTemplate = async () => {
    await getCloneVmTemplateApi().then(res => {
        const link = document.createElement("a")
        link.style.display = "none";
        link.href = URL.createObjectURL(res.data)
        link.setAttribute("download", "cloneVmTemplate.xlsx")
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)

    })
}

export const postCloneVmRun = async (data: []) => {
    const res = await postCloneVmRunApi(data)
    console.log(res.data)
    try {
        if (res.data.code !== 0) {
            console.warn(res.data.message)
        }
    } catch (error) {
        console.error(error)
    }
}

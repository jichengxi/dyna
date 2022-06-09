export default [
    {
        url: '/api/vmware/clusters',
        method: 'get',
        response: () => {
            return {
                code: 0,
                message: 'ok',
                data: [
                    {
                        // key: 0,
                        name: '21viacloud',
                        ip: '222.73.18.26',
                        user: 'administrator@vsphere.local',
                        password: '111111',
                        tags: ['nice', 'developer']
                    },
                    {
                        // key: 1,
                        name: 'gz-zsc',
                        ip: '172.24.32.51',
                        user: 'administrator@vsphere.local',
                        password: '111111',
                        tags: ['nice', 'developer']
                    },
                    {
                        // key: 2,
                        name: 'gz-mix',
                        ip: '172.24.32.53',
                        user: 'administrator@vsphere.local',
                        password: '111111',
                        tags: ['nice', 'developer']
                    }
                ]
            }
        }
    }
]



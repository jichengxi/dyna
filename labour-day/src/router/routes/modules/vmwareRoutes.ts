import Layout from '@/layout/index.vue'
import {markRaw} from "vue";
import {IconMenuTable} from "@/components/AppIcons";

export default [
    {
        name: 'Vmware',
        path: '/vmware',
        component: Layout,
        redirect: '/vmware/home/view',
        meta: {
            title: '虚拟化',
            role: ['admin'],
        },
        children: [
            {
                name: 'VmwareHome',
                path: 'home',
                component: () => import('@/views2/vmware/index.vue'),
                redirect: '/vmware/home/view',
                meta: {
                    title: 'vmware',
                    role: ['admin'],
                    icon: markRaw(IconMenuTable),
                },
                children: [
                    {
                        name: 'VmwareView',
                        path: 'view',
                        component: () => import('@/views2/vmware/view/index.vue'),
                        meta: {
                            title: '详情',
                            role: ['admin'],
                        }
                    },
                    {
                        name: 'VmwareClone',
                        path: 'clone',
                        component: () => import('@/views2/vmware/index.vue'),
                        redirect: '/vmware/home/clone/jobs',
                        meta: {
                            title: '克隆',
                            role: ['admin'],
                        },
                        children: [
                            {
                                name: 'VmwareCloneJobsList',
                                path: 'jobs',
                                component: () => import('@/views2/vmware/clone/clone-vm-jobs-list.vue'),
                                meta: {
                                    title: '克隆任务列表',
                                    role: ['admin'],
                                },
                            },
                            {
                                name: 'VmwareCloneJobsAdd',
                                path: 'jobs-add',
                                component: () => import('@/views2/vmware/clone/clone-vm-jobs-add.vue'),
                                meta: {
                                    title: '添加克隆任务',
                                    role: ['admin'],
                                },
                            }

                        ]
                    },
                    {
                        name: 'VmwareRelocate',
                        path: 'relocate',
                        component: () => import('@/views2/vmware/relocate/index.vue'),
                        meta: {
                            title: '迁移',
                            role: ['admin'],
                        },
                    },
                    {
                        name: 'VmwareHostInfo',
                        path: 'host-info',
                        isHidden: true,
                        component: () => import('@/views2/vmware/view/host-info.vue'),
                        meta: {
                            title: '主机详情',
                            role: ['admin'],
                        },
                    }
                ],
            }
        ]
    }
]

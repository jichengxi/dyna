import Layout from '@/layout/index.vue'
import {markRaw} from "vue";
import {IconMenuTable} from "@/components/AppIcons";

export default [
    {
        name: 'Admin',
        path: '/admin',
        component: Layout,
        redirect: '/admin/manage/vc-manage',
        meta: {
            title: '集群管理',
            role: ['admin'],
        },
        children: [
            {
                name: 'Manage',
                path: 'manage',
                component: () => import('@/views2/management/index.vue'),
                meta: {
                    title: '集群管理',
                    role: ['admin'],
                    icon: markRaw(IconMenuTable),
                },
                children: [
                    {
                        name: 'VcManage',
                        path: 'vc-manage',
                        component: () => import('@/views2/management/vmware/index.vue'),
                        meta: {
                            title: 'vc集群管理',
                            role: ['admin'],
                        },
                    },
                    {
                        name: 'ContainerManage',
                        path: 'container-manage',
                        component: () => import('@/views2/management/vmware/test.vue'),
                        meta: {
                            title: '容器集群管理',
                            role: ['admin'],
                        },
                    },
                    {
                        name: 'VcAdd',
                        path: 'vc-add',
                        component: () => import('@/views2/management/vmware/vc-add.vue'),
                        meta: {
                            title: '录入vc集群',
                            role: ['admin'],
                        },
                    },
                ]
            }
        ]
    }
]


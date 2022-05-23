/**
 * @file helm router 配置
 */

// Helm应用列表
const helms = () => import(/* webpackChunkName: 'helm' */'@/views/helm')

// Helm模板列表
const helmTplList = () => import(/* webpackChunkName: 'helm' */'@/views/helm/tpl-list.vue')

// Helm模板详情
const helmTplDetail = () => import(/* webpackChunkName: 'helm' */'@/views/helm/tpl-detail.vue')

// Helm实例化
const helmTplInstance = () => import(/* webpackChunkName: 'helm' */'@/views/helm/tpl-instance.vue')

// Helm app详情
const helmAppDetail = () => import(/* webpackChunkName: 'helm' */'@/views/helm/app-detail.vue')

// helm repo
const helmRepo = () => import(/* webpackChunkName: 'helm' */'@/views/helm/repo.vue')
const helmRepoEdit = () => import(/* webpackChunkName: 'helm' */'@/views/helm/repo-edit.vue')

const childRoutes = [
    {
        path: ':projectCode/helm',
        name: 'helms',
        component: helms
    },
    {
        path: ':projectCode/helm/list',
        name: 'helmTplList',
        component: helmTplList
    },
    {
        path: ':projectCode/helm/tpl/:tplId',
        name: 'helmTplDetail',
        component: helmTplDetail,
        meta: {
            menuId: 'helmTplList'
        }
    },
    {
        path: ':projectCode/helm/instance/:tplId',
        name: 'helmTplInstance',
        component: helmTplInstance,
        meta: {
            menuId: 'helmTplList'
        }
    },
    {
        path: ':projectCode/helm/app/:appId',
        name: 'helmAppDetail',
        component: helmAppDetail,
        meta: {
            menuId: 'helms'
        }
    },
    {
        path: ':projectCode/helm/repo',
        name: 'helmRepo',
        component: helmRepo,
        meta: {
            title: window.i18n.t('仓库'),
            hideBack: true
        }
    },
    {
        path: ':projectCode/helm/repo/edit/:name?',
        name: 'helmEdit',
        props: true,
        component: helmRepoEdit
    }
]

export default childRoutes

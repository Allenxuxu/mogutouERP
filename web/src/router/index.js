import Vue from 'vue'
import Router from 'vue-router'

// in development-env not use lazy-loading, because lazy-loading too many pages will cause webpack hot update too slow. so only in production use lazy-loading;
// detail: https://panjiachen.github.io/vue-element-admin-site/#/lazy-loading

Vue.use(Router)

/* Layout */
import Layout from '../views/layout/Layout'

/**
* hidden: true                   if `hidden:true` will not show in the sidebar(default is false)
* alwaysShow: true               if set true, will always show the root menu, whatever its child routes length
*                                if not set alwaysShow, only more than one route under the children
*                                it will becomes nested mode, otherwise not show the root menu
* redirect: noredirect           if `redirect:noredirect` will no redirect in the breadcrumb
* name:'router-name'             the name is used by <keep-alive> (must set!!!)
* meta : {
    title: 'title'               the name show in submenu and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    breadcrumb: false            if false, the item will hidden in breadcrumb(default is true)
  }
**/
export const constantRouterMap = [
  { path: '/login', component: () => import('@/views/login/index'), hidden: true },
  { path: '/404', component: () => import('@/views/404'), hidden: true },

  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    name: 'Dashboard',
    hidden: true,
    children: [{
      path: 'dashboard',
      component: () => import('@/views/index')
    }]
  },

  {
    path: '/company',
    component: Layout,
    redirect: '/company/staff',
    name: 'Company',
    meta: { title: '公司', icon: 'company' },
    children: [
      {
        path: 'staff',
        name: 'CStaff',
        component: () => import('@/views/company/staff'),
        meta: { title: '所有员工' }
      }
    ]
  },

  {
    path: '/stock',
    component: Layout,
    redirect: '/stock/index',
    children: [
      {
        path: 'index',
        name: 'Stock',
        component: () => import('@/views/stock/index'),
        meta: { title: '库存查询', icon: 'warehouse' }
      }
    ]
  }
]

export default new Router({
  mode: 'hash',
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRouterMap
})

export const asyncRouterMap = [
  {
    path: '/order',
    component: Layout,
    redirect: '/order/sales',
    alwaysShow: true, // will always show the root menu
    meta: {
      title: '订单',
      icon: 'order'
    },
    children: [
      {
        path: 'sales',
        component: () => import('@/views/order/sales'),
        redirect: '/order/sales/edit',
        name: 'Sales',
        meta: {
          title: '客户订单'
        },
        children: [
          {
            path: 'edit',
            component: () => import('@/views/order/sales/list'),
            name: 'EditSales',
            meta: {
              title: '所有订单'
            }
          },
          {
            path: 'create',
            component: () => import('@/views/order/sales/create'),
            name: 'CreateSales',
            meta: {
              title: '新增订单'
            }
          }
        ]
      },
      {
        path: 'purchase',
        component: () => import('@/views/order/purchase'),
        name: 'Purchase',
        redirect: '/order/purchase/edit',
        meta: {
          title: '采购订单',
          roles: ['PM'] // or you can only set roles in sub nav
        },
        children: [
          {
            path: 'edit',
            component: () => import('@/views/order/purchase/list'),
            name: 'EditPurchase',
            meta: {
              title: '所有订单'
            }
          },
          {
            path: 'create',
            component: () => import('@/views/order/purchase/create'),
            name: 'CreatePurchase',
            meta: {
              title: '新增订单'
            }
          }
        ]
      }
    ]
  },

  {
    path: '/commodity',
    component: Layout,
    redirect: '/commodity/edit',
    name: 'Commodity',
    meta: {
      roles: ['CM'],
      title: '商品管理',
      icon: 'commodity'
    },
    children: [
      {
        path: 'edit',
        name: 'EditCommodity',
        component: () => import('@/views/commodity/edit'),
        meta: { title: '信息删改' }
      },
      {
        path: 'create',
        name: 'CreateCommodity',
        component: () => import('@/views/commodity/create'),
        meta: { title: '增加商品' }
      }
    ]
  },

  {
    path: '/staff',
    component: Layout,
    redirect: '/staff/index',
    meta: {
      roles: ['admin']
    },
    children: [
      {
        path: 'index',
        name: 'Staff',
        component: () => import('@/views/staff/index'),
        meta: { title: '员工管理', icon: 'staff' }
      }
    ]
  },

  {
    path: '/finance',
    component: Layout,
    redirect: '/finance/index',
    meta: {
      roles: ['admin']
    },
    children: [
      {
        path: 'index',
        name: 'Finance',
        component: () => import('@/views/finance/index'),
        meta: { title: '财务管理', icon: 'finance' }
      }
    ]
  },

  { path: '*', redirect: '/404', hidden: false }
]

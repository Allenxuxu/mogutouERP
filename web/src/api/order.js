import request from '@/utils/request'

export function getCustormerOrders() {
  return request({
    url: '/orders/custormer',
    method: 'get'
  })
}

export function getPurchaseOrders() {
  return request({
    url: '/orders/purchase',
    method: 'get'
  })
}

export function createPurchaseOrder(data) {
  return request({
    url: '/order/purchase',
    method: 'post',
    data: data
  })
}

export function createCustormerOrder(data) {
  return request({
    url: '/order/custormer',
    method: 'post',
    data: data
  })
}

export function deletePurchaseOrder(id) {
  return request({
    url: '/order/purchase/' + id,
    method: 'delete'
  })
}

export function deleteCustormerOrder(id) {
  return request({
    url: '/order/custormer/' + id,
    method: 'delete'
  })
}

export function confirmPurchaseOrder(id, data) {
  return request({
    url: '/order/purchase/' + id + '/confirm',
    method: 'patch',
    data: data
  })
}

export function confirmCustormerOrder(id, data) {
  return request({
    url: '/order/custormer/' + id + '/confirm',
    method: 'patch',
    data: data
  })
}

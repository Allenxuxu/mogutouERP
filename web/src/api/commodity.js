import request from '@/utils/request'

export function getCommodities() {
  return request({
    url: '/commodities',
    method: 'get'
  })
}

export function createCommodity(data) {
  return request({
    url: '/commodity',
    method: 'post',
    data: data
  })
}

export function deleteCommodity(id) {
  return request({
    url: '/commodity/' + id,
    method: 'delete'
  })
}

export function updateCommodity(id, data) {
  return request({
    url: '/commodity/' + id,
    method: 'patch',
    data: data
  })
}

export function getCommoditiesAsAdmin() {
  return request({
    url: '/admin/commodities',
    method: 'get'
  })
}

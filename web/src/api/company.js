import request from '@/utils/request'

export function getStaff() {
  return request({
    url: '/users',
    method: 'get'
  })
}

export function createStaff(data) {
  return request({
    url: '/user',
    method: 'post',
    data: data
  })
}

export function updateStaff(id, data) {
  return request({
    url: '/user/' + id,
    method: 'patch',
    data: data
  })
}

export function deleteStaff(id) {
  return request({
    url: '/user/' + id,
    method: 'delete'
  })
}

export function updatePassword(data) {
  return request({
    url: '/userPassword',
    method: 'patch',
    data: data
  })
}

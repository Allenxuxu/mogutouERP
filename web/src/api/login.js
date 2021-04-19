import request from '@/utils/request'
import md5 from 'js-md5'

export function login(tel, password) {
  return request({
    url: '/login',
    method: 'post',
    data: {
      tel: tel,
      password: md5(password)
    }
  })
}

export function getInfo() {
  return request({
    url: '/user',
    method: 'get'
  })
}

export function logout() {
  return request({
    url: '/logout',
    method: 'get'
  })
}

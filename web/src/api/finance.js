import request from '@/utils/request'

export function getFinance(year) {
  return request({
    url: '/admin/finance/' + year,
    method: 'get'
  })
}

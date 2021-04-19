import Vue from 'vue'

import 'normalize.css/normalize.css' // A modern alternative to CSS resets

import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
// import locale from 'element-ui/lib/locale/lang/en' // lang i18n

import '@/styles/index.scss' // global css

import App from './App'
import router from './router'
import store from './store'

import '@/icons' // icon
import '@/permission' // permission control

Vue.use(ElementUI, {})

Vue.config.productionTip = false
Vue.directive('enterToNext', {
  inserted: function(el) {
    console.log('enterToNext...')
    // let frm = el.querySelector('.el-form');
    const inputs = el.querySelectorAll('input')
    console.log(inputs)
    // 绑定回写事件
    for (var i = 0; i < inputs.length; i++) {
      inputs[i].setAttribute('keyFocusIndex', i)
      inputs[i].addEventListener('keyup', (ev) => {
        if (ev.keyCode === 13) {
          const targetTo = ev.srcElement.getAttribute('keyFocusTo')
          if (targetTo) {
            this.$refs[targetTo].$el.focus()
          } else {
            var attrIndex = ev.srcElement.getAttribute('keyFocusIndex')
            var ctlI = parseInt(attrIndex)
            if (ctlI < inputs.length - 1) { inputs[ctlI + 1].focus() }
          }
        }
      })
    }
  }
})
new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App)
})

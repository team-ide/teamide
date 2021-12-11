import Vue from 'vue'
import App from '@/App.vue'
import router from '@/router'

Vue.config.productionTip = false

import tool from '@/tool'
import server from '@/server'

Vue.prototype.tool = tool
Vue.prototype.server = server

import '@mdi/font/css/materialdesignicons.css'

import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'

import '@/app.scss'

// Import Bootstrap an BootstrapVue CSS files (order is important)
// import 'bootstrap/dist/css/bootstrap.css'
// import 'bootstrap-vue/dist/bootstrap-vue.css'

// Make BootstrapVue available throughout your project
Vue.use(BootstrapVue)
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin)

import { } from '@/component'

new Vue({
  router,
  render: h => h(App),
  created() {
    this.tool.$bvToast = this.$bvToast;
  },
  mounted() {
    this.tool.init()
  }
}).$mount('#app')
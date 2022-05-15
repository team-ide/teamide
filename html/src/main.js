import Vue from 'vue'
import App from '@/App.vue'
import router from '@/router'


Vue.config.productionTip = false

import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
Vue.use(ElementUI)


import tool from '@/tool'
import server from '@/server'

Vue.prototype.tool = tool
Vue.prototype.server = server

import '@mdi/font/css/materialdesignicons.css'


import tm from 'teamide-ui'
Vue.use(tm)

import { } from '@/component'

import form from "@/form";

Vue.prototype.form = form;

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
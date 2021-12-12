import Vue from 'vue'
import Router from 'vue-router'

import Index from '@/views/Index'
import Toolbox from '@/views/Toolbox'
import Workspace from '@/views/Workspace'

Vue.use(Router)

export default new Router({
  routes: [{
    path: '/',
    component: Index
  },
  {
    path: '/toolbox',
    component: Toolbox
  },
  {
    path: '/workspace',
    component: Workspace
  },
  ]
})
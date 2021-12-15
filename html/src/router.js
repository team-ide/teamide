import Vue from 'vue'
import Router from 'vue-router'

import Toolbox from '@/views/Toolbox'
import Workspace from '@/views/Workspace'

import Manage from '@/views/manage/Index'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/toolbox',
      component: Toolbox
    },
    {
      path: '/workspace',
      component: Workspace
    },
    {
      path: '/manage',
      component: Manage
    },
  ]
})
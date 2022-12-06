import Vue from 'vue'
import Router from 'vue-router'


import Index from '@/views/Index'

// push相同地址 会报错，使用该方式规避
let push = Router.prototype.push;
Router.prototype.push = function (location) {
  return push.call(this, location).catch(e => e);
};

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      component: Index
    },
  ]
})
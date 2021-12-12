import Vue from 'vue'

import InfoBox from './InfoBox.vue'
Vue.component('InfoBox', InfoBox);

import SystemInfoBox from './SystemInfoBox.vue'
Vue.component('SystemInfoBox', SystemInfoBox);

import AlertBox from './AlertBox.vue'
Vue.component('AlertBox', AlertBox);

import Login from './Login.vue'
Vue.component('Login', Login);

import Register from './Register.vue'
Vue.component('Register', Register);

import Workspace from './workspace/Workspace.vue'
Vue.component('Workspace', Workspace);

import Console from './console/Console.vue'
Vue.component('Console', Console);

export default {};
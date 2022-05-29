import Vue from 'vue'

import './tm';

import './iconfont/iconfont.css'

import TMLayout from './tm/Layout.vue'
Vue.component(TMLayout.name, TMLayout);
import TMLayoutBar from './tm/LayoutBar.vue'
Vue.component(TMLayoutBar.name, TMLayoutBar);

import InfoBox from './message/InfoBox.vue'
Vue.component('InfoBox', InfoBox);

import SystemInfoBox from './message/SystemInfoBox.vue'
Vue.component('SystemInfoBox', SystemInfoBox);

import AlertBox from './message/AlertBox.vue'
Vue.component('AlertBox', AlertBox);

import Login from './Login.vue'
Vue.component('Login', Login);

import Register from './Register.vue'
Vue.component('Register', Register);

import WaterfallLayout from './WaterfallLayout.vue'
Vue.component('WaterfallLayout', WaterfallLayout);

import IconFont from './icon/IconFont.vue'
Vue.component('IconFont', IconFont);
import Icon from './icon/Icon.vue'
Vue.component('Icon', Icon);

import Contextmenu from './contextmenu/Contextmenu.vue'
Vue.component('Contextmenu', Contextmenu);
import ContextmenuMenus from './contextmenu/ContextmenuMenus.vue'
Vue.component('ContextmenuMenus', ContextmenuMenus);

import Form from './form/Form.vue'
Vue.component('Form', Form);
import FormBox from './form/FormBox.vue'
Vue.component('FormBox', FormBox);
import FormDialog from './form/FormDialog.vue'
Vue.component('FormDialog', FormDialog);

import ShouldLogin from './error/ShouldLogin.vue'
Vue.component('ShouldLogin', ShouldLogin);

import NoPower from './error/NoPower.vue'
Vue.component('NoPower', NoPower);

import TabEditor from './tab-editor/Index.vue'
Vue.component('TabEditor', TabEditor);

import MenuBox from './menu/MenuBox.vue'
Vue.component('MenuBox', MenuBox);

import MenuSubBox from './menu/MenuSubBox.vue'
Vue.component('MenuSubBox', MenuSubBox);

import MenuItem from './menu/MenuItem.vue'
Vue.component('MenuItem', MenuItem);


import ToolboxEditor from './toolbox/Index.vue'
Vue.component('ToolboxEditor', ToolboxEditor);

import ToolboxRedisEditor from './toolbox/redis/Index.vue'
Vue.component('ToolboxRedisEditor', ToolboxRedisEditor);

import ToolboxKafkaEditor from './toolbox/kafka/Index.vue'
Vue.component('ToolboxKafkaEditor', ToolboxKafkaEditor);

import ToolboxZookeeperEditor from './toolbox/Zookeeper.vue'
Vue.component('ToolboxZookeeperEditor', ToolboxZookeeperEditor);

import ToolboxElasticsearchEditor from './toolbox/elasticsearch/Index.vue'
Vue.component('ToolboxElasticsearchEditor', ToolboxElasticsearchEditor);

import ToolboxDatabaseEditor from './toolbox/database/Index.vue'
Vue.component('ToolboxDatabaseEditor', ToolboxDatabaseEditor);

import ToolboxSSHEditor from './toolbox/ssh/Index.vue'
Vue.component('ToolboxSSHEditor', ToolboxSSHEditor);

import ToolboxOtherEditor from './toolbox/other/Index.vue'
Vue.component('ToolboxOtherEditor', ToolboxOtherEditor);

import ScriptInfo from './script/Info.vue'
Vue.component('ScriptInfo', ScriptInfo);

// import Editor from './Editor.vue'
// Vue.component('Editor', Editor);

export default {};
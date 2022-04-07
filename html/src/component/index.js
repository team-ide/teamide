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

import Contextmenu from './contextmenu/Contextmenu.vue'
Vue.component('Contextmenu', Contextmenu);
import ContextmenuMenus from './contextmenu/ContextmenuMenus.vue'
Vue.component('ContextmenuMenus', ContextmenuMenus);

import Form from './Form.vue'
Vue.component('Form', Form);

import ShouldLogin from './ShouldLogin.vue'
Vue.component('ShouldLogin', ShouldLogin);

import NoPower from './NoPower.vue'
Vue.component('NoPower', NoPower);

import ModelEditor from './modelEditor/Index.vue'
Vue.component('ModelEditor', ModelEditor);

import ModelEditorFieldInput from './modelEditor/FieldInput.vue'
Vue.component('ModelEditorFieldInput', ModelEditorFieldInput);

import ModelEditorFields from './modelEditor/Fields.vue'
Vue.component('ModelEditorFields', ModelEditorFields);

import ModelEditorField from './modelEditor/Field.vue'
Vue.component('ModelEditorField', ModelEditorField);

import ModelEditorFieldList from './modelEditor/FieldList.vue'
Vue.component('ModelEditorFieldList', ModelEditorFieldList);

import ModelEditorAction from './modelEditor/Action.vue'
Vue.component('ModelEditorAction', ModelEditorAction);

import ModelEditorStep from './modelEditor/Step.vue'
Vue.component('ModelEditorStep', ModelEditorStep);

import ModelEditorSteps from './modelEditor/Steps.vue'
Vue.component('ModelEditorSteps', ModelEditorSteps);

import MenuBox from './menu/MenuBox.vue'
Vue.component('MenuBox', MenuBox);

import MenuSubBox from './menu/MenuSubBox.vue'
Vue.component('MenuSubBox', MenuSubBox);

import MenuItem from './menu/MenuItem.vue'
Vue.component('MenuItem', MenuItem);


import ToolboxEditor from './toolbox/Index.vue'
Vue.component('ToolboxEditor', ToolboxEditor);

import ToolboxRedisEditor from './toolbox/Redis.vue'
Vue.component('ToolboxRedisEditor', ToolboxRedisEditor);

import ToolboxKafkaEditor from './toolbox/kafka/Index.vue'
Vue.component('ToolboxKafkaEditor', ToolboxKafkaEditor);

import ToolboxKafkaTopic from './toolbox/kafka/Topic.vue'
Vue.component('ToolboxKafkaTopic', ToolboxKafkaTopic);

import ToolboxKafkaTabs from './toolbox/kafka/Tabs.vue'
Vue.component('ToolboxKafkaTabs', ToolboxKafkaTabs);

import ToolboxKafkaTopicData from './toolbox/kafka/TopicData.vue'
Vue.component('ToolboxKafkaTopicData', ToolboxKafkaTopicData);

import ToolboxZookeeperEditor from './toolbox/Zookeeper.vue'
Vue.component('ToolboxZookeeperEditor', ToolboxZookeeperEditor);

import ToolboxElasticsearchEditor from './toolbox/Elasticsearch.vue'
Vue.component('ToolboxElasticsearchEditor', ToolboxElasticsearchEditor);

import ToolboxDatabaseEditor from './toolbox/database/Index.vue'
Vue.component('ToolboxDatabaseEditor', ToolboxDatabaseEditor);

import ToolboxDatabaseDatabase from './toolbox/database/Database.vue'
Vue.component('ToolboxDatabaseDatabase', ToolboxDatabaseDatabase);

import ToolboxDatabaseTable from './toolbox/database/Table.vue'
Vue.component('ToolboxDatabaseTable', ToolboxDatabaseTable);

import ToolboxDatabaseTableData from './toolbox/database/TableData.vue'
Vue.component('ToolboxDatabaseTableData', ToolboxDatabaseTableData);

import ToolboxDatabaseTabs from './toolbox/database/Tabs.vue'
Vue.component('ToolboxDatabaseTabs', ToolboxDatabaseTabs);

import ToolboxDatabaseSql from './toolbox/database/Sql.vue'
Vue.component('ToolboxDatabaseSql', ToolboxDatabaseSql);

import ToolboxSSHEditor from './toolbox/ssh/Index.vue'
Vue.component('ToolboxSSHEditor', ToolboxSSHEditor);
import ToolboxSSH from './toolbox/ssh/SSH.vue'
Vue.component('ToolboxSSH', ToolboxSSH);
import ToolboxFTP from './toolbox/ssh/FTP.vue'
Vue.component('ToolboxFTP', ToolboxFTP);
import ToolboxFTPFiles from './toolbox/ssh/Files.vue'
Vue.component('ToolboxFTPFiles', ToolboxFTPFiles);
export default {};
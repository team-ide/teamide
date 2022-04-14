<template>
  <div class="toolbox-database-tabs">
    <template v-if="ready">
      <TabEditor
        ref="TabEditor"
        :source="source"
        :onRemoveTab="onRemoveTab"
        :onActiveTab="onActiveTab"
      >
        <template v-slot:body="{ tab }">
          <template v-if="tab.type == 'data'">
            <ToolboxDatabaseTableData
              :source="source"
              :toolbox="toolbox"
              :toolboxType="toolboxType"
              :data="data"
              :wrap="wrap"
              :database="tab.data.database"
              :table="tab.data"
            >
            </ToolboxDatabaseTableData>
          </template>
          <template v-else-if="tab.type == 'sql'">
            <ToolboxDatabaseSql
              :source="source"
              :toolbox="toolbox"
              :toolboxType="toolboxType"
              :data="data"
              :wrap="wrap"
              :sqlData="tab.data"
            >
            </ToolboxDatabaseSql>
          </template>
        </template>
      </TabEditor>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "data", "toolboxType", "toolbox", "option", "wrap"],
  data() {
    return {
      ready: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.wrap.doActiveTab = this.doActiveTab;
      this.wrap.createTab = this.createTab;
      this.wrap.createTabByData = this.createTabByData;
      this.wrap.getTabByData = this.getTabByData;
      this.wrap.addTab = this.addTab;
      this.ready = true;
    },
    getTab(tab) {
      return this.$refs.TabEditor.getTab(tab);
    },
    onRemoveTab(tab) {},
    onActiveTab(tab) {},
    addTab(tab) {
      return this.$refs.TabEditor.addTab(tab);
    },
    doActiveTab(tab) {
      return this.$refs.TabEditor.doActiveTab(tab);
    },
    getTabKeyByData(type, data) {
      let key;
      if (type == "data") {
        key = "" + data.database.name + ":" + data.name;
      } else {
        key = data.key;
      }

      return key;
    },
    getTabByData(type, data) {
      let key = this.getTabKeyByData(type, data);
      let tab = this.getTab(key);
      return tab;
    },
    createTabByData(type, data) {
      let key = this.getTabKeyByData(type, data);

      let tab = this.getTab(key);
      if (tab == null) {
        let title = "";
        let name = "";
        if (type == "data") {
          title = data.database.name + "." + data.name;
          name = data.database.name + "." + data.name;
        } else {
          title = data.title || data.name;
          name = data.name || data.title;
        }
        tab = {
          key,
          data,
          type,
          title,
          name,
        };
        tab.active = false;
        tab.changed = false;
      }
      return tab;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-tabs {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>

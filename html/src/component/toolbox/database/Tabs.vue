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
          <template v-if="tab.extend.type == 'data'">
            <ToolboxDatabaseTableData
              :source="source"
              :toolbox="toolbox"
              :wrap="wrap"
              :database="tab.extend.database"
              :table="tab.extend.table"
            >
            </ToolboxDatabaseTableData>
          </template>
          <template v-else-if="tab.extend.type == 'sql'">
            <ToolboxDatabaseSql
              :source="source"
              :toolbox="toolbox"
              :toolboxType="toolboxType"
              :wrap="wrap"
            >
            </ToolboxDatabaseSql>
          </template>
          <template v-else-if="tab.extend.type == 'ddl'">
            <DDL
              :source="source"
              :wrap="wrap"
              :database="tab.extend.database"
              :table="tab.extend.table"
            >
            </DDL>
          </template>
        </template>
      </TabEditor>
    </template>
  </div>
</template>


<script>
import DDL from "./DDL";

export default {
  components: {DDL},
  props: ["source", "toolboxType", "toolbox", "option", "wrap"],
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
      this.wrap.addTab = this.addTab;
      this.wrap.getTab = this.getTab;
      this.ready = true;
    },
    getTab(tab) {
      return this.$refs.TabEditor.getTab(tab);
    },
    onRemoveTab(tab) {
      this.wrap.onRemoveTab(tab);
    },
    onActiveTab(tab) {
      this.wrap.onActiveTab(tab);
    },
    addTab(tab) {
      return this.$refs.TabEditor.addTab(tab);
    },
    doActiveTab(tab) {
      return this.$refs.TabEditor.doActiveTab(tab);
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

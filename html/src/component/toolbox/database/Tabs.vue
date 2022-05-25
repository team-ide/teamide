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
            <TableData
              :source="source"
              :toolbox="toolbox"
              :wrap="wrap"
              :database="tab.extend.database"
              :table="tab.extend.table"
              :extend="tab.extend"
              :tab="tab"
            >
            </TableData>
          </template>
          <template v-else-if="tab.extend.type == 'sql'">
            <Sql
              :source="source"
              :wrap="wrap"
              :extend="tab.extend"
              :databases="databases"
              :tab="tab"
            >
            </Sql>
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
          <template v-if="tab.extend.type == 'table'">
            <Table
              :source="source"
              :toolbox="toolbox"
              :wrap="wrap"
              :database="tab.extend.database"
              :table="tab.extend.table"
              :extend="tab.extend"
            >
            </Table>
          </template>
          <template v-if="tab.extend.type == 'export'">
            <Export
              :source="source"
              :toolbox="toolbox"
              :wrap="wrap"
              :database="tab.extend.database"
              :table="tab.extend.table"
              :extend="tab.extend"
            >
            </Export>
          </template>
          <template v-if="tab.extend.type == 'import'">
            <Import
              :source="source"
              :toolbox="toolbox"
              :wrap="wrap"
              :database="tab.extend.database"
              :table="tab.extend.table"
              :extend="tab.extend"
            >
            </Import>
          </template>
        </template>
      </TabEditor>
    </template>
  </div>
</template>


<script>
import DDL from "./DDL";
import Sql from "./Sql";
import Table from "./Table";
import TableData from "./TableData";
import Export from "./Export";
import Import from "./Import";

export default {
  components: { DDL, Sql, Table, TableData, Export, Import },
  props: ["source", "toolboxType", "toolbox", "option", "wrap", "databases"],
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
      return this.$refs.TabEditor && this.$refs.TabEditor.getTab(tab);
    },
    onRemoveTab(tab) {
      this.wrap.onRemoveTab(tab);
    },
    onActiveTab(tab) {
      this.wrap.onActiveTab(tab);
    },
    addTab(tab) {
      return this.$refs.TabEditor && this.$refs.TabEditor.addTab(tab);
    },
    doActiveTab(tab) {
      return this.$refs.TabEditor && this.$refs.TabEditor.doActiveTab(tab);
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

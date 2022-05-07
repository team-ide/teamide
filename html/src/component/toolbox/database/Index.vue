<template>
  <div class="toolbox-database-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout width="400px" class="">
          <ToolboxDatabaseDatabase
            ref="ToolboxDatabaseDatabase"
            :source="source"
            :toolbox="toolbox"
            :toolboxType="toolboxType"
            :wrap="wrap"
          >
          </ToolboxDatabaseDatabase>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <ToolboxDatabaseTabs
            :source="source"
            :toolbox="toolbox"
            :toolboxType="toolboxType"
            :wrap="wrap"
          >
          </ToolboxDatabaseTabs>
        </tm-layout>
      </tm-layout>
      <ShowDatabaseCreate :source="source" :wrap="wrap"> </ShowDatabaseCreate>
      <ShowTableCreate :source="source" :wrap="wrap"> </ShowTableCreate>
      <ShowSqlForInsert :source="source" :wrap="wrap"> </ShowSqlForInsert>
    </template>
  </div>
</template>


<script>
import ShowDatabaseCreate from "./ShowDatabaseCreate";
import ShowTableCreate from "./ShowTableCreate";
import ShowSqlForInsert from "./ShowSqlForInsert";

export default {
  components: {
    ShowDatabaseCreate,
    ShowTableCreate,
    ShowSqlForInsert,
  },
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
      this.ready = true;
    },
    refresh() {
      this.$refs.ToolboxDatabaseDatabase.refresh();
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-editor {
  width: 100%;
  height: 100%;
}
</style>

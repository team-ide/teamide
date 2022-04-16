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
            :data="data"
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
            :data="data"
            :wrap="wrap"
          >
          </ToolboxDatabaseTabs>
        </tm-layout>
      </tm-layout>
      <ShowDatabaseCreate :source="source" :wrap="wrap"> </ShowDatabaseCreate>
      <ShowTableCreate :source="source" :wrap="wrap"> </ShowTableCreate>
    </template>
  </div>
</template>


<script>
import ShowDatabaseCreate from "./ShowDatabaseCreate";
import ShowTableCreate from "./ShowTableCreate";

export default {
  components: { ShowDatabaseCreate, ShowTableCreate },
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

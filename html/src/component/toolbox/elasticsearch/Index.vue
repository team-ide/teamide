<template>
  <div class="toolbox-elasticsearch-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout width="400px" class="scrollbar">
          <IndexNames
            ref="ToolboxElasticsearchIndexNames"
            :source="source"
            :toolbox="toolbox"
            :toolboxType="toolboxType"
            :wrap="wrap"
          >
          </IndexNames>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <Tabs
            :source="source"
            :toolbox="toolbox"
            :toolboxType="toolboxType"
            :wrap="wrap"
          >
          </Tabs>
        </tm-layout>
      </tm-layout>
      <IndexForm :source="source" :toolbox="toolbox" :wrap="wrap"> </IndexForm>
      <DataForm :source="source" :toolbox="toolbox" :wrap="wrap"> </DataForm>
      <ShowData :source="source" :toolbox="toolbox" :wrap="wrap"> </ShowData>
    </template>
  </div>
</template>


<script>
import IndexNames from "./IndexNames";
import Tabs from "./Tabs";
import IndexForm from "./IndexForm";
import DataForm from "./DataForm";
import ShowData from "./ShowData";
export default {
  components: { IndexNames, Tabs, IndexForm, DataForm, ShowData },
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
      this.$refs.ToolboxElasticsearchIndexNames.refresh();
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-elasticsearch-editor {
  width: 100%;
  height: 100%;
  user-select: text;
}
</style>

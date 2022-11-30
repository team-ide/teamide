<template>
  <div class="toolbox-elasticsearch-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout width="400px" class="">
          <IndexNames
            ref="ToolboxElasticsearchIndexNames"
            :source="source"
            :toolboxWorker="toolboxWorker"
            :extend="extend"
          >
          </IndexNames>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <Tabs :source="source" :toolboxWorker="toolboxWorker"> </Tabs>
        </tm-layout>
      </tm-layout>
      <IndexForm :source="source" :toolboxWorker="toolboxWorker"> </IndexForm>
      <ReindexForm :source="source" :toolboxWorker="toolboxWorker">
      </ReindexForm>
      <DataForm :source="source" :toolboxWorker="toolboxWorker"> </DataForm>
      <MappingForm :source="source" :toolboxWorker="toolboxWorker">
      </MappingForm>
      <ShowInfo :source="source" :toolboxWorker="toolboxWorker"> </ShowInfo>
    </template>
  </div>
</template>


<script>
import IndexNames from "./IndexNames";
import Tabs from "./Tabs";
import IndexForm from "./IndexForm";
import ReindexForm from "./ReindexForm";
import MappingForm from "./MappingForm";
import DataForm from "./DataForm";
import ShowInfo from "./ShowInfo";
export default {
  components: {
    IndexNames,
    Tabs,
    IndexForm,
    ReindexForm,
    MappingForm,
    DataForm,
    ShowInfo,
  },
  props: ["source", "toolboxWorker", "extend"],
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

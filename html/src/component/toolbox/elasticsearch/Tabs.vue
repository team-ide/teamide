<template>
  <div class="toolbox-elasticsearch-tabs">
    <div class="default-tabs-container">
      <WorkspaceTabs :source="source" :itemsWorker="toolboxWorker.itemsWorker">
      </WorkspaceTabs>
    </div>
    <div class="default-spans-container">
      <WorkspaceSpans :source="source" :itemsWorker="toolboxWorker.itemsWorker">
        <template v-slot:span="{ item }">
          <IndexData
            :source="source"
            :toolboxWorker="toolboxWorker"
            :indexName="item.extend.indexName"
          >
          </IndexData>
        </template>
      </WorkspaceSpans>
    </div>
  </div>
</template>


<script>
import IndexData from "./IndexData";
export default {
  components: { IndexData },
  props: ["source", "toolboxWorker"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    init() {},
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-elasticsearch-tabs {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>

<template>
  <div class="toolbox-elasticsearch-tabs">
    <div class="default-tabs-container">
      <WorkspaceTabs :source="source" :itemsWorker="toolboxWorker.itemsWorker">
      </WorkspaceTabs>
    </div>
    <div class="default-spans-container">
      <WorkspaceSpans :source="source" :itemsWorker="toolboxWorker.itemsWorker">
        <template v-slot:span="{ item }">
          <template v-if="item.extend == null || item.extend.type == 'data'">
            <IndexData
              :source="source"
              :toolboxWorker="toolboxWorker"
              :indexName="item.extend.indexName"
            >
            </IndexData>
          </template>
          <template v-if="item.extend.type == 'import'">
            <Import
              :source="source"
              :toolboxWorker="toolboxWorker"
              :indexName="item.extend.indexName"
              :extend="item.extend"
              :tabId="item.tabId"
              :actived="
                toolboxWorker.itemsWorker.activeItem &&
                item.key == toolboxWorker.itemsWorker.activeItem.key
              "
            >
            </Import>
          </template>
        </template>
      </WorkspaceSpans>
    </div>
  </div>
</template>


<script>
import IndexData from "./IndexData";
import Import from "./Import";

export default {
  components: { Import, IndexData },
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

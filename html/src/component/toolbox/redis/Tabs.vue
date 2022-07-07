<template>
  <div class="toolbox-redis-tabs">
    <div class="default-tabs-container">
      <WorkspaceTabs :source="source" :itemsWorker="toolboxWorker.itemsWorker">
      </WorkspaceTabs>
    </div>
    <div class="default-spans-container">
      <WorkspaceSpans :source="source" :itemsWorker="toolboxWorker.itemsWorker">
        <template v-slot:span="{ item }">
          <template v-if="item.extend.type == 'data'">
            <Data
              :source="source"
              :toolboxWorker="toolboxWorker"
              :extend="item.extend"
            >
            </Data>
          </template>
          <template v-else-if="item.extend.type == 'import'">
            <Import
              :source="source"
              :toolboxWorker="toolboxWorker"
              :extend="item.extend"
            >
            </Import>
          </template>
          <template v-else-if="item.extend.type == 'export'">
            <Export
              :source="source"
              :toolboxWorker="toolboxWorker"
              :extend="item.extend"
            >
            </Export>
          </template>
        </template>
      </WorkspaceSpans>
    </div>
  </div>
</template>


<script>
import Data from "./Data.vue";
import Import from "./Import.vue";
import Export from "./Export.vue";

export default {
  components: { Data, Import, Export },
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
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-redis-tabs {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>

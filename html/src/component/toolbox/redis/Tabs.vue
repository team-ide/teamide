<template>
  <div class="toolbox-redis-tabs">
    <template v-if="ready">
      <TabEditor
        ref="TabEditor"
        :source="source"
        :onRemoveTab="onRemoveTab"
        :onActiveTab="onActiveTab"
      >
        <template v-slot:body="{ tab }">
          <template v-if="tab.extend.type == 'data'">
            <Data
              :source="source"
              :toolbox="toolbox"
              :wrap="wrap"
              :extend="tab.extend"
              :tab="tab"
            >
            </Data>
          </template>
          <template v-else-if="tab.extend.type == 'import'">
            <Import
              :source="source"
              :toolbox="toolbox"
              :wrap="wrap"
              :extend="tab.extend"
              :tab="tab"
            >
            </Import>
          </template>
          <template v-else-if="tab.extend.type == 'export'">
            <Export
              :source="source"
              :toolbox="toolbox"
              :wrap="wrap"
              :extend="tab.extend"
              :tab="tab"
            >
            </Export>
          </template>
        </template>
      </TabEditor>
    </template>
  </div>
</template>


<script>
import Data from "./Data.vue";
import Import from "./Import.vue";
import Export from "./Export.vue";

export default {
  components: { Data, Import, Export },
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
.toolbox-redis-tabs {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>

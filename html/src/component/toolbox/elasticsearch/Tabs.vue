<template>
  <div class="toolbox-elasticsearch-tabs">
    <template v-if="ready">
      <TabEditor
        ref="TabEditor"
        :source="source"
        :onRemoveTab="onRemoveTab"
        :onActiveTab="onActiveTab"
      >
        <template v-slot:body="{ tab }">
          <IndexData
            :source="source"
            :toolbox="toolbox"
            :toolboxType="toolboxType"
            :wrap="wrap"
            :indexName="tab.extend.indexName"
          >
          </IndexData>
        </template>
      </TabEditor>
    </template>
  </div>
</template>


<script>
import IndexData from "./IndexData";
export default {
  components: { IndexData },
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
.toolbox-elasticsearch-tabs {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>

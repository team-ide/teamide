<template>
  <div class="toolbox-kafka-tabs">
    <template v-if="ready">
      <TabEditor
        ref="TabEditor"
        :source="source"
        :onRemoveTab="onRemoveTab"
        :onActiveTab="onActiveTab"
      >
        <template v-slot:body="{ tab }">
          <ToolboxKafkaTopicData
            :source="source"
            :toolbox="toolbox"
            :toolboxType="toolboxType"
            :data="data"
            :wrap="wrap"
            :topic="tab.data"
          >
          </ToolboxKafkaTopicData>
        </template>
      </TabEditor>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
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
      this.wrap.doActiveTab = this.doActiveTab;
      this.wrap.createTab = this.createTab;
      this.wrap.createTabByData = this.createTabByData;
      this.wrap.getTabByData = this.getTabByData;
      this.wrap.addTab = this.addTab;
      this.ready = true;
    },
    getTab(tab) {
      return this.$refs.TabEditor.getTab(tab);
    },
    onRemoveTab(tab) {},
    onActiveTab(tab) {},
    addTab(tab) {
      return this.$refs.TabEditor.addTab(tab);
    },
    doActiveTab(tab) {
      return this.$refs.TabEditor.doActiveTab(tab);
    },
    getTabKeyByData(data) {
      let key = "" + data.name;

      return key;
    },
    getTabByData(data) {
      let key = this.getTabKeyByData(data);
      let tab = this.getTab(key);
      return tab;
    },
    createTabByData(data) {
      let key = this.getTabKeyByData(data);

      let tab = this.getTab(key);
      if (tab == null) {
        let title = data.name;
        let name = data.name;
        tab = {
          key,
          data,
          title,
          name,
        };
        tab.active = false;
        tab.changed = false;
      }
      return tab;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-kafka-tabs {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>

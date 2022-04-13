<template>
  <div class="toolbox-main">
    <TabEditor
      ref="TabEditor"
      :source="source"
      :onRemoveTab="onRemoveTab"
      :onActiveTab="onActiveTab"
    >
      <template v-slot:body="{ tab }">
        <ToolboxEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="tab.toolboxType"
          :data="tab.data"
          :extend="tab.extend"
          :active="tab.active"
        ></ToolboxEditor>
      </template>
    </TabEditor>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "context"],
  data() {
    return {};
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    getTab(tab) {
      return this.$refs.TabEditor.getTab(tab);
    },
    onRemoveTab(tab) {
      this.toolbox.closeOpen(tab.openId);
    },
    onActiveTab(tab) {
      this.toolbox.activeOpen(tab.openId);
    },
    addTab(tab) {
      return this.$refs.TabEditor.addTab(tab);
    },
    doActiveTab(tab) {
      return this.$refs.TabEditor.doActiveTab(tab);
    },
    getTabKeyByData(openData) {
      let key = "" + openData.openId;
      return key;
    },
    getTabByData(openData) {
      let key = this.getTabKeyByData(openData);
      let tab = this.getTab(key);
      return tab;
    },
    createTabByData(openData) {
      let key = this.getTabKeyByData(openData);

      let tab = this.getTab(key);
      if (tab == null) {
        let toolboxType = openData.toolboxType;
        let data = openData.data;
        let extend = openData.extend;
        let title = toolboxType.text + " : " + data.name;
        let name = data.name;

        extend = extend || {};
        this.toolbox.formatExtend(toolboxType, data, extend);
        if (extend.isFTP) {
          title = "FTP : " + data.name;
          name = "FTP : " + data.name;
        } else if (toolboxType.name == "ssh") {
          title = "SSH : " + data.name;
          name = "SSH : " + data.name;
        }
        tab = {
          key,
          data,
          title,
          name,
          toolboxType,
          extend,
          openId: openData.openId,
        };
        tab.active = false;
      }
      return tab;
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.toolbox.doActiveTab = this.doActiveTab;
    this.toolbox.createTab = this.createTab;
    this.toolbox.createTabByData = this.createTabByData;
    this.toolbox.getTabByData = this.getTabByData;
    this.toolbox.addTab = this.addTab;
  },
};
</script>

<style>
.toolbox-main {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>

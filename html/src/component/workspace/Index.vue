<template>
  <div
    class="workspace-container"
    :class="{
      'workspace-theme-dark': theme.isDark,
    }"
    :style="{
      backgroundColor: theme.backgroundColor,
    }"
  >
    <div class="workspace-header"></div>
    <div class="workspace-main">
      <div class="workspace-main-tabs-container">
        <WorkspaceTabs :source="source" :itemsWorker="mainItemsWorker">
          <slot name="mainTabLeftExtend" slot="leftExtend"></slot>
          <slot name="mainTabRightExtend" slot="rightExtend"></slot>
        </WorkspaceTabs>
      </div>
      <div class="workspace-main-spans-container">
        <WorkspaceSpans :source="source" :itemsWorker="mainItemsWorker">
          <template v-slot:span="{ item }">
            <template v-if="item.isToolbox">
              <ToolboxEditor
                :source="source"
                :toolboxType="item.toolboxType"
                :toolboxId="item.toolboxId"
                :openId="item.openId"
                :extend="item.extend"
                :item="item"
              >
              </ToolboxEditor>
            </template>
          </template>
        </WorkspaceSpans>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "onMainRemoveItem", "onMainActiveItem", "toMainCopyItem"],
  data() {
    let mainItemsWorker = this.tool.newItemsWorker({
      onRemoveItem: this.onMainRemoveItem,
      onActiveItem: this.onMainActiveItem,
      toCopyItem: this.toMainCopyItem,
    });
    return {
      mainItemsWorker: mainItemsWorker,
      theme: {
        isDark: true,
        backgroundColor: "#383838",
      },
    };
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
.workspace-container {
  width: 100%;
  height: 100%;
  margin: 0px;
  padding: 0px;
  position: relative;
}
.workspace-container.workspace-theme-dark {
  color: #d9d9d9;
}
.workspace-header {
  width: 100%;
  height: 0px;
  margin: 0px;
  padding: 0px;
  position: relative;
  border-bottom: 0px solid #4e4e4e;
}
.workspace-main {
  width: 100%;
  height: calc(100% - 0px);
  margin: 0px;
  padding: 0px;
  position: relative;
}
.workspace-main-tabs-container {
  width: 100%;
  height: 30px;
  position: relative;
  border-bottom: 1px solid #4e4e4e;
}
.workspace-main-spans-container {
  width: 100%;
  height: calc(100% - 30px);
  position: relative;
}
.default-tabs-container {
  width: 100%;
  height: 30px;
  position: relative;
  border-bottom: 1px solid #4e4e4e;
}
.default-spans-container {
  width: 100%;
  height: calc(100% - 30px);
  position: relative;
}
</style>

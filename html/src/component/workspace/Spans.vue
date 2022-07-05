<template>
  <div class="workspace-spans">
    <div class="workspace-spans-body">
      <template v-for="one in tabs">
        <div
          :key="one.key"
          class="workspace-spans-one"
          :class="{ active: one == activeTab }"
        >
          <template v-if="one.isToolbox">
            <ToolboxEditor
              :source="source"
              :tab="tab"
              :toolboxType="one.toolboxType"
              :toolboxId="one.toolboxId"
              :openId="one.openId"
              :extend="one.extend"
            >
            </ToolboxEditor>
          </template>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "tabs", "activeTab"],
  data() {
    return {};
  },
  computed: {},
  watch: {
    activeTab() {
      this.onActiveTabFocue();
    },
  },
  methods: {
    init() {},
    onActiveTabFocue() {
      if (this.activeTab) {
        let slot = this.getTabSpanSlot(this.activeTab);
        if (slot == null) {
          return;
        }
        slot.onFocus && slot.onFocus();
      }
    },
    getTabSpanSlot(tab) {
      let refs = this.$refs[`span-${tab == null ? "" : tab.key}`];
      if (refs != null && refs.length > 0) {
        return refs[0];
      }
      return refs;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  destroyed() {},
};
</script>

<style >
.workspace-spans {
  width: 100%;
  height: 100%;
  position: relative;
}

.workspace-spans-body {
  width: 100%;
  height: 100%;
  position: relative;
}
.workspace-spans-one {
  position: absolute;
  width: 100%;
  height: 100%;
  left: 0px;
  right: 0px;
  transform: scale(0);
}
.workspace-spans-one.active {
  transform: scale(1);
}
</style>

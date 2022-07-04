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
    <div class="workspace-main">
      <div class="workspace-main-tabs-container">
        <Tabs
          :source="source"
          :activeTab="mainActiveItem"
          :tabs="mainItems"
          :onActiveTab="toMainActiveItem"
        ></Tabs>
      </div>
      <div class="workspace-main-spans-container">
        <div class="workspace-spans"></div>
      </div>
    </div>
  </div>
</template>

<script>
import Tabs from "./Tabs.vue";

export default {
  components: { Tabs },
  props: ["source"],
  data() {
    return {
      headerItems: [],
      leftItems: [],
      mainItems: [],
      mainActiveItem: null,
      footerItems: [],
      theme: {
        isDark: true,
        backgroundColor: "#383838",
      },
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      for (let i = 0; i < 20; i++) {
        let item = {};
        item.key = "key:item:" + i;
        item.name = "name:item:" + i;
        item.title = "title:item:" + i;
        item.show = true;
        this.mainItems.push(item);
      }
    },
    toMainActiveItem(item) {
      this.mainActiveItem = item;
    },
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
.workspace-main {
  width: 100%;
  height: 100%;
  margin: 0px;
  padding: 0px;
  position: relative;
}
.workspace-main-tabs-container {
  width: 100%;
  height: 25px;
  position: relative;
  border-bottom: 1px solid #4e4e4e;
}
</style>

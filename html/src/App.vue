<template>
  <div id="app">
    <InfoBox :source="source"></InfoBox>
    <SystemInfoBox :source="source"></SystemInfoBox>
    <AlertBox :source="source"></AlertBox>
    <Frame
      :source="source"
      v-if="!source.removeFrame"
      v-show="source.showFrame"
    ></Frame>
    <Login
      :source="source"
      v-if="!source.removeLogin"
      v-show="source.showLogin"
    ></Login>
    <Register
      :source="source"
      v-if="!source.removeRegister"
      v-show="source.showRegister"
    ></Register>
    <Workspace
      :source="source"
      v-if="!source.removeWorkspace"
      v-show="source.showWorkspace"
    ></Workspace>
    <Console
      :source="source"
      v-if="!source.removeConsole"
      v-show="source.showConsole"
    ></Console>
  </div>
</template>

<script>
import source from "@/source";

import Frame from "@/views/frame/Index.vue";

export default {
  components: { Frame },
  props: [],
  data() {
    return { source };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {},
  // 在实例创建完成后被立即调用
  created() {
    this.tool.$bvToast = this.$bvToast;
  },
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {},
  destroyed() {},
};
</script>

<style>
html,
body {
  height: 100%;
  width: 100%;
}
#app {
  height: 100%;
  width: 100%;
}
</style>

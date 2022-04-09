<template>
  <div class="workspace-box">
    <template v-if="ready"> Wellcome Workspace! </template>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {
      ready: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    "$route.path"() {
      this.init();
    },
  },
  methods: {
    init() {
      if (this.ready) {
        return;
      }
      if (!this.tool.isWorkspacePage(this.$route.path)) {
        return;
      }
      this.ready = true;
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {},
};
</script>

<style>
.workspace-box {
  width: 100%;
  height: 100%;
  margin: 0px;
  padding: 0px;
  position: relative;
}
</style>

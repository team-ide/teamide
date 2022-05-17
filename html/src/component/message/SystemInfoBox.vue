<template>
  <div></div>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {};
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {},
  // 在实例创建完成后被立即调用
  created() {
    this.tool.systemInfo = (info) => {
      this.$notify({
        title: "系统提示",
        message: info,
        type: "info",
      });
    };
    this.tool.systemSuccess = (success) => {
      this.$notify({
        title: "系统提示",
        message: success,
        type: "success",
      });
    };

    this.tool.systemWarn = (warn) => {
      this.$notify({
        title: "系统警告",
        message: warn,
        type: "warning",
      });
    };

    this.tool.systemError = (error) => {
      this.$notify({
        title: "系统错误",
        message: error,
        type: "error",
      });
    };
  },
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {},
};
</script>

<style>
</style>

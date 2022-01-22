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
  methods: {
    toast(title, info, variant) {
      let config = {
        title: title,
        toaster: "b-toaster-top-right",
        variant: variant,
        autoHideDelay: 3000,
        appendToast: true,
      };
      this.$bvToast.toast(info, config);
    },
  },
  // 在实例创建完成后被立即调用
  created() {
    this.tool.systemInfo = (info) => {
      this.toast("系统提示", info, "info");
    };

    this.tool.systemWarn = (warn) => {
      this.toast("系统警告", warn, "warning");
    };

    this.tool.systemError = (error) => {
      this.toast("系统错误", error, "danger");
    };
  },
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {},
};
</script>

<style>
</style>

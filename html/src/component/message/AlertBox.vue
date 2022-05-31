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
    this.tool.alert = (msg, okTitle) => {
      return new Promise((resolve, reject) => {
        this.$alert(msg, "提示", {
          confirmButtonText: okTitle || "确定",
          callback: (action) => {
            resolve(action);
          },
        });
      });
    };
    this.tool.confirm = (msg, okTitle, cancelTitle) => {
      return this.$confirm(msg, "提示", {
        confirmButtonText: okTitle || "确定",
        cancelButtonText: cancelTitle || "取消",
      });
    };
    this.tool.message = (msg, okTitle, cancelTitle) => {
      return new Promise((resolve, reject) => {
        this.$msgbox({
          title: "提示",
          message: msg,
          showCancelButton: true,
          confirmButtonText: okTitle || "确定",
          cancelButtonText: cancelTitle || "取消",
          beforeClose: (action, instance, done) => {
            if (action === "confirm") {
              done();
              resolve(action);
            } else {
              done();
              reject();
            }
          },
        });
      });
    };
  },
  mounted() {},
};
</script>

<style>
</style>

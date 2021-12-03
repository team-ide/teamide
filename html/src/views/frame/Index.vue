<template>
  <div class="frame-box">{{ source }}</div>
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
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.tool
      .alert("这是提示")
      .then(() => {
        console.log("确认");
      })
      .catch(() => {
        console.log("取消");
      });
    this.tool
      .confirm("这是确认")
      .then(() => {
        console.log("确认");
      })
      .catch(() => {
        console.log("取消");
      });
  },
};
</script>

<style>
.frame-box {
  width: 100%;
  height: 100%;
}
</style>

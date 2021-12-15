<template>
  <div class="frame-body-box">
    <template v-if="tool.hasFrame($route.path)">
      <div class="frame-body">
        <div class="frame-body-left"></div>
        <div class="frame-body-center">
          <router-view :source="source"></router-view>
        </div>
        <div class="frame-body-right"></div>
      </div>
    </template>
    <template v-else>
      <router-view :source="source"></router-view>
    </template>
  </div>
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
  mounted() {},
};
</script>

<style>
.frame-body-box {
  width: 100%;
  height: calc(100% - 56px);
  background: #ddd;
  display: inline-block;
}
.frame-body {
  min-width: 1024px;
  max-width: 1400px;
  margin: 10px auto;
  border: 1px solid #ddd;
  min-height: calc(100% - 30px);
  background: #fff;
  display: flex;
}
.frame-body-left {
  width: 280px;
  border-right: 1px solid #ddd;
}
.frame-body-center {
  flex: 1;
}
.frame-body-right {
  width: 280px;
  border-left: 1px solid #ddd;
}
.page-box {
  width: 100%;
  height: 100%;
}
.page-header {
  padding: 10px 20px;
  border-bottom: 1px solid #ddd;
  font-weight: 600;
  font-size: 22px;
}
.page-body {
  padding: 10px 20px;
}
</style>

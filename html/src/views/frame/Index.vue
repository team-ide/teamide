<template>
  <div class="frame-box">
    <Header :source="source"></Header>
    <Body :source="source"></Body>
    <Footer :source="source"></Footer>
  </div>
</template>

<script>
import Header from "@/views/frame/Header.vue";
import Body from "@/views/frame/Body.vue";
import Footer from "@/views/frame/Footer.vue";

export default {
  components: { Header, Body, Footer },
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
.frame-box {
  width: 100%;
  height: 100%;
}
</style>

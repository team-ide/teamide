<script>
export default {
  components: {},
  props: ["size", "pointer"],
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

<template  >
  <i
    class="mdi-icon mdi"
    :style="{ 'font-size': `${size}px`, cursor: pointer ? 'auto' : 'pointer' }"
  >
    <slot></slot>
  </i>
</template>

<style>
</style>

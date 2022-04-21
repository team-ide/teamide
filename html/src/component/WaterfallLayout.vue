<template>
  <div class="waterfall-layout-box">
    <slot></slot>
  </div>
</template>

<script>
export default {
  components: {},
  props: {
    gap: {
      type: Number,
      default: 12,
    },
    columns: {
      type: Number,
      default: 2,
    },
    startLeft: {
      type: Number,
      default: 0,
    },
    startTop: {
      type: Number,
      default: 0,
    },
  },
  data() {
    return {};
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    // this.doLayout();
  },
  methods: {
    calcPositions({ elements }) {
      if (!elements || !elements.length) {
        return [];
      }
      const y = [];
      const positions = [];
      elements.forEach((item, index) => {
        if (y.length < this.columns) {
          y.push(item.offsetHeight);
          positions.push({
            left:
              (index % this.columns) * (item.offsetWidth + this.gap) +
              this.startLeft,
            top: 0 + this.startTop,
          });
        } else {
          const min = Math.min(...y);
          const idx = y.indexOf(min);
          y.splice(idx, 1, min + this.gap + item.offsetHeight);
          positions.push({
            left: idx * (item.offsetWidth + this.gap) + this.startLeft,
            top: min + this.gap + this.startTop,
          });
        }
      });
      return {
        positions,
        containerHeight:
          positions[positions.length - 1].top +
          elements[elements.length - 1].offsetHeight,
      };
    },
    doLayout() {
      const children = [...this.$el.querySelectorAll(".waterfall-layout-item")];
      if (children.length === 0) {
        return;
      }
      const { positions, containerHeight } = this.calcPositions({
        elements: children,
      });
      children.forEach((item, index) => {
        item.style.cssText = `left:${positions[index].left}px;top:${positions[index].top}px;`;
      });
      this.$el.style.height = `${containerHeight}px`;
    },
  },
};
</script>

<style>
.waterfall-layout-box {
  position: relative;
}
.waterfall-layout-item {
  position: absolute;
}
</style>

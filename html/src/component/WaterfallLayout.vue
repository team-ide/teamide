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
    marginBoxBottom: {
      type: Number,
      default: 0,
    },
  },
  data() {
    return {};
  },
  computed: {},
  watch: {},
  created() {},
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
          elements[elements.length - 1].offsetHeight +
          this.marginBoxBottom,
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

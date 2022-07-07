<template>
  <div class="workspace-spans">
    <div class="workspace-spans-body">
      <template v-for="one in itemsWorker.items">
        <div
          :key="one.key"
          class="workspace-spans-one"
          :class="{ active: one == itemsWorker.activeItem }"
        >
          <slot name="span" :item="one" :ref="`span-${one.key}`"></slot>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "itemsWorker"],
  data() {
    return {};
  },
  computed: {},
  watch: {
    "itemsWorker.activeItem"() {
      this.$nextTick(() => {
        this.onActiveItemFocue();
      });
    },
  },
  methods: {
    init() {},
    onActiveItemFocue() {
      if (this.itemsWorker.activeItem) {
        let slot = this.getItemSpanSlot(this.itemsWorker.activeItem);
        if (slot == null) {
          return;
        }
        slot.onFocus && slot.onFocus();
      }
    },
    getItemSpanSlot(item) {
      let index = this.itemsWorker.items.indexOf(item);
      let $vue = this.$children[index];
      return $vue;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  destroyed() {},
};
</script>

<style >
.workspace-spans {
  width: 100%;
  height: 100%;
  position: relative;
}

.workspace-spans-body {
  width: 100%;
  height: 100%;
  position: relative;
}
.workspace-spans-one {
  position: absolute;
  width: 100%;
  height: 100%;
  left: 0px;
  right: 0px;
  transform: scale(0);
}
.workspace-spans-one.active {
  transform: scale(1);
}
</style>

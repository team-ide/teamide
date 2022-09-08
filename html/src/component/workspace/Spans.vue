<template>
  <div class="workspace-spans">
    <div class="workspace-spans-body">
      <template v-for="one in itemsWorker.items">
        <div
          :key="`span-${one.key}`"
          class="workspace-spans-one"
          :class="{
            active:
              itemsWorker.activeItem && one.key == itemsWorker.activeItem.key,
          }"
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
        let find = null;
        this.$children.forEach((one) => {
          let el = this.tool.jQuery(one.$el).parent();
          if (el.hasClass("active")) {
            find = one;
          }
        });
        if (find == null) {
          return;
        }
        find.onFocus && find.onFocus();
      }
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {},
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

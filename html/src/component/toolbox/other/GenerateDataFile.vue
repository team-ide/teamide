<template>
  <div class="toolbox-generateDataFile-editor">
    <template v-if="ready">
      <div class="text-center color-orange pdtb-20">功能开发中，敬请期待！</div>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxType", "toolbox", "extend", "wrap"],
  data() {
    return {
      ready: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.ready = true;
    },
    refresh() {
      this.$nextTick(() => {});
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-generateDataFile-editor {
  width: 100%;
  height: 100%;
}
</style>

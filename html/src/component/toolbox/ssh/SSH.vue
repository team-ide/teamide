<template>
  <div class="toolbox-ssh-editor">
    <Terminal
      :source="source"
      :toolboxWorker="toolboxWorker"
      place="ssh"
      :placeId="`${toolboxWorker.toolboxId}`"
    ></Terminal>
  </div>
</template>


<script>
import Terminal from "../terminal/Terminal.vue";

export default {
  components: { Terminal },
  props: ["source", "extend", "toolboxWorker"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    async init() {},
    onFocus() {
      this.$children.forEach((one) => {
        one.onFocus && one.onFocus();
      });
    },
    refresh() {
      this.$children.forEach((one) => {
        one.refresh && one.refresh();
      });
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    this.isDestroyed = true;
  },
};
</script>

<style>
.toolbox-ssh-editor {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>

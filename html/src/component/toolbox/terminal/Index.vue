<template>
  <div class="toolbox-terminal-editor">
    <Terminal
      :source="source"
      :toolboxWorker="toolboxWorker"
      :place="place"
      :placeId="placeId"
      :extend="extend"
    >
    </Terminal>
  </div>
</template>


<script>
import Terminal from "./Terminal.vue";

export default {
  components: { Terminal },
  props: ["source", "toolboxWorker", "extend"],
  data() {
    let extend = this.extend || {};
    return {
      ready: false,
      place: extend.place,
      placeId: extend.placeId,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {},
    refresh() {},
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-terminal-editor {
  width: 100%;
  height: 100%;
}
</style>

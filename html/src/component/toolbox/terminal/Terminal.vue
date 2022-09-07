<template>
  <div class="toolbox-terminal-box">
    <tm-layout height="100%"> </tm-layout>
  </div>
</template>


<script>
import _worker from "./worker.js";

export default {
  components: {},
  props: ["source", "toolboxWorker", "place", "placeId"],
  data() {
    let worker = _worker.newWorker({
      workerId: this.toolboxWorker.workerId,
      place: this.place,
      placeId: this.placeId,
    });
    return {
      worker: worker,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
    },
    refresh() {},
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-terminal-box {
  width: 100%;
  height: 100%;
}
</style>

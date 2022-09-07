<template>
  <div class="toolbox-terminal-editor">
    <tm-layout height="100%"> </tm-layout>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      place: null,
      placeId: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.place = this.extend.place;
      this.placeId = this.extend.placeId;
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
.toolbox-terminal-editor {
  width: 100%;
  height: 100%;
}
</style>

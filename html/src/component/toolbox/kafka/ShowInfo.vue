<template>
  <div></div>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolboxWorker"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    async show() {
      let infoData = await this.loadInfo();
      this.toolboxWorker.showJSONData(infoData);
    },
    hide() {
      this.showDialog = false;
    },
    async loadInfo() {
      let param = {};
      let res = await this.toolboxWorker.work("info", param);
      res.data = res.data || {};
      return res.data.info;
    },
    init() {},
  },
  created() {},
  mounted() {
    this.toolboxWorker.showInfo = this.show;
    this.init();
  },
};
</script>

<style>
</style>

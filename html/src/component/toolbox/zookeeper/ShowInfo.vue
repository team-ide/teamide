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
    async loadInfo() {
      let param = this.toolboxWorker.getWorkParam({});
      let res = await this.server.zookeeper.info(param);
      res.data = res.data || {};
      return res.data;
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

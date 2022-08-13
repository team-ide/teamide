<template>
  <el-dialog
    ref="modal"
    title="信息查看"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="1000px"
  >
    <div class="ft-15">
      <el-input
        type="textarea"
        v-model="info"
        :autosize="{ minRows: 10, maxRows: 25 }"
      >
      </el-input>
    </div>
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolboxWorker"],
  data() {
    return {
      showDialog: false,
      info: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async show() {
      this.info = await this.loadInfo();
      this.showDialog = true;
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
    this.toolboxWorker.hideInfo = this.hide;
    this.init();
  },
};
</script>

<style>
</style>

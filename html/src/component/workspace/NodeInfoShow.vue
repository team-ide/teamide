<template>
  <el-dialog
    ref="modal"
    :title="`节点查看`"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :show-close="true"
    :append-to-body="true"
    :visible="showBox"
    :before-close="hide"
    :fullscreen="true"
    width="100%"
    class="node-info-dialog"
  >
    <tm-layout height="100%">
      <tm-layout height="150px" class="scrollbar">
        <div class="node-info-dialog-header">
          {{ node }}
        </div>
      </tm-layout>
      <tm-layout height="auto" class="">
        <div class="node-info-dialog-body">
          {{ node }}
        </div>
      </tm-layout>
    </tm-layout>
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {
      showBox: false,
      node: null,
      loading: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    show(node) {
      this.node = node;
      this.showBox = true;
    },
    hide() {
      this.showBox = false;
    },
    init() {},
  },
  created() {},
  updated() {
    this.tool.showNodeInfo = this.show;
    this.tool.hideNodeInfo = this.hide;
  },
  mounted() {
    this.init();
    this.tool.showNodeInfo = this.show;
    this.tool.hideNodeInfo = this.hide;
  },
};
</script>

<style>
.node-info-dialog {
  width: 100%;
  height: 100%;
  user-select: text;
}
.node-info-dialog .el-dialog {
  background: #0f1b26;
  color: #ffffff;
  position: absolute;
  top: 30px;
  bottom: 30px;
  left: 30px;
  right: 30px;
  width: auto;
  height: auto;
}
.node-info-dialog .el-dialog__title {
  color: #ffffff;
}
.node-info-dialog .el-dialog__body {
  position: relative;
  width: 100%;
  height: calc(100% - 55px);
  padding: 0px;
}
.node-info-dialog-header {
  position: relative;
}
.node-info-dialog-body {
  position: relative;
  width: 100%;
  height: 100%;
}
</style>

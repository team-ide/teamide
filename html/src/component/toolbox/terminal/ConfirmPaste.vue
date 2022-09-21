<template>
  <el-dialog
    ref="modal"
    :title="`确认粘贴`"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="900px"
  >
    <div>
      <el-input
        type="textarea"
        v-model="text"
        :readonly="true"
        :autosize="{ minRows: 10, maxRows: 25 }"
      >
      </el-input>
    </div>
    <div class="pdt-20">
      <div class="tm-btn bg-green ft-13 mgr-10" @click="doConfirm">确认</div>
      <div class="tm-btn bg-grey ft-13" @click="hide">取消</div>
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
      text: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    show(text, confirm, cancel) {
      this.text = text;
      this.confirm = confirm;
      this.cancel = cancel;
      this.showDialog = true;
    },
    doConfirm() {
      this.showDialog = false;
      this.confirm && this.confirm();
    },
    doCancel() {
      this.hide();
    },
    hide() {
      this.showDialog = false;
      this.cancel && this.cancel();
    },
  },
  created() {},
  mounted() {
    this.toolboxWorker.showConfirmPaste = this.show;
    this.toolboxWorker.hideConfirmPaste = this.show;
  },
};
</script>

<style>
</style>

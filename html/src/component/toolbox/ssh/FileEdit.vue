<template>
  <el-dialog
    ref="modal"
    :title="`编辑文件：${path}`"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="1200px"
  >
    <div class="mgt--20 toolbox-ftp-edit-file">
      <template v-if="error != null">
        <div class="bg-red ft-12 pd-5">{{ error }}</div>
      </template>
      <template v-else>
        <textarea v-model="text" class="toolbox-ftp-edit-file-textarea">
        </textarea>
      </template>
      <div class="pdtb-10">
        <div
          class="tm-btn bg-green ft-13"
          @click="toSave"
          :class="{ 'tm-disabled': saveIng }"
        >
          保存
        </div>
      </div>
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
      path: null,
      text: null,
      place: null,
      error: null,
      saveIng: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async show(place, file) {
      this.place = place;
      this.file = file;
      this.path = file.path;
      await this.toLoad();
      this.showDialog = true;
    },
    hide() {
      this.showDialog = false;
    },
    async toLoad() {
      let param = {
        place: this.place,
        path: this.path,
      };
      this.error = null;
      this.text = null;
      let res = await this.toolboxWorker.work("readText", param);
      if (res.code != 0) {
        this.error = res.msg || res.error;
        return;
      }
      res.data = res.data || {};
      let response = res.data.response || {};
      this.text = response.text;
    },
    async toSave() {
      this.saveIng = true;
      let param = {
        place: this.place,
        path: this.path,
        text: this.text,
      };
      let res = await this.toolboxWorker.work("saveText", param);
      if (res.code == 0) {
        this.tool.success("保存成功!");
      }
      this.saveIng = false;
    },
    init() {},
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-ftp-edit-file-textarea {
  width: 100%;
  height: 500px;
  letter-spacing: 1px;
  word-spacing: 5px;
  word-break: break-all;
  font-size: 12px;
  border: 1px solid #ddd;
  padding: 0px 5px;
  outline: none;
  user-select: none;
  resize: none;
}
</style>

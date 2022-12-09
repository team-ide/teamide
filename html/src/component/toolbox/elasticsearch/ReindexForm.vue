<template>
  <el-dialog
    ref="modal"
    :title="`索引[${indexName}]迁移`"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="700px"
    top="40px"
  >
    <div class="">
      <el-form ref="form" size="mini" @submit.native.prevent>
        <el-form-item label="源索引">
          <el-input v-model="sourceIndexName"> </el-input>
        </el-form-item>
        <el-form-item label="迁移至索引">
          <el-input v-model="destIndexName"> </el-input>
        </el-form-item>
      </el-form>
    </div>
    <div class="">
      <div
        class="tm-btn bg-teal-8 ft-18 pdtb-5"
        :class="{ 'tm-disabled': saveBtnDisabled }"
        @click="doSave"
      >
        迁移
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
      indexName: null,
      sourceIndexName: null,
      destIndexName: null,
      saveBtnDisabled: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    show(data, callback) {
      data = data || {};

      this.indexName = data.indexName;
      this.sourceIndexName = data.indexName;
      this.destIndexName = data.indexName + "-xxx";

      this.callback = callback;
      this.showDialog = true;
    },
    hide() {
      this.showDialog = false;
    },
    async doSave() {
      let mapping = null;
      try {
        mapping = JSON.parse(this.mappingValue);
      } catch (e) {
        try {
          mapping = eval("(" + this.mappingValue + ")");
        } catch (error2) {
          this.tool.error("请输入有效JSON:" + e.toString());
          return;
        }
      }
      if (this.tool.isEmpty(this.sourceIndexName)) {
        this.tool.error("请输入源索引名称");
        return;
      }
      if (this.tool.isEmpty(this.destIndexName)) {
        this.tool.error("请输入迁移至索引名称");
        return;
      }

      this.saveBtnDisabled = true;

      let param = {
        sourceIndexName: this.sourceIndexName,
        destIndexName: this.destIndexName,
      };
      let flag = await this.callback(param);
      this.saveBtnDisabled = false;
      if (flag) {
        this.hide();
      }
    },
    init() {},
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.toolboxWorker.showReindexForm = this.show;
    this.toolboxWorker.hideReindexForm = this.hide;
    this.init();
  },
};
</script>

<style>
</style>

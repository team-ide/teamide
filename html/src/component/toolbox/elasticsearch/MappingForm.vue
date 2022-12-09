<template>
  <el-dialog
    ref="modal"
    :title="`索引[${indexName}]结构`"
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
        <el-form-item label="结构">
          <el-input
            type="textarea"
            v-model="mappingValue"
            :autosize="{ minRows: 5, maxRows: 10 }"
          >
          </el-input>
        </el-form-item>
        <template v-if="mappingJSON != null">
          <el-form-item label="结构JSON预览">
            <el-input
              type="textarea"
              v-model="mappingJSON"
              :autosize="{ minRows: 5, maxRows: 10 }"
            >
            </el-input>
          </el-form-item>
        </template>
      </el-form>
    </div>
    <div class="">
      <div
        class="tm-btn bg-teal-8 ft-18 pdtb-5 tm-btn-block"
        :class="{ 'tm-disabled': saveBtnDisabled }"
        @click="doSave"
      >
        保存
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
      mappingValue: null,
      mappingJSON: null,
      saveBtnDisabled: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    mappingValue(value) {
      this.mappingJSON = null;
      if (this.tool.isNotEmpty(value)) {
        try {
          let data = null;
          try {
            data = JSON.parse(value);
          } catch (error) {
            try {
              data = eval("(" + value + ")");
            } catch (error2) {
              throw error;
            }
          }
          this.mappingJSON = JSON.stringify(data, null, "    ");
        } catch (e) {
          this.mappingJSON = e.toString();
        }
      }
    },
  },
  methods: {
    show(data, callback) {
      data = data || {};

      this.indexName = data.indexName;
      this.mapping = data.mapping || {};
      this.mappingValue = JSON.stringify(this.mapping, null, "    ");

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

      this.saveBtnDisabled = true;

      let param = {
        indexName: this.indexName,
        mapping: mapping,
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
    this.toolboxWorker.showMappingForm = this.show;
    this.toolboxWorker.hideMappingForm = this.hide;
    this.init();
  },
};
</script>

<style>
</style>

<template>
  <el-dialog
    ref="modal"
    title="创建索引"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="900px"
  >
    <Form :source="source" ref="formBuild"> </Form>
    <el-form ref="form" size="mini" @submit.native.prevent>
      <el-form-item label="结构">
        <el-input
          type="textarea"
          v-model="mappingValue"
          :autosize="{ minRows: 5, maxRows: 20 }"
        >
        </el-input>
      </el-form-item>
      <template v-if="mappingJSON != null">
        <el-form-item label="结构JSON预览">
          <el-input
            type="textarea"
            v-model="mappingJSON"
            :autosize="{ minRows: 5, maxRows: 20 }"
          >
          </el-input>
        </el-form-item>
      </template>
    </el-form>
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
      formBuild: null,
      formData: null,
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
          this.mappingJSON = JSON.stringify(data, null, "  ");
        } catch (e) {
          this.mappingJSON = e;
        }
      }
    },
  },
  methods: {
    show(data, callback) {
      data = data || {};
      let mapping = data.mapping || {
        settings: {
          number_of_shards: 1,
          number_of_replicas: 0,
        },
        mappings: {
          properties: {
            title: { type: "text" },
          },
        },
      };
      this.mappingValue = JSON.stringify(mapping, null, "  ");

      this.formBuild = this.form.build(
        this.form.toolboxOption.elasticsearch.index
      );
      let formData = this.formBuild.newDefaultData();
      for (let key in formData) {
        formData[key] = data[key];
      }
      this.formData = formData;

      this.callback = callback;
      this.showDialog = true;

      this.$nextTick(() => {
        this.$refs.formBuild.build(this.formBuild, this.formData);
      });
    },
    hide() {
      this.showDialog = false;
    },
    doSave() {
      let mapping = null;
      try {
        mapping = JSON.parse(this.mappingValue);
      } catch (e) {
        try {
          mapping = eval("(" + this.mappingValue + ")");
        } catch (error2) {
          this.tool.error("请输入有效JSON:" + e);
          return;
        }
      }

      this.saveBtnDisabled = true;
      this.formBuild.validate(this.formData).then(async (res) => {
        if (res.valid) {
          let param = {};
          Object.assign(param, this.formData);
          param.mapping = mapping;
          let flag = await this.callback(param);
          this.saveBtnDisabled = false;
          if (flag) {
            this.hide();
          }
        } else {
          this.saveBtnDisabled = false;
        }
      });
    },
    init() {},
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.toolboxWorker.showIndexForm = this.show;
    this.toolboxWorker.hideIndexForm = this.hide;
    this.init();
  },
};
</script>

<style>
</style>

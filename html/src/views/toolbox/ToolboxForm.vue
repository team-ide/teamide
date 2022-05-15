<template>
  <el-dialog
    ref="modal"
    :title="toolboxType == null ? '' : toolboxType.text"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="800px"
  >
    <Form :source="source" ref="formBuild" class=""> </Form>
    <div class="pdb-10 ft-16 color-grey">配置</div>
    <Form :source="source" ref="formOptionBuild" class=""> </Form>
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
  props: ["source", "toolbox"],
  data() {
    return {
      showDialog: false,
      formBuild: null,
      formData: null,
      formOptionBuild: null,
      formOptionData: null,
      saveBtnDisabled: false,
      toolboxType: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    show(toolboxType, data, callback) {
      data = data || {};
      this.formBuild = this.form.build(this.form.toolbox);
      let formData = this.formBuild.newDefaultData();
      if (Object.keys(data).length > 0) {
        for (let key in formData) {
          formData[key] = data[key];
        }
      }
      this.formData = formData;

      let option = {};
      if (this.tool.isNotEmpty(data.option)) {
        option = JSON.parse(data.option);
      }

      this.formOptionBuild = this.form.build(toolboxType.configForm);
      let formOptionData = this.formOptionBuild.newDefaultData();
      if (Object.keys(option).length > 0) {
        for (let key in formOptionData) {
          formOptionData[key] = option[key];
        }
      }
      this.formOptionData = formOptionData;

      this.toolboxType = toolboxType;
      this.callback = callback;

      this.showDialog = true;

      this.$nextTick(() => {
        this.$refs.formBuild.build(this.formBuild, this.formData);
        this.$refs.formOptionBuild.build(
          this.formOptionBuild,
          this.formOptionData
        );
      });
    },
    hide() {
      this.showDialog = false;
    },
    doSave() {
      this.saveBtnDisabled = true;
      this.formBuild.validate(this.formData).then(async (res) => {
        if (res.valid) {
          this.formOptionBuild
            .validate(this.formOptionData)
            .then(async (res) => {
              if (res.valid) {
                let param = {};
                Object.assign(param, this.formData);
                param.option = JSON.stringify(this.formOptionData);
                let flag = await this.callback(this.toolboxType, param);
                this.saveBtnDisabled = false;
                if (flag) {
                  this.hide();
                }
              } else {
                this.saveBtnDisabled = false;
              }
            });
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
    this.toolbox.showToolboxForm = this.show;
    this.toolbox.hideToolboxForm = this.hide;
    this.init();
  },
};
</script>

<style>
</style>

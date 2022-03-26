<template>
  <b-modal
    ref="modal"
    :title="toolboxType == null ? '' : toolboxType.text"
    :hide-header-close="false"
    :no-close-on-backdrop="true"
    :no-close-on-esc="true"
    :hide-backdrop="true"
    hide-footer
  >
    <Form
      v-if="formBuild != null"
      :form="formBuild"
      :formData="formData"
      class=""
    >
    </Form>
    <div class="pdb-10 ft-16 color-grey">配置</div>
    <Form
      v-if="formOptionBuild != null"
      :form="formOptionBuild"
      :formData="formOptionData"
      class=""
    >
    </Form>
    <div class="">
      <div
        class="tm-btn bg-teal-8 ft-18 pdtb-5 tm-btn-block"
        :class="{ 'tm-disabled': saveBtnDisabled }"
        @click="doSave"
      >
        保存
      </div>
    </div>
  </b-modal>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox"],
  data() {
    return {
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
      for (let key in formData) {
        formData[key] = data[key];
      }
      this.formData = formData;

      let option = {};
      if (this.tool.isNotEmpty(data.option)) {
        option = JSON.parse(data.option);
      }

      this.formOptionBuild = this.form.build(
        this.form.toolboxOption[toolboxType.name]
      );
      let formOptionData = this.formOptionBuild.newDefaultData();
      if (Object.keys(option).length > 0) {
        for (let key in formOptionData) {
          formOptionData[key] = option[key];
        }
      }
      this.formOptionData = formOptionData;

      this.toolboxType = toolboxType;
      this.callback = callback;
      this.$refs["modal"].show();
    },
    hide() {
      this.$refs["modal"].hide();
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
                  this.$refs["modal"].hide();
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

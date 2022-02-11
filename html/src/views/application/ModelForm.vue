<template>
  <b-modal
    ref="modal"
    :title="modelType == null ? '' : modelType.text"
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
      class="pd-10"
    >
      <div class="pdtb-10">
        <div
          class="tm-btn bg-teal-8 ft-18 pdtb-5 tm-btn-block"
          :class="{ 'tm-disabled': saveBtnDisabled }"
          @click="doSave"
        >
          保存
        </div>
      </div>
    </Form>
  </b-modal>
</template>

<script>
export default {
  components: {},
  props: ["source", "application"],
  data() {
    return {
      formBuild: null,
      formData: null,
      saveBtnDisabled: false,
      modelType: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    show(modelType, data, callback) {
      this.formBuild = this.form.build(this.form.model);
      let formData = this.formBuild.newDefaultData();
      data = data || {};
      for (let key in formData) {
        formData[key] = data[key];
      }
      this.formData = formData;
      this.modelType = modelType;
      this.callback = callback;
      this.$refs["modal"].show();
    },
    hide() {
      this.$refs["modal"].hide();
    },
    doSave() {
      this.saveBtnDisabled = true;
      this.formBuild.validate(this.formData).then((res) => {
        if (res.valid) {
          let param = {};
          Object.assign(param, this.formData);
          let flag = this.callback(this.modelType, param);
          this.saveBtnDisabled = false;
          if (flag) {
            this.$refs["modal"].hide();
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
    this.application.showModelForm = this.show;
    this.application.hideModelForm = this.hide;
    this.init();
  },
};
</script>

<style>
</style>

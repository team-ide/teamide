<template>
  <b-modal
    ref="modal"
    :title="'Topic:' + topic + ':推送'"
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
    <div class="">
      <div
        class="tm-btn bg-teal-8 ft-18 pdtb-5 tm-btn-block"
        :class="{ 'tm-disabled': saveBtnDisabled }"
        @click="doSave"
      >
        推送
      </div>
    </div>
  </b-modal>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "wrap"],
  data() {
    return {
      topic: null,
      formBuild: null,
      formData: null,
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
      this.topic = data.topic;

      this.formBuild = this.form.build(
        this.form.toolboxOption.elasticsearch.data
      );
      let formData = this.formBuild.newDefaultData();
      formData.topic = data.topic;
      if (data.keyType) {
        formData.keyType = data.keyType;
      }
      if (data.valueType) {
        formData.valueType = data.valueType;
      }
      this.formData = formData;

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
          let param = {};
          Object.assign(param, this.formData);
          let flag = await this.callback(param);
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
    this.wrap.showPushForm = this.show;
    this.wrap.hidePushForm = this.hide;
    this.init();
  },
};
</script>

<style>
</style>

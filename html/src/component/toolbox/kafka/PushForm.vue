<template>
  <el-dialog
    ref="modal"
    :title="'Topic:' + topic + ':推送'"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="700px"
    top="40px"
  >
    <Form :source="source" ref="formBuild"> </Form>
    <div class="">
      <div
        class="tm-btn bg-teal-8 ft-18 pdtb-5 tm-btn-block"
        :class="{ 'tm-disabled': saveBtnDisabled }"
        @click="doSave"
      >
        推送
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

      this.formBuild = this.form.build(this.form.toolboxOption.kafka.push);
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
      this.showDialog = true;

      this.$nextTick(() => {
        this.$refs.formBuild.build(this.formBuild, this.formData);
      });
    },
    hide() {
      this.showDialog = false;
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
    this.toolboxWorker.showPushForm = this.show;
    this.toolboxWorker.hidePushForm = this.hide;
    this.init();
  },
};
</script>

<style>
</style>

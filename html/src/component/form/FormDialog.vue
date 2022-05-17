<template>
  <el-dialog
    ref="modal"
    :title="`${title_ || title}`"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    :width="width || '700px'"
  >
    <FormBox
      :source="source"
      ref="FormBox"
      :onSave="onSave"
      :saveText="saveText"
      :onSuccess="onSuccess"
    >
    </FormBox>
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: [
    "source",
    "title",
    "width",
    "wrap",
    "onSave",
    "saveText",
    "showName",
    "hideName",
  ],
  data() {
    return {
      showDialog: false,
      saveBtnDisabled: false,
      title_: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      if (this.wrap) {
        if (this.showName) {
          this.wrap[this.showName] = this.show;
        }
        if (this.hideName) {
          this.wrap[this.hideName] = this.hide;
        }
      }
    },
    show(options) {
      this.showDialog = true;
      options = options || {};

      this.title_ = options.title;

      let formConfigList = [];
      let formList = options.form || [];
      if (!Array.isArray(formList)) {
        formList = [formList];
      }
      let dataList = options.data || [];
      if (!Array.isArray(dataList)) {
        dataList = [dataList];
      }
      formList.forEach((form, index) => {
        let formConfig = {};
        formConfig.form = form;
        formConfig.data = dataList[index];

        formConfigList.push(formConfig);
      });

      this.$nextTick(() => {
        this.$refs.FormBox.build(formConfigList, options);
      });
    },
    hide() {
      this.showDialog = false;
    },
    onSuccess() {
      this.hide();
    },
    onError() {},
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
</style>

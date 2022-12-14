<template>
  <el-dialog
    ref="modal"
    :title="title || '文案'"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="1200px"
    top="40px"
  >
    <div class="mgt--20">
      <div style="height: 620px !important">
        <Editor
          ref="Editor"
          :source="source"
          :value="text"
          language="html"
        ></Editor>
      </div>
    </div>
    <div class="mgt-10">
      <div
        v-if="onSave != null"
        class="tm-btn bg-teal-8 ft-18 pdtb-5"
        @click="doSave"
      >
        {{ saveText || "保存" }}
      </div>
    </div>
  </el-dialog>
</template>

<script>
var JSONbig = require("json-bigint");
var JSONbigString = JSONbig({});

export default {
  components: {},
  props: ["source"],
  data() {
    return {
      showDialog: false,
      text: null,
      title: null,
      onSave: null,
      onCancel: null,
      saveText: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    async show(data, options) {
      options = options || {};
      this.onSave = options.onSave;
      this.onCancel = options.onCancel;
      this.title = options.title;
      this.saveText = options.saveText;
      this.showDialog = true;
      this.$nextTick(() => {
        this.text = "";
        if (data != null) {
          try {
            if (typeof data == "string") {
              if (data != "") {
                try {
                  let json = JSONbigString.parse(data);
                  this.text = JSONbigString.stringify(json, null, "    ");
                } catch (e) {
                  this.text = data;
                }
              }
            } else {
              this.text = JSONbigString.stringify(data, null, "    ");
            }
          } catch (e) {
            this.text = e.toString();
          }
        }
        this.$refs.Editor.setValue(this.text);
      });
    },
    hide() {
      this.showDialog = false;
      this.onCancel && this.onCancel();
    },
    doSave() {
      this.showDialog = false;
      let text = this.$refs.Editor.getValue();
      var jsonData = null;
      var jsonError = null;
      try {
        jsonData = JSONbigString.parse(text);
      } catch (e) {
        jsonError = e.toString();
      }
      this.onSave &&
        this.onSave({
          text: text,
          jsonData: jsonData,
          jsonError: jsonError,
        });
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {},
};
</script>

<style>
</style>

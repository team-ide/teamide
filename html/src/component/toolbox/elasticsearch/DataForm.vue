<template>
  <el-dialog
    ref="modal"
    :title="`索引[${indexName}]数据`"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="1200px"
    top="40px"
  >
    <div class="mgt--25">
      <el-form ref="form" size="mini" @submit.native.prevent>
        <el-form-item label="ID">
          <el-input type="input" v-model="id" style="width: 800px"> </el-input>
        </el-form-item>
        <div class="tm-row">
          <div class="col-6 pdr-5">
            <el-form-item label="数据">
              <div
                style="
                  height: 520px !important;
                  display: inline-block;
                  width: 100%;
                "
              >
                <Editor
                  ref="docValue"
                  :source="source"
                  language="json"
                  :change="onChangeDocValue"
                ></Editor>
              </div>
            </el-form-item>
          </div>
          <div class="col-6 pdl-5">
            <el-form-item label="数据JSON预览">
              <div
                style="
                  height: 520px !important;
                  display: inline-block;
                  width: 100%;
                "
              >
                <Editor ref="docJSON" :source="source" language="json"></Editor>
              </div>
            </el-form-item>
          </div>
        </div>
      </el-form>
    </div>
    <div class="">
      <div
        class="tm-btn bg-teal-8 ft-18 pdtb-5"
        :class="{ 'tm-disabled': saveBtnDisabled }"
        @click="doSave"
      >
        保存
      </div>
    </div>
  </el-dialog>
</template>

<script>
var JSONbig = require("json-bigint");
var JSONbigString = JSONbig({});

export default {
  components: {},
  props: ["source", "toolboxWorker"],
  data() {
    return {
      showDialog: false,
      indexName: null,
      id: null,
      docValue: null,
      docJSON: null,
      saveBtnDisabled: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    docValue(value) {
      this.docJSON = null;
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
          this.docJSON = JSONbigString.stringify(data, null, "  ");
        } catch (e) {
          this.docJSON = e.toString();
        }
      }
      this.$refs.docJSON.setValue(this.docJSON);
    },
  },
  methods: {
    onChangeDocValue(value) {
      this.docValue = value;
    },
    show(data, mapping, callback) {
      data = data || {};

      this.indexName = data.indexName;
      this.id = data.id;
      let doc = data.doc || {};

      if (mapping && mapping.mappings && mapping.mappings.properties) {
        for (let name in mapping.mappings.properties) {
          if (doc[name] != null) {
            continue;
          }
          let property = mapping.mappings.properties[name];
          doc[name] = null;
          if (property.type == "text") {
            doc[name] = null;
          }
          if (property.properties) {
            let v = {};
            for (let k in property.properties) {
              v[k] = null;
            }
            doc[name] = v;
          }
        }
      }

      this.callback = callback;
      this.showDialog = true;
      this.$nextTick(() => {
        this.docValue = JSONbigString.stringify(doc, null, "  ");
        this.$refs.docValue.setValue(this.docValue);
      });
    },
    hide() {
      this.showDialog = false;
    },
    async doSave() {
      let doc = null;
      try {
        doc = JSON.parse(this.docValue);
      } catch (e) {
        try {
          doc = eval("(" + this.docValue + ")");
        } catch (error2) {
          this.tool.error("请输入有效JSON:" + e);
          return;
        }
      }
      let id = this.id;
      if (this.tool.isEmpty(id)) {
        this.tool.error("请输入ID");
        return;
      }

      this.saveBtnDisabled = true;

      let param = {
        indexName: this.indexName,
        doc: doc,
        id: id,
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
    this.toolboxWorker.showDataForm = this.show;
    this.toolboxWorker.hideDataForm = this.hide;
    this.init();
  },
};
</script>

<style>
</style>

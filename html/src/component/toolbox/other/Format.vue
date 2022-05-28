<template>
  <div class="toolbox-format-editor">
    <tm-layout height="100%">
      <tm-layout width="50%" class="">
        <el-form class="pdt-10 pdlr-10" size="mini" @submit.native.prevent>
          <el-form-item label="格式" class="mgb-5">
            <el-select v-model="fromType" style="width: 100px" @change="change">
              <el-option
                v-for="(one, index) in types"
                :key="index"
                :value="one.value"
                :label="one.text"
                :disabled="one.disabled"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="输入" class="mgb-5">
            <el-input
              type="textarea"
              v-model="from"
              @change="change"
              @input="change"
            >
            </el-input>
          </el-form-item>
        </el-form>
        <template v-if="error != null">
          <div class="color-error pdlr-10">
            异常： <span>{{ error }}</span>
          </div>
        </template>
      </tm-layout>
      <tm-layout-bar right></tm-layout-bar>
      <tm-layout width="auto">
        <el-form class="pdt-10 pdlr-10" size="mini" @submit.native.prevent>
          <el-form-item label="格式" class="mgb-5">
            <el-checkbox-group v-model="toTypes" @change="change">
              <el-checkbox
                v-for="(one, index) in types"
                :key="index"
                :label="one.value"
                :disabled="one.disabled"
              >
                {{ one.text }}
              </el-checkbox>
            </el-checkbox-group>
          </el-form-item>
          <template v-for="(one, index) in tos">
            <el-form-item :key="index" :label="one.toType" class="mgb-5">
              <el-input type="textarea" v-model="one.value"> </el-input>
              <template v-if="one.error != null">
                <div class="color-error pdlr-10">
                  异常： <span>{{ one.error }}</span>
                </div>
              </template>
            </el-form-item>
          </template>
        </el-form>
      </tm-layout>
    </tm-layout>
  </div>
</template>


<script>
import jsYaml from "js-yaml";

export default {
  components: {},
  props: ["source", "toolboxType", "toolbox", "extend", "wrap"],
  data() {
    return {
      from: null,
      fromType: null,
      tos: [],
      toTypes: [],
      types: [
        { text: "JSON", value: "json" },
        { text: "YAML", value: "yaml" },
        { text: "URL", value: "url" },
        { text: "XML", value: "xml", disabled: true },
        { text: "TOML", value: "toml", disabled: true },
      ],
      error: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      let extend = this.extend || {};
      this.from = extend.from;
      this.fromType = extend.fromType;
      this.toTypes = extend.toTypes || [];
      this.format();
    },
    refresh() {
      this.$nextTick(() => {
        this.format();
      });
    },
    change() {
      let extend = this.extend || {};
      extend.from = this.from;
      extend.fromType = this.fromType || "json";
      extend.toTypes = this.toTypes;
      this.wrap.updateExtend(extend);
      this.format();
    },
    format() {
      let tos = [];
      if (this.toTypes) {
        this.toTypes.forEach((toType) => {
          let to = {};
          to.from = this.from;
          to.fromType = this.fromType;
          to.toType = toType;
          to.value = null;
          to.error = null;
          this.formatValue(to);
          tos.push(to);
        });
      }
      this.tos = tos;
    },
    formatValue(data) {
      this.error = null;
      if (data == null || this.tool.isEmpty(data.from)) {
        return;
      }
      let fromData = this.getFromData(data);
      data.value = this.getToValue(data, fromData);
    },
    getFromData(data) {
      this.error = null;

      let fromData = null;

      if (this.tool.isNotEmpty(data.from)) {
        let fromType = ("" + data.fromType).toLowerCase();
        try {
          if (fromType == "json") {
            let json = null;
            try {
              json = JSON.parse(data.from);
            } catch (error) {
              try {
                json = eval("(" + data.from + ")");
              } catch (error2) {
                throw error;
              }
            }
            fromData = json;
          } else if (fromType == "yaml") {
            fromData = jsYaml.load(data.from);
          }
        } catch (e) {
          this.error = e;
        }
      }

      return fromData;
    },
    getToValue(data, fromData) {
      data.error = null;
      if (data == null || this.tool.isEmpty(data.toType) || fromData == null) {
        return null;
      }
      let value = null;
      let toType = ("" + data.toType).toLowerCase();

      try {
        if (toType == "json") {
          value = JSON.stringify(fromData, null, "  ");
        } else if (toType == "yaml") {
          value = jsYaml.dump(fromData);
        }
      } catch (e) {
        data.error = e;
      }
      return value;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-format-editor {
  width: 100%;
  height: 100%;
}

.toolbox-format-editor textarea {
  width: 100%;
  height: 300px !important;
  letter-spacing: 1px;
  word-spacing: 5px;
  padding: 5px;
}
</style>

<template>
  <select
    class="model-select"
    v-model="value"
    @change="onChange"
    :style="{ width: width + 'px' }"
  >
    <option value="">{{ placeholder ? placeholder : "请选择" }}</option>
    <template v-if="items != null">
      <option
        v-for="(item, index) in items"
        :key="index"
        :value="item.value || item.name"
      >
        {{ item.text || item.name }}
        <template v-if="tool.isNotEmpty(item.comment)">
          ({{ item.comment }})
        </template>
      </option>
    </template>
  </select>
</template>

<script>
export default {
  props: [
    "source",
    "context",
    "wrap",
    "bean",
    "name",
    "options",
    "placeholder",
    "isDataTypeOption",
    "isStructOption",
    "isActionOption",
    "isColumnTypeOption",
    "isIndexTypeOption",
    "isDatabaseTypeOption",
  ],
  components: {},
  data() {
    return {
      width: null,
      value: null,
      items: null,
    };
  },
  watch: {
    bean() {
      this.init();
    },
    value() {
      this.initWidth(this.value);
    },
    options() {
      this.initOptions();
    },
    context() {
      if (this.isDataTypeOption) {
        this.initDataTypeOptions();
      } else if (this.isStructOption) {
        this.initStructOptions();
      } else if (this.isActionOption) {
        this.initActionOptions();
      }
    },
  },
  methods: {
    initWidth(value) {
      let text = this.placeholder ? this.placeholder : "请选择";
      let items = this.items || [];
      items.forEach((item) => {
        if (item.value == value || item.name == value) {
          text = item.text || item.name;
          if (this.tool.isNotEmpty(item.comment)) {
            text += "(" + item.comment + ")";
          }
        }
      });

      let compute = this.tool.computeFontSize("" + text, "12px");
      this.width = compute.width + 30;
    },
    onChange() {
      this.$emit("change", this.value);
      this.wrap && this.wrap.onChange(this.bean, this.name, this.value);
    },
    init() {
      this.initOptions();
      let value = this.getBeanValue();
      if (this.tool.isEmpty(value)) {
        value = "";
      }
      this.value = value;
      this.initWidth(this.value);
    },
    getBeanValue() {
      this.bean = this.bean || {};
      if (this.bean[this.name] == null) {
        this.bean[this.name] = null;
      }
      return this.bean[this.name];
    },
    initOptions() {
      if (this.options != null) {
        this.items = this.options;
      } else if (this.isDataTypeOption) {
        this.initDataTypeOptions();
      } else if (this.isColumnTypeOption) {
        this.initColumnTypeOptions();
      } else if (this.isIndexTypeOption) {
        this.initIndexTypeOptions();
      } else if (this.isDatabaseTypeOption) {
        this.initIndexTypeOptions();
      } else if (this.isStructOption) {
        this.initStructOptions();
      } else if (this.isActionOption) {
        this.initActionOptions();
      }
    },
    initActionOptions() {
      this.items = this.source.application.getActionOptions(this.context);
    },
    initStructOptions() {
      this.items = this.source.application.getStructOptions(this.context);
    },
    initDataTypeOptions() {
      this.items = this.source.application.getDataTypeOptions(this.context);
    },
    initColumnTypeOptions() {
      this.items = this.source.application.getColumnTypeOptions(this.context);
    },
    initIndexTypeOptions() {
      this.items = this.source.application.getIndexTypeOptions(this.context);
    },
    initDatabaseTypeOptions() {
      this.items = this.source.application.getDatabaseTypeOptions(this.context);
    },
  },
  mounted() {
    this.init();
  },
  beforeCreate() {},
};
</script>

<style >
.model-select {
  padding: 0px 0px 0px 5px;
  border-color: transparent;
  outline: none;
  background: transparent;
  color: #f9f9f9;
  border-bottom: 1px groove #f9f9f9;
  min-width: 40px;
  text-align: center;
  -moz-appearance: auto;
  -webkit-appearance: auto;
  font-size: 12px;
  line-height: 22px;
  height: 22px;
}
.model-select option {
  background-color: #ffffff;
  color: #3e3e3e;
  text-align: left;
}
</style>

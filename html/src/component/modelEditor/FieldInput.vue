<template>
  <div v-if="ready" class="input">
    <template v-if="field.type == 'switch'">
      <Switch_
        :source="source"
        :context="context"
        :bean="bean"
        :name="field.name"
        :readonly="field.readonly"
        :wrap="wrap"
      ></Switch_>
    </template>
    <template v-else-if="field.type == 'select'">
      <Select_
        :source="source"
        :context="context"
        :bean="bean"
        :name="field.name"
        :options="options"
        :readonly="field.readonly"
        :wrap="wrap"
      ></Select_>
    </template>
    <template v-else>
      <Input_
        :source="source"
        :context="context"
        :bean="bean"
        :name="field.name"
        :readonly="field.readonly"
        :isNumber="field.isNumber"
        :wrap="wrap"
      ></Input_>
    </template>
  </div>
</template>


<script>
import Input_ from "./Input.vue";
import Switch_ from "./Switch.vue";
import Select_ from "./Select.vue";

export default {
  components: { Input_, Switch_, Select_ },
  props: ["source", "context", "wrap", "bean", "field"],
  data() {
    return {
      options: null,
      ready: false,
    };
  },
  computed: {},
  watch: {
    context() {
      if (this.field.isDataTypeOption) {
        this.initDataTypeOptions();
      } else if (this.field.isStructOption) {
        this.initStructOptions();
      } else if (this.field.isActionOption) {
        this.initActionOptions();
      }
    },
  },
  methods: {
    init() {
      this.initOptions();
      this.ready = true;
    },
    initOptions() {
      if (this.field.options != null) {
        this.options = this.field.options;
      } else if (this.field.isDataTypeOption) {
        this.initDataTypeOptions();
      } else if (this.field.isColumnTypeOption) {
        this.initColumnTypeOptions();
      } else if (this.field.isIndexTypeOption) {
        this.initIndexTypeOptions();
      } else if (this.field.isDatabaseTypeOption) {
        this.initIndexTypeOptions();
      } else if (this.field.isStructOption) {
        this.initStructOptions();
      } else if (this.field.isActionOption) {
        this.initActionOptions();
      }
    },
    initActionOptions() {
      let options = [];
      if (this.context.actions) {
        this.context.actions.forEach((one) => {
          options.push({
            value: one.name,
            text: one.name,
            comment: one.comment,
          });
        });
      }
      this.options = options;
    },
    initStructOptions() {
      let options = [];
      if (this.context.structs) {
        this.context.structs.forEach((one) => {
          options.push({
            value: one.name,
            text: one.name,
            comment: one.comment,
          });
        });
      }
      this.options = options;
    },
    initDataTypeOptions() {
      let options = [];
      options.push({ value: "string", text: "字符串" });
      options.push({ value: "int", text: "整形" });
      options.push({ value: "long", text: "长整型" });
      options.push({ value: "boolean", text: "布尔型" });
      options.push({ value: "byte", text: "字节型" });
      options.push({ value: "date", text: "日期" });
      options.push({ value: "short", text: "短整型" });
      options.push({ value: "double", text: "双精度浮点型" });
      options.push({ value: "float", text: "浮点型" });
      options.push({ value: "map", text: "集合" });

      if (this.context.structs) {
        this.context.structs.forEach((one) => {
          options.push({
            value: one.name,
            text: one.name,
            comment: one.comment,
          });
        });
      }
      this.options = options;
    },
    initColumnTypeOptions() {
      let options = [];
      options.push({ value: "varchar", text: "varchar" });
      options.push({ value: "bigint", text: "bigint" });
      options.push({ value: "int", text: "int" });
      options.push({ value: "datetime", text: "datetime" });
      options.push({ value: "number", text: "number" });
      this.options = options;
    },
    initIndexTypeOptions() {
      let options = [];
      options.push({ value: "unique", text: "普通索引" });
      this.options = options;
    },
    initDatabaseTypeOptions() {
      let options = [];
      options.push({ value: "Mysql", text: "MySql" });
      this.options = options;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
</style>

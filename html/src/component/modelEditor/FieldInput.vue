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
      >
      </Switch_>
    </template>
    <template v-else-if="field.type == 'select'">
      <Select_
        :source="source"
        :context="context"
        :bean="bean"
        :name="field.name"
        :options="options"
        :readonly="field.readonly"
        :isDataTypeOption="field.isDataTypeOption"
        :isStructOption="field.isStructOption"
        :isActionOption="field.isActionOption"
        :isColumnTypeOption="field.isColumnTypeOption"
        :isIndexTypeOption="field.isIndexTypeOption"
        :isDatabaseTypeOption="field.isDatabaseTypeOption"
        :wrap="wrap"
      >
      </Select_>
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
      >
      </Input_>
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
  watch: {},
  methods: {
    init() {
      this.ready = true;
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

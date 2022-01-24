<template>
  <select
    class="model-select"
    v-model="value"
    @change="onChange"
    :style="{ width: width + 'px' }"
  >
    <option value="">{{ placeholder ? placeholder : "请选择" }}</option>
    <template v-if="options != null">
      <option
        v-for="(item, index) in options"
        :key="index"
        :value="item.value || item.name"
      >
        {{ item.text || item.name }}
        <template v-if="item.comment"> ( {{ item.comment }} ) </template>
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
  ],
  components: {},
  data() {
    return {
      width: null,
      value: null,
    };
  },
  watch: {
    bean() {
      this.init();
    },
    value() {
      this.initWidth(this.value);
    },
  },
  methods: {
    initWidth(value) {
      let text = this.placeholder ? this.placeholder : "请选择";
      let list = this.options || [];
      list.forEach((option) => {
        if (option.value == value || option.name == value) {
          text = option.text || option.name;
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
  },
  mounted() {
    this.init();
  },
  beforeCreate() {},
};
</script>

<style >
.model-select {
  padding: 2px 0px 1px 5px;
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
}
.model-select option {
  background-color: #ffffff;
  color: #3e3e3e;
  text-align: left;
}
</style>

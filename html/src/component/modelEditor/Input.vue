<template>
  <input
    class="model-input"
    v-model="value"
    :style="{ width: width + 'px' }"
    @change="onChange"
    @input="onInput"
    autocomplete="off"
  />
</template>

<script>
export default {
  props: ["source", "context", "wrap", "bean", "name", "validate", "isNumber"],
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
      if (this.tool.isEmpty(value)) {
        value = "";
      }
      let compute = this.tool.computeFontSize(value, "12px");
      this.width = compute.width + 15;
    },
    onInput() {
      this.initWidth(this.$el.value);
    },
    onChange() {
      if (this.validate) {
        let error = this.validate(this.bean, this.name, this.value);
        if (this.tool.isNotEmpty(error)) {
          this.tool.error(error);
          this.$el.focus();
          return;
        }
      }
      if (this.isNumber) {
        this.$emit("change", Number(this.value));
        this.wrap &&
          this.wrap.onChange(this.bean, this.name, Number(this.value));
      } else {
        this.$emit("change", this.value);
        this.wrap && this.wrap.onChange(this.bean, this.name, this.value);
      }
    },
    init() {
      this.value = this.getBeanValue();
      this.initWidth(this.value);
    },
    bindEvent() {
      this.$el.addEventListener("keydown", (event) => {
        //ctrl+s
        if (event.keyCode === 83 && event.ctrlKey) {
          this.$el.blur();
        }
      });
    },
    getBeanValue() {
      this.bean = this.bean || {};
      return this.bean[this.name];
    },
  },
  mounted() {
    this.init();
    this.bindEvent();
  },
  beforeCreate() {},
};
</script>

<style >
.model-input {
  padding: 0px 5px;
  border-color: transparent;
  outline: none;
  /* background: transparent; */
  color: #f9f9f9;
  /* border-bottom: 1px groove #f9f9f9; */
  min-width: 40px;
  text-align: center;
  font-size: 12px;
  line-height: 22px;
  height: 22px;
  background: #343434;
}
.model-table .model-input {
  width: 100% !important;
}
</style>

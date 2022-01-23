<template>
  <input
    v-model="value"
    :style="{ width: width + 'px' }"
    @change="onChange"
    @input="onInput"
    autocomplete="off"
  />
</template>

<script>
export default {
  props: ["wrap", "bean", "name", "validate"],
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
      let compute = this.tool.computeFontSize(value, "13px");
      this.width = compute.width + 10;
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
      this.wrap && this.wrap.onChange(this.bean, this.name, this.value);
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
      if (this.bean[this.name] == null) {
        this.bean[this.name] = null;
      }
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
</style>

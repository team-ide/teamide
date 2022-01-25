<template>
  <input
    type="checkbox"
    class="model-switch"
    v-model="value"
    @change="onChange"
    @input="onInput"
    autocomplete="off"
    :value="true"
    :checked="tool.isTrue(value)"
  />
</template>

<script>
export default {
  props: ["wrap", "bean", "name", "validate"],
  components: {},
  data() {
    return {
      value: null,
    };
  },
  watch: {
    bean() {
      this.init();
    },
    value() {},
  },
  methods: {
    onInput() {},
    onChange() {
      if (this.validate) {
        let error = this.validate(this.bean, this.name, this.value);
        if (this.tool.isNotEmpty(error)) {
          this.tool.error(error);
          this.$el.focus();
          return;
        }
      }

      this.$emit("change", this.value);
      this.wrap && this.wrap.onChange(this.bean, this.name, this.value);
    },
    init() {
      this.value = this.getBeanValue();
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
.model-switch {
  background: transparent;
  width: 12px;
  height: 22px;
}
</style>

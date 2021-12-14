<template>
  <b-form>
    <template v-if="key != null && formData != null">
      <template v-for="(field, index) in form.fields">
        <b-form-group
          :key="`key-${key}-${field.name}-${index}`"
          :id="`form-${key}-group-${field.name}-${index}`"
          :label="field.label"
          :label-for="`key-${key}-${field.name}-input`"
          :description="field.description"
          :disabled="field.disabled"
          :state="field.valid"
        >
          <b-form-input
            :id="`key-${key}-${field.name}-input`"
            v-model="formData[field.name]"
            :type="field.type"
            :placeholder="field.placeholder"
            :required="field.required"
            :state="field.valid"
          ></b-form-input>
          <b-form-invalid-feedback v-if="field.validMessage">{{
            field.validMessage
          }}</b-form-invalid-feedback>
        </b-form-group>
      </template>
    </template>
    <slot></slot>
    <b-button v-if="_saveShow" variant="primary" @click="doSave">
      {{ saveText || "保存" }}
    </b-button>
  </b-form>
</template>


<script>
export default {
  components: {},
  props: ["source", "form", "formData", "saveShow", "saveText"],
  data() {
    return {
      key: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {
    _saveShow() {
      if (this.saveShow == undefined || this.saveShow == null) {
        return true;
      }
      return this.saveShow;
    },
  },
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    doSave() {
      this.form.validate(this.formData).then((res) => {
        console.log(res);
      });
    },
    init() {
      this.key = this.tool.getNumber();
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.init();
  },
};
</script>

<style>
</style>

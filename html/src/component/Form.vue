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
      <b-button v-if="saveShow" variant="primary" @click="doSave">
        {{ saveText }}
      </b-button>
    </template>
  </b-form>
</template>


<script>
export default {
  components: {},
  props: ["source", "form"],
  data() {
    return {
      formData: null,
      saveText: "保存",
      saveShow: true,
      key: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    doSave() {
      this.form.validate(this.formData).then((res) => {
        console.log(res);
      });
    },
    init() {
      let formData = this.form.newDefaultData();
      this.formData = formData;
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

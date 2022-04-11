<template>
  <b-form v-if="ready">
    <template v-if="key != null && formData != null && fileData != null">
      <template v-for="(field, index) in form.fields">
        <b-form-group
          :key="`key-${key}-${field.name}-${index}`"
          :id="`form-${key}-group-${field.name}-${index}`"
          :label="field.label"
          :description="field.description"
          :disabled="field.disabled"
          :state="field.valid"
        >
          <template v-if="field.type == 'select'">
            <b-form-select
              v-model="formData[field.name]"
              :placeholder="field.placeholder"
              :required="field.required"
              :options="field.options"
              :state="field.valid"
            >
            </b-form-select>
          </template>
          <template v-else-if="field.type == 'switch'"> </template>
          <template v-else-if="field.type == 'file'">
            <div
              class="ft-12 pdb-5"
              v-if="tool.isNotEmpty(formData[field.name])"
            >
              <span class="color-grey">文件：</span>
              <a
                class="tm-link color-green"
                :href="source.filesUrl + formData[field.name]"
                >{{ formData[field.name] }}</a
              >
            </div>
            <b-form-file
              v-model="fileData[field.name]"
              :placeholder="field.placeholder"
              :required="field.required"
              :state="field.valid"
              browse-text="选择文件"
              @input="fileInput(field.name, $event)"
            >
            </b-form-file>
          </template>
          <template v-else-if="field.type == 'textarea'">
            <b-form-textarea
              v-model="formData[field.name]"
              :placeholder="field.placeholder"
              :required="field.required"
              :state="field.valid"
            >
            </b-form-textarea>
          </template>
          <template v-else>
            <b-form-input
              v-model="formData[field.name]"
              :type="field.type"
              :placeholder="field.placeholder"
              :required="field.required"
              :state="field.valid"
            >
            </b-form-input>
          </template>
          <b-form-invalid-feedback v-if="field.validMessage">
            {{ field.validMessage }}
          </b-form-invalid-feedback>
        </b-form-group>
      </template>
    </template>
    <slot></slot>
  </b-form>
</template>


<script>
export default {
  components: {},
  props: ["source", "form", "formData"],
  data() {
    return {
      key: null,
      ready: false,
      fileData: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    init() {
      this.key = this.tool.getNumber();
      let fileData = {};
      this.form.fields.forEach((one) => {
        if (one.type == "file") {
          fileData[one.name] = null;
        }
      });
      this.fileData = fileData;
      this.ready = true;
    },
    async fileInput(name, event) {
      let file = this.fileData[name];
      console.log(file);
      let form = new FormData();
      form.append("place", "other");
      form.append("file", file);
      let res = await this.server.upload(form);
      if (res.code != 0) {
        this.tool.error(res.msg);
        return false;
      }
      this.formData[name] = res.data.files[0].path;

      return true;
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

<template>
  <el-form v-if="ready" size="mini">
    <template v-if="key != null && formData != null && fileObjectMap != null">
      <template v-for="field in form.fields">
        <el-form-item :key="`key-${key}-${field.name}`" :label="field.label">
          <template v-if="field.type == 'select'">
            <el-select
              v-model="formData[field.name]"
              :placeholder="field.placeholder"
              :required="field.required"
              style="width: 100%"
            >
              <el-option
                v-for="(one, index) in field.options"
                :key="index"
                :value="one.value"
                :label="one.text"
              >
              </el-option>
            </el-select>
          </template>
          <template v-else-if="field.type == 'switch'">
            <el-switch v-model="form.name"> </el-switch>
          </template>
          <template v-else-if="field.type == 'file'">
            <div
              class="ft-12 pdb-5"
              v-if="tool.isNotEmpty(formData[field.name])"
            >
              <span class="color-grey">文件：</span>
              <a
                class="tm-link color-green"
                :href="source.filesUrl + formData[field.name]"
              >
                {{ formData[field.name] }}
              </a>
            </div>
            <el-upload
              class="upload-file"
              :action="source.api + 'upload'"
              :limit="1"
              :data="{ place: 'other' }"
              :headers="{ JWT: tool.getJWT() }"
              name="file"
              :on-success="fileObjectMap[field.name].success"
              :show-file-list="false"
            >
              <div class="tm-link color-teal-8">点击上传</div>
            </el-upload>
          </template>
          <template v-else>
            <el-input
              v-model="formData[field.name]"
              :type="field.type"
              :placeholder="field.placeholder"
              :required="field.required"
            >
            </el-input>
          </template>
          <div class="color-red" v-if="field.validMessage">
            {{ field.validMessage }}
          </div>
        </el-form-item>
      </template>
    </template>
    <slot></slot>
  </el-form>
</template>


<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {
      key: null,
      form: null,
      formData: null,
      ready: false,
      fileObjectMap: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    build(form, formData) {
      this.form = form;
      this.formData = formData;
      this.key = this.tool.getNumber();
      let fileObjectMap = {};
      this.form.fields.forEach((one) => {
        if (one.type == "file") {
          fileObjectMap[one.name] = {
            success: (response, file, fileList) => {
              if (response.code != 0) {
                this.tool.error(response.msg);
                return false;
              }
              this.formData[one.name] = response.data.files[0].path;
            },
          };
        }
      });
      this.fileObjectMap = fileObjectMap;
      this.ready = true;
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    // this.init();
  },
};
</script>

<style>
</style>

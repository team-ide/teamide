<template>
  <el-form v-if="ready" size="mini" @submit.native.prevent>
    <template
      v-if="
        key != null &&
        formBuild != null &&
        formData != null &&
        fileObjectMap != null
      "
    >
      <template v-for="field in formBuild.fields">
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
          <template v-else-if="field.type == 'textarea'">
            <el-input
              type="textarea"
              v-model="formData[field.name]"
              :autosize="{ minRows: 5, maxRows: 10 }"
            >
            </el-input>
          </template>
          <template v-else-if="field.type == 'json'">
            <el-input
              type="textarea"
              v-model="jsonStringMap[field.name].value"
              :autosize="{ minRows: 5, maxRows: 20 }"
              @input="jsonStringChange(jsonStringMap[field.name])"
              @change="jsonStringChange(jsonStringMap[field.name])"
            >
            </el-input>
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
  props: ["source", "formBuild", "formData"],
  data() {
    return {
      key: null,
      ready: false,
      fileObjectMap: null,
      jsonStringMap: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    init() {
      this.build();
    },
    build() {
      this.key = this.tool.getNumber();
      let fileObjectMap = {};
      let jsonStringMap = {};
      this.formBuild.fields.forEach((one) => {
        let name = one.name;
        let type = one.type;
        if (type == "json") {
          let json = this.formData[name];
          let jsonString = null;
          if (json != null) {
            if (typeof json == "object") {
              jsonString = JSON.stringify(json, null, "  ");
            }
          }
          jsonStringMap[name] = {
            field: one,
            value: jsonString,
            onChange: () => {
              console.log(jsonStringMap[name]);
            },
          };
        } else if (type == "file") {
          fileObjectMap[name] = {
            success: (response, file, fileList) => {
              if (response.code != 0) {
                this.tool.error(response.msg);
                return false;
              }
              this.formData[name] = response.data.files[0].path;
            },
          };
        }
      });
      this.jsonStringMap = jsonStringMap;
      this.fileObjectMap = fileObjectMap;
      this.ready = true;
    },
    async validate(data) {
      let validateResult = await this.form.validate(data || this.formData);
      return validateResult;
    },
    jsonStringChange(bean) {
      let field = bean.field;
      try {
        field.jsonStringValue = bean.value;
        this.tool.stringToJSON(bean.value);
        field.validMessage = null;
      } catch (error) {
        field.validMessage = error;
      }
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

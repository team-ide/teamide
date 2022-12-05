<template>
  <el-form
    v-if="ready"
    size="mini"
    @submit.native.prevent
    :label-width="`${formBuild.labelWidth}`"
  >
    <template
      v-if="
        key != null &&
        formBuild != null &&
        formData != null &&
        fileObjectMap != null
      "
    >
      <template v-for="field in formBuild.fields">
        <template v-if="tool.isEmpty(field.vIf) || exec(field.vIf, formData)">
          <el-form-item :key="`key-${key}-${field.name}`" :label="field.label">
            <template v-if="field.type == 'select'">
              <el-select
                v-model="formData[field.name]"
                :placeholder="field.placeholder"
                :required="field.required"
                style="width: 100%"
                clearable
                filterable
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
                @input="valueChange(field)"
                @change="valueChange(field)"
              >
              </el-input>
            </template>
            <template v-else-if="field.type == 'json'">
              <el-input
                type="textarea"
                v-model="jsonStringMap[field.name].value"
                :autosize="{ minRows: 5, maxRows: 20 }"
                @input="
                  valueChange(field) &&
                    jsonStringChange(jsonStringMap[field.name])
                "
                @change="
                  valueChange(field) &&
                    jsonStringChange(jsonStringMap[field.name])
                "
              >
              </el-input>
            </template>
            <template v-else-if="field.type == 'list'">
              <el-table :data="listObjectMap[field.name].list">
                <template
                  v-for="(listField, listFieldIndex) in listObjectMap[
                    field.name
                  ].fields"
                >
                  <el-table-column
                    :key="listFieldIndex"
                    :label="listField.label"
                    fixed
                  >
                    <template slot-scope="scope">
                      <el-input
                        v-model="scope.row[listField.name]"
                        type="text"
                      />
                    </template>
                  </el-table-column>
                </template>
                <el-table-column label="操作" width="200px">
                  <template slot="header" s>
                    <div
                      class="tm-link color-green mgl-10"
                      @click="listObjectMap[field.name].add({})"
                    >
                      新增
                    </div>
                  </template>
                  <template slot-scope="scope">
                    <div
                      class="tm-link color-grey mglr-5"
                      @click="listObjectMap[field.name].up(scope.row)"
                    >
                      上移
                    </div>
                    <div
                      class="tm-link color-grey mglr-5"
                      @click="listObjectMap[field.name].down(scope.row)"
                    >
                      下移
                    </div>
                    <div
                      class="tm-link color-red mglr-5"
                      @click="listObjectMap[field.name].remove(scope.row)"
                    >
                      删除
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </template>
            <template v-else-if="field.type == 'jsonView'">
              <el-input
                type="textarea"
                v-model="jsonViewMap[field.bindName].value"
                :autosize="{ minRows: 5, maxRows: 20 }"
              >
              </el-input>
            </template>
            <template v-else>
              <el-input
                v-model="formData[field.name]"
                :type="field.type"
                :placeholder="field.placeholder"
                :required="field.required"
                @input="valueChange(field)"
                @change="valueChange(field)"
              >
              </el-input>
            </template>
            <div class="color-red" v-if="field.validMessage">
              {{ field.validMessage }}
            </div>
          </el-form-item>
        </template>
      </template>
    </template>
    <slot></slot>
  </el-form>
</template>


<script>
var JSONbig = require("json-bigint");
var JSONbigString = JSONbig({});

export default {
  components: {},
  props: ["source", "formBuild", "formData"],
  data() {
    return {
      key: null,
      ready: false,
      fileObjectMap: null,
      jsonStringMap: null,
      jsonViewMap: null,
      listObjectMap: null,
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
      if (this.formBuild == null) {
        return;
      }
      this.key = this.tool.getNumber();
      let fileObjectMap = {};
      let jsonStringMap = {};
      let jsonViewMap = {};
      let listObjectMap = {};
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
            onChange: () => {},
          };
          one.jsonStringValue = jsonString;
        } else if (type == "jsonView") {
          let json = this.formData[one.bindName];
          let jsonString = null;
          if (json != null) {
            if (typeof json == "object") {
              jsonString = JSON.stringify(json, null, "  ");
            }
          }
          jsonViewMap[one.bindName] = {
            field: one,
            value: jsonString,
            onChange: () => {},
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
        } else if (type == "list") {
          let listFields = one.fields || [];
          let list = this.formData[name] || [];
          one.fields = one.fields || [];
          let listObject = {
            list: list,
            fields: listFields,
            fullData: (data) => {
              data = data || {};
              listFields.forEach((listField) => {
                data[listField.name] = data[listField.name];
              });
            },
            add: (data) => {
              listObject.fullData(data);
              list.push(data);
            },

            up: (data) => {
              this.tool.up(listObject, "list", data);
            },
            down: (data) => {
              this.tool.down(listObject, "list", data);
            },
            remove: (data) => {
              let findIndex = list.indexOf(data);
              if (findIndex >= 0) {
                list.splice(findIndex, 1);
              }
            },
          };
          listObjectMap[name] = listObject;
          this.formData[name] = list;
        }
      });
      this.listObjectMap = listObjectMap;
      this.jsonViewMap = jsonViewMap;
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
    exec(vIf, data) {
      if (this.tool.isEmpty(vIf)) {
        return true;
      }
      try {
        var script = ``;
        for (let key in data) {
          script += `var ` + key + `= data['` + key + `'];`;
        }
        script += vIf;
        var res = eval("" + script + "");
        return res;
      } catch (error) {
        console.log(error);
      }
      return false;
    },
    valueChange(field) {
      let value = this.formData[field.name];
      let jsonView = this.jsonViewMap[field.name];
      if (jsonView != null) {
        let jsonString = null;
        if (this.tool.isJSONString(value)) {
          try {
            let json = null;
            try {
              json = JSONbigString.parse(value);
            } catch (error) {
              try {
                json = eval("(" + value + ")");
              } catch (error2) {
                throw error;
              }
            }
            jsonString = JSON.stringify(json, null, "  ");
          } catch (e) {
            jsonString = e;
          }
        }
        jsonView.value = jsonString;
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

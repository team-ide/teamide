<template>
  <div class="toolbox-redis-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="100%">
          <tm-layout :width="style.left.width">
            <tm-layout height="60px">
              <b-form inline class="pdt-20 pdlr-10">
                <label label-size="sm" class="pdr-5"
                  >Key（支持 * 模糊搜索）</label
                >
                <b-form-input
                  size="sm"
                  id="inline-form-input-name"
                  class="mb-2 mr-sm-2 mb-sm-0"
                  v-model="searchForm.pattern"
                  placeholder="key*"
                >
                </b-form-input>
                <label label-size="sm" class="pdr-5">数量</label>
                <b-form-input
                  size="sm"
                  id="inline-form-input-name"
                  class="mb-2 mr-sm-2 mb-sm-0"
                  v-model="searchForm.size"
                >
                </b-form-input>

                <div class="tm-btn bg-teal-8 ft-14" @click="toSearch">搜索</div>
                <div class="tm-btn bg-green ft-14" @click="toInsert">新增</div>
                <div
                  class="tm-btn bg-orange ft-14"
                  @click="toDeletePattern(searchForm.pattern)"
                >
                  删除
                </div>
              </b-form>
            </tm-layout>
            <tm-layout height="auto" class="scrollbar">
              <div class="pd-10" style="o">
                <template v-if="searchResult != null">
                  <table>
                    <thead>
                      <tr>
                        <th>Key （共 {{ searchResult.count }} 个）</th>
                        <th width="120">操作</th>
                      </tr>
                    </thead>
                    <tbody>
                      <template v-if="searchResult.keys.length == 0">
                        <tr>
                          <td colspan="2">
                            <div class="text-center ft-13 pdtb-10">
                              暂无匹配数据!
                            </div>
                          </td>
                        </tr>
                      </template>
                      <template v-else>
                        <template v-for="(one, index) in searchResult.keys">
                          <tr :key="index">
                            <td>{{ one }}</td>
                            <td>
                              <div
                                class="tm-btn color-blue tm-btn-xs"
                                @click="toUpdate(one)"
                              >
                                修改
                              </div>
                              <div
                                class="tm-btn color-orange tm-btn-xs"
                                @click="toDelete(one)"
                              >
                                删除
                              </div>
                            </td>
                          </tr>
                        </template>
                      </template>
                    </tbody>
                  </table>
                </template>
              </div>
            </tm-layout>
          </tm-layout>
          <tm-layout-bar right></tm-layout-bar>
          <tm-layout width="auto">
            <b-form class="pd-10">
              <b-form-group label="类型" label-size="sm">
                <b-form-select
                  v-model="form.type"
                  placeholder="请选择类型"
                  :options="[{ text: 'string', value: 'string' }]"
                >
                </b-form-select>
              </b-form-group>
              <b-form-group label="Key" label-size="sm">
                <b-form-input size="sm" v-model="form.key"> </b-form-input>
              </b-form-group>
              <template v-if="form.type == 'string'">
                <b-form-group label="Value" label-size="sm">
                  <b-form-textarea
                    size="sm"
                    rows="5"
                    max-rows="10"
                    v-model="form.value"
                  >
                  </b-form-textarea>
                </b-form-group>
                <template v-if="form.valueJson != null">
                  <b-form-group label="值JSON预览" label-size="sm">
                    <b-form-textarea
                      size="sm"
                      rows="5"
                      max-rows="10"
                      v-model="form.valueJson"
                    >
                    </b-form-textarea>
                  </b-form-group>
                </template>
              </template>
              <template v-else>
                <div class="text-center ft-13 pdtb-10">
                  暂不支持[{{ form.type }}]类型的值编辑!
                </div>
              </template>
              <div class="pdtb-20">
                <div class="tm-btn bg-teal-8" @click="toSave">保存</div>
              </div>
            </b-form>
          </tm-layout>
        </tm-layout>
      </tm-layout>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "data", "toolboxType", "toolbox", "option", "wrap"],
  data() {
    return {
      style: {
        left: {
          width: "800px",
        },
        main: {},
      },
      searchForm: {
        pattern: "*",
        size: 50,
      },
      form: {
        type: "string",
        key: null,
        value: null,
        valueJson: null,
      },
      searchResult: null,
      ready: false,
    };
  },
  computed: {},
  watch: {
    "form.value"(value) {
      this.form.valueJson = null;
      if (this.tool.isNotEmpty(value)) {
        try {
          if (
            (value.startsWith("{") && value.endsWith("}")) ||
            (value.startsWith("[") && value.endsWith("]"))
          ) {
            let data = JSON.parse(value);
            this.form.valueJson = JSON.stringify(data, null, "    ");
          }
        } catch (e) {
          this.form.valueJson = e;
        }
      }
    },
  },
  methods: {
    init() {
      this.ready = true;
      this.loadKeys();
    },
    toSearch() {
      this.loadKeys();
    },
    toSave() {
      this.doSave();
    },
    toInsert() {
      this.form.type = "string";
      this.form.key = null;
      this.form.value = null;
    },
    async toUpdate(one) {
      let data = await this.get(one);
      if (data == null) {
        data = {};
      }
      this.form.key = one;
      this.form.type = data.type || "string";
      this.form.value = data.value;
    },
    toDelete(key) {
      this.tool
        .confirm("确认删除[" + key + "]？")
        .then(async () => {
          this.doDelete(key);
        })
        .catch((e) => {});
    },
    toDeletePattern(pattern) {
      this.tool
        .confirm("将删除所有匹配[" + pattern + "]的Key，确定删除？")
        .then(async () => {
          this.doDeletePattern(pattern);
        })
        .catch((e) => {});
    },
    async loadKeys() {
      let param = {};
      Object.assign(param, this.searchForm);
      if (this.tool.isEmpty(param.size)) {
        param.size = 50;
      }
      param.size = Number(param.size);
      let res = await this.wrap.work("keys", param);
      this.searchResult = res.data;
    },
    async get(key) {
      let param = {
        key: key,
      };
      let res = await this.wrap.work("get", param);
      return res.data;
    },
    async doSave() {
      let param = {};
      Object.assign(param, this.form);
      param.type = "set";
      let res = await this.wrap.work("do", param);
      if (res.code == 0) {
        this.tool.info("保存成功!");
        this.toSearch();
      }
    },
    async doDelete(key) {
      let param = {
        key: key,
      };
      let res = await this.wrap.work("delete", param);
      if (res.code == 0) {
        this.tool.info("删除成功!");
        this.toSearch();
      }
    },
    async doDeletePattern(pattern) {
      let param = {
        pattern: pattern,
      };
      let res = await this.wrap.work("deletePattern", param);
      if (res.code == 0) {
        this.tool.info("删除成功!");
        this.toSearch();
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
.toolbox-redis-editor {
  width: 100%;
  height: 100%;
}
</style>

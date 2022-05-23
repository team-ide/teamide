<template>
  <div class="toolbox-redis-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout :width="style.left.width">
          <tm-layout height="120px">
            <el-form
              class="pdt-10"
              label-width="150px"
              size="mini"
              :inline="true"
            >
              <el-form-item label="Database" label-width="70px">
                <el-input v-model="searchForm.database" style="width: 50px" />
              </el-form-item>
              <el-form-item label="Key(支持*模糊搜索)">
                <el-input v-model="searchForm.pattern" style="width: 120px" />
              </el-form-item>
              <el-form-item label="数量" label-width="50px">
                <el-input v-model="searchForm.size" style="width: 50px" />
              </el-form-item>
              <div class="pdlr-10">
                <div class="tm-btn tm-btn-sm bg-teal-8 ft-13" @click="toSearch">
                  搜索
                </div>
                <div class="tm-btn tm-btn-sm bg-green ft-13" @click="toInsert">
                  新增
                </div>
                <div
                  class="tm-btn tm-btn-sm bg-orange ft-13"
                  @click="
                    toDeletePattern(searchForm.database, searchForm.pattern)
                  "
                >
                  删除
                </div>
              </div>
            </el-form>
          </tm-layout>
          <tm-layout height="auto" class="scrollbar">
            <div class="pd-10" style="o">
              <table>
                <thead>
                  <tr>
                    <th>
                      Key
                      <template v-if="searchResult != null">
                        （共
                        {{ searchResult.count }}
                        个）
                      </template>
                    </th>
                    <th width="120">操作</th>
                  </tr>
                </thead>
                <tbody>
                  <template v-if="searchResult == null">
                    <tr>
                      <td colspan="2">
                        <div class="text-center ft-13 pdtb-10">
                          数据加载中，请稍后!
                        </div>
                      </td>
                    </tr>
                  </template>
                  <template
                    v-else-if="
                      searchResult.dataList == null ||
                      searchResult.dataList.length == 0
                    "
                  >
                    <tr>
                      <td colspan="2">
                        <div class="text-center ft-13 pdtb-10">
                          暂无匹配数据!
                        </div>
                      </td>
                    </tr>
                  </template>
                  <template v-else>
                    <template v-for="(one, index) in searchResult.dataList">
                      <tr :key="index" @click="rowClick(one)">
                        <td>{{ one.key }}</td>
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
            </div>
          </tm-layout>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <div>
            <el-form
              ref="form"
              size="mini"
              class="pd-10"
              @submit.native.prevent
              v-loading="get_loading"
            >
              <el-form-item label="Database">
                <el-input v-model="form.database"> </el-input>
              </el-form-item>
              <el-form-item label="值类型">
                <el-select
                  placeholder="请选择类型"
                  v-model="form.type"
                  style="width: 100%"
                >
                  <el-option label="string" value="string"> </el-option>
                  <el-option label="list" value="list"></el-option>
                  <el-option label="set" value="set"></el-option>
                  <el-option label="hash" value="hash"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="Key">
                <el-input v-model="form.key"> </el-input>
              </el-form-item>
              <template
                v-if="
                  form.type == 'list' ||
                  form.type == 'set' ||
                  form.type == 'hash'
                "
              >
                <el-form-item label="ValueSize（加载值数量）">
                  <el-input v-model="form.valueSize" @change="valueSizeChange">
                  </el-input>
                </el-form-item>
              </template>
              <template v-if="form.type == 'string'">
                <el-form-item label="Value">
                  <el-input
                    type="textarea"
                    v-model="form.value"
                    :autosize="{ minRows: 5, maxRows: 10 }"
                  >
                  </el-input>
                </el-form-item>
                <template v-if="form.valueJson != null">
                  <el-form-item label="值JSON预览">
                    <el-input
                      type="textarea"
                      v-model="form.valueJson"
                      :autosize="{ minRows: 5, maxRows: 10 }"
                    >
                    </el-input>
                  </el-form-item>
                </template>
                <div class="pdtb-20">
                  <div class="tm-btn bg-teal-8" @click="toSave">保存</div>
                </div>
              </template>
              <template v-else-if="form.type == 'list'">
                <div class="">
                  List
                  <a
                    class="tm-link color-green ft-15 mgl-20"
                    @click="
                      addOne('list', {
                        value: null,
                        isNew: true,
                        index: form.list.length,
                      })
                    "
                    title="添加"
                  >
                    <i class="mdi mdi-plus"></i>
                  </a>
                </div>
                <div
                  v-if="form.list.length > 0"
                  class="pd-5 mgt-5 bd-1 bd-grey-7"
                >
                  <template v-for="(one, index) in form.list">
                    <el-row :key="index">
                      <el-col :span="20">
                        <el-form-item
                          :label="'Value' + '（' + one.index + '）'"
                        >
                          <el-input v-model="one.value"></el-input>
                        </el-form-item>
                      </el-col>
                      <el-col :span="4">
                        <a
                          class="tm-link color-green ft-15 mgt-30"
                          @click="
                            toDo(one.isNew ? 'rpush' : 'lset', {
                              value: one.value,
                              index: one.index,
                            })
                          "
                          title="保存"
                        >
                          <i class="mdi mdi-content-save-outline"></i>
                        </a>
                        <a
                          class="tm-link color-red ft-15 mgt-30"
                          @click="toDo('lrem', { value: one.value, count: 0 })"
                          title="删除"
                          v-if="one.canDel"
                        >
                          <i class="mdi mdi-delete-outline"></i>
                        </a>
                      </el-col>
                    </el-row>
                  </template>
                </div>
              </template>
              <template v-else-if="form.type == 'set'">
                <div class="">
                  Set
                  <a
                    class="tm-link color-green ft-15 mgl-20"
                    @click="
                      addOne('set', {
                        value: null,
                        isNew: true,
                        index: form.set.length,
                      })
                    "
                    title="添加"
                  >
                    <i class="mdi mdi-plus"></i>
                  </a>
                </div>
                <div
                  v-if="form.set.length > 0"
                  class="pd-5 mgt-5 bd-1 bd-grey-7"
                >
                  <template v-for="(one, index) in form.set">
                    <el-row :key="index">
                      <el-col :span="20">
                        <el-form-item
                          :label="'Value' + '（' + one.index + '）'"
                        >
                          <el-input v-model="one.value"></el-input>
                        </el-form-item>
                      </el-col>
                      <el-col :span="4">
                        <a
                          class="tm-link color-green ft-15 mgt-30"
                          @click="
                            toDo('sadd', {
                              value: one.value,
                              index: one.index,
                            })
                          "
                          title="保存"
                        >
                          <i class="mdi mdi-content-save-outline"></i>
                        </a>
                        <a
                          class="tm-link color-red ft-15 mgt-30"
                          @click="toDo('srem', { value: one.value })"
                          title="删除"
                          v-if="!one.isNew"
                        >
                          <i class="mdi mdi-delete-outline"></i>
                        </a>
                      </el-col>
                    </el-row>
                  </template>
                </div>
              </template>
              <template v-else-if="form.type == 'hash'">
                <div class="">
                  Hash
                  <a
                    class="tm-link color-green ft-15 mgl-20"
                    @click="
                      addOne('hash', {
                        field: null,
                        value: null,
                        isNew: true,
                        index: form.hash.length,
                      })
                    "
                    title="添加"
                  >
                    <i class="mdi mdi-plus"></i>
                  </a>
                </div>
                <div
                  v-if="form.hash.length > 0"
                  class="pd-5 mgt-5 bd-1 bd-grey-7"
                >
                  <template v-for="(one, index) in form.hash">
                    <el-row :key="index">
                      <el-col :span="10">
                        <el-form-item
                          :label="'Field' + '（' + one.index + '）'"
                        >
                          <el-input v-model="one.field"></el-input>
                        </el-form-item>
                      </el-col>
                      <el-col :span="10">
                        <el-form-item :label="'Value'">
                          <el-input v-model="one.value"></el-input>
                        </el-form-item>
                      </el-col>
                      <el-col :span="4">
                        <a
                          class="tm-link color-green ft-15 mgt-30"
                          @click="
                            toDo('hset', {
                              field: one.field,
                              value: one.value,
                              index: one.index,
                            })
                          "
                          title="保存"
                        >
                          <i class="mdi mdi-content-save-outline"></i>
                        </a>
                        <a
                          class="tm-link color-red ft-15 mgt-30"
                          @click="toDo('hdel', { field: one.field })"
                          title="删除"
                          v-if="!one.isNew"
                        >
                          <i class="mdi mdi-delete-outline"></i>
                        </a>
                      </el-col>
                    </el-row>
                  </template>
                </div>
              </template>
              <template v-else>
                <div class="text-center ft-13 pdtb-10">
                  暂不支持[{{ form.type }}]类型的值编辑!
                </div>
              </template>
            </el-form>
          </div>
        </tm-layout>
      </tm-layout>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxType", "toolbox", "option", "wrap"],
  data() {
    return {
      style: {
        left: {
          width: "600px",
        },
        main: {},
      },
      get_loading: false,
      searchForm: {
        database: 0,
        pattern: "*",
        size: 50,
      },
      form: {
        database: 0,
        type: "string",
        valueSize: 10,
        key: null,
        value: null,
        valueJson: null,
        set: [],
        list: [],
        hash: [],
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
    valueSizeChange() {
      if (this.tool.isNotEmpty(this.form.key)) {
        this.get(this.form.database, this.form.key);
      }
    },
    addOne(type, one) {
      this.form[type].push(one);
    },
    init() {
      this.ready = true;
      this.loadKeys();
    },
    refresh() {
      this.toSearch();
    },
    toSearch() {
      this.loadKeys();
    },
    toSave() {
      this.doSave();
    },
    toInsert() {
      this.form.database = Number(this.searchForm.database);
      this.form.type = "string";
      this.form.key = null;
      this.form.value = null;
      this.form.set = [];
      this.form.list = [];
      this.form.hash = [];
    },
    async toDo(type, data) {
      data = data || {};
      data.database = Number(this.form.database);
      data.key = this.form.key;
      data.type = type;
      if (this.tool.isEmpty(type)) {
        this.tool.error("操作类型不能为空！");
        return;
      }
      if (this.tool.isEmpty(data.key)) {
        this.tool.error("Key不能为空！");
        return;
      }

      let res = await this.wrap.work("do", data);
      if (res.code != 0) {
        this.tool.error(res.msg);
      } else {
        this.tool.success("操作成功");
        let find = false;

        if (this.searchResult && this.searchResult.dataList) {
          this.searchResult.dataList.forEach((one) => {
            if (one.key == data.key) {
              find = true;
            }
          });
        }
        if (!find) {
          this.toSearch();
        }

        this.get(data.database, data.key);
      }
    },
    rowClick(data) {
      this.rowClickTimeCache = this.rowClickTimeCache || {};
      let nowTime = new Date().getTime();
      let clickTime = this.rowClickTimeCache[data];
      this.rowClickTimeCache[data] = nowTime;
      if (clickTime) {
        let timeout = nowTime - clickTime;
        if (timeout < 300) {
          delete this.rowClickTimeCache[data];
          this.rowDbClick(data);
        }
      }
    },
    rowDbClick(data) {
      this.toUpdate(data);
    },
    async toUpdate(one) {
      let data = await this.get(one.database, one.key);
      if (data == null) {
        data = {};
      }
    },
    initFormData(data) {
      this.form.database = Number(data.database);
      this.form.key = data.key;
      this.form.type = data.type || "string";

      this.form.value = null;
      this.form.set = [];
      this.form.list = [];
      this.form.hash = [];
      if (data.type == "list") {
        data.value = data.value || [];
        data.value.forEach((one, index) => {
          let canDel = true;
          data.value.forEach((one_, index_) => {
            if (one_ == one && index != index_) {
              canDel = false;
            }
          });
          this.form.list.push({
            value: one,
            index: index,
            canDel: canDel,
          });
        });
      } else if (data.type == "set") {
        data.value = data.value || [];
        data.value.forEach((one, index) => {
          this.form.set.push({
            value: one,
            index: index,
          });
        });
      } else if (data.type == "hash") {
        data.value = data.value || {};
        let index = 0;
        for (var field in data.value) {
          this.form.hash.push({
            field: field,
            value: data.value[field],
            index: index,
          });
          index++;
        }
      } else {
        this.form.value = data.value;
      }
    },
    toDelete(data) {
      this.tool
        .confirm("确认删除[" + data.key + "]？")
        .then(async () => {
          this.doDelete(data.database, data.key);
        })
        .catch((e) => {});
    },
    toDeletePattern(database, pattern) {
      this.tool
        .confirm("将删除所有匹配[" + pattern + "]的Key，确定删除？")
        .then(async () => {
          this.doDeletePattern(database, pattern);
        })
        .catch((e) => {});
    },
    async loadKeys() {
      this.searchResult = null;
      let param = {};
      this.searchForm.database = Number(this.searchForm.database);
      Object.assign(param, this.searchForm);
      if (this.tool.isEmpty(param.size)) {
        param.size = 50;
      }
      param.size = Number(param.size);
      let res = await this.wrap.work("keys", param);
      this.searchResult = res.data;
    },
    async get(database, key) {
      this.get_loading = true;
      let param = {
        database: Number(database),
        key: key,
        valueSize: Number(this.form.valueSize),
      };
      let res = await this.wrap.work("get", param);
      this.get_loading = false;
      let data = res.data || {};
      this.initFormData(data);
      return data;
    },
    async doSave() {
      let param = {};
      this.form.valueSize = Number(this.form.valueSize);
      Object.assign(param, this.form);
      param.type = "set";
      let res = await this.wrap.work("do", param);
      if (res.code == 0) {
        this.tool.success("保存成功!");

        let find = false;

        if (this.searchResult && this.searchResult.dataList) {
          this.searchResult.dataList.forEach((one) => {
            if (one.key == param.key) {
              find = true;
            }
          });
        }
        if (!find) {
          this.toSearch();
        }

        this.get(param.database, param.key);
      }
    },
    async doDelete(database, key) {
      let param = {
        database: Number(database),
        key: key,
      };
      let res = await this.wrap.work("delete", param);
      if (res.code == 0) {
        this.tool.success("删除成功!");
        this.toSearch();
      }
    },
    async doDeletePattern(database, pattern) {
      let param = {
        database: Number(database),
        pattern: pattern,
      };
      let res = await this.wrap.work("deletePattern", param);
      if (res.code == 0) {
        this.tool.success("删除成功!");
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
  user-select: text;
}
</style>

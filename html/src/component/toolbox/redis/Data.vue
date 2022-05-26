<template>
  <div class="toolbox-redis-data" v-loading="loading">
    <template v-if="ready">
      <div class="toolbox-redis-form-box">
        <el-form ref="form" size="mini" @submit.native.prevent>
          <el-form-item label="Database">
            <el-input v-model="form.database"> </el-input>
          </el-form-item>
          <el-form-item label="值类型">
            <el-select
              placeholder="请选择类型"
              v-model="form.valueType"
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
              form.valueType == 'list' ||
              form.valueType == 'set' ||
              form.valueType == 'hash'
            "
          >
            <el-form-item label="ValueSize（加载值数量）">
              <el-input v-model="form.valueSize" @change="valueSizeChange">
              </el-input>
            </el-form-item>
          </template>
          <template v-if="form.valueType == 'string'">
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
            <div class="pdlr-10">
              <div class="tm-btn bg-teal-8" @click="toSave">保存</div>
            </div>
          </template>
          <template v-else-if="form.valueType == 'list'">
            <div class="pdlr-10">
              List
              <a
                class="tm-link color-green ft-15 mgl-20"
                @click="
                  addOne('listList', {
                    value: null,
                    index: listList.length,
                  })
                "
                title="添加"
              >
                <i class="mdi mdi-plus"></i>
              </a>
            </div>
            <div
              v-if="listList.length > 0"
              class="mglr-10 mgt-10 bd-1 bd-grey-7"
            >
              <template v-for="(one, index) in listList">
                <el-row :key="index">
                  <el-col :span="20">
                    <el-form-item :label="'Value' + '（' + one.index + '）'">
                      <el-input v-model="one.value"></el-input>
                    </el-form-item>
                  </el-col>
                  <el-col :span="4">
                    <a
                      class="tm-link color-green ft-15 mgt-30"
                      @click="
                        toDo(one.isNew ? 'RPush' : 'LSet', {
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
                      @click="
                        toDo('LRem', {
                          value: one.value,
                          count: 0,
                          isNew: one.isNew,
                        })
                      "
                      title="删除"
                    >
                      <i class="mdi mdi-delete-outline"></i>
                    </a>
                  </el-col>
                </el-row>
              </template>
            </div>
          </template>
          <template v-else-if="form.valueType == 'set'">
            <div class="pdlr-10">
              Set
              <a
                class="tm-link color-green ft-15 mgl-20"
                @click="
                  addOne('setList', {
                    value: null,
                    index: setList.length,
                  })
                "
                title="添加"
              >
                <i class="mdi mdi-plus"></i>
              </a>
            </div>
            <div
              v-if="setList.length > 0"
              class="mglr-10 mgt-10 bd-1 bd-grey-7"
            >
              <template v-for="(one, index) in setList">
                <el-row :key="index">
                  <el-col :span="20">
                    <el-form-item :label="'Value' + '（' + one.index + '）'">
                      <el-input v-model="one.value"></el-input>
                    </el-form-item>
                  </el-col>
                  <el-col :span="4">
                    <a
                      class="tm-link color-green ft-15 mgt-30"
                      @click="
                        toDo('SAdd', {
                          value: one.value,
                          index: one.index,
                          isNew: one.isNew,
                        })
                      "
                      title="保存"
                    >
                      <i class="mdi mdi-content-save-outline"></i>
                    </a>
                    <a
                      class="tm-link color-red ft-15 mgt-30"
                      @click="toDo('SRem', { value: one.value })"
                      title="删除"
                    >
                      <i class="mdi mdi-delete-outline"></i>
                    </a>
                  </el-col>
                </el-row>
              </template>
            </div>
          </template>
          <template v-else-if="form.valueType == 'hash'">
            <div class="pdlr-10">
              Hash
              <a
                class="tm-link color-green ft-15 mgl-20"
                @click="
                  addOne('hashList', {
                    field: null,
                    value: null,
                    index: hashList.length,
                  })
                "
                title="添加"
              >
                <i class="mdi mdi-plus"></i>
              </a>
            </div>
            <div
              v-if="hashList.length > 0"
              class="mglr-10 mgt-10 bd-1 bd-grey-7"
            >
              <template v-for="(one, index) in hashList">
                <el-row :key="index">
                  <el-col :span="10">
                    <el-form-item :label="'Field' + '（' + one.index + '）'">
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
                        toDo('HSet', {
                          field: one.field,
                          value: one.value,
                          index: one.index,
                          isNew: one.isNew,
                        })
                      "
                      title="保存"
                    >
                      <i class="mdi mdi-content-save-outline"></i>
                    </a>
                    <a
                      class="tm-link color-red ft-15 mgt-30"
                      @click="toDo('HDel', { field: one.field })"
                      title="删除"
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
              暂不支持[{{ form.valueType }}]类型的值编辑!
            </div>
          </template>
        </el-form>
      </div>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxType", "toolbox", "option", "extend", "wrap"],
  data() {
    return {
      ready: false,
      loading: false,
      form: {
        database: 0,
        valueType: "string",
        valueSize: 10,
        key: null,
        value: null,
        valueJson: null,
      },
      setList: [],
      listList: [],
      hashList: [],
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
      let extend = this.extend || {};
      let database = extend.database || 0;
      let key = extend.key;
      this.form.database = Number(database);
      this.form.key = key;
      this.initForm();
      this.ready = true;
    },
    refresh() {},
    async initForm() {
      this.loading = true;
      this.form.value = null;
      this.setList = [];
      this.listList = [];
      this.hashList = [];

      if (this.tool.isNotEmpty(this.form.key)) {
        let keyData = await this.load(this.form.database, this.form.key);
        this.initFormData(keyData);
      }
      if (this.tool.isEmpty(this.form.valueType)) {
        this.form.valueType = "string";
      }
      this.loading = false;
    },
    async load(database, key) {
      let param = {
        database: Number(database),
        key: key,
        valueSize: Number(this.form.valueSize),
      };
      let res = await this.wrap.work("get", param);
      let data = res.data || {};
      return data;
    },
    initFormData(data) {
      this.form.database = Number(data.database);
      this.form.key = data.key;
      if (data.valueType != "none") {
        this.form.valueType = data.valueType || "string";
      }

      this.form.value = null;
      this.setList = [];
      this.listList = [];
      this.hashList = [];
      if (data.valueType == "list") {
        data.value = data.value || [];
        data.value.forEach((one, index) => {
          this.listList.push({
            value: one,
            index: index,
            isNew: false,
          });
        });
      } else if (data.valueType == "set") {
        data.value = data.value || [];
        data.value.forEach((one, index) => {
          this.setList.push({
            value: one,
            index: index,
            isNew: false,
          });
        });
      } else if (data.valueType == "hash") {
        data.value = data.value || {};
        let index = 0;
        for (var field in data.value) {
          this.hashList.push({
            field: field,
            value: data.value[field],
            index: index,
            isNew: false,
          });
          index++;
        }
      } else {
        this.form.value = data.value;
      }
    },
    valueSizeChange() {
      if (this.tool.isNotEmpty(this.form.key)) {
        this.get(this.form.database, this.form.key);
      }
    },
    addOne(type, one) {
      one.isNew = true;
      this[type].push(one);
    },
    toSave() {
      this.doSave();
    },
    async toDo(doType, data) {
      data = data || {};
      data.database = Number(this.form.database);
      data.key = this.form.key;
      data.doType = doType;
      if (this.tool.isEmpty(doType)) {
        this.tool.error("操作类型不能为空！");
        return;
      }
      if (this.tool.isEmpty(data.key)) {
        this.tool.error("Key不能为空！");
        return;
      }
      if (data.doType == "LRem") {
        if (data.isNew) {
          this.listList.splice(this.listList.indexOf(data), 1);
          return;
        }
        this.tool
          .confirm("确认删除[" + data.key + "]下所有值[" + data.value + "]？")
          .then(async () => {
            await this.doDo(data);
          })
          .catch((e) => {});
      } else if (data.doType == "SRem") {
        if (data.isNew) {
          this.setList.splice(this.setList.indexOf(data), 1);
          return;
        }
        this.tool
          .confirm("确认删除[" + data.key + "]下键[" + data.value + "]？")
          .then(async () => {
            await this.doDo(data);
          })
          .catch((e) => {});
      } else if (data.doType == "HDel") {
        if (data.isNew) {
          this.hashList.splice(this.hashList.indexOf(data), 1);
          return;
        }
        this.tool
          .confirm("确认删除[" + data.key + "]下键[" + data.value + "]？")
          .then(async () => {
            await this.doDo(data);
          })
          .catch((e) => {});
      } else {
        await this.doDo(data);
      }
    },
    async doDo(data) {
      let res = await this.wrap.work("do", data);
      if (res.code != 0) {
        this.tool.error(res.msg);
      } else {
        this.tool.success("操作成功");

        this.initForm();
      }
    },
    async doSave() {
      let param = {};
      this.form.valueSize = Number(this.form.valueSize);
      Object.assign(param, this.form);
      param.doType = "set";
      let res = await this.wrap.work("do", param);
      if (res.code == 0) {
        this.tool.success("保存成功!");

        this.initForm();
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
.toolbox-redis-data {
  width: 100%;
  height: 100%;
}
.toolbox-redis-form-box {
  width: 100%;
  height: 100%;
}
.toolbox-redis-form-box .el-form {
  width: 100%;
  height: 100%;
}
.toolbox-redis-form-box .el-form .el-form-item {
  padding: 0px 10px;
}
</style>

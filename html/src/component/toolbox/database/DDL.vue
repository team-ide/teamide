<template>
  <div class="toolbox-database-ddl" v-loading="loading">
    <template v-if="ready">
      <el-form
        class="pdt-10"
        ref="form"
        :model="form"
        label-width="90px"
        size="mini"
        :inline="true"
      >
        <el-form-item label="数据库类型">
          <el-select
            placeholder="当前库类型"
            v-model="form.databaseType"
            @change="toLoad"
            style="width: 120px"
          >
            <el-option label="MySql" value="mysql"> </el-option>
            <el-option label="Sqlite" value="sqlite"> </el-option>
            <el-option label="Oracle" value="oracle"> </el-option>
            <el-option label="达梦" value="dameng"> </el-option>
            <el-option label="神通" value="shentong"> </el-option>
            <el-option label="金仓" value="kingbase"> </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="生成建库">
          <el-switch v-model="form.generateDatabase" @change="toLoad">
          </el-switch>
        </el-form-item>
        <el-form-item label="追加库名">
          <el-switch v-model="form.appendDatabase" @change="toLoad">
          </el-switch>
        </el-form-item>
        <template v-if="form.appendDatabase">
          <el-form-item label="库名包装">
            <el-select
              placeholder="不包装"
              v-model="form.databasePackingCharacter"
              @change="toLoad"
              style="width: 100px"
            >
              <el-option
                v-for="(one, index) in packingCharacters"
                :key="index"
                :value="one.value"
              >
                {{ one.text }}
              </el-option>
            </el-select>
          </el-form-item>
        </template>
        <el-form-item label="表名包装">
          <el-select
            placeholder="不包装"
            v-model="form.tablePackingCharacter"
            @change="toLoad"
            style="width: 100px"
          >
            <el-option
              v-for="(one, index) in packingCharacters"
              :key="index"
              :value="one.value"
            >
              {{ one.text }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="字段包装">
          <el-select
            placeholder="不包装"
            v-model="form.columnPackingCharacter"
            @change="toLoad"
            style="width: 100px"
          >
            <el-option
              v-for="(one, index) in packingCharacters"
              :key="index"
              :value="one.value"
            >
              {{ one.text }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="字符值包装">
          <el-select
            placeholder="不包装"
            v-model="form.stringPackingCharacter"
            @change="toLoad"
            style="width: 100px"
          >
            <el-option
              v-for="(one, index) in stringPackingCharacters"
              :key="index"
              :value="one.value"
            >
              {{ one.text }}
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <textarea v-model="showDDL"> </textarea>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "wrap", "database", "table"],
  data() {
    return {
      ready: false,
      showDDL: null,
      loading: false,
      packingCharacters: [
        { value: "", text: "不包装" },
        { value: "'", text: "'" },
        { value: '"', text: '"' },
        { value: "`", text: "`" },
      ],
      stringPackingCharacters: [
        { value: "'", text: "'" },
        { value: '"', text: '"' },
      ],
      form: {
        databaseType: "",
        generateDatabase: false,
        appendDatabase: false,
        databasePackingCharacter: null,
        tablePackingCharacter: null,
        columnPackingCharacter: null,
        stringPackingCharacter: "'",
      },
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.ready = true;
      this.toLoad();
    },
    async toLoad() {
      this.loading = true;
      await this.load();
      this.loading = false;
    },
    async load() {
      let sqls = await this.loadDDL(this.database, this.table);
      let ddl = "";
      sqls.forEach((sql, index) => {
        if (index > 0) {
          ddl += "\n";
        }
        ddl += sql + ";";
      });
      this.showDDL = ddl;
    },
    async loadDDL(database, table) {
      let param = Object.assign({}, this.form);
      param.database = database;
      param.table = table;
      let res = await this.wrap.work("ddl", param);
      res.data = res.data || {};
      return res.data.sqlList || [];
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-ddl {
  width: 100%;
  height: 100%;
}
.toolbox-database-ddl textarea {
  width: 100%;
  height: calc(100% - 140px) !important;
  margin-top: 23px;
  letter-spacing: 1px;
  word-spacing: 5px;
}
</style>

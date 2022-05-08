<template>
  <div class="toolbox-database-ddl">
    <template v-if="ready">
      <el-form
        class="mgt-10"
        ref="form"
        :model="form"
        label-width="110px"
        size="mini"
        :inline="true"
      >
        <el-form-item label="数据库类型">
          <el-select
            placeholder="当前库类型"
            v-model="form.databaseType"
            @change="toLoad"
          >
            <el-option label="MySql" value="mysql"> </el-option>
            <el-option label="Oracle" value="oracle"> </el-option>
            <el-option label="达梦" value="dameng"> </el-option>
            <el-option label="神通" value="shentong"> </el-option>
            <el-option label="金仓" value="kingbase"> </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <b-form-textarea
        size="sm"
        rows="10"
        max-rows="30"
        v-model="showDDL"
        class="toolbox-database-textarea"
      >
      </b-form-textarea>
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
      form: {
        databaseType: "",
      },
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.ready = true;
      this.load();
    },
    async toLoad() {
      await this.load();
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
      let param = {
        database: database,
        table: table,
        databaseType: this.form.databaseType,
      };
      let res = await this.wrap.work("ddl", param);
      res.data = res.data || {};
      return res.data.sqls || [];
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
.toolbox-database-textarea {
  width: 100%;
  height: calc(100% - 70px) !important;
  margin-top: 23px;
  letter-spacing: 1px;
  word-spacing: 5px;
}
</style>

<template>
  <b-modal
    ref="modal"
    :title="tableDetail == null ? '' : tableDetail.name"
    :hide-header-close="false"
    :no-close-on-backdrop="true"
    :no-close-on-esc="true"
    :hide-backdrop="true"
    hide-footer
    size="lg"
  >
    <div class="ft-15">
      <b-form-textarea size="sm" rows="10" max-rows="30" v-model="sql" readonly>
      </b-form-textarea>
    </div>
  </b-modal>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "wrap"],
  data() {
    return {
      sql: null,
      tableDetail: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    async show(tableDetail, datas) {
      datas = datas || [];
      this.tableDetail = tableDetail;
      this.sql = await this.initSql(datas);
      this.$refs["modal"].show();
    },
    hide() {
      this.$refs["modal"].hide();
    },
    async initSql(datas) {
      let res = "";
      if (this.tableDetail == null) {
        return res;
      }
      datas.forEach((one) => {
        let sql = "INSERT INTO " + this.tableDetail.name;
        sql += " (";
        this.tableDetail.columns.forEach((column) => {
          sql += "" + column.name + ",";
        });
        if (sql.endsWith(",")) {
          sql = sql.substring(0, sql.length - 1);
        }
        sql += " )";

        sql += " VALUES (";
        this.tableDetail.columns.forEach((column) => {
          let value = one[column.name];
          if (value == null) {
            sql += "null,";
          } else {
            if (
              column.type == "int" ||
              column.type == "bigint" ||
              column.type == "bit" ||
              column.type == "number"
            ) {
              sql += "" + value + ",";
            } else {
              sql += "'" + value + "',";
            }
          }
        });
        if (sql.endsWith(",")) {
          sql = sql.substring(0, sql.length - 1);
        }
        sql += " )";

        res += sql + ";\n";
      });
      return res;
    },
    init() {},
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.wrap.showSqlForInsert = this.show;
    this.init();
  },
};
</script>

<style>
</style>

<template>
  <div class="toolbox-database-table-data">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="120px" style="overflow: hidden">
          <tm-layout width="400px">
            <ul class="part-box scrollbar mg-0" v-if="tableDetail != null">
              <template v-for="(one, index) in form.wheres">
                <li :key="index">
                  <input v-model="one.checked" type="checkbox" />
                  <template v-if="one.sqlConditionalOperation == 'custom'">
                  </template>
                  <template v-else>
                    <select
                      v-model="one.name"
                      @change="initInputWidth"
                      class="part-form-input"
                    >
                      <option :value="null" text="请选择">请选择</option>
                      <template v-for="(one, index) in tableDetail.columns">
                        <option
                          :key="index"
                          :value="one.name"
                          :text="one.name + '&nbsp;'"
                        >
                          {{ one.name }}
                          <template v-if="tool.isNotEmpty(one.comment)">
                            （{{ one.comment }}）
                          </template>
                        </option>
                      </template>
                    </select>
                  </template>
                  <select
                    v-model="one.sqlConditionalOperation"
                    @change="initInputWidth"
                    class="part-form-input"
                  >
                    <template
                      v-for="(one, index) in source.sqlConditionalOperations"
                    >
                      <option :key="index" :value="one.value" :text="one.text">
                        {{ one.text }}
                        <template v-if="tool.isNotEmpty(one.comment)">
                          （{{ one.comment }}）
                        </template>
                      </option>
                    </template>
                  </select>
                  <template
                    v-if="
                      one.sqlConditionalOperation == 'is null' ||
                      one.sqlConditionalOperation == 'is not null' ||
                      one.sqlConditionalOperation == 'is empty' ||
                      one.sqlConditionalOperation == 'is not empty'
                    "
                  >
                  </template>
                  <template
                    v-else-if="
                      one.sqlConditionalOperation == 'between' ||
                      one.sqlConditionalOperation == 'not between'
                    "
                  >
                    <input
                      v-model="one.before"
                      type="text"
                      @input="initInputWidth"
                      @change="initInputWidth"
                      class="part-form-input"
                    />
                    <span class="mglr-5">,</span>
                    <input
                      v-model="one.after"
                      type="text"
                      @input="initInputWidth"
                      @change="initInputWidth"
                      class="part-form-input"
                    />
                  </template>
                  <template v-else-if="one.sqlConditionalOperation == 'custom'">
                    <input
                      v-model="one.customSql"
                      type="text"
                      @input="initInputWidth"
                      @change="initInputWidth"
                      class="part-form-input"
                    />
                  </template>
                  <template v-else>
                    <input
                      v-model="one.value"
                      type="text"
                      @input="initInputWidth"
                      @change="initInputWidth"
                      class="part-form-input"
                    />
                  </template>

                  <select
                    v-model="one.andOr"
                    @change="initInputWidth"
                    class="part-form-input"
                  >
                    <option value="AND">AND</option>
                    <option value="OR">OR</option>
                  </select>
                </li>
              </template>
              <li class="pdl-5">
                <div @click="addWhere" class="color-green tm-link mgr-10">
                  添加条件
                </div>
                <div @click="doSearch" class="color-green tm-link mgr-10">
                  查询
                </div>
              </li>
            </ul>
          </tm-layout>
          <tm-layout-bar right></tm-layout-bar>
          <tm-layout width="400px">
            <ul class="part-box scrollbar mg-0" v-if="tableDetail != null">
              <li></li>
            </ul>
          </tm-layout>
          <tm-layout-bar right></tm-layout-bar>
          <tm-layout>
            <ul class="part-box scrollbar mg-0" v-if="tableDetail != null">
              <li></li>
            </ul>
          </tm-layout>
        </tm-layout>
        <tm-layout-bar bottom></tm-layout-bar>
        <tm-layout height="20px">
          <div class="ft-12 pdl-10" v-if="tableDetail != null">
            <div class="color-red tm-link mgr-10">删除</div>
            <div class="color-green tm-link mgr-10">保存修改</div>
            <div class="color-blue tm-link mgr-10">新增</div>
            <div @click="exportDataForInsert" class="color-grey tm-link mgr-10">
              导出（Insert）
            </div>
          </div>
        </tm-layout>
        <tm-layout height="auto" v-loading="datas_loading">
          <div
            class="toolbox-database-table-data-table"
            v-if="tableDetail != null"
          >
            <el-table
              :data="datas"
              :border="true"
              height="100%"
              style="width: 100%"
              size="mini"
            >
              <el-table-column width="70">
                <template slot-scope="scope" label="">
                  <!-- :checked="selects.indexOf(scope.row) >= 0" -->
                  <input
                    type="checkbox"
                    class="mgl-5"
                    style="width: auto"
                    v-model="selects"
                    :value="scope.row"
                  />
                  <span class="mgl-5">{{ scope.$index }}</span>
                  <template v-if="updates.indexOf(scope.row) >= 0">
                    <i class="mgl-5 mdi mdi-text-box-outline"></i>
                  </template>
                  <template v-if="inserts.indexOf(scope.row) >= 0">
                    <i class="mgl-5 mdi mdi-text-box-plus-outline"></i>
                  </template>
                </template>
              </el-table-column>
              <template v-for="(column, index) in form.columns">
                <template v-if="column.checked">
                  <el-table-column
                    :key="index"
                    :prop="column.name"
                    :label="column.name"
                    width="150"
                  >
                    <template slot-scope="scope">
                      <div class="">
                        <input
                          v-model="scope.row[column.name]"
                          :placeholder="
                            scope.row[column.name] == null ? 'null' : ''
                          "
                          type="text"
                        />
                      </div>
                    </template>
                  </el-table-column>
                </template>
              </template>
            </el-table>
          </div>
        </tm-layout>
        <tm-layout height="30px">
          <div
            class="
              ft-12
              pdt-2
              text-center
              toolbox-database-table-data-pagination
            "
            v-if="tableDetail != null && sql != null"
          >
            <el-pagination
              small
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              :current-page="pageIndex"
              :page-sizes="[10, 50, 100, 200, 500]"
              :page-size="pageSize"
              layout="total, sizes, prev, pager, next, jumper"
              :total="total"
              :disabled="total <= 0"
            >
            </el-pagination>
          </div>
        </tm-layout>
        <tm-layout-bar top></tm-layout-bar>
        <tm-layout height="50px" class="scrollbar">
          <div class="ft-12 pdlr-10" v-if="tableDetail != null && sql != null">
            <div style="line-height: 20px">
              <span style="word-break: break-all">{{ sql }}</span>
            </div>
          </div>
        </tm-layout>
      </tm-layout>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: [
    "source",
    "toolboxType",
    "toolbox",
    "option",
    "database",
    "table",
    "wrap",
  ],
  data() {
    return {
      ready: false,
      tableDetail: null,
      datas_loading: false,
      datas: null,
      sql: null,
      params: null,
      inserts: [],
      updates: [],
      selects: [],
      pageSize: 10,
      pageIndex: 1,
      total: 0,
      form: {
        wheres: [],
        orders: [],
        columns: [],
      },
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      await this.initTable();
      if (this.tableDetail == null) {
        return;
      }
      this.result = null;
      this.form.wheres = [];
      this.form.orders = [];
      this.form.columns = [];
      if (this.tableDetail && this.tableDetail.columns) {
        this.tableDetail.columns.forEach((one) => {
          let column = Object.assign({}, one);
          column.checked = true;
          this.form.columns.push(column);
        });
      }
      this.ready = true;
      this.doSearch();
    },
    async initTable() {
      this.tableDetail = await this.wrap.getTableDetail(
        this.database,
        this.table
      );
    },
    handleSizeChange(pageSize) {
      this.pageSize = pageSize;
      this.doSearch();
    },
    handleCurrentChange(pageIndex) {
      this.pageIndex = pageIndex;
      this.doSearch();
    },
    initInputWidth() {
      this.$nextTick(() => {
        let es = this.$el.getElementsByClassName("part-form-input");
        if (es) {
          Array.prototype.forEach.call(es, (one) => {
            this.initWidth(one);
          });
        }
      });
    },
    initWidth(event) {
      let target = event;
      if (event.target) {
        target = event.target;
      }
      let value = target.value;
      this.tool.initInputWidth(event);
    },
    addWhere() {
      let where = {
        checked: true,
        name: null,
        value: null,
        before: null,
        after: null,
        sqlConditionalOperation: "=",
        andOr: "and",
      };
      let column = null;
      if (this.tableDetail && this.tableDetail.columns) {
        if (this.tableDetail.columns.length > 0) {
          this.tableDetail.columns.forEach((one) => {
            if (column != null) {
              return;
            }
            let find = false;
            this.form.wheres.forEach((w) => {
              if (w.name == one.name) {
                find = true;
              }
            });
            if (find) {
              return;
            }
            column = one;
          });
        }
      }
      if (column != null) {
        where.name = column.name;
      }

      this.form.wheres.push(where);
      this.initInputWidth();
    },
    async doSearch() {
      let wheres = [];
      let orders = [];
      let columns = [];
      this.form.wheres.forEach((where) => {
        if (!where.checked) {
          return;
        }
        if (this.tool.isEmpty(where.name)) {
          return;
        }
        wheres.push(where);
      });
      this.form.orders.forEach((order) => {
        if (!order.checked) {
          return;
        }
        orders.push(order);
      });
      this.form.columns.forEach((column) => {
        if (!column.checked) {
          return;
        }
        columns.push(column);
      });
      let data = {};
      data.database = this.database;
      data.table = this.table;
      data.wheres = wheres;
      data.orders = orders;
      data.columns = columns;
      data.pageIndex = this.pageIndex;
      data.pageSize = this.pageSize;
      this.datas_loading = true;

      this.datas = null;
      this.sql = null;
      this.total = 0;
      this.params = null;
      this.updates = [];
      this.inserts = [];
      this.selects = [];

      let res = await this.wrap.work("datas", data);
      res.data = res.data || {};

      this.datas = res.data.datas;
      this.sql = res.data.sql;
      this.total = Number(res.data.total || 0);
      this.params = res.data.params;
      this.datas_loading = false;
    },
    exportDataForInsert() {
      this.wrap.showSqlForInsert(this.tableDetail, this.selects);
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-table-data {
  width: 100%;
  height: 100%;
}
.toolbox-database-table-data-table {
  width: 100%;
  height: 100%;
}

.toolbox-database-table-data .el-table,
.toolbox-database-table-data .el-table__expanded-cell {
  background: transparent;
}
.toolbox-database-table-data .el-table th.el-table__cell {
  background: transparent;
}
.toolbox-database-table-data .el-table td.el-table__cell,
.toolbox-database-table-data .el-table th.el-table__cell.is-leaf {
  border-bottom: 1px solid #473939;
}
.toolbox-database-table-data .el-table th,
.toolbox-database-table-data .el-table tr {
  background: transparent;
}
.toolbox-database-table-data .el-table .el-table__row:hover td.el-table__cell {
  background-color: #473939;
}
.toolbox-database-table-data .el-table {
  color: unset;
}
.toolbox-database-table-data .el-table--mini td,
.toolbox-database-table-data .el-table--mini th {
  padding: 3px 0;
}
.toolbox-database-table-data .el-table--border::after,
.toolbox-database-table-data .el-table--group::after,
.toolbox-database-table-data .el-table::before {
  background-color: transparent;
}
.toolbox-database-table-data .el-table--border,
.toolbox-database-table-data .el-table--group {
  border: 1px solid #404040;
}
.toolbox-database-table-data .el-table td,
.toolbox-database-table-data .el-table th.is-leaf {
  border-bottom: 1px solid #404040;
}
.toolbox-database-table-data .el-table--border .el-table__cell,
.toolbox-database-table-data
  .el-table__body-wrapper
  .el-table--border.is-scrolling-left
  ~ .el-table__fixed {
  border-right: 1px solid #404040;
}
.toolbox-database-table-data
  .el-table--border
  th.el-table__cell.gutter:last-of-type {
  border-bottom: 1px solid #404040;
}
.toolbox-database-table-data .el-divider {
  background-color: #404040;
}
.toolbox-database-table-data .el-loading-mask {
  background-color: rgba(255, 255, 255, 0.1);
}

.toolbox-database-table-data .el-table .el-table__cell {
  padding: 0px;
}
.toolbox-database-table-data .el-table .cell {
  white-space: nowrap;
}
.toolbox-database-table-data .el-table tr,
.toolbox-database-table-data .el-table th.el-table__cell {
  background-color: #2d2d2d;
}
.toolbox-database-table-data .el-table tbody .cell {
  padding-left: 0px !important;
}
.toolbox-database-table-data .el-table input {
  color: #fff;
  border: 0px dashed transparent;
  background-color: transparent;
  width: 100%;
  padding: 0px 5px;
  box-sizing: border-box;
  outline: none;
  font-size: 12px;
}
.toolbox-database-table-data-pagination .el-pagination {
  color: #929292;
}
.toolbox-database-table-data-pagination .el-pagination .el-input__inner {
  height: 22px !important;
  line-height: 22px !important;
  background-color: transparent;
  border: 0px solid #dcdfe6;
}
.toolbox-database-table-data-pagination .el-pagination .el-input__icon {
  line-height: 22px !important;
}
.toolbox-database-table-data-pagination .el-pagination button:disabled {
  background-color: transparent;
}
.toolbox-database-table-data-pagination .el-input.is-disabled .el-input__inner {
  background-color: transparent;
}
.toolbox-database-table-data-pagination .el-pagination .btn-next,
.toolbox-database-table-data-pagination .el-pagination .btn-prev {
  background-color: transparent;
}
.toolbox-database-table-data-pagination .el-pagination .btn-next,
.toolbox-database-table-data-pagination .el-pagination .btn-prev {
  color: #929292;
}
.toolbox-database-table-data-pagination .el-pager li.btn-quicknext,
.toolbox-database-table-data-pagination .el-pager li.btn-quickprev {
  color: #929292;
}
.toolbox-database-table-data-pagination .el-pager li {
  background: transparent;
}
</style>

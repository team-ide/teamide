<template>
  <div class="database-table-detail">
    <template v-if="tableDetail != null">
      <el-form ref="form" inline size="mini">
        <el-form-item label="表名" class="mgb-0">
          <template v-if="!isInsert">
            <el-input v-model="tableName" readonly> </el-input>
          </template>
          <template v-else>
            <el-input v-model="tableDetail.tableName" @change="change">
            </el-input>
          </template>
        </el-form-item>
        <el-form-item label="注释" class="mgb-0">
          <el-input v-model="tableDetail.tableComment" @change="change">
          </el-input>
        </el-form-item>
      </el-form>
      <div class="database-table-detail-tabs-box">
        <el-tabs v-model="activeName">
          <el-tab-pane label="字段" name="column">
            <div style="height: 100%">
              <div class="pdtb-5 ft-13">
                <span class="color-grey">字段列表</span>
                <div class="tm-link color-green mgl-10" @click="addColumn({})">
                  新增
                </div>
              </div>
              <el-table
                :data="columnList"
                :border="true"
                style="width: 100%"
                height="calc(100% - 30px)"
                size="mini"
              >
                <el-table-column label="字段名" fixed>
                  <template slot-scope="scope">
                    <div class="">
                      <el-input
                        v-model="scope.row.columnName"
                        type="text"
                        @change="change"
                      />
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="类型" width="120">
                  <template slot-scope="scope">
                    <div class="">
                      <el-select
                        placeholder="选中类型"
                        v-model="scope.row.columnDataType"
                        @change="change"
                        filterable
                      >
                        <el-option
                          v-for="(one, index) in columnTypeInfoList"
                          :key="index"
                          :value="one.name"
                          :label="one.name"
                        >
                        </el-option>
                      </el-select>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="长度" width="70">
                  <template slot-scope="scope">
                    <div class="">
                      <el-input
                        v-model="scope.row.columnLength"
                        type="text"
                        @change="change"
                      />
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="小数点" width="70">
                  <template slot-scope="scope">
                    <div class="">
                      <el-input
                        v-model="scope.row.columnDecimal"
                        type="text"
                        @change="change"
                      />
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="必填" width="60">
                  <template slot-scope="scope">
                    <div class="">
                      <el-switch
                        v-model="scope.row.columnNotNull"
                        @change="change"
                      />
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="主键" width="60">
                  <template slot-scope="scope">
                    <div class="">
                      <el-switch
                        v-model="scope.row.primaryKey"
                        @change="change"
                      />
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="默认值">
                  <template slot-scope="scope">
                    <div class="">
                      <el-input
                        v-model="scope.row.columnDefault"
                        type="text"
                        @change="change"
                      />
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="注释">
                  <template slot-scope="scope">
                    <div class="">
                      <el-input
                        v-model="scope.row.columnComment"
                        type="text"
                        @change="change"
                      />
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="200px">
                  <template slot-scope="scope">
                    <div
                      class="tm-link color-grey mglr-5"
                      @click="upColumn(scope.row)"
                    >
                      上移
                    </div>
                    <div
                      class="tm-link color-grey mglr-5"
                      @click="downColumn(scope.row)"
                    >
                      下移
                    </div>
                    <div
                      class="tm-link color-grey mglr-5"
                      @click="addColumn({}, scope.row)"
                    >
                      插入字段
                    </div>
                    <div
                      class="tm-link color-red mglr-5"
                      @click="removeColumn(scope.row)"
                    >
                      删除
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-tab-pane>
          <el-tab-pane label="索引" name="index">
            <div style="height: 100%">
              <div class="pdtb-5 ft-13">
                <span class="color-grey">索引列表</span>
                <div class="tm-link color-green mgl-10" @click="addIndex({})">
                  新增
                </div>
              </div>
              <el-table
                :data="indexList"
                :border="true"
                style="width: 100%"
                height="calc(100% - 30px)"
                size="mini"
              >
                <el-table-column label="索引名称" fixed width="200">
                  <template slot-scope="scope">
                    <div class="">
                      <el-input
                        v-model="scope.row.indexName"
                        type="text"
                        @change="change"
                      />
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="类型" width="120">
                  <template slot-scope="scope">
                    <div class="">
                      <el-select
                        placeholder="普通索引"
                        v-model="scope.row.indexType"
                        @change="change"
                      >
                        <el-option value="" label="普通索引"></el-option>

                        <el-option
                          v-for="(one, index) in indexTypeInfoList"
                          :key="index"
                          :value="one.name"
                          :label="one.name"
                        >
                        </el-option>
                      </el-select>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="字段">
                  <template slot-scope="scope">
                    <div class="">
                      <el-select
                        placeholder="选中字段"
                        v-model="scope.row.columnNames"
                        @change="change"
                        filterable
                        multiple
                        style="width: 100%"
                      >
                        <el-option
                          v-for="(one, index) in tableDetail.columnList"
                          :key="index"
                          :value="one.columnName"
                          :label="one.columnName"
                        >
                        </el-option>
                      </el-select>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="注释" width="200">
                  <template slot-scope="scope">
                    <div class="">
                      <el-input
                        v-model="scope.row.indexComment"
                        type="text"
                        @change="change"
                      />
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="80">
                  <template slot-scope="scope">
                    <div
                      class="tm-link color-red mglr-5"
                      @click="removeIndex(scope.row)"
                    >
                      删除
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-tab-pane>
          <el-tab-pane label="SQL（ 只作展示 ）" name="sql">
            <div class="" style="height: 100%">
              <el-form class="" size="mini" inline>
                <Pack
                  :source="source"
                  :toolboxWorker="toolboxWorker"
                  :form="formData"
                  :change="toLoad"
                >
                </Pack>
              </el-form>
              <div style="height: calc(100% - 60px)">
                <Editor
                  ref="Editor"
                  :source="source"
                  :value="showSQL"
                  language="sql"
                ></Editor>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </template>
  </div>
</template>

<script>
import Pack from "./Pack";

export default {
  components: { Pack },
  props: [
    "source",
    "toolboxWorker",
    "onChange",
    "onError",
    "formData",
    "isInsert",
    "getFormData",
    "columnTypeInfoList",
    "indexTypeInfoList",
  ],
  data() {
    return {
      tableDetail: null,
      columnList: [],
      indexList: [],
      activeName: "column",
      showSQL: null,
      tableName: null,
    };
  },
  computed: {},
  watch: {
    "formData.targetDatabaseType"() {
      this.toLoad();
    },
    "formData.tableName"() {
      this.toLoad();
    },
    "formData.ownerName"() {
      this.toLoad();
    },
    isInsert() {
      this.toLoad();
    },
  },
  methods: {
    init(ownerName, tableName, tableDetail) {
      this.ownerName = ownerName;
      this.tableName = tableName;
      this.tableDetail = tableDetail;
      this.initData();
      this.$nextTick(() => {
        this.canLoadSql = true;
      });
    },
    initData() {
      this.columnList = [];
      this.indexList = [];
      let tableDetail = this.tableDetail;
      tableDetail.columnList.forEach((column, i) => {
        if (column.default == "") {
          column.default = null;
        }
        if (column.deleted) {
          if (this.columnList.indexOf(column) >= 0) {
            this.columnList.splice(this.columnList.indexOf(column), 1);
          }
        } else {
          if (this.columnList.indexOf(column) < 0) {
            this.columnList.push(column);
          }
        }
      });
      tableDetail.indexList.forEach((index) => {
        if (index.deleted) {
          if (this.indexList.indexOf(index) >= 0) {
            this.indexList.splice(this.indexList.indexOf(index), 1);
          }
        } else {
          if (this.indexList.indexOf(index) < 0) {
            this.indexList.push(index);
          }
        }
      });
    },
    change() {
      this.callChange();
    },
    upColumn(column) {
      this.tool.up(this, "columnList", column);
      this.tool.up(this.tableDetail, "columnList", column);

      this.initData();
      this.callChange();
    },
    downColumn(column) {
      this.tool.down(this, "columnList", column);
      this.tool.down(this.tableDetail, "columnList", column);

      this.initData();
      this.callChange();
    },

    addColumn(column, after) {
      column = column || {};
      column.columnName = column.columnName || "";
      column.columnDataType = column.columnDataType || "varchar";
      column.columnLength = column.columnLength || 250;
      column.columnDecimal = column.columnDecimal || 0;
      column.primaryKey = column.primaryKey || false;
      column.columnNotNull = column.columnNotNull || false;
      column.columnDefault = column.columnDefault || "";
      column.columnComment = column.columnComment || "";

      let appendIndex = this.tableDetail.columnList.indexOf(after);
      if (appendIndex < 0) {
        appendIndex = this.tableDetail.columnList.length;
      } else {
        appendIndex++;
      }
      this.tableDetail.columnList.splice(appendIndex, 0, column);
      this.initData();
      this.callChange();
    },
    removeColumn(column) {
      if (column.oldColumn == null) {
        let findIndex = this.tableDetail.columnList.indexOf(column);
        if (findIndex >= 0) {
          this.tableDetail.columnList.splice(findIndex, 1);
        }
      } else {
        column.deleted = true;
      }
      this.initData();
      this.callChange();
    },
    addIndex(index, after) {
      index = index || {};
      index.indexName = index.indexName || "";
      index.indexType = index.indexType || "";
      index.indexComment = index.indexComment || "";
      index.columnNames = index.columnNames || [];

      let appendIndex = this.tableDetail.indexList.indexOf(after);
      if (appendIndex < 0) {
        appendIndex = this.tableDetail.indexList.length;
      } else {
        appendIndex++;
      }
      this.tableDetail.indexList.splice(appendIndex, 0, index);
      this.initData();
      this.callChange();
    },
    removeIndex(index) {
      if (index.oldIndex == null) {
        let findIndex = this.tableDetail.indexList.indexOf(index);
        if (findIndex >= 0) {
          this.tableDetail.indexList.splice(findIndex, 1);
        }
      } else {
        index.deleted = true;
      }
      this.initData();
      this.callChange();
    },
    callChange() {
      this.onChange(this.tableDetail);
      this.toLoad();
    },
    async toLoad() {
      if (!this.canLoadSql) {
        return;
      }
      this.showSQL = "";
      let res = await this.loadSqls();
      let sqlList = res.sqlList || [];
      sqlList.forEach((sql) => {
        this.showSQL += sql + ";\n\n";
      });
      this.$refs.Editor.setValue(this.showSQL);
    },
    async loadSqls() {
      let data = this.getFormData();
      let res = null;
      if (this.isInsert) {
        res = await this.toolboxWorker.work("tableCreateSql", data);
      } else {
        data.ownerName = this.ownerName;
        data.tableName = this.tableName;
        res = await this.toolboxWorker.work("tableUpdateSql", data);
      }
      this.error = null;
      if (res.code != 0) {
        this.onError(res.msg);
        return;
      }
      return res.data || {};
    },
  },
  created() {},
  mounted() {},
};
</script>

<style>
.database-table-detail {
  width: 100%;
  height: 100%;
}
.database-table-detail-tabs-box {
  width: 100%;
  height: calc(100% - 40px);
}
.database-table-detail-tabs-box .el-tabs {
  height: 100%;
}
.database-table-detail-tabs-box .el-tabs__content {
  height: calc(100% - 40px);
}
.database-table-detail-tabs-box .el-tab-pane {
  height: 100%;
}
.database-table-detail-tabs-box .el-tabs__item {
  color: darkgrey;
}
.database-table-detail-tabs-box .el-tabs__item.is-active {
  color: inherit;
}
.database-table-detail-tabs-box .el-tabs__active-bar {
  background-color: #ffffff;
}
.database-table-detail-tabs-box .el-tabs__nav-wrap::after {
  background-color: transparent;
}
</style>

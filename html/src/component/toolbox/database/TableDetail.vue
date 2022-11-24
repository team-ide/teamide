<template>
  <div class="database-table-detail">
    <div class="" v-if="tableDetail != null">
      <el-form ref="form" inline size="mini">
        <el-form-item label="表名" class="mgb-0">
          <el-input v-model="tableDetail.tableName" @change="change">
          </el-input>
        </el-form-item>
        <el-form-item label="注释" class="mgb-0">
          <el-input v-model="tableDetail.tableComment" @change="change">
          </el-input>
        </el-form-item>
      </el-form>
      <div class="">
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
          size="mini"
        >
          <el-table-column label="字段名">
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
                <el-switch v-model="scope.row.columnNotNull" @change="change" />
              </div>
            </template>
          </el-table-column>
          <el-table-column label="主键" width="60">
            <template slot-scope="scope">
              <div class="">
                <el-switch v-model="scope.row.primaryKey" @change="change" />
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
      <div class="">
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
          size="mini"
        >
          <el-table-column label="索引名称" width="200">
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
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: [
    "source",
    "toolboxWorker",
    "onChange",
    "columnTypeInfoList",
    "indexTypeInfoList",
  ],
  data() {
    return {
      tableDetail: null,
      columnList: [],
      indexList: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    init(tableDetail) {
      this.tableDetail = tableDetail;
      this.initData();
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
      this.columnList.forEach((column, i) => {
        delete column.beforeColumn;
        if (
          column.beforeColumn_ != null &&
          column.beforeColumn_ != column.oldBeforeColumn &&
          this.columnList.indexOf(column.beforeColumn_) >= 0
        ) {
          column.beforeColumn = column.beforeColumn_.name;
        } else {
          delete column.beforeColumn_;
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
      this.onChange(this.tableDetail);
    },
    upColumn(column) {
      this.tool.up(this, "columnList", column);
      this.tool.up(this.tableDetail, "columnList", column);

      let findIndex = this.columnList.indexOf(column);
      if (findIndex == 0) {
        column.beforeColumn_ = null;
        if (this.columnList.length > 1) {
          this.columnList[1].beforeColumn_ = column;
        }
      } else {
        column.beforeColumn_ = this.columnList[findIndex - 1];
      }

      this.initData();
      this.onChange(this.tableDetail);
    },
    downColumn(column) {
      this.tool.down(this, "columnList", column);
      this.tool.down(this.tableDetail, "columnList", column);

      let findIndex = this.columnList.indexOf(column);
      if (findIndex == 0) {
        column.beforeColumn_ = null;
        if (this.columnList.length > 1) {
          this.columnList[1].beforeColumn_ = column;
        }
      } else {
        column.beforeColumn_ = this.columnList[findIndex - 1];
      }

      this.initData();
      this.onChange(this.tableDetail);
    },

    addColumn(column, after) {
      column = column || {};
      column.columnName = column.columnName || "";
      column.columnType = column.columnType || "varchar";
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
      column.beforeColumn_ = this.tableDetail.columnList[appendIndex - 1];
      this.initData();
      this.onChange(this.tableDetail);
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
      this.onChange(this.tableDetail);
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
      this.onChange(this.tableDetail);
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
      this.onChange(this.tableDetail);
    },
  },
  created() {},
  mounted() {},
};
</script>

<style>
</style>

<template>
  <div class="toolbox-database-export">
    <div class="app-scroll-bar pd-10" style="height: calc(100% - 120px)">
      <el-form ref="form" :model="form" size="mini" inline>
        <el-form-item label="导出类型">
          <el-select v-model="form.exportType" style="width: 100px">
            <el-option
              v-for="(one, index) in exportTypes"
              :key="index"
              :value="one.value"
              :label="one.text"
              :disabled="one.disabled"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <template v-if="form.exportType == 'csv' || form.exportType == 'txt'">
          <el-form-item label="列分割字符">
            <el-input v-model="form.separator" style="width: 100px"> </el-input>
          </el-form-item>
          <el-form-item label="换行符转换">
            <el-input v-model="form.linefeed" style="width: 100px"> </el-input>
          </el-form-item>
        </template>
        <template v-if="form.exportType == 'sql'">
          <Pack
            :source="source"
            :toolboxWorker="toolboxWorker"
            :form="form"
            :change="toLoad"
          >
          </Pack>
        </template>
      </el-form>

      <template v-for="(owner, ownerIndex) in form.owners">
        <div :key="ownerIndex">
          <div>
            <el-form ref="form" :model="form" size="mini" inline>
              <el-form-item label="库名称">
                <el-select v-model="owner.ownerName" style="width: 150px">
                  <el-option
                    v-for="(one, index) in ownerList"
                    :key="index"
                    :value="one.ownerName"
                    :label="one.ownerName"
                  >
                  </el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="导出名称">
                <el-input v-model="owner.exportName" style="width: 150px">
                </el-input>
              </el-form-item>
            </el-form>
          </div>

          <template v-for="(table, tableIndex) in owner.tables">
            <div :key="tableIndex">
              <div>
                <el-form ref="form" :model="form" size="mini" inline>
                  <el-form-item label="库名称">
                    <el-select v-model="table.tableName" style="width: 150px">
                      <el-option
                        v-for="(one, index) in owner.tableList"
                        :key="index"
                        :value="one.tableName"
                        :label="one.tableName"
                      >
                      </el-option>
                    </el-select>
                  </el-form-item>
                  <el-form-item label="导出名称">
                    <el-input v-model="table.exportName" style="width: 150px">
                    </el-input>
                  </el-form-item>
                </el-form>
              </div>
            </div>
          </template>
        </div>
      </template>

      <template
        v-if="
          form.exportType == 'excel' ||
          form.exportType == 'sql' ||
          form.exportType == 'csv'
        "
      >
        <div v-if="tableDetail != null" class="mgt-20">
          <div class="mgb-10">
            <div class="tm-link color-grey" @click="addExportColumn({})">
              添加
            </div>
          </div>
          <el-table
            :data="exportColumnList"
            border
            style="width: 100%"
            size="mini"
          >
            <el-table-column label="字段">
              <template slot-scope="scope">
                <div class="">
                  <el-select v-model="scope.row.columnName" style="width: 100%">
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
            <el-table-column label="导出名称（列名，字段名）">
              <template slot-scope="scope">
                <div class="">
                  <el-input v-model="scope.row.exportName" type="text" />
                </div>
              </template>
            </el-table-column>
            <el-table-column label="导出固定值（函数脚本，默认为查询出的值）">
              <template slot-scope="scope">
                <div class="">
                  <el-input v-model="scope.row.value" type="text" />
                </div>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200px">
              <template slot-scope="scope">
                <div
                  class="tm-link color-grey mglr-5"
                  @click="upExportColumn(scope.row)"
                >
                  上移
                </div>
                <div
                  class="tm-link color-grey mglr-5"
                  @click="downExportColumn(scope.row)"
                >
                  下移
                </div>
                <div
                  class="tm-link color-grey mglr-5"
                  @click="addExportColumn({}, scope.row)"
                >
                  插入
                </div>
                <div
                  class="tm-link color-red mglr-5"
                  @click="removeExportColumn(scope.row)"
                >
                  删除
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </template>
    </div>

    <div class="pdlr-10 mgt-10" v-if="taskKey == null">
      <div class="tm-btn bg-green" @click="toExport">开始</div>
    </div>
  </div>
</template>


<script>
import Pack from "./Pack";

export default {
  components: { Pack },
  props: [
    "source",
    "toolboxWorker",
    "extend",
    "ownerName",
    "tableName",
    "owners",
    "columnTypeInfoList",
    "indexTypeInfoList",
  ],
  data() {
    return {
      ready: false,
      exportTypes: [
        { text: "SQL", value: "sql" },
        { text: "Excel", value: "excel" },
        { text: "CSV", value: "csv" },
        { text: "Txt", value: "txt" },
      ],
      form: {
        exportType: "excel",
        exportOwnerName: "",
        exportTableName: "",

        separator: "|:-:|",
        linefeed: "|:-n-:|",
        targetDatabaseType: "",
        appendOwnerName: true,
        ownerNamePackChar: "",
        tableNamePackChar: "",
        columnNamePackChar: "",
        sqlValuePackChar: "",

        owners: [],
      },
      ownerList: [],
      exportColumnList: null,
      tableDetail: null,
      taskKey: null,
      task: null,
    };
  },
  computed: {},
  watch: {
    "form.exportType"() {
      if (this.form.exportType == "txt") {
        this.form.separator = "|:-:|";
      } else if (this.form.exportType == "csv") {
        this.form.separator = ",";
      }
    },
  },
  methods: {
    async init() {
      let ownerList = [];
      if (this.tool.isNotEmpty(this.ownerName)) {
        let owner = {
          ownerName: this.ownerName,
          tables: null,
        };

        if (this.tool.isNotEmpty(this.tableName)) {
          let tableDetail = await this.toolboxWorker.getTableDetail(
            owner.ownerName,
            this.tableName
          );
          if (tableDetail == null) {
            tableDetail = {
              tableName: this.tableName,
              columnList: [],
            };
          }
          owner.tables = [];
          owner.tables.push(tableDetail);
        }
        ownerList.push(owner);
      } else {
        ownerList = await this.toolboxWorker.loadOwners();
      }
      // this.exportColumnList = [];

      // if (this.tableDetail) {
      //   this.tableDetail.columnList.forEach((column) => {
      //     let exportColumn = {};
      //     exportColumn.columnName = column.columnName;
      //     exportColumn.exportName = column.columnName;
      //     exportColumn.value = null;
      //     this.exportColumnList.push(exportColumn);
      //   });
      // }

      this.ready = true;
    },
    async newExportTable(owner, table) {
      let exportTable = {
        tableName: null,
        exportName: null,
        columnList: [],
        exportColumnList: [],
      };
      if (table.columnList == null) {
        let detail = await this.toolboxWorker.getTableDetail(
          owner.ownerName,
          table.tableName
        );
        if (detail) {
          Object.assign(table, detail);
        }
        table.columnList = table.columnList || [];
      }
      exportTable.columnList = table.columnList || [];
      table.columnList.forEach((column) => {
        let exportColumn = {};
        exportColumn.columnName = column.columnName;
        exportColumn.exportName = column.columnName;
        exportColumn.value = null;
        exportTable.exportColumnList.push(exportColumn);
      });
    },

    upExportColumn(exportColumn) {
      this.tool.up(this, "exportColumnList", exportColumn);
    },
    downExportColumn(exportColumn) {
      this.tool.down(this, "exportColumnList", exportColumn);
    },
    addExportColumn(exportColumn, after) {
      exportColumn = exportColumn || {};
      exportColumn.columnName = exportColumn.columnName || "";
      exportColumn.exportName = exportColumn.exportName || "";
      exportColumn.value = exportColumn.value || "";

      let appendIndex = this.exportColumnList.indexOf(after);
      if (appendIndex < 0) {
        appendIndex = this.exportColumnList.length;
      } else {
        appendIndex++;
      }
      this.exportColumnList.splice(appendIndex, 0, exportColumn);
    },
    removeExportColumn(exportColumn) {
      let findIndex = this.exportColumnList.indexOf(exportColumn);
      if (findIndex >= 0) {
        this.exportColumnList.splice(findIndex, 1);
      }
    },
    async toExport() {
      if (this.task != null) {
        this.cleanTask(this.task.key);
      }
      this.task = null;
      this.taskKey = null;
      let res = await this.doExport();
      this.taskKey = res.taskKey;
      this.loadStatus();
    },
    async doExport() {
      let param = Object.assign({}, this.form);
      this.toolboxWorker.formatParam(param);

      param.ownerName = this.ownerName;
      param.tableName = this.tableName;

      if (this.tableDetail) {
        param.columnList = this.tableDetail.columnList;
      }
      param.exportColumnList = this.exportColumnList;

      let res = await this.toolboxWorker.work("export", param);
      res.data = res.data || {};
      return res.data;
    },
    async loadStatus() {
      if (this.taskKey == null) {
        return;
      }
      if (this.task != null && this.task.isEnd) {
        this.taskKey = null;
        this.cleanTask();
        return;
      }
      let param = {
        taskKey: this.taskKey,
      };
      let res = await this.toolboxWorker.work("exportStatus", param);
      res.data = res.data || {};
      this.task = res.data.task;
      setTimeout(this.loadStatus, 100);
    },
    async stopTask() {
      if (this.taskKey == null) {
        return;
      }
      let param = {
        taskKey: this.taskKey,
      };
      await this.toolboxWorker.work("exportStop", param);
    },
    async cleanTask(taskKey) {
      let param = {
        taskKey: taskKey,
      };
      await this.toolboxWorker.work("exportClean", param);
    },
    toDownload() {
      if (this.task == null) {
        this.tool.error("任务数据丢失");
        return;
      }
      let url =
        this.source.api +
        "api/toolbox/database/export/download?taskKey=" +
        encodeURIComponent(this.task.key) +
        "&jwt=" +
        encodeURIComponent(this.tool.getJWT());
      window.location.href = url;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-export {
  width: 100%;
  height: 100%;
}
</style>

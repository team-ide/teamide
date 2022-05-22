<template>
  <div class="toolbox-database-export">
    <div class="scrollbar pd-10" style="height: calc(100% - 120px)">
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
        <template v-if="form.exportType == 'sql'">
          <el-form-item label="追加库名">
            <el-switch v-model="form.appendDatabase"> </el-switch>
          </el-form-item>
          <template v-if="form.appendDatabase">
            <el-form-item label="导出库名（导出SQL文件拼接的库名）">
              <el-input v-model="form.exportDatabase" style="width: 120px">
              </el-input>
            </el-form-item>
            <el-form-item label="库名包装">
              <el-select
                placeholder="不包装"
                v-model="form.databasePackingCharacter"
                style="width: 90px"
              >
                <el-option
                  v-for="(one, index) in packingCharacters"
                  :key="index"
                  :value="one.value"
                  :label="one.text"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </template>
          <el-form-item label="导出表名（导出SQL文件拼接的表名）">
            <el-input v-model="form.exportTable" style="width: 120px">
            </el-input>
          </el-form-item>
          <el-form-item label="表名包装">
            <el-select
              placeholder="不包装"
              v-model="form.tablePackingCharacter"
              style="width: 90px"
            >
              <el-option
                v-for="(one, index) in packingCharacters"
                :key="index"
                :value="one.value"
                :label="one.text"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="字段包装">
            <el-select
              placeholder="不包装"
              v-model="form.columnPackingCharacter"
              style="width: 90px"
            >
              <el-option
                v-for="(one, index) in packingCharacters"
                :key="index"
                :value="one.value"
                :label="one.text"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="字符值包装">
            <el-select
              v-model="form.stringPackingCharacter"
              style="width: 60px"
            >
              <el-option
                v-for="(one, index) in stringPackingCharacters"
                :key="index"
                :value="one.value"
                :label="one.text"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="日期函数">
            <el-select v-model="form.dateFunction" style="width: 120px">
              <el-option
                v-for="(one, index) in dateFunctions"
                :key="index"
                :value="one.value"
                :label="one.text"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </template>
      </el-form>
      <template v-if="form.exportType == 'excel' || form.exportType == 'sql'">
        <div
          v-if="tableDetail != null"
          class="mgt-20 toolbox-database-table-data-table"
        >
          <div class="mgb-10">
            <div class="tm-link color-grey" @click="addExportColumn">添加</div>
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
                  <el-select
                    v-model="scope.row.column"
                    style="width: 100%"
                    size="mini"
                  >
                    <el-option
                      v-for="(one, index) in tableDetail.columnList"
                      :key="index"
                      :value="one.name"
                      :label="one.name"
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
    <div class="pdlr-10 mgt-10" style="user-select: text">
      <div class="ft-12">
        <span class="color-grey">任务状态：</span>
        <template v-if="task == null">
          <span class="color-orange pdr-10">暂未开始</span>
        </template>
        <template v-else>
          <template v-if="!task.isEnd">
            <span class="color-orange pdr-10"> 处理中 </span>
          </template>
          <template v-if="task.isStop">
            <span class="color-red pdr-10"> 已停止 </span>
          </template>
          <template v-else>
            <span class="color-green pdr-10"> 执行完成 </span>
          </template>
          <span class="color-grey pdr-10">
            开始：
            <span>
              {{
                tool.formatDate(new Date(task.startTime), "yyyy-MM-dd hh:mm:ss")
              }}
            </span>
          </span>
          <template v-if="task.isEnd">
            <span class="color-grey pdr-10">
              结束：
              <span>
                {{
                  tool.formatDate(new Date(task.endTime), "yyyy-MM-dd hh:mm:ss")
                }}
              </span>
            </span>
            <span class="color-grey pdr-10">
              耗时： <span>{{ task.useTime }} 毫秒</span>
            </span>
          </template>
          <template v-if="!task.isEnd">
            <div @click="stopTask" class="color-red tm-link mgr-10">
              停止执行
            </div>
          </template>
          <div class="mgt-5">
            <span class="color-grey pdr-10">
              预计导出： <span>{{ task.dataCount }}</span>
            </span>
            <span class="color-grey pdr-10">
              已准备数据： <span>{{ task.readyDataCount }}</span>
            </span>
            <span class="color-success pdr-10">
              成功： <span>{{ task.successCount }}</span>
            </span>
            <span class="color-error pdr-10">
              异常： <span>{{ task.errorCount }}</span>
            </span>
            <template v-if="task.isEnd">
              <div class="tm-link color-green mgl-50" @click="toDownload">
                下载
              </div>
            </template>
          </div>
          <template v-if="task.error != null">
            <div class="mgt-5 color-error pdr-10">
              异常： <span>{{ task.error }}</span>
            </div>
          </template>
        </template>
      </div>
    </div>
    <div class="pdlr-10 mgt-10" v-if="taskKey == null">
      <div class="tm-btn bg-green" @click="toExport">导出</div>
    </div>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolbox", "wrap", "extend", "database", "table"],
  data() {
    return {
      ready: false,
      exportTypes: [
        { text: "SQL", value: "sql" },
        { text: "Excel", value: "excel" },
      ],
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
      dateFunctions: [
        {
          value: "",
          text: "无函数",
        },
        {
          value: "to_date('$value','yyyy-mm-dd hh24:mi:ss')",
          text: "to_date('$value','yyyy-mm-dd hh24:mi:ss')",
        },
        {
          value: "to_timestamp('$value','yyyy-mm-dd hh24:mi:ss')",
          text: "to_timestamp('$value','yyyy-mm-dd hh24:mi:ss')",
        },
      ],
      form: {
        exportType: "excel",
        appendDatabase: true,
        exportDatabase: "",
        databasePackingCharacter: "`",
        exportTable: "",
        tablePackingCharacter: "`",
        columnPackingCharacter: "`",
        stringPackingCharacter: "'",
        dateFunction: "",
      },
      exportColumnList: null,
      tableDetail: null,
      taskKey: null,
      task: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      if (this.tool.isNotEmpty(this.table)) {
        this.tableDetail = await this.wrap.getTableDetail(
          this.database,
          this.table
        );
      }
      this.form.exportDatabase = this.database;
      this.form.exportTable = this.table;
      this.exportColumnList = [];

      this.tableDetail.columnList.forEach((column) => {
        let exportColumn = {};
        exportColumn.column = column.name;
        exportColumn.exportName = column.name;
        exportColumn.value = null;
        this.exportColumnList.push(exportColumn);
      });
      this.ready = true;
    },

    upExportColumn(exportColumn) {
      this.tool.up(this, "exportColumnList", exportColumn);
    },
    downExportColumn(exportColumn) {
      this.tool.down(this, "exportColumnList", exportColumn);
    },
    addExportColumn(exportColumn, after) {
      exportColumn = exportColumn || {};
      exportColumn.column = exportColumn.column || "";
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
      param.database = this.database;
      param.table = this.table;
      param.columnList = this.tableDetail.columnList;
      param.exportColumnList = this.exportColumnList;

      let res = await this.wrap.work("export", param);
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
      let res = await this.wrap.work("exportStatus", param);
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
      await this.wrap.work("exportStop", param);
    },
    async cleanTask(taskKey) {
      let param = {
        taskKey: taskKey,
      };
      await this.wrap.work("exportClean", param);
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

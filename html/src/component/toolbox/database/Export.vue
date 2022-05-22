<template>
  <div class="toolbox-database-export">
    <div class="pd-10">
      <el-form ref="form" :model="form" size="mini" inline>
        <el-form-item label="导出类型">
          <el-select v-model="form.exportType" style="width: 100px">
            <el-option
              v-for="(one, index) in exportTypes"
              :key="index"
              :value="one.value"
              :label="one.text"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <template v-if="form.exportType == 'sql'">
          <el-form-item label="追加库名">
            <el-switch v-model="form.appendDatabase"> </el-switch>
          </el-form-item>
          <template v-if="form.appendDatabase">
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
      <div class="mgt-10">
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
                  tool.formatDate(
                    new Date(task.startTime),
                    "yyyy-MM-dd hh:mm:ss"
                  )
                }}
              </span>
            </span>
            <template v-if="task.isEnd">
              <span class="color-grey pdr-10">
                结束：
                <span>
                  {{
                    tool.formatDate(
                      new Date(task.endTime),
                      "yyyy-MM-dd hh:mm:ss"
                    )
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
            </div>
            <template v-if="task.error != null">
              <div class="mgt-5 color-error pdr-10">
                异常： <span>{{ task.error }}</span>
              </div>
            </template>
          </template>
        </div>
      </div>
      <div class="mgt-20" v-if="taskKey == null">
        <div class="tm-btn bg-green" @click="toExport">导出</div>
      </div>
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
        exportType: "sql",
        appendDatabase: true,
        databasePackingCharacter: "`",
        tablePackingCharacter: "`",
        columnPackingCharacter: "`",
        stringPackingCharacter: "'",
        dateFunction: "",
      },
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
      this.ready = true;
    },
    async toExport() {
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
    async cleanTask() {
      if (this.taskKey == null) {
        return;
      }
      let param = {
        taskKey: this.taskKey,
      };
      await this.wrap.work("exportClean", param);
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

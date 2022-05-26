<template>
  <div class="toolbox-redis-export">
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
      </el-form>
      <template v-if="form.exportType == 'excel'"> </template>
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
  props: ["source", "toolbox", "wrap", "extend"],
  data() {
    return {
      ready: false,
      exportTypes: [
        { text: "SQL", value: "sql" },
        { text: "Excel", value: "excel" },
      ],
      form: {
        exportType: "excel",
      },
      exportColumnList: null,
      taskKey: null,
      task: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      this.exportColumnList = [];

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
      param.redis = this.redis;
      param.table = this.table;
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
        "api/toolbox/redis/export/download?taskKey=" +
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
.toolbox-redis-export {
  width: 100%;
  height: 100%;
}
</style>

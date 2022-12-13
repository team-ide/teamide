<template>
  <div class="toolbox-elasticsearch-import">
    <div class="app-scroll-bar pd-10" style="height: calc(100% - 120px)">
      <el-form size="mini" inline>
        <el-form-item label="类型">
          <el-select v-model="formData.importType" style="width: 100px">
            <el-option
              v-for="(one, index) in importTypes"
              :key="index"
              :value="one.value"
              :label="one.text"
              :disabled="one.disabled"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="索引">
          <el-input v-model="formData.indexName" style="width: 200px">
          </el-input>
        </el-form-item>
        <el-form-item label="导入数量">
          <el-input v-model="formData.count" style="width: 100px"> </el-input>
        </el-form-item>
        <el-form-item label="批量导入">
          <el-input v-model="formData.batchNumber" style="width: 100px">
          </el-input>
        </el-form-item>
        <el-form-item label="ID">
          <el-input v-model="formData.id" style="width: 300px"> </el-input>
        </el-form-item>
      </el-form>
      <div>
        <div class="tm-link color-green mgr-5" @click="addImportColumn({})">
          添加字段
        </div>
      </div>
      <el-table :data="columnList" border style="width: 100%" size="mini">
        <el-table-column label="字段名称">
          <template slot-scope="scope">
            <div class="">
              <el-input v-model="scope.row.name" type="text" />
            </div>
          </template>
        </el-table-column>
        <el-table-column label="导入固定值（函数脚本，默认为查询出的值）">
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
              @click="upImportColumn(scope.row)"
            >
              上移
            </div>
            <div
              class="tm-link color-grey mglr-5"
              @click="downImportColumn(scope.row)"
            >
              下移
            </div>
            <div
              class="tm-link color-grey mglr-5"
              @click="addImportColumn({}, scope.row)"
            >
              插入
            </div>
            <div
              class="tm-link color-red mglr-5"
              @click="removeImportColumn(scope.row)"
            >
              删除
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <div class="mglr-10 mgt-10" style="user-select: text">
      <div class="ft-12">
        <span class="color-grey">任务状态：</span>
        <template v-if="task == null">
          <span class="color-orange pdr-10">暂未开始</span>
        </template>
        <template v-else>
          <template v-if="!task.isEnd">
            <span class="color-orange pdr-10"> 处理中 </span>
          </template>
          <template v-else-if="task.isStop">
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
            <div @click="stopTask()" class="color-red tm-link mgr-10">
              停止执行
            </div>
          </template>
          <div class="mgt-5">
            <span class="color-grey pdr-10">
              数据 总数/成功/失败：
              <span>
                {{ task.dataCount }}
                /
                <span class="color-green">
                  {{ task.dataSuccessCount }}
                </span>
                /
                <span class="color-red">
                  {{ task.dataErrorCount }}
                </span>
              </span>
            </span>
            <template v-if="task.isEnd">
              <div @click="taskClean()" class="color-orange tm-link mgr-10">
                删除
              </div>
            </template>
          </div>
          <template v-if="tool.isNotEmpty(task.error)">
            <div class="mgt-5 color-error pdr-10">
              异常： <span>{{ task.error }}</span>
            </div>
          </template>
        </template>
      </div>
    </div>
    <div class="pdlr-10 mgt-10" v-if="taskId == null">
      <div class="tm-btn bg-green" @click="toDo">开始</div>
    </div>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "actived", "extend", "indexName"],
  data() {
    return {
      ready: false,
      importTypes: [{ text: "数据策略", value: "strategy" }],
      formData: {
        importType: "strategy",
        indexName: null,
        count: 1,
        batchNumber: 200,
        id: "",
      },
      columnList: [],
      taskId: null,
      task: null,
    };
  },
  computed: {},
  watch: {
    "formData.importType"() {},
  },
  methods: {
    onFocus() {
      if (this.inited) {
        return;
      }
      this.$nextTick(async () => {
        this.init();
      });
    },
    packChange() {},
    async init() {
      this.inited = true;
      this.formData.indexName = this.indexName;
      this.formData.id = "'id-' + (_$index + 0)";
      this.ready = true;
    },
    upImportColumn(importColumn) {
      this.tool.up(this, "columnList", importColumn);
    },
    downImportColumn(importColumn) {
      this.tool.down(this, "columnList", importColumn);
    },
    addImportColumn(importColumn, after) {
      importColumn = importColumn || {};
      importColumn.name = importColumn.name || "name";
      importColumn.value = importColumn.value || "'name-' + (_$index + 0)";

      let appendIndex = this.columnList.indexOf(after);
      if (appendIndex < 0) {
        appendIndex = this.columnList.length;
      } else {
        appendIndex++;
      }
      this.columnList.splice(appendIndex, 0, importColumn);
    },
    removeImportColumn(importColumn) {
      let findIndex = this.columnList.indexOf(importColumn);
      if (findIndex >= 0) {
        this.columnList.splice(findIndex, 1);
      }
    },
    async toDo() {
      if (this.task != null) {
        this.taskClean();
      }
      this.task = null;
      this.taskId = null;
      let res = await this.start();
      if (res) {
        this.taskId = res.taskId;
        this.loadStatus();
      }
    },
    async start() {
      let param = this.toolboxWorker.getWorkParam(Object.assign(this.formData));

      param.batchNumber = Number(param.batchNumber);
      param.count = Number(param.count);
      if (this.tool.isEmpty(param.id)) {
        this.tool.error("请输入id值策略");
        return;
      }
      param.columnList = [];
      this.columnList.forEach((one) => {
        param.columnList.push(one);
      });

      let res = await this.server.elasticsearch.import(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      return res.data;
    },
    async loadStatus() {
      if (this.taskId == null) {
        return;
      }
      if (this.task != null && this.task.isEnd) {
        this.taskId = null;
        return;
      }
      if (this.isDestroyed) {
        return;
      }
      let param = this.toolboxWorker.getWorkParam({
        taskId: this.taskId,
      });
      let res = await this.server.elasticsearch.taskStatus(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      this.task = res.data;
      setTimeout(this.loadStatus, 100);
    },
    async stopTask() {
      if (this.task == null) {
        return;
      }
      let param = this.toolboxWorker.getWorkParam({
        taskId: this.task.taskId,
      });
      await this.server.elasticsearch.taskStop(param);
    },
    async taskClean() {
      if (this.task == null) {
        return;
      }
      let param = this.toolboxWorker.getWorkParam({
        taskId: this.task.taskId,
      });
      this.task = null;
      this.taskId = null;
      await this.server.elasticsearch.taskClean(param);
    },
  },
  created() {},
  mounted() {
    if (this.actived) {
      this.init();
    }
  },
  beforeDestroy() {
    this.isDestroyed = true;
    this.taskClean();
  },
};
</script>

<style>
.toolbox-elasticsearch-import {
  width: 100%;
  height: 100%;
}
</style>

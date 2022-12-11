<template>
  <div class="toolbox-kafka-topic">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="80px">
          <el-form class="pdt-10 pdlr-10" size="mini" inline>
            <el-form-item label="搜索" class="mgb-5">
              <el-input v-model="searchForm.pattern" style="width: 300px" />
            </el-form-item>
            <el-form-item label="" class="mgb-5">
              <div class="tm-btn tm-btn-xs bg-grey-6" @click="loadTopics">
                刷新
              </div>
              <div class="tm-btn tm-btn-xs bg-teal-8" @click="toInsert">
                新建主题
              </div>
              <div class="tm-btn tm-btn-xs bg-grey" @click="toImport()">
                导入
              </div>
              <div class="tm-btn tm-btn-xs bg-grey" @click="toExport()">
                导出
              </div>
              <div
                class="tm-btn tm-btn-xs bg-grey"
                @click="toolboxWorker.showInfo()"
              >
                信息
              </div>
            </el-form-item>
          </el-form>
        </tm-layout>
        <tm-layout height="auto">
          <template v-if="topics == null">
            <div class="text-center ft-13 pdtb-10">数据加载中，请稍后!</div>
          </template>
          <template v-else-if="topics.length == 0">
            <div class="text-center ft-13 pdtb-10">暂无匹配数据!</div>
          </template>
          <template v-else>
            <div class="text-center ft-13 pdtb-10" style="height: 40px">
              Topics （{{ topics.length }}）
              <span class="color-orange">双击查看Topic数据</span>
            </div>
            <div
              class="data-list-box app-scroll-bar"
              style="height: calc(100% - 40px); user-select: text"
            >
              <template v-for="(one, index) in topics">
                <div
                  :key="index"
                  v-if="
                    tool.isEmpty(searchForm.pattern) ||
                    one.topic
                      .toLowerCase()
                      .indexOf(searchForm.pattern.toLowerCase()) >= 0
                  "
                  class="data-list-one"
                  @click="rowClick(one)"
                  @contextmenu="dataContextmenu(one)"
                >
                  <div class="data-list-one-text">
                    {{ one.topic }}
                  </div>
                </div>
              </template>
            </div>
          </template>
        </tm-layout>
      </tm-layout>
    </template>
    <FormDialog
      ref="TopicForm"
      :source="source"
      title="主题"
      :onSave="doInsert"
    ></FormDialog>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      ready: false,
      topics: null,
      searchForm: {
        pattern: null,
      },
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.ready = true;
      this.loadTopics();
    },
    refresh() {
      this.loadTopics();
    },
    rowClick(data) {
      this.rowClickTimeCache = this.rowClickTimeCache || {};
      let nowTime = new Date().getTime();
      let clickTime = this.rowClickTimeCache[data];
      this.rowClickTimeCache[data] = nowTime;
      if (clickTime) {
        let timeout = nowTime - clickTime;
        if (timeout < 300) {
          delete this.rowClickTimeCache[data];
          this.rowDbClick(data);
        }
      }
    },
    rowDbClick(data) {
      this.toOpenTopic(data);
    },
    toOpenTopic(data) {
      let extend = {
        name: data.topic,
        title: data.topic,
        type: "data",
        topic: data.topic,
        onlyOpenOneKey: "kafka:data:topic" + data.topic,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    dataContextmenu(data) {
      let menus = [];
      menus.push({
        header: data.topic,
      });
      menus.push({
        text: "数据",
        onClick: () => {
          this.toOpenTopic(data);
        },
      });
      menus.push({
        text: "导入",
        onClick: () => {
          this.toImport(data);
        },
      });
      menus.push({
        text: "导出",
        onClick: () => {
          this.toExport(data);
        },
      });
      menus.push({
        text: "删除",
        onClick: () => {
          this.toDelete(data);
        },
      });
      if (menus.length > 0) {
        this.tool.showContextmenu(menus);
      }
    },
    async toExport(data) {
      this.tool.warn("功能还未完善，敬请期待！");
      return;
    },
    async toImport(data) {
      this.tool.warn("功能还未完善，敬请期待！");
      return;
    },
    toInsert() {
      let data = {};

      this.$refs.TopicForm.show({
        title: `创建主题`,
        form: [this.form.toolbox.kafka.topic],
        data: [data],
      });
    },
    async doInsert(dataList) {
      let data = dataList[0];
      let param = this.toolboxWorker.getWorkParam({
        topic: data.topic,
        numPartitions: Number(data.numPartitions),
        replicationFactor: Number(data.replicationFactor),
      });
      let res = await this.server.kafka.createTopic(param);
      if (res.code == 0) {
        await this.loadTopics();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    toDelete(data) {
      let msg = "确认删除";
      msg += "主题[" + data.topic + "]";
      msg += "?";
      this.tool
        .confirm(msg)
        .then(async () => {
          this.doDelete(data.topic);
        })
        .catch((e) => {});
    },
    async doDelete(topic) {
      let param = this.toolboxWorker.getWorkParam({
        topic: topic,
      });
      let res = await this.server.kafka.deleteTopic(param);
      if (res.code == 0) {
        this.tool.success("删除成功!");
        this.loadTopics();
      } else {
        this.tool.error(res.msg);
      }
    },
    async loadTopics() {
      this.topics = null;
      let param = this.toolboxWorker.getWorkParam({});
      let res = await this.server.kafka.topics(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      this.topics = res.data || [];
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-kafka-topic {
  width: 100%;
  height: 100%;
}
</style>

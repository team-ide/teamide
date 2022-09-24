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
              <div
                class="tm-btn tm-btn-xs bg-grey"
                @click="toolboxWorker.showInfo()"
              >
                新建主题
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
                    one.name
                      .toLowerCase()
                      .indexOf(searchForm.pattern.toLowerCase()) >= 0
                  "
                  class="data-list-one"
                  @click="rowClick(one)"
                  @contextmenu="dataContextmenu(one)"
                >
                  <div class="data-list-one-text">
                    {{ one.name }}
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
        name: data.name,
        title: data.name,
        type: "data",
        topic: data.name,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    dataContextmenu(data) {
      let menus = [];
      menus.push({
        header: data.name,
      });
      menus.push({
        text: "数据",
        onClick: () => {
          this.toOpenTopic(data);
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
      let param = {
        topic: data.topic,
        numPartitions: Number(data.numPartitions),
        replicationFactor: Number(data.replicationFactor),
      };
      let res = await this.toolboxWorker.work("createTopic", param);
      if (res.code == 0) {
        await this.loadTopics();
        return true;
      } else {
        return false;
      }
    },
    toDelete(data) {
      let msg = "确认删除";
      msg += "主题[" + data.name + "]";
      msg += "?";
      this.tool
        .confirm(msg)
        .then(async () => {
          this.doDelete(data.name);
        })
        .catch((e) => {});
    },
    async doDelete(topic) {
      let param = {
        topic: topic,
      };
      let res = await this.toolboxWorker.work("deleteTopic", param);
      if (res.code == 0) {
        this.tool.success("删除成功!");
        this.loadTopics();
      }
    },
    async loadTopics() {
      this.topics = null;
      let param = {};
      let res = await this.toolboxWorker.work("topics", param);
      res.data = res.data || {};
      res.data.topics = res.data.topics || [];
      let topics = [];
      res.data.topics.forEach((one) => {
        let topic = {};
        topic.name = one;

        topics.push(topic);
      });
      this.topics = topics;
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

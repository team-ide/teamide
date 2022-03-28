<template>
  <div class="toolbox-kafka-topic">
    <template v-if="ready">
      <div class="pd-10">
        <table>
          <thead>
            <tr>
              <th>Topic</th>
              <th width="120">
                <a
                  class="tm-link color-grey-3 ft-14 mglr-2"
                  @click="loadTopics()"
                >
                  <i class="mdi mdi-reload"></i>
                </a>
                <a
                  class="tm-link color-green-3 ft-14 mglr-2"
                  @click="toInsert()"
                >
                  <i class="mdi mdi-plus"></i>
                </a>
              </th>
            </tr>
          </thead>
          <tbody>
            <template v-if="topics == null">
              <tr>
                <td colspan="2">
                  <div class="text-center ft-13 pdtb-10">加载中...</div>
                </td>
              </tr>
            </template>
            <template v-if="topics.length == 0">
              <tr>
                <td colspan="2">
                  <div class="text-center ft-13 pdtb-10">暂无匹配数据!</div>
                </td>
              </tr>
            </template>
            <template v-else>
              <template v-for="(one, index) in topics">
                <tr :key="index" @click="rowClick(one)" class="tm-pointer">
                  <td>{{ one.name }}</td>
                  <td>
                    <div
                      class="tm-btn color-blue tm-btn-xs"
                      @click="toOpenTopic(one)"
                    >
                      数据
                    </div>
                    <div
                      class="tm-btn color-orange tm-btn-xs"
                      @click="toDelete(one)"
                    >
                      删除
                    </div>
                  </td>
                </tr>
              </template>
            </template>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "data", "toolboxType", "toolbox", "option", "wrap"],
  data() {
    return {
      ready: false,
      topics: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.ready = true;
      this.loadTopics();
    },
    rowClick(data) {
      let nowTime = new Date().getTime();
      let clickTime = data.clickTime;
      data.clickTime = nowTime;
      if (clickTime) {
        let timeout = nowTime - clickTime;
        if (timeout < 300) {
          data.clickTime = null;
          this.rowDbClick(data);
        }
      }
    },
    rowDbClick(data) {
      this.toOpenTopic(data);
    },
    toOpenTopic(data) {
      let tab = this.wrap.createTabByData(data);
      this.wrap.addTab(tab);
      this.wrap.doActiveTab(tab);
    },
    toInsert() {
      let data = {};
      this.wrap.showTopicForm(data, (m) => {
        let flag = this.doInsert(m);
        return flag;
      });
    },
    async doInsert(data) {
      let param = {
        topic: data.topic,
        numPartitions: Number(data.numPartitions),
        replicationFactor: Number(data.replicationFactor),
      };
      let res = await this.wrap.work("createTopic", param);
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
      let res = await this.wrap.work("deleteTopic", param);
      if (res.code == 0) {
        this.tool.info("删除成功!");
        this.loadTopics();
      }
    },
    async loadTopics() {
      this.topics = null;
      let param = {};
      let res = await this.wrap.work("topics", param);
      res.data = res.data || {};
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

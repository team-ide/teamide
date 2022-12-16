<template>
  <div class="toolbox-node-editor">
    <template v-if="type == 'node-context'">
      <NodeContext
        :source="source"
        :toolboxWorker="toolboxWorker"
        :nodeContext="nodeContext"
      >
      </NodeContext>
    </template>
    <template v-if="type == 'net-proxy-context'">
      <NetProxyContext
        :source="source"
        :toolboxWorker="toolboxWorker"
        :nodeContext="nodeContext"
      >
      </NetProxyContext>
    </template>
    <template v-else-if="type == 'node-info'">
      <Info :source="source" :serverId="serverId"></Info>
    </template>
    <NodeForm
      :source="source"
      :toolboxWorker="toolboxWorker"
      :nodeContext="nodeContext"
    >
    </NodeForm>
    <NodeNetProxyForm
      :source="source"
      :toolboxWorker="toolboxWorker"
      :nodeContext="nodeContext"
    >
    </NodeNetProxyForm>
  </div>
</template>


<script>
import Info from "./Info.vue";
import NodeContext from "./NodeContext.vue";
import NetProxyContext from "./NetProxyContext.vue";
import NodeForm from "./NodeForm.vue";
import NodeNetProxyForm from "./NodeNetProxyForm.vue";

export default {
  components: {
    Info,
    NodeContext,
    NetProxyContext,
    NodeForm,
    NodeNetProxyForm,
  },
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      type: null,
      serverId: null,
      nodeContext: {
        nodeLocalList: [],
        nodeList: [],
        nodeNetProxyList: [],
        localIpList: [],
        nodeOptionMap: {},
      },
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      if (this.extend) {
        this.type = this.extend.type;
        this.serverId = this.extend.serverId;
      }
      this.toolboxWorker.initNodeContext = this.initNodeContext;
      await this.initNodeContext();
      this.server.addListenOnEvent("node-data-change", this.onNodeDataChange);
    },
    onNodeDataChange(data) {
      try {
        if (data.type == "count") {
        } else if (data.type == "node-list") {
          this.initNodeList(data.nodeList);
        } else if (data.type == "net-proxy-list") {
          this.initNodeNetProxyList(data.netProxyList);
        }
      } catch (error) {}
    },
    refresh() {},
    async initNodeContext() {
      let data = {};
      if (this.source.login.user != null) {
        let res = await this.server.node.context({});
        if (res.code != 0) {
          this.tool.error(res.msg);
        }
        data = res.data || {};
      }

      this.nodeContext.localIpList = data.localIpList || [];
      this.initNodeList(data.nodeList);
      this.initNodeNetProxyList(data.netProxyList);
    },
    initNodeList(nodeList) {
      nodeList = nodeList || [];
      this.form.node.nodeOptions.splice(0, this.form.node.nodeOptions.length);
      let nodeOptionMap = {};
      var nodeLocalList = [];
      nodeList.forEach((one) => {
        let option = {};
        option.isStarted = one.isStarted;
        option.value = one.serverId;
        option.text = one.name;

        this.form.node.nodeOptions.push(option);
        nodeOptionMap[option.value] = option;
        if (one.isLocal == 1) {
          nodeLocalList.push(one);
        }
      });
      this.nodeContext.nodeLocalList = nodeLocalList;
      this.nodeContext.nodeList = nodeList;
      this.nodeContext.nodeOptionMap = nodeOptionMap;
      this.source.nodeList = nodeList;
    },
    initNodeNetProxyList(nodeNetProxyList) {
      nodeNetProxyList = nodeNetProxyList || [];
      nodeNetProxyList.forEach((one) => {
        one.innerMonitorData = one.innerMonitorData || {
          readSize: null,
          readSizeUnit: null,
          readLastSleep: null,
          readLastSleepUnit: null,
          readLastTimestamp: null,
          writeSize: null,
          writeSizeUnit: null,
          writeLastSleep: null,
          writeLastSleepUnit: null,
          writeLastTimestamp: null,
        };
        one.outerMonitorData = one.outerMonitorData || {
          readSize: null,
          readSizeUnit: null,
          readLastSleep: null,
          readLastSleepUnit: null,
          readLastTimestamp: null,
          writeSize: null,
          writeSizeUnit: null,
          writeLastSleep: null,
          writeLastSleepUnit: null,
          writeLastTimestamp: null,
        };
      });
      this.nodeContext.nodeNetProxyList = nodeNetProxyList;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  deactivated() {
    this.server.removeListenOnEvent("node-data-change", this.onNodeDataChange);
  },
};
</script>

<style>
.toolbox-node-editor {
  width: 100%;
  height: 100%;
}
</style>

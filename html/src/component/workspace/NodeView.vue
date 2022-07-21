<template>
  <div class="node-view-box">
    <div ref="container" class="node-view-container"></div>
  </div>
</template>

<script>
import { Graph } from "@antv/x6";
import "@antv/x6-vue-shape";

import NodeInfo from "./NodeInfo.vue";

export default {
  components: {},
  props: ["source", "nodeList", "onNodeMoved"],
  data() {
    return {
      nodeWrapList: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    nodeList() {
      this.initData();
      this.initDataView();
    },
  },
  methods: {
    init() {
      // 创建 Graph 的实例
      this.graph = new Graph({
        container: this.$refs.container,
        grid: true,
      });

      Graph.unregisterNode("node-info");
      Graph.registerNode("node-info", {
        inherit: "vue-shape",
        component: {
          template: `<node-info />`,
          components: {
            NodeInfo,
          },
        },
      });
      this.initData();
      this.initDataView();
    },

    initData() {
      let nodeList = this.nodeList || [];
      let nodeWrapList = [];
      let lastX = 300;
      let lastY = 0;
      nodeList.forEach((one) => {
        let info = one.info;
        let nodeModel = one.nodeModel;
        let nodeWrap = { info, nodeModel };
        if (nodeModel) {
          nodeWrap.id = nodeModel.serverId;
          nodeWrap.text = nodeModel.name;
        } else {
          nodeWrap.id = info.id;
          nodeWrap.text = info.id;
        }
        nodeWrap.status = 2;
        nodeWrap.statusError = null;
        if (info) {
          nodeWrap.status = info.status;
          nodeWrap.statusError = info.statusError;
          nodeWrap.connIdList = info.connNodeIdList || [];
        }
        lastY += 50;
        let compute = this.tool.computeFontSize(nodeWrap.text, "15px", "600");

        let width = compute.width + 40;
        if (width < 120) {
          width = 120;
        }
        nodeWrap.x = lastX;
        nodeWrap.y = lastY;
        if (nodeModel) {
          let option = this.tool.getOptionJSON(nodeModel.option);
          if (option.x) {
            nodeWrap.x = option.x;
          }
          if (option.y) {
            nodeWrap.y = option.y;
          }
        }

        nodeWrap.width = width;
        nodeWrap.height = 70;
        lastY += nodeWrap.height;

        nodeWrapList.push(nodeWrap);
      });

      this.nodeWrapList = nodeWrapList;
    },
    initDataView() {
      var nodeWrapList = this.nodeWrapList || [];
      let newNodeWrapMap = {};
      nodeWrapList.forEach((one) => {
        newNodeWrapMap[one.id] = one;
      });

      let graph = this.graph;
      let edges = graph.getEdges();
      edges = edges || [];
      edges.forEach((one) => {
        graph.removeEdge(one);
      });
      let gNodeList = graph.getNodes();
      gNodeList = gNodeList || [];
      let oldNodeMap = {};
      gNodeList.forEach((one) => {
        // let position = one.getPosition();
        // oldNodeMap[one.id] = { x: position.x, y: position.y };
        oldNodeMap[one.id] = one;
        if (newNodeWrapMap[one.id] == null) {
          graph.removeNode(one);
        }
      });

      let nodeMap = {};
      nodeWrapList.forEach((one) => {
        let oldNode = oldNodeMap[one.id];
        let x = one.x;
        let y = one.y;
        if (oldNode) {
          oldNode.setData(one);
          nodeMap[one.id] = oldNode;
        } else {
          let gNode = graph.addNode({
            id: one.id,
            shape: "node-info",
            x: x,
            y: y,
            width: one.width,
            height: one.height,
            source: this.source,
            data: one,
          });
          nodeMap[one.id] = gNode;
        }
      });

      nodeWrapList.forEach((one) => {
        let sourceId = one.id;
        let source = nodeMap[sourceId];
        if (source == null) {
          return;
        }
        let connIdList = one.connIdList || [];
        connIdList.forEach((connId) => {
          let targetId = connId;
          let target = nodeMap[targetId];
          if (target == null) {
            return;
          }
          // 渲染边
          graph.addEdge({
            source: sourceId,
            target: targetId,
          });
        });
      });

      graph.on("node:moved", ({ e, x, y, node, view }) => {
        this.onNodeMoved && this.onNodeMoved(node.data, { x: x, y: y });
      });
    },
  },
  created() {},
  updated() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.node-view-box {
  position: relative;
  width: 100%;
  height: 100%;
  background: #09131c;
}
.node-view-container {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>

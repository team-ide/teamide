<template>
  <div class="node-view-box">
    <div ref="container" class="node-view-container"></div>
  </div>
</template>

<script>
import { Graph } from "@antv/x6";
import "@antv/x6-vue-shape";

import NodeInfo from "./NodeInfo.vue";
import float from "js-yaml/lib/type/float";

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
      let lastX = 0;
      let lastY = 0;
      nodeList.forEach((one) => {
        let info = one.info;
        let nodeModel = one.nodeModel;
        let nodeWrap = { info, nodeModel };
        nodeWrap.isRoot = false;
        var connServerIdList = [];
        var historyConnServerIdList = [];
        if (nodeModel) {
          if (nodeModel.isRoot == 1) {
            nodeWrap.isRoot = true;
          }
          nodeWrap.id = nodeModel.serverId;
          nodeWrap.text = nodeModel.name;
          nodeWrap.serverId = nodeModel.serverId;
          if (this.tool.isNotEmpty(nodeModel.connServerIds)) {
            try {
              connServerIdList = JSON.parse(nodeModel.connServerIds);
            } catch (e) {}
          }
          if (this.tool.isNotEmpty(nodeModel.historyConnServerIds)) {
            try {
              historyConnServerIdList = JSON.parse(
                nodeModel.historyConnServerIds
              );
            } catch (e) {}
          }
        } else {
          nodeWrap.id = info.id;
          nodeWrap.text = info.id;
          nodeWrap.serverId = info.id;
        }
        nodeWrap.connServerIdList = connServerIdList;
        nodeWrap.historyConnServerIdList = historyConnServerIdList;
        nodeWrap.status = 2;
        nodeWrap.statusError = null;
        if (info) {
          nodeWrap.status = info.status;
          nodeWrap.statusError = info.statusError;
          nodeWrap.connIdList = info.connNodeIdList || [];
        }
        let compute = this.tool.computeFontSize(nodeWrap.text, "15px", "600");

        let width = compute.width + 40;
        if (width < 120) {
          width = 120;
        }
        lastY += 50;
        lastX += 50;
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
        // lastY += nodeWrap.height;

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
        let allConnIdList = [];
        if (one.connServerIdList) {
          one.connServerIdList.forEach((one) => {
            if (allConnIdList.indexOf(one) < 0) {
              allConnIdList.push(one);
            }
          });
        }
        if (one.historyConnServerIdList) {
          one.historyConnServerIdList.forEach((one) => {
            if (allConnIdList.indexOf(one) < 0) {
              allConnIdList.push(one);
            }
          });
        }
        if (one.connIdList) {
          one.connIdList.forEach((one) => {
            if (allConnIdList.indexOf(one) < 0) {
              allConnIdList.push(one);
            }
          });
        }
        allConnIdList.forEach((connId) => {
          let targetId = connId;
          let target = nodeMap[targetId];
          if (target == null) {
            return;
          }
          var isConn = one.connIdList.indexOf(connId) >= 0;
          var isConfig = one.connServerIdList.indexOf(connId) >= 0;
          var isHistory = one.historyConnServerIdList.indexOf(connId) >= 0;
          var stroke = "#8b8b8b";
          var strokeWidth = 1;
          if (isConn) {
            stroke = "#a5a5a5";
            strokeWidth = 1;
          } else if (isConfig) {
            stroke = "#626262";
          } else if (isHistory) {
            stroke = "#4a4a4a";
          }
          // 渲染边
          graph.addEdge({
            source: sourceId,
            target: targetId,
            router: "metro",
            connector: "rounded",
            attrs: {
              line: {
                stroke: stroke,
                strokeWidth: strokeWidth,
              },
            },
          });
        });
      });

      graph.on("node:moved", ({ e, x, y, node, view }) => {
        this.onNodeMoved && this.onNodeMoved(node.data, { x: x, y: y });
      });
      graph.on("node:contextmenu", ({ e, node, view }) => {
        let data = node.data;
        let menus = [];
        menus.push({
          header: data.text,
        });
        menus.push({
          text: "连接节点",
          onClick: () => {
            this.tool.toInsertConnNode(data);
          },
        });
        menus.push({
          text: "查看",
          onClick: () => {
            this.tool.showNodeInfo(data);
          },
        });

        if (menus.length > 0) {
          this.tool.showContextmenu(menus);
        }
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

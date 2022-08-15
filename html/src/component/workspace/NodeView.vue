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
  computed: {},
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
        width: "100%",
        height: "100%",
        grid: {
          visible: true,
          args: [],
        },
        scroller: {
          enabled: true,
          pannable: true,
          padding: {
            left: 0,
            top: 0,
          },
        },
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
        let model = one.model;
        let nodeWrap = { info, model };
        nodeWrap.isLocal = false;
        var connServerIdList = [];
        var historyConnServerIdList = [];
        if (model) {
          if (model.isLocal == 1) {
            nodeWrap.isLocal = true;
          }
          nodeWrap.nodeId = model.nodeId;
          nodeWrap.id = model.serverId;
          nodeWrap.text = model.name;
          nodeWrap.serverId = model.serverId;
          if (this.tool.isNotEmpty(model.connServerIds)) {
            try {
              connServerIdList = JSON.parse(model.connServerIds);
            } catch (e) {}
          }
          if (this.tool.isNotEmpty(model.historyConnServerIds)) {
            try {
              historyConnServerIdList = JSON.parse(model.historyConnServerIds);
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
        if (model) {
          let option = this.tool.getOptionJSON(model.option);
          if (option.x) {
            nodeWrap.x = option.x;
          }
          if (option.y) {
            nodeWrap.y = option.y;
          }
        }
        if (nodeWrap.x == null || nodeWrap.y == null) {
          lastX += 50;
          nodeWrap.x = lastX;
          lastY += 50;
          nodeWrap.y = lastY;
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
          target.isConn = true;
          target.isConfig = true;
          target.isHistory = true;

          var stroke = "#8b8b8b";
          var strokeWidth = 1;
          if (isConn) {
            stroke = "#a5a5a5";
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
        let targetNode = nodeMap[data.id] || {};
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
        if (this.tool.openByExtend && data.model) {
          menus.push({
            text: "查看",
            onClick: () => {
              this.tool.openByExtend({
                toolboxType: "node",
                type: "info",
                title: "查看节点-" + data.model.name,
                serverId: data.model.serverId,
              });
              this.tool.hideNodeDialog();
            },
          });
        }

        if (data.model && !data.isLocal) {
          if (data.model.enabled == 1) {
            menus.push({
              text: "停用",
              onClick: () => {
                this.tool.toDisableNode(data);
              },
            });
          } else {
            menus.push({
              text: "启用",
              onClick: () => {
                this.tool.toEnableNode(data);
              },
            });
          }
        }
        if (
          (!data.isLocal && (targetNode.isConfig || targetNode.isHistory)) ||
          !targetNode.isConn
        ) {
          menus.push({
            text: "删除",
            onClick: () => {
              this.tool.toDeleteNode(data);
            },
          });
        }

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
  position: relative;
}

/* 滚动条样式*/
.x6-graph-scroller:hover::-webkit-scrollbar-thumb {
  box-shadow: inset 0 0 5px #333333;
  background: #333333;
}
.x6-graph-scroller:hover::-webkit-scrollbar-track {
  box-shadow: inset 0 0 5px #262626;
  background: #262626;
}
.x6-graph-scroller:hover::-webkit-scrollbar-corner {
  background: #262626;
}

.x6-graph-scroller::-webkit-scrollbar {
  width: 5px;
  height: 5px;
}
.x6-graph-scroller:hover::-webkit-scrollbar {
  width: 5px;
  height: 5px;
}
.x6-graph-scroller::-webkit-scrollbar-thumb {
  border-radius: 0px;
}
.x6-graph-scroller::-webkit-scrollbar-track {
  border-radius: 0;
}
.x6-graph-scroller::-webkit-scrollbar-corner {
  background: transparent;
}
</style>

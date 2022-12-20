<template>
  <div>
    <FormDialog
      ref="InsertNode"
      :source="source"
      title="新增节点"
      :onSave="doInsert"
    ></FormDialog>
    <FormDialog
      ref="UpdateNode"
      :source="source"
      title="编辑Node"
      :onSave="doUpdate"
    ></FormDialog>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "nodeContext"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    init() {},
    toInsertLocal() {
      this.tool.stopEvent();
      let data = {
        name: this.source.login.user.name + "-本地节点",
        bindAddress: ":21090",
        bindToken: this.tool.md5("bindToken" + new Date().getTime()),
      };

      this.$refs.InsertNode.show({
        title: `设置本地节点`,
        form: [this.form.node.local],
        isLocal: true,
        serverId: this.tool.md5("serverId" + new Date().getTime()),
        data: [data],
      });
    },
    toInsertToNode(parentNode) {
      this.tool.stopEvent();
      let data = {};
      let parentServerId = null;
      if (parentNode && parentNode.serverId) {
        parentServerId = parentNode.serverId;
      } else {
        this.tool.warn("父节点ID丢失");
        return;
      }
      this.$refs.InsertNode.show({
        title: `子节点`,
        form: [this.form.node.toNode],
        data: [data],
        parentServerId: parentServerId,
      });
    },
    toInsertFromNode(toNode) {
      this.tool.stopEvent();
      let data = {};
      let toServerId = null;
      if (toNode && toNode.serverId) {
        toServerId = toNode.serverId;
      } else {
        this.tool.warn("子节点ID丢失");
        return;
      }
      this.$refs.InsertNode.show({
        title: `父节点`,
        form: [this.form.node.fromNode],
        data: [data],
        toServerId: toServerId,
      });
    },
    async doInsert(dataList, config) {
      let data = dataList[0];
      data.parentServerId = config.parentServerId;
      data.toServerId = config.toServerId;
      if (config.isLocal) {
        data.isLocal = 1;
        data.serverId = config.serverId;
      } else {
        data.isLocal = 0;
      }
      let res = await this.server.node.insert(data);
      if (res.code == 0) {
        this.tool.success("新增成功");
        if (config.isLocal) {
          this.toolboxWorker.initNodeContext();
        }
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    toUpdate(data) {
      this.tool.stopEvent();
      this.$refs.InsertNode.show({
        title: `编辑[${data.name}]`,
        nodeId: data.nodeId,
        form: [this.form.node.toNode],
        data: [data],
      });
    },
    async doUpdate(dataList, config) {
      let data = dataList[0];
      data.nodeId = config.nodeId;
      let res = await this.server.node.update(data);
      if (res.code == 0) {
        this.tool.success("修改成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    toEnable(data) {
      this.tool.stopEvent();
      if (!data || !data.nodeId) {
        this.tool.warn("节点ID丢失");
        return;
      }
      return this.doEnable(data.nodeId);
    },
    async doEnable(nodeId) {
      let res = await this.server.node.enable({ nodeId: nodeId });
      if (res.code == 0) {
        this.tool.success("启用成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    toDisable(data) {
      this.tool.stopEvent();
      if (!data || !data.nodeId) {
        this.tool.warn("节点ID丢失");
        return;
      }
      this.tool
        .confirm("禁用[" + data.text + "]节点，相关功能将无法使用，确定禁用？")
        .then(async () => {
          return this.doDisable(data.nodeId);
        })
        .catch((e) => {});
    },
    async doDisable(nodeId) {
      let res = await this.server.node.disable({ nodeId: nodeId });
      if (res.code == 0) {
        this.tool.success("禁用成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    toDelete(data) {
      this.tool.stopEvent();
      if (!data || !data.nodeId) {
        this.tool.warn("节点ID丢失");
        return;
      }
      this.tool
        .confirm(
          "删除[" +
            data.name +
            "]节点，将删除所有关联数据且无法恢复，确定删除？"
        )
        .then(async () => {
          return this.doDelete(data.nodeId);
        })
        .catch((e) => {});
    },
    async doDelete(nodeId) {
      let res = await this.server.node.delete({ nodeId: nodeId });
      if (res.code == 0) {
        this.tool.success("删除成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doDeleteByServerId(serverId) {
      let res = await this.server.node.delete({ serverId: serverId });
      if (res.code == 0) {
        this.tool.success("删除成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    onNodeMoved(node, position) {
      if (node.nodeId == null) {
        return;
      }
      let option = this.tool.getOptionJSON(node.option);
      option.x = position.x;
      option.y = position.y;
      let optionStr = JSON.stringify(option);
      if (this.nodeContext.nodeList) {
        this.nodeContext.nodeList.forEach((one) => {
          if (one && one.nodeId == node.nodeId) {
            one.option = optionStr;
          }
        });
      }
      this.server.node.updateOption({
        nodeId: node.nodeId,
        option: optionStr,
      });
    },
  },
  created() {},
  updated() {
    this.toolboxWorker.toInsertToNode = this.toInsertToNode;
    this.toolboxWorker.toInsertFromNode = this.toInsertFromNode;
    this.toolboxWorker.toInsertLocalNode = this.toInsertLocal;
    this.toolboxWorker.toDeleteNode = this.toDelete;
    this.toolboxWorker.doDeleteNode = this.doDelete;
    this.toolboxWorker.toEnableNode = this.toEnable;
    this.toolboxWorker.toDisableNode = this.toDisable;
    this.toolboxWorker.onNodeMoved = this.onNodeMoved;
  },
  mounted() {
    this.init();
    this.toolboxWorker.toInsertToNode = this.toInsertToNode;
    this.toolboxWorker.toInsertFromNode = this.toInsertFromNode;
    this.toolboxWorker.toInsertLocalNode = this.toInsertLocal;
    this.toolboxWorker.toDeleteNode = this.toDelete;
    this.toolboxWorker.doDeleteNode = this.doDelete;
    this.toolboxWorker.toEnableNode = this.toEnable;
    this.toolboxWorker.toDisableNode = this.toDisable;
    this.toolboxWorker.onNodeMoved = this.onNodeMoved;
  },
};
</script>

<style>
</style>

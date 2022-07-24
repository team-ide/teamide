<template>
  <el-dialog
    ref="modal"
    :title="`节点`"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showBox"
    :before-close="hide"
    :fullscreen="true"
    width="100%"
    class="node-context-dialog"
  >
    <div class="node-context-box">
      <tm-layout height="100%">
        <tm-layout height="100px" class="scrollbar">
          <div class="node-context-box-header">
            <div class="node-context-box-header pdlr-20 ft-12">
              <div class="color-grey">
                节点程序下载地址:
                <span class="color-green pdlr-10">
                  https://gitee.com/teamide/teamide/releases
                </span>
                或
                <span class="color-green pdlr-10">
                  https://github.com/team-ide/teamide/releases
                </span>
              </div>
              <div class="color-grey">
                节点启动:
                <span class="color-green pdlr-10">
                  ./node -id node1 -address :21090 -token xxx
                </span>
                <span class="color-grey">
                  -id 节点ID,必须唯一 -address 节点启动绑定地址 -token
                  节点Token,用于节点直接连接鉴权
                </span>
              </div>
              <div class="color-grey">
                节点启动连接到某一个节点:
                <span class="color-green pdlr-10">
                  ./node -id node1 -address :21090 -token xxx -connAddress
                  ip:port -connToken xxx
                </span>
                <span class="color-grey">
                  -connAddress 目标节点的ip:port -connToken 目标节点的Token
                </span>
              </div>
              <template v-if="nodeRoot != null">
                <div class="color-grey">
                  当前节点连接到其它节点:
                  <span class="color-orange pdlr-10"> 右击节点进行操作 </span>
                </div>
                <div class="color-grey">
                  其它节点连接到当前节点:
                  <template v-for="(address, index) in rootAddressList">
                    <div :key="index" class="color-green">
                      ./node -id node1 -address :21090 -token xxx -connAddress
                      {{ address }} -connToken {{ nodeRoot.bindToken }}
                    </div>
                  </template>
                </div>
              </template>
            </div>
          </div>
        </tm-layout>
        <tm-layout height="auto" class="">
          <div class="node-context-body" v-if="ready">
            <template v-if="nodeRoot == null">
              <div class="text-center pdt-50">
                <div class="tm-btn bg-green tm-btn-lg" @click="toInsertRoot">
                  设置根节点
                </div>
              </div>
            </template>
            <template v-else>
              <NodeView
                :source="source"
                :nodeList="nodeList"
                :onNodeMoved="onNodeMoved"
              ></NodeView>
            </template>
          </div>
        </tm-layout>
      </tm-layout>

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
  </el-dialog>
</template>

<script>
import NodeView from "./NodeView.vue";

export default {
  components: { NodeView },
  props: ["source"],
  data() {
    return {
      showBox: false,
      nodeRoot: null,
      nodeList: [],
      loading: false,
      ready: false,
      rootAddressList: [],
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    showBox() {
      if (this.showBox) {
        this.initData();
      }
    },
    "source.nodeList"() {
      if (this.showBox) {
        this.initView();
      }
    },
  },
  methods: {
    show() {
      this.showBox = true;
    },
    showSwitch() {
      this.showBox = !this.showBox;
    },
    hide() {
      this.showBox = false;
    },
    init() {},
    async initData() {
      this.ready = false;
      this.initView();
      this.ready = true;
    },
    initView() {
      this.nodeRoot = this.source.nodeRoot;
      this.nodeList = this.source.nodeList;

      let rootAddressList = [];
      if (this.nodeRoot != null) {
        let address = this.nodeRoot.bindAddress;
        if (this.tool.isNotEmpty(address) && address.indexOf(":") >= 0) {
          let lastIndex = address.lastIndexOf(":");
          let ip = address.substring(0, lastIndex);
          let port = address.substring(lastIndex + 1);
          if (this.tool.isEmpty(ip) || ip == "0.0.0.0") {
            this.source.localIpList.forEach((localIp) => {
              rootAddressList.push(localIp + ":" + port);
            });
          } else {
            rootAddressList.push(address);
          }
        }
      }
      this.rootAddressList = rootAddressList;
    },
    nodeContextmenu(node) {
      let menus = [];
      menus.push({
        header: node.name,
      });
      menus.push({
        text: "修改",
        onClick: () => {
          this.toUpdate(node);
        },
      });
      menus.push({
        text: "删除",
        onClick: () => {
          this.toDelete(node);
        },
      });

      if (menus.length > 0) {
        this.tool.showContextmenu(menus);
      }
    },
    toInsertRoot() {
      this.tool.stopEvent();
      let data = {
        name: "根节点",
        bindAddress: ":21090",
        bindToken: this.tool.md5("bindToken" + new Date().getTime()),
      };

      this.$refs.InsertNode.show({
        title: `设置根节点`,
        form: [this.form.node.root],
        isLocal: true,
        serverId: this.tool.md5("serverId" + new Date().getTime()),
        data: [data],
      });
    },
    toInsertConnNode(parentNode) {
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
        title: `连接节点`,
        form: [this.form.node.connNode],
        data: [data],
        parentServerId: parentServerId,
      });
    },
    async doInsert(dataList, config) {
      let data = dataList[0];
      data.parentServerId = config.parentServerId;
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
          this.source.initNodeContext();
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
        form: [this.form.node.connNode],
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
      if (!data || !data.model || !data.model.nodeId) {
        this.tool.warn("节点ID丢失");
        return;
      }
      this.tool
        .confirm(
          "删除[" +
            data.model.name +
            "]节点，将删除所有关联数据且无法恢复，确定删除？"
        )
        .then(async () => {
          return this.doDelete(data.model.nodeId);
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
      if (node.model == null) {
        return;
      }
      let option = this.tool.getOptionJSON(node.model.option);
      option.x = position.x;
      option.y = position.y;
      let optionStr = JSON.stringify(option);
      this.server.node.updateOption({
        nodeId: node.model.nodeId,
        option: optionStr,
      });
    },
  },
  created() {},
  updated() {
    this.tool.showNodeDialog = this.show;
    this.tool.showSwitchNodeDialog = this.showSwitch;
    this.tool.hideNodeDialog = this.hide;
    this.tool.showNodeInfo = this.showNodeInfo;
    this.tool.toInsertConnNode = this.toInsertConnNode;
    this.tool.toDeleteNode = this.toDelete;
    this.tool.doDeleteNode = this.doDelete;
    this.tool.toEnableNode = this.toEnable;
    this.tool.toDisableNode = this.toDisable;
  },
  mounted() {
    this.init();
    this.tool.showNodeDialog = this.show;
    this.tool.showSwitchNodeDialog = this.showSwitch;
    this.tool.hideNodeDialog = this.hide;
    this.tool.showNodeInfo = this.showNodeInfo;
    this.tool.toInsertConnNode = this.toInsertConnNode;
    this.tool.toDeleteNode = this.toDelete;
    this.tool.doDeleteNode = this.doDelete;
    this.tool.toEnableNode = this.toEnable;
    this.tool.toDisableNode = this.toDisable;
  },
};
</script>

<style>
.node-context-dialog {
  top: 30px !important;
}
.node-context-dialog .el-dialog {
  background: #0f1b26;
  color: #ffffff;
}
.node-context-dialog .el-dialog__title {
  color: #ffffff;
}
.node-context-dialog .el-dialog__body {
  position: relative;
  width: 100%;
  height: calc(100% - 55px);
  padding: 0px;
}
.node-context-box {
  position: relative;
  width: 100%;
  height: 100%;
}
.node-context-box-header {
  position: relative;
  user-select: text;
  width: 100%;
  height: 100%;
}
.node-context-body {
  position: relative;
  width: 100%;
  height: 100%;
}
</style>

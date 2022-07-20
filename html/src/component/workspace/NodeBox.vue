<template>
  <div class="node-context-box" :class="{ 'node-context-box-show': showBox }">
    <div class="node-context-box-header">
      <div class="color-white ft-16 pdl-10 pdt-5">
        节点
        <span class="color-orange mgl-20 ft-12">请右击进行操作</span>
      </div>
      <div
        style="display: inline-block; position: absolute; top: 0px; right: 0px"
      >
        <span title="关闭" class="tm-link color-write mgr-0" @click="hide">
          <i class="mdi mdi-close ft-21"></i>
        </span>
      </div>
    </div>
    <div class="node-context-body">
      <template v-if="source.nodeRoot == null">
        <div class="text-center pdt-50">
          <div class="tm-btn bg-green tm-btn-lg" @click="toInsertRoot">
            设置根节点
          </div>
        </div>
      </template>
      <template v-else> </template>
    </div>
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
  props: ["source"],
  data() {
    return {
      showBox: false,
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
      await this.source.initNodeContext();
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
        bindToken: this.tool.md5("TokenTime" + new Date().getTime()),
      };

      this.$refs.InsertNode.show({
        title: `设置根节点`,
        form: [this.form.node.root],
        isRoot: true,
        serverId: this.tool.md5("TokenTime" + new Date().getTime()),
        data: [data],
      });
    },
    toInsert() {
      this.tool.stopEvent();
      let data = {
        connToken: this.tool.md5("TokenTime" + new Date().getTime()),
      };

      this.$refs.InsertNode.show({
        title: `设置根节点`,
        form: [this.form.node.connNode],
        data: [data],
      });
    },
    async doInsert(dataList, config) {
      let data = dataList[0];
      data.parentServerId = config.parentServerId;
      if (config.isRoot) {
        data.isRoot = 1;
        data.serverId = config.serverId;
      } else {
        data.isRoot = 0;
      }
      let res = await this.server.node.insert(data);
      if (res.code == 0) {
        this.tool.success("新增成功");
        this.initData();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    toUpdate(data) {
      this.tool.stopEvent();
      this.$refs.UpdateNodeGroup.show({
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
        this.initData();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
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
.node-context-box {
  position: absolute;
  top: 30px;
  width: 100%;
  background: #0f1b26;
  transition: all 0s;
  transform: scale(0);
  height: calc(100% - 30px);
  color: #ffffff;
}
.node-context-box.node-context-box-show {
  transform: scale(1);
}
.node-context-box-header {
  height: 40px;
}
.node-context-body {
  width: 100%;
  height: calc(100% - 40px);
}
</style>

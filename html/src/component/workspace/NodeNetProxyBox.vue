<template>
  <div
    class="node-net-proxy-context-box"
    :class="{ 'node-net-proxy-context-box-show': showBox }"
  >
    <div class="node-net-proxy-context-box-header">
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
    <div class="node-net-proxy-context-body"></div>
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
    toInsert() {
      this.tool.stopEvent();
      let data = {};
      let optionsJSON = {};

      this.$refs.InsertNodeGroup.show({
        title: `新增工具分组`,
        form: [this.form.node.group, this.form.node.group.option],
        data: [data, optionsJSON],
      });
    },
    toUpdate(data) {
      this.tool.stopEvent();
      this.updateGroupData = data;

      let optionsJSON = this.tool.getOptionJSON(data.option);

      this.$refs.UpdateNodeGroup.show({
        title: `编辑[${data.name}]工具分组`,
        form: [this.form.node.group, this.form.node.group.option],
        data: [data, optionsJSON],
      });
    },
    toDelete(data) {
      this.tool.stopEvent();
      this.tool
        .confirm(
          "删除工具分组[" + data.name + "]并将该分组下工具移出该组，确定删除？"
        )
        .then(async () => {
          return this.doDeleteGroup(data);
        })
        .catch((e) => {});
    },
    async doDelete(data) {
      let res = await this.server.node.group.delete(data);
      if (res.code == 0) {
        this.tool.success("删除分组成功");
        this.initData();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doUpdate(dataList, config) {
      let data = dataList[0];
      let optionJSON = dataList[1];
      data.groupId = this.updateGroupData.groupId;
      data.option = JSON.stringify(optionJSON);
      let res = await this.server.node.group.update(data);
      if (res.code == 0) {
        this.tool.success("修改分组成功");
        this.initData();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doInsert(dataList, config) {
      let data = dataList[0];
      let optionJSON = dataList[1];
      data.option = JSON.stringify(optionJSON);
      let res = await this.server.node.group.insert(data);
      if (res.code == 0) {
        this.tool.success("新增分组成功");
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
.node-net-proxy-context-box {
  position: absolute;
  top: 30px;
  width: 100%;
  background: #0f1b26;
  transition: all 0s;
  transform: scale(0);
  height: calc(100% - 30px);
  color: #ffffff;
}
.node-net-proxy-context-box.node-net-proxy-context-box-show {
  transform: scale(1);
}
.node-net-proxy-context-box-header {
  height: 40px;
}
.node-net-proxy-context-body {
  width: 100%;
  height: calc(100% - 40px);
  display: flex;
}
</style>

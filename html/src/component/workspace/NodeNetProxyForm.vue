<template>
  <div>
    <FormDialog
      ref="InsertNodeNetProxy"
      :source="source"
      :onSave="doInsert"
    ></FormDialog>
    <FormDialog
      ref="UpdateNodeNetProxy"
      :source="source"
      :onSave="doUpdate"
    ></FormDialog>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    init() {},
    toCopy(data) {
      this.tool.stopEvent();

      this.$refs.InsertNodeNetProxy.show({
        title: `设置网络代理`,
        form: [this.form.node.netProxy],
        data: [data.model],
      });
    },
    toInsert() {
      this.tool.stopEvent();
      let data = {};

      this.$refs.InsertNodeNetProxy.show({
        title: `设置网络代理`,
        form: [this.form.node.netProxy],
        data: [data],
      });
    },
    async doInsert(dataList, config) {
      let data = dataList[0];
      let res = await this.server.node.netProxy.insert(data);
      if (res.code == 0) {
        this.tool.success("新增成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    toUpdate(data) {
      this.tool.stopEvent();
      this.$refs.InsertNodeNetProxy.show({
        title: `编辑[${data.name}]网络代理`,
        netProxyId: data.netProxyId,
        form: [this.form.node.netProxy],
        data: [data],
      });
    },
    async doUpdate(dataList, config) {
      let data = dataList[0];
      data.netProxyId = config.netProxyId;
      let res = await this.server.node.netProxy.update(data);
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
      if (!data || !data.model || !data.model.netProxyId) {
        this.tool.warn("代理ID丢失");
        return;
      }
      return this.doEnable(data.model.netProxyId);
    },
    async doEnable(netProxyId) {
      let res = await this.server.node.netProxy.enable({
        netProxyId: netProxyId,
      });
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
      if (!data || !data.model || !data.model.netProxyId) {
        this.tool.warn("代理ID丢失");
        return;
      }
      this.tool
        .confirm(
          "禁用[" + data.model.name + "]代理，相关功能将无法使用，确定禁用？"
        )
        .then(async () => {
          return this.doDisable(data.model.netProxyId);
        })
        .catch((e) => {});
    },
    async doDisable(netProxyId) {
      let res = await this.server.node.netProxy.disable({
        netProxyId: netProxyId,
      });
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
      if (!data || !data.model || !data.model.netProxyId) {
        this.tool.warn("代理ID丢失");
        return;
      }
      this.tool
        .confirm(
          "删除[" +
            data.model.name +
            "]代理，将删除关联数据且无法恢复，确定删除？"
        )
        .then(async () => {
          return this.doDelete(data.model.netProxyId);
        })
        .catch((e) => {});
    },
    async doDelete(netProxyId) {
      let res = await this.server.node.netProxy.delete({
        netProxyId: netProxyId,
      });
      if (res.code == 0) {
        this.tool.success("删除成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
  },
  created() {},
  updated() {
    this.tool.showNodeNetProxyInfo = this.showNodeNetProxyInfo;
    this.tool.toInsertNodeNetProxy = this.toInsert;
    this.tool.toDeleteNodeNetProxy = this.toDelete;
    this.tool.doDeleteNodeNetProxy = this.doDelete;
    this.tool.toEnableNodeNetProxy = this.toEnable;
    this.tool.toDisableNodeNetProxy = this.toDisable;
  },
  async mounted() {
    this.init();
    this.tool.showNodeNetProxyInfo = this.showNodeNetProxyInfo;
    this.tool.toInsertNodeNetProxy = this.toInsert;
    this.tool.toCopyNodeNetProxy = this.toCopy;
    this.tool.toDeleteNodeNetProxy = this.toDelete;
    this.tool.doDeleteNodeNetProxy = this.doDelete;
    this.tool.toEnableNodeNetProxy = this.toEnable;
    this.tool.toDisableNodeNetProxy = this.toDisable;
  },
};
</script>

<style>
</style>

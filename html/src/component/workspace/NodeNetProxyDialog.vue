<template>
  <el-dialog
    ref="model"
    :title="`网络代理`"
    :close-on-click-model="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showBox"
    :before-close="hide"
    :fullscreen="true"
    width="100%"
    class="node-net-proxy-dialog"
  >
    <div class="node-net-proxy-box">
      <tm-layout height="100%">
        <tm-layout height="60px">
          <div class="node-net-proxy-box-header">
            <el-form class="pdt-10 pdlr-20" size="mini" inline>
              <el-form-item label="" class="mgb-5">
                <div class="tm-btn tm-btn-sm bg-green ft-13" @click="toInsert">
                  新增
                </div>
              </el-form-item>
            </el-form>
          </div>
        </tm-layout>
        <tm-layout height="auto" class="">
          <div
            class="node-net-proxy-body scrollbar toolbox-editor"
            v-if="ready"
          >
            <el-table
              :data="source.nodeNetProxyList"
              :border="true"
              height="100%"
              style="width: 100%"
              size="mini"
            >
              <el-table-column width="70" label="序号" fixed>
                <template slot-scope="scope">
                  <span class="mgl-5">{{ scope.$index + 1 }}</span>
                </template>
              </el-table-column>
              <el-table-column label="名称" width="200">
                <template slot-scope="scope">
                  <div class="">
                    {{ scope.row.model.name }}
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="输入端">
                <template slot-scope="scope">
                  <div class="text-left pd-5">
                    <div class="">
                      <template
                        v-if="
                          source.nodeOptionMap[scope.row.model.innerServerId] !=
                          null
                        "
                      >
                        节点:<span class="pdl-10">
                          {{
                            source.nodeOptionMap[scope.row.model.innerServerId]
                              .text
                          }}
                        </span>
                      </template>
                      <template v-else>
                        节点:<span class="pdl-10">
                          {{ scope.row.model.innerServerId }}
                        </span>
                      </template>
                    </div>
                    <div class="">
                      类型:<span class="pdl-10">{{
                        scope.row.model.innerType
                      }}</span>
                    </div>
                    <div class="">
                      地址:<span class="pdl-10">{{
                        scope.row.model.innerAddress
                      }}</span>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="输出端">
                <template slot-scope="scope">
                  <div class="text-left pd-5">
                    <div class="">
                      <template
                        v-if="
                          source.nodeOptionMap[scope.row.model.outerServerId] !=
                          null
                        "
                      >
                        节点:<span class="pdl-10">
                          {{
                            source.nodeOptionMap[scope.row.model.outerServerId]
                              .text
                          }}
                        </span>
                      </template>
                      <template v-else>
                        节点:<span class="pdl-10">
                          {{ scope.row.model.outerServerId }}
                        </span>
                      </template>
                    </div>
                    <div class="">
                      类型:<span class="pdl-10">{{
                        scope.row.model.outerType
                      }}</span>
                    </div>
                    <div class="">
                      地址:<span class="pdl-10">{{
                        scope.row.model.outerAddress
                      }}</span>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="启动状态" width="100">
                <template slot-scope="scope">
                  <div class="text-left pd-5">
                    <div class="">
                      <template v-if="scope.row.isStarted">
                        <span class="color-green"> 启动中 </span>
                      </template>
                      <template v-else>
                        <span class="color-red"> 已停止 </span>
                      </template>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="状态" width="100">
                <template slot-scope="scope">
                  <div class="text-left pd-5">
                    <div class="">
                      <template v-if="scope.row.model.enabled == 1">
                        <span class="color-green"> 启用 </span>
                      </template>
                      <template v-else>
                        <span class="color-red"> 停用 </span>
                      </template>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="200">
                <template slot-scope="scope">
                  <div
                    class="tm-link color-orange mgl-5"
                    v-if="scope.row.model.enabled == 1"
                    @click="toDisable(scope.row)"
                  >
                    停用
                  </div>
                  <div
                    class="tm-link color-green mgl-5"
                    v-if="scope.row.model.enabled != 1"
                    @click="toEnable(scope.row)"
                  >
                    启用
                  </div>
                  <div
                    class="tm-link color-grey mgl-5"
                    @click="toCopy(scope.row)"
                  >
                    复制
                  </div>
                  <div
                    class="tm-link color-red mgl-5"
                    @click="toDelete(scope.row)"
                  >
                    删除
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </tm-layout>
      </tm-layout>
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
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {
      showBox: false,
      ready: false,
      dataList: null,
    };
  },
  computed: {},
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
      this.ready = false;
      this.dataList = this.source.nodeNetProxyList || [];
      this.initView();
      this.ready = true;
    },
    initView() {},
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
    this.tool.showNodeNetProxyDialog = this.show;
    this.tool.showSwitchNodeNetProxyDialog = this.showSwitch;
    this.tool.hideNodeNetProxyDialog = this.hide;
    this.tool.showNodeNetProxyInfo = this.showNodeNetProxyInfo;
    this.tool.toDeleteNodeNetProxy = this.toDelete;
    this.tool.doDeleteNodeNetProxy = this.doDelete;
    this.tool.toEnableNodeNetProxy = this.toEnable;
    this.tool.toDisableNodeNetProxy = this.toDisable;
  },
  mounted() {
    this.init();
    this.tool.showNodeNetProxyDialog = this.show;
    this.tool.showSwitchNodeNetProxyDialog = this.showSwitch;
    this.tool.hideNodeNetProxyDialog = this.hide;
    this.tool.showNodeNetProxyInfo = this.showNodeNetProxyInfo;
    this.tool.toDeleteNodeNetProxy = this.toDelete;
    this.tool.doDeleteNodeNetProxy = this.doDelete;
    this.tool.toEnableNodeNetProxy = this.toEnable;
    this.tool.toDisableNodeNetProxy = this.toDisable;
  },
};
</script>

<style>
.node-net-proxy-dialog .el-dialog {
  background: #0f1b26;
  color: #ffffff;
  user-select: text;
}
.node-net-proxy-dialog .el-dialog__title {
  color: #ffffff;
}
.node-net-proxy-dialog .el-dialog__body {
  position: relative;
  width: 100%;
  height: calc(100% - 55px);
  padding: 0px;
}
.node-net-proxy-box {
  position: relative;
  width: 100%;
  height: 100%;
}
.node-net-proxy-box-header {
  position: relative;
  width: 100%;
  height: 100%;
  user-select: text;
}
.node-net-proxy-body {
  position: relative;
  width: 100%;
  height: 100%;
}
.node-net-proxy-body {
  position: relative;
  width: 100%;
  height: 100%;
}
</style>

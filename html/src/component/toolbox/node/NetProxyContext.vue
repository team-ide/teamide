<template>
  <div class="node-net-proxy-box">
    <tm-layout height="100%">
      <tm-layout height="60px">
        <div class="node-net-proxy-box-header">
          <el-form class="pdt-10 pdlr-20" size="mini" inline>
            <el-form-item label="" class="mgb-5">
              <div
                class="tm-btn tm-btn-sm bg-green ft-13"
                @click="tool.toInsertNodeNetProxy()"
              >
                新增
              </div>
            </el-form-item>
          </el-form>
        </div>
      </tm-layout>
      <tm-layout height="auto" class="">
        <div class="node-net-proxy-body app-scroll-bar toolbox-editor" v-if="ready">
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
                  {{ scope.row.name }}
                </div>
              </template>
            </el-table-column>
            <el-table-column label="输入端">
              <template slot-scope="scope">
                <div class="text-left pd-5">
                  <div class="">
                    节点:
                    <span
                      v-if="
                        source.nodeOptionMap[scope.row.innerServerId] != null
                      "
                      class="pdl-10"
                      :class="{
                        'color-green':
                          source.nodeOptionMap[scope.row.innerServerId]
                            .isStarted,

                        'color-red':
                          !source.nodeOptionMap[scope.row.innerServerId]
                            .isStarted,
                      }"
                    >
                      {{ source.nodeOptionMap[scope.row.innerServerId].text }}
                    </span>
                    <span v-else class="pdl-10">
                      {{ scope.row.innerServerId }}
                    </span>
                  </div>
                  <div class="">
                    状态:
                    <template v-if="scope.row.innerIsStarted">
                      <span class="pdl-10 color-green"> 启动中 </span>
                    </template>
                    <template v-else>
                      <span class="pdl-10 color-red"> 已停止 </span>
                    </template>
                  </div>
                  <div class="">
                    类型:<span class="pdlr-10 color-blue">
                      {{ scope.row.innerType }}
                    </span>
                    地址:<span class="pdl-10 color-blue">
                      {{ scope.row.innerAddress }}
                    </span>
                  </div>
                  <div class="">
                    输入:<span class="pdlr-10 color-blue">
                      {{ scope.row.innerMonitorData.readSize }}
                      {{ scope.row.innerMonitorData.readSizeUnit }}
                    </span>
                    输入速度:<span class="pdlr-10 color-blue">
                      {{ scope.row.innerMonitorData.readLastSleep }}
                      {{ scope.row.innerMonitorData.readLastSleepUnit }}
                    </span>
                    最后时间:<span class="pdlr-10 color-blue">
                      {{
                        tool.formatDateByTime(
                          scope.row.innerMonitorData.readLastTimestamp
                        )
                      }}
                    </span>
                  </div>
                  <div class="ft-12">
                    输出:<span class="pdlr-10 color-blue">
                      {{ scope.row.innerMonitorData.writeSize }}
                      {{ scope.row.innerMonitorData.writeSizeUnit }}
                    </span>
                    输出速度:<span class="pdlr-10 color-blue">
                      {{ scope.row.innerMonitorData.writeLastSleep }}
                      {{ scope.row.innerMonitorData.writeLastSleepUnit }}
                    </span>
                    最后时间:<span class="pdlr-10 color-blue">
                      {{
                        tool.formatDateByTime(
                          scope.row.innerMonitorData.writeLastTimestamp
                        )
                      }}
                    </span>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="输出端">
              <template slot-scope="scope">
                <div class="text-left pd-5">
                  <div class="">
                    节点:<span
                      v-if="
                        source.nodeOptionMap[scope.row.outerServerId] != null
                      "
                      class="pdl-10"
                      :class="{
                        'color-green':
                          source.nodeOptionMap[scope.row.outerServerId]
                            .isStarted,

                        'color-red':
                          !source.nodeOptionMap[scope.row.outerServerId]
                            .isStarted,
                      }"
                    >
                      {{ source.nodeOptionMap[scope.row.outerServerId].text }}
                    </span>
                    <span v-else class="pdl-10">
                      {{ scope.row.outerServerId }}
                    </span>
                  </div>
                  <div class="">
                    状态:
                    <template v-if="scope.row.outerIsStarted">
                      <span class="pdl-10 color-green"> 启动中 </span>
                    </template>
                    <template v-else>
                      <span class="pdl-10 color-red"> 已停止 </span>
                    </template>
                  </div>
                  <div class="">
                    类型:<span class="pdlr-10 color-blue">
                      {{ scope.row.outerType }}
                    </span>
                    地址:<span class="pdl-10 color-blue">
                      {{ scope.row.outerAddress }}
                    </span>
                  </div>
                  <div class="ft-12">
                    输入:<span class="pdlr-10 color-blue">
                      {{ scope.row.outerMonitorData.readSize }}
                      {{ scope.row.outerMonitorData.readSizeUnit }}
                    </span>
                    输入速度:<span class="pdlr-10 color-blue">
                      {{ scope.row.outerMonitorData.readLastSleep }}
                      {{ scope.row.outerMonitorData.readLastSleepUnit }}
                    </span>
                    最后时间:<span class="pdlr-10 color-blue">
                      {{
                        tool.formatDateByTime(
                          scope.row.outerMonitorData.readLastTimestamp
                        )
                      }}
                    </span>
                  </div>
                  <div class="ft-12">
                    输出:<span class="pdlr-10 color-blue">
                      {{ scope.row.outerMonitorData.writeSize }}
                      {{ scope.row.outerMonitorData.writeSizeUnit }}
                    </span>
                    输出速度:<span class="pdlr-10 color-blue">
                      {{ scope.row.outerMonitorData.writeLastSleep }}
                      {{ scope.row.outerMonitorData.writeLastSleepUnit }}
                    </span>
                    最后时间:<span class="pdlr-10 color-blue">
                      {{
                        tool.formatDateByTime(
                          scope.row.outerMonitorData.writeLastTimestamp
                        )
                      }}
                    </span>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="60">
              <template slot-scope="scope">
                <div class="pd-5">
                  <div class="">
                    <template v-if="scope.row.enabled == 1">
                      <span class="color-green"> 启用 </span>
                    </template>
                    <template v-else>
                      <span class="color-red"> 停用 </span>
                    </template>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150">
              <template slot-scope="scope">
                <div
                  class="tm-link color-orange mgl-5"
                  v-if="scope.row.enabled == 1"
                  @click="tool.toDisableNodeNetProxy(scope.row)"
                >
                  停用
                </div>
                <div
                  class="tm-link color-green mgl-5"
                  v-if="scope.row.enabled != 1"
                  @click="tool.toEnableNodeNetProxy(scope.row)"
                >
                  启用
                </div>
                <div
                  class="tm-link color-grey mgl-5"
                  @click="tool.toCopyNodeNetProxy(scope.row)"
                >
                  复制
                </div>
                <div
                  class="tm-link color-red mgl-5"
                  @click="tool.toDeleteNodeNetProxy(scope.row)"
                >
                  删除
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </tm-layout>
    </tm-layout>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {
      isDestroyed: false,
      ready: false,
      dataList: null,
      loadMonitorDataInit: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      if (!this.loadMonitorDataInit) {
        this.loadMonitorDataInit = true;
        await this.loadMonitorData();
      }
      this.initData();
    },
    async initData() {
      this.ready = false;
      this.dataList = this.source.nodeNetProxyList || [];
      this.initView();
      this.ready = true;
    },
    initView() {},
    async loadMonitorData() {
      if (this.isDestroyed) {
        return;
      }
      let idList = [];
      let list = this.source.nodeNetProxyList || [];
      list.forEach((one) => {
        idList.push(one.code);
      });
      if (idList.length > 0) {
        let res = await this.server.node.netProxy.monitorData({
          idList: idList,
        });
        res.data = res.data || {};
        let netProxyMonitorDataList = res.data.netProxyMonitorDataList || [];
        netProxyMonitorDataList.forEach((netProxyMonitorData) => {
          list.forEach((one) => {
            if (one.code == netProxyMonitorData.id) {
              one.innerMonitorData = netProxyMonitorData.innerMonitorData;
              one.outerMonitorData = netProxyMonitorData.outerMonitorData;
            }
          });
        });
      }
      window.setTimeout(() => {
        this.loadMonitorData();
      }, 1000 * 5);
    },
  },
  created() {},
  updated() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    this.isDestroyed = true;
  },
};
</script>

<style>
.node-net-proxy-box {
  position: relative;
  width: 100%;
  height: 100%;
  user-select: text;
}
.node-net-proxy-box-header {
  position: relative;
  width: 100%;
  height: 100%;
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

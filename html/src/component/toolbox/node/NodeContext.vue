<template>
  <div class="node-context-box">
    <tm-layout height="100%">
      <tm-layout height="150px" class="app-scroll-bar">
        <div class="node-context-box-header">
          <div class="pdlr-20 ft-12 pdtb-10">
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
                ./node -id node1 -address :21090 -token xxx -connAddress ip:port
                -connToken xxx
              </span>
              <span class="color-grey">
                -connAddress 目标节点的ip:port -connToken 目标节点的Token
              </span>
            </div>
            <template v-if="localNodeConnList.length > 0">
              <div class="color-grey">
                当前节点连接到其它节点:
                <span class="color-orange pdlr-10"> 右击节点进行操作 </span>
              </div>
              <div class="color-grey">
                其它节点连接到当前节点:
                <template v-for="(conn, index) in localNodeConnList">
                  <div :key="index" class="color-green">
                    ./node -id node1 -address :21090 -token xxx -connAddress
                    {{ conn.address }} -connToken {{ conn.bindToken }}
                  </div>
                </template>
              </div>
            </template>
          </div>
        </div>
      </tm-layout>
      <tm-layout height="auto" class="">
        <div class="node-context-body" v-if="ready">
          <template v-if="localNodeConnList.length == 0">
            <div class="text-center pdt-50">
              <div
                class="tm-btn bg-green tm-btn-lg"
                @click="tool.toInsertLocalNode()"
              >
                设置本地节点
              </div>
            </div>
          </template>
          <template v-else>
            <NodeView
              :source="source"
              :nodeList="nodeList"
              :onNodeMoved="tool.onNodeMoved"
            ></NodeView>
          </template>
        </div>
      </tm-layout>
    </tm-layout>
  </div>
</template>

<script>
import NodeView from "./NodeView.vue";

export default {
  components: { NodeView },
  props: ["source"],
  data() {
    return {
      isDestroyed: false,
      nodeList: [],
      loading: false,
      ready: false,
      localNodeConnList: [],
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    "source.nodeList"() {
      if (!this.isDestroyed) {
        this.initView();
      }
    },
  },
  methods: {
    init() {
      this.initData();
    },
    async initData() {
      this.ready = false;
      this.initView();
      this.ready = true;
    },
    initView() {
      let nodeLocalList = this.source.nodeLocalList || [];

      let localNodeConnList = [];
      nodeLocalList.forEach((one) => {
        let address = one.bindAddress;
        if (this.tool.isNotEmpty(address) && address.indexOf(":") >= 0) {
          let lastIndex = address.lastIndexOf(":");
          let ip = address.substring(0, lastIndex);
          let port = address.substring(lastIndex + 1);
          if (this.tool.isEmpty(ip) || ip == "0.0.0.0") {
            this.source.localIpList.forEach((localIp) => {
              localNodeConnList.push({
                address: localIp + ":" + port,
                bindToken: one.bindToken,
              });
            });
          } else {
            localNodeConnList.push({
              address: address,
              bindToken: one.bindToken,
            });
          }
        }
      });
      this.localNodeConnList = localNodeConnList;
      this.nodeList = this.source.nodeList;
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
.node-context-box {
  position: relative;
  width: 100%;
  height: 100%;
}
.node-context-box-header {
  position: relative;
  user-select: text;
}
.node-context-body {
  position: relative;
  width: 100%;
  height: 100%;
}
</style>

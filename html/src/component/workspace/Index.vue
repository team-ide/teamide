<template>
  <div
    class="workspace-container"
    :class="{
      'workspace-theme-dark': theme.isDark,
    }"
    :style="{
      backgroundColor: theme.backgroundColor,
    }"
  >
    <div class="workspace-header">
      <div class="workspace-header-nav-box">
        <div class="workspace-header-nav" @click="showSwitchToolboxContext()">
          工具箱
          <span class="color-green mgl-2">({{ source.toolboxCount }})</span>
        </div>
        <div class="workspace-header-nav" @click="showSwitchNodeNetProxyBox()">
          隧道|代理|透传
          <span class="color-green mgl-2">
            (
            {{ source.nodeNetProxyCount }}
            /
            {{ source.nodeNetProxySuccessCount }}
            )
          </span>
        </div>
        <div class="workspace-header-nav tm-disabled">
          监控
          <span class="color-green mgl-2">({{ 0 }}/{{ 0 }})</span>
        </div>
        <div class="workspace-header-nav" @click="showSwitchNodeBox()">
          节点
          <span class="color-green mgl-2">
            (
            {{ source.nodeCount }}
            /
            {{ source.nodeSuccessCount }}
            )
          </span>
        </div>
        <div style="flex: 1"></div>
        <template v-if="source.login.user == null">
          <div
            class="workspace-header-nav"
            v-if="source.hasPower('login')"
            @click="tool.toLogin()"
          >
            登录
          </div>
        </template>
        <template v-else>
          <div class="workspace-header-nav">
            {{ source.login.user.name }}
          </div>
        </template>
        <div class="workspace-header-nav" @click="tool.showUpdateCheck()">
          <template v-if="source.hasNewVersion"> 有新版本 </template>
          <template v-else> 检测新版本 </template>
        </div>
      </div>
    </div>
    <div class="workspace-main">
      <div class="workspace-main-tabs-container">
        <WorkspaceTabs :source="source" :itemsWorker="mainItemsWorker">
          <!-- <el-dropdown
        slot="leftExtend"
        trigger="click"
        class="workspace-tabs-nav-dropdown"
      >
        <div class="workspace-tabs-nav tm-pointer pdlr-5 ft-12">
          全部
          <i class="mdi mdi-menu-down"></i>
        </div>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item> 全部 </el-dropdown-item>

          <template v-for="(one, index) in source.toolboxGroups">
            <el-dropdown-item :key="index" :command="one">
              {{ one.name || one.title }}
            </el-dropdown-item>
          </template>
        </el-dropdown-menu>
      </el-dropdown> -->
          <!-- <div
            slot="rightExtend"
            class="workspace-tabs-nav tm-pointer color-green pdlr-2"
            @click="showSwitchToolboxContext()"
          >
            <i class="mdi mdi-plus"></i>
          </div> -->
        </WorkspaceTabs>
      </div>
      <div class="workspace-main-spans-container">
        <WorkspaceSpans :source="source" :itemsWorker="mainItemsWorker">
          <template v-slot:span="{ item }">
            <template v-if="item.isToolbox">
              <ToolboxEditor
                :source="source"
                :toolboxType="item.toolboxType"
                :toolboxId="item.toolboxId"
                :openId="item.openId"
                :extend="item.extend"
                :item="item"
              >
              </ToolboxEditor>
            </template>
          </template>
        </WorkspaceSpans>
      </div>
      <ToolboxContext
        ref="ToolboxContext"
        :source="source"
        :openByToolboxId="openByToolboxId"
      ></ToolboxContext>
      <NodeBox ref="NodeBox" :source="source"></NodeBox>
      <NodeNetProxyBox ref="NodeNetProxyBox" :source="source"></NodeNetProxyBox>
    </div>
  </div>
</template>

<script>
import ToolboxContext from "./ToolboxContext";
import NodeBox from "./NodeBox";
import NodeNetProxyBox from "./NodeNetProxyBox";

export default {
  components: { ToolboxContext, NodeBox, NodeNetProxyBox },
  props: ["source"],
  data() {
    let mainItemsWorker = this.tool.newItemsWorker({
      onRemoveItem: this.onMainRemoveItem,
      onActiveItem: this.onMainActiveItem,
      toCopyItem: this.toMainCopyItem,
    });
    return {
      mainItemsWorker: mainItemsWorker,
      theme: {
        isDark: true,
        backgroundColor: "#383838",
      },
    };
  },
  computed: {},
  watch: {
    "source.login.user"() {
      this.initData();
    },
  },
  methods: {
    init() {
      this.initData();
    },
    initData() {
      if (this.source.login.user == null) {
        return;
      }
      if (this.initDataed) {
        return;
      }
      this.initDataed = true;
      this.source.initUserToolboxData();
      this.source.initNodeContext();
      this.initOpens();
    },
    showSwitchNodeBox() {
      this.$refs.NodeBox.showSwitch();
    },
    showNodeBox() {
      this.$refs.NodeBox.show();
    },
    hideNodeBox() {
      this.$refs.NodeBox.hide();
    },
    showSwitchNodeNetProxyBoy() {
      this.$refs.NodeNetProxyBox.showSwitch();
    },
    showNodeNetProxyBox() {
      this.$refs.NodeNetProxyBox.show();
    },
    hideNodeNetProxyBox() {
      this.$refs.NodeNetProxyBox.hide();
    },
    showSwitchToolboxContext() {
      this.$refs.ToolboxContext.showSwitch();
    },
    showToolboxContext() {
      this.$refs.ToolboxContext.show();
    },
    hideToolboxContext() {
      this.$refs.ToolboxContext.hide();
    },
    addMainItem(item, fromItem) {
      this.mainItemsWorker.addItem(item, fromItem);
    },
    toMainActiveItem(item) {
      this.mainItemsWorker.toActiveItem(item);
    },
    getMainItems() {
      return this.mainItemsWorker.items || [];
    },
    toMainCopyItem(item) {
      let extend = item.extend;
      this.openByToolboxId(item.toolboxId, extend, item, item.createTime);
    },
    async openByToolboxId(toolboxId, extend, fromItem, createTime) {
      let param = {
        toolboxId: toolboxId,
        extend: JSON.stringify(extend || {}),
      };
      if (createTime) {
        param.createTime = createTime;
      }
      let res = await this.server.toolbox.open(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      } else {
        let openData = res.data.open;
        let item = this.addMainItemByOpen(openData, fromItem);
        if (item != null) {
          this.$nextTick(() => {
            this.toMainActiveItem(item);
          });
        }
      }
    },
    addMainItemByOpen(open, fromItem) {
      let item = {};
      item.key = open.openId;
      item.name = open.toolboxName;
      item.title = open.toolboxName;
      item.show = true;
      item.isToolbox = true;
      item.toolboxType = open.toolboxType;
      item.toolboxId = open.toolboxId;
      item.toolboxGroupId = open.toolboxGroupId;
      item.createTime = open.createTime;
      item.extend = this.tool.getOptionJSON(open.extend);

      if (item.toolboxType == "ssh" || item.extend.isFTP) {
        item.extend.local = item.extend.local || {};
        item.extend.remote = item.extend.remote || {};
        item.extend.local.dir = item.extend.local.dir || "";
        item.extend.remote.dir = item.extend.remote.dir || "";
      }
      item.openId = open.openId;
      switch (item.toolboxType) {
        case "database":
          item.iconFont = "teamide-database";
          break;
        case "redis":
          item.iconFont = "teamide-redis";
          break;
        case "elasticsearch":
          item.iconFont = "teamide-elasticsearch";
          break;
        case "zookeeper":
          item.iconFont = "teamide-zookeeper";
          break;
        case "kafka":
          item.iconFont = "teamide-kafka";
          break;
        case "ssh":
          item.iconFont = "teamide-ssh";
          if (item.extend && item.extend.isFTP) {
            item.iconFont = "teamide-ftp";
          }
          break;
        case "other":
          break;
      }
      this.addMainItem(item, fromItem);
      open.item = item;
      return item;
    },
    async initOpens() {
      let opens = [];
      let res = await this.server.toolbox.queryOpens({});
      if (res.code != 0) {
        this.tool.error(res.msg);
      } else {
        opens = res.data.opens || [];
      }
      opens.forEach((one) => {
        this.addMainItemByOpen(one);
      });
      // 激活最后
      let activeOpen = null;
      opens.forEach(async (one) => {
        if (activeOpen == null) {
          activeOpen = one;
        } else {
          if (
            new Date(one.openTime).getTime() >
            new Date(activeOpen.openTime).getTime()
          ) {
            activeOpen = one;
          }
        }
      });
      if (activeOpen != null) {
        this.toMainActiveItem(activeOpen.item);
      } else {
        this.showToolboxContext();
      }
    },
    async onMainActiveItem(item) {
      if (item == null || this.tool.isEmpty(item.openId)) {
        return;
      }
      let res = await this.server.toolbox.open({
        openId: item.openId,
      });
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      this.hideToolboxContext();
    },
    async onMainRemoveItem(item) {
      if (item == null || this.tool.isEmpty(item.openId)) {
        return;
      }
      let res = await this.server.toolbox.close({
        openId: item.openId,
      });
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      if (this.getMainItems().length == 0) {
        this.showToolboxContext();
      }
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.workspace-container {
  width: 100%;
  height: 100%;
  margin: 0px;
  padding: 0px;
  position: relative;
}
.workspace-container.workspace-theme-dark {
  color: #d9d9d9;
}
.workspace-header {
  width: 100%;
  height: 30px;
  margin: 0px;
  padding: 0px;
  position: relative;
  border-bottom: 1px solid #4e4e4e;
}
.workspace-main {
  width: 100%;
  height: calc(100% - 30px);
  margin: 0px;
  padding: 0px;
  position: relative;
}
.workspace-main-tabs-container {
  width: 100%;
  height: 30px;
  position: relative;
  border-bottom: 1px solid #4e4e4e;
}
.workspace-main-spans-container {
  width: 100%;
  height: calc(100% - 30px);
  position: relative;
}
.workspace-header-nav-box {
  height: 100%;
  align-items: center;
  display: flex;
  white-space: nowrap;
  padding: 0px 0px;
}
.workspace-header-nav {
  height: 100%;
  font-size: 12px;
  padding: 0px 10px;
  cursor: pointer;
  white-space: nowrap;
  align-items: center;
  display: flex;
}
.workspace-header-nav:hover {
  background-color: #505050;
}
.default-tabs-container {
  width: 100%;
  height: 30px;
  position: relative;
  border-bottom: 1px solid #4e4e4e;
}
.default-spans-container {
  width: 100%;
  height: calc(100% - 30px);
  position: relative;
}
</style>

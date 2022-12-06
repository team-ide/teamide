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
        <div
          v-if="source.hasPower('toolbox')"
          class="workspace-header-nav"
          @click="tool.showSwitchToolboxContext()"
        >
          工具箱
          <span class="color-green mgl-2">({{ source.toolboxCount }})</span>
        </div>
        <div
          v-if="source.hasPower('node')"
          class="workspace-header-nav"
          @click="openNodeContext()"
        >
          节点
        </div>
        <div
          v-if="source.hasPower('node_net_proxy')"
          class="workspace-header-nav"
          @click="openNodeNetProxyContext()"
        >
          网络代理|透传
        </div>
        <div v-if="source.hasPower('terminal')" class="workspace-header-nav">
          <el-dropdown
            trigger="click"
            class="terminal-dropdown"
            ref="terminalDropdown"
          >
            <span class="el-dropdown-link">
              终端<i class="el-icon-arrow-down el-icon--right"></i>
            </span>
            <el-dropdown-menu slot="dropdown" class="terminal-dropdown-menu">
              <MenuBox>
                <MenuItem @click="openTerminal('local')">本地</MenuItem>
                <template
                  v-if="
                    source.sshToolboxList && source.sshToolboxList.length > 0
                  "
                >
                  <MenuItem>
                    SSH
                    <MenuSubBox
                      slot="MenuSubBox"
                      style="max-height: 500px; overflow-y: scroll"
                    >
                      <template v-for="(one, index) in source.sshToolboxList">
                        <MenuItem
                          :key="index"
                          @click="openTerminal('ssh', one)"
                        >
                          {{ one.name }}
                        </MenuItem>
                      </template>
                    </MenuSubBox>
                  </MenuItem>
                </template>
                <template v-if="source.nodeList && source.nodeList.length > 0">
                  <MenuItem>
                    节点
                    <MenuSubBox slot="MenuSubBox">
                      <template v-for="(one, index) in source.nodeList">
                        <MenuItem
                          :key="index"
                          @click="openTerminal('node', one)"
                        >
                          {{ one.name }}
                        </MenuItem>
                      </template>
                    </MenuSubBox>
                  </MenuItem>
                </template>
              </MenuBox>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
        <div
          v-if="source.hasPower('file_manager')"
          class="workspace-header-nav"
        >
          <el-dropdown
            trigger="click"
            class="file-manager-dropdown"
            ref="fileManagerDropdown"
          >
            <span class="el-dropdown-link">
              文件管理器<i class="el-icon-arrow-down el-icon--right"></i>
            </span>
            <el-dropdown-menu
              slot="dropdown"
              class="file-manager-dropdown-menu"
            >
              <MenuBox>
                <MenuItem @click="openFileManager('local')">本地</MenuItem>
                <template
                  v-if="
                    source.sshToolboxList && source.sshToolboxList.length > 0
                  "
                >
                  <MenuItem>
                    SSH
                    <MenuSubBox
                      slot="MenuSubBox"
                      style="max-height: 500px; overflow-y: scroll"
                    >
                      <template v-for="(one, index) in source.sshToolboxList">
                        <MenuItem
                          :key="index"
                          @click="openFileManager('ssh', one)"
                        >
                          {{ one.name }}
                        </MenuItem>
                      </template>
                    </MenuSubBox>
                  </MenuItem>
                </template>
                <template v-if="source.nodeList && source.nodeList.length > 0">
                  <MenuItem>
                    节点
                    <MenuSubBox slot="MenuSubBox">
                      <template v-for="(one, index) in source.nodeList">
                        <MenuItem
                          :key="index"
                          @click="openFileManager('node', one)"
                        >
                          {{ one.name }}
                        </MenuItem>
                      </template>
                    </MenuSubBox>
                  </MenuItem>
                </template>
              </MenuBox>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
        <div style="flex: 1"></div>
        <template v-if="source.login.user == null">
          <div
            class="workspace-header-nav"
            v-if="source.hasPower('login')"
            @click="tool.toLogin()"
          >
            登 录
          </div>
          <div
            class="workspace-header-nav"
            v-if="source.hasPower('register')"
            @click="tool.toRegister()"
          >
            注 册
          </div>
        </template>
        <template v-else>
          <div class="workspace-header-nav">
            <el-dropdown
              trigger="click"
              class="user-dropdown"
              ref="userDropdown"
            >
              <span class="el-dropdown-link">
                {{ source.login.user.name }}
                <i class="el-icon-arrow-down el-icon--right"></i>
              </span>
              <el-dropdown-menu slot="dropdown" class="user-dropdown-menu">
                <MenuBox>
                  <MenuItem @click="openPage('userCenter', '个人中心')">
                    个人中心
                  </MenuItem>
                  <MenuItem
                    v-if="source.hasPower('logout')"
                    @click="tool.toLogout()"
                  >
                    退出登录
                  </MenuItem>
                </MenuBox>
              </el-dropdown-menu>
            </el-dropdown>
          </div>
        </template>
        <div
          v-if="source.hasPower('update_check')"
          class="workspace-header-nav"
          @click="tool.showUpdateCheck()"
        >
          <template v-if="source.hasNewVersion"> 有新版本 </template>
          <template v-else> 检测新版本 </template>
        </div>
      </div>
    </div>
    <div class="workspace-main">
      <div class="workspace-main-tabs-container">
        <WorkspaceTabs :source="source" :itemsWorker="mainItemsWorker">
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
    </div>
    <ToolboxContext
      :source="source"
      :openByToolboxId="openByToolboxId"
    ></ToolboxContext>
  </div>
</template>

<script>
import ToolboxContext from "./ToolboxContext";

export default {
  components: { ToolboxContext },
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
    "source.login.userId"() {
      this.initUserData();
    },
  },
  methods: {
    async initUserData() {
      await this.source.initUserToolboxData();
      await this.initOpens();
    },
    init() {
      this.initUserData();
      this.server.addServerSocketOnOpen(() => {});
    },
    openNodeContext() {
      this.tool.openByExtend({
        toolboxType: "node",
        type: "node-context",
        title: "节点",
        onlyOpenOneKey: "node:node-context",
      });
    },
    openNodeNetProxyContext() {
      this.tool.openByExtend({
        toolboxType: "node",
        type: "net-proxy-context",
        title: "网络透传",
        onlyOpenOneKey: "node:net-proxy-context",
      });
    },
    openPage(page, title) {
      this.tool.openByExtend({
        toolboxType: "page",
        page: page,
        title: title,
        onlyOpenOneKey: "page:" + page,
      });
    },
    openTerminal(place, placeData) {
      let extend = {
        toolboxType: "terminal",
        place: place,
        title: null,
        placeId: null,
      };
      if (place == "local") {
        extend.title = "终端-本地";
      } else if (place == "ssh") {
        extend.title = "" + placeData.name;
        extend.placeId = "" + placeData.toolboxId;
      } else if (place == "node") {
        extend.title = "" + placeData.name;
        extend.placeId = "" + placeData.serverId;
      } else {
        this.tool.error("暂不支持该配置作为终端");
        return;
      }
      this.tool.openByExtend(extend);
      this.$refs.terminalDropdown && this.$refs.terminalDropdown.hide();
    },
    openFileManager(place, placeData) {
      let extend = {
        toolboxType: "file-manager",
        place: place,
        title: null,
        placeId: null,
      };
      if (place == "local") {
        extend.title = "文件管理器-本地";
      } else if (place == "ssh") {
        extend.title = "" + placeData.name;
        extend.placeId = "" + placeData.toolboxId;
      } else if (place == "node") {
        extend.title = "" + placeData.name;
        extend.placeId = "" + placeData.serverId;
      } else {
        this.tool.error("暂不支持该配置作为文件管理器");
        return;
      }
      this.tool.openByExtend(extend);
      this.$refs.fileManagerDropdown && this.$refs.fileManagerDropdown.hide();
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
      if (this.tool.isEmpty(item.toolboxId)) {
        this.openByExtend(extend, item, item.createTime);
      } else {
        this.openByToolboxId(item.toolboxId, extend, item, item.createTime);
      }
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
    async openByExtend(extend, fromItem, createTime) {
      if (
        extend == null ||
        Object.keys(extend) == 0 ||
        this.tool.isEmpty(extend.toolboxType)
      ) {
        this.tool.error("根据扩展打开需要配置类型");
        return;
      }
      if (this.tool.isEmpty(extend.title)) {
        this.tool.error("根据扩展打开需要配置标题");
        return;
      }
      if (this.tool.isNotEmpty(extend.onlyOpenOneKey)) {
        let items = this.getMainItems() || [];
        let find = null;
        items.forEach((item) => {
          if (
            item &&
            item.extend &&
            item.extend.onlyOpenOneKey == extend.onlyOpenOneKey
          ) {
            find = item;
          }
        });
        if (find != null) {
          this.toMainActiveItem(find);
          return;
        }
      }
      let param = {
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

      if (this.tool.isEmpty(item.title)) {
        item.name = item.extend.name;
        item.title = item.extend.title;
      }
      if (this.tool.isEmpty(item.toolboxType)) {
        item.toolboxType = item.extend.toolboxType;
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
          item.icon = "mdi-console";
          if (item.extend && item.extend.isFTP) {
            item.icon = "mdi-folder";
          }
          break;
        case "other":
          break;
        default:
          if (item.extend) {
            if (item.extend.toolboxType == "file-manager") {
              item.icon = "mdi-folder";
            } else if (item.extend.toolboxType == "terminal") {
              item.icon = "mdi-console";
            }
          }
      }
      this.addMainItem(item, fromItem);
      open.item = item;
      return item;
    },
    async initOpens() {
      let opens = [];
      if (this.source.login.user != null) {
        let res = await this.server.toolbox.queryOpens({});
        if (res.code != 0) {
          this.tool.error(res.msg);
        } else {
          opens = res.data.opens || [];
        }
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
        if (this.source.login.user != null) {
          this.tool.showToolboxContext();
        }
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
      this.tool.hideToolboxContext();
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
        if (this.source.login.user != null) {
          this.tool.showToolboxContext();
        }
      }
    },
  },
  created() {},
  updated() {
    this.tool.openByExtend = this.openByExtend;
    this.tool.openPage = this.openPage;
    this.tool.openTerminal = this.openTerminal;
    this.tool.openFileManager = this.openFileManager;
  },
  mounted() {
    this.tool.openByExtend = this.openByExtend;
    this.tool.openPage = this.openPage;
    this.tool.openTerminal = this.openTerminal;
    this.tool.openFileManager = this.openFileManager;

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

.file-manager-dropdown.el-dropdown {
  color: unset;
  font-size: unset;
  display: flex;
  white-space: nowrap;
  align-items: center;
}
.file-manager-dropdown-menu.el-dropdown-menu {
  padding: 0;
  margin: 0;
  border: 0;
  border-radius: 4px;
  box-shadow: 0 0 0;
  background: transparent;
  top: 35px !important;
}
.file-manager-dropdown-menu .menu-box a {
  cursor: pointer;
}

.terminal-dropdown.el-dropdown {
  color: unset;
  font-size: unset;
  display: flex;
  white-space: nowrap;
  align-items: center;
}
.terminal-dropdown-menu.el-dropdown-menu {
  padding: 0;
  margin: 0;
  border: 0;
  border-radius: 4px;
  box-shadow: 0 0 0;
  background: transparent;
  top: 35px !important;
}
.terminal-dropdown-menu .menu-box a {
  cursor: pointer;
}

.user-dropdown.el-dropdown {
  color: unset;
  font-size: unset;
  display: flex;
  white-space: nowrap;
  align-items: center;
}
.user-dropdown-menu.el-dropdown-menu {
  padding: 0;
  margin: 0;
  border: 0;
  border-radius: 4px;
  box-shadow: 0 0 0;
  background: transparent;
  top: 35px !important;
}
.user-dropdown-menu .menu-box a {
  cursor: pointer;
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

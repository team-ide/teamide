<template>
  <div class="toolbox-main">
    <TabEditor
      ref="TabEditor"
      :source="source"
      :onRemoveTab="onRemoveTab"
      :onActiveTab="onActiveTab"
      :slotTab="true"
      :copyTab="toCopyTab"
      :hasOpenNewWindow="hasOpenNewWindow"
      :openNewWindow="openNewWindow"
    >
      <template v-slot:tab="{ tab }">
        <span class="toolbox-tab-title">
          <template v-if="tab.toolboxType.name == 'database'">
            <IconFont class="teamide-database"> </IconFont>
          </template>
          <template v-else-if="tab.toolboxType.name == 'redis'">
            <IconFont class="teamide-redis"> </IconFont>
          </template>
          <template v-else-if="tab.toolboxType.name == 'elasticsearch'">
            <IconFont class="teamide-elasticsearch"> </IconFont>
          </template>
          <template v-else-if="tab.toolboxType.name == 'kafka'">
            <IconFont class="teamide-kafka"> </IconFont>
          </template>
          <template v-else-if="tab.toolboxType.name == 'zookeeper'">
            <IconFont class="teamide-zookeeper"> </IconFont>
          </template>
          <template
            v-else-if="tab.toolboxType.name == 'ssh' && tab.extend.isFTP"
          >
            <IconFont class="teamide-ftp"> </IconFont>
          </template>
          <template v-else-if="tab.toolboxType.name == 'ssh'">
            <IconFont class="teamide-ssh"> </IconFont>
          </template>
          <span>{{ tab.name }}</span>
          <template v-if="tool.isNotEmpty(tab.comment)">
            <span>>{{ tab.comment }}</span>
          </template>
        </span>
      </template>
      <template v-slot:body="{ tab }">
        <ToolboxEditor
          :source="source"
          :toolbox="toolbox"
          :tab="tab"
          :toolboxType="tab.toolboxType"
          :openId="tab.openId"
          :toolboxData="tab.toolboxData"
          :extend="tab.extend"
          :active="tab.active"
          :updateOpenExtend="updateOpenExtend"
          :updateOpenComment="updateOpenComment"
        >
        </ToolboxEditor>
      </template>
      <div slot="extend" class="tab-header-extend">
        <div
          class="tab-header-nav tm-pointer"
          @click="toolbox.showSwitchToolboxType()"
        >
          <i class="mdi mdi-plus"></i>
        </div>
      </div>
    </TabEditor>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "context"],
  data() {
    return {
      hasOpenNewWindow: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    "$route.path"() {
      if (!this.tool.isToolboxPage(this.$route.path)) {
        return;
      }
      this.$refs.TabEditor.onFocus();
    },
  },
  methods: {
    init() {
      if (
        window.electron &&
        window.electron.ipcRenderer &&
        window.electron.ipcRenderer.sendMessage
      ) {
        this.hasOpenNewWindow = true;
      }
      this.onIpcRendererOnce();

      this.initOpens();
    },
    onIpcRendererOnce() {
      if (
        window.electron &&
        window.electron.ipcRenderer &&
        window.electron.ipcRenderer.once
      ) {
        window.electron.ipcRenderer.once("close-open-window", (config) => {
          this.onIpcRendererOnce();
          if (config == null) {
            return;
          }
          let key = this.getTabKeyByData(config);
          if (key == null) {
            return;
          }
          let tab = this.getTab(key);
          if (tab == null) {
            return;
          }
          this.updateOpenExtend(tab.openId, { openNewWindow: false });
          this.showTab(tab);
          this.doActiveTab(tab.openId);
        });
      }
    },
    openNewWindow(tab) {
      let url = this.source.url + "#/open/toolbox/" + tab.openId;
      if (this.hasOpenNewWindow) {
        this.updateOpenExtend(tab.openId, { openNewWindow: true });
        this.hideTab(tab);
        window.electron.ipcRenderer.sendMessage("open-new-window", {
          url: url,
          openId: tab.openId,
          type: "toolbox",
        });
      } else {
        window.open(url);
      }
    },
    getTabs() {
      return this.$refs.TabEditor && this.$refs.TabEditor.getTabs();
    },
    getTab(tab) {
      return this.$refs.TabEditor && this.$refs.TabEditor.getTab(tab);
    },
    showTab(tab) {
      return this.$refs.TabEditor && this.$refs.TabEditor.showTab(tab);
    },
    hideTab(tab) {
      return this.$refs.TabEditor && this.$refs.TabEditor.hideTab(tab);
    },
    onRemoveTab(tab) {
      this.source.toolbox.closeOpen(tab.openId);
      if (this.getTabs().length == 0) {
        this.source.toolbox.showToolboxType();
      }
    },
    onActiveTab(tab) {
      this.source.toolbox.activeOpen(tab.openId);
      this.source.toolbox.hideToolboxType();
    },
    addTab(tab, fromTab) {
      return this.$refs.TabEditor.addTab(tab, fromTab);
    },
    doActiveTab(tab) {
      return this.$refs.TabEditor.doActiveTab(tab);
    },
    getTabKeyByData(openData) {
      let key = "" + openData.openId;
      return key;
    },
    getTabByData(openData) {
      let key = this.getTabKeyByData(openData);
      let tab = this.getTab(key);
      return tab;
    },
    updateOpenComment(openId, comment) {
      let tab = this.getTab("" + openId);
      if (tab == null) {
        return;
      }
      tab.comment = comment;
    },
    async updateOpenExtend(openId, keyValueMap) {
      let tab = this.getTab("" + openId);
      if (tab == null) {
        return;
      }
      if (keyValueMap == null) {
        return;
      }
      if (Object.keys(keyValueMap) == 0) {
        return;
      }
      let obj = tab.extend;
      for (let key in keyValueMap) {
        let value = keyValueMap[key];
        let names = key.split(".");
        names.forEach((name, index) => {
          if (index < names.length - 1) {
            obj[name] = obj[name] || {};
            obj = obj[name];
          } else {
            obj[name] = value;
          }
        });
      }
      let param = {
        openId: openId,
        extend: JSON.stringify(tab.extend),
      };
      let res = await this.server.toolbox.updateOpenExtend(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
    },
    toolboxDataOpen(toolboxData, fromTab) {
      this.tool.stopEvent();
      this.source.toolbox.hideToolboxType();

      let extend = {};

      if (toolboxData && toolboxData.toolboxType) {
        let toolboxType = this.source.toolbox.getToolboxType(
          toolboxData.toolboxType
        );
        if (toolboxType && toolboxType.name == "other") {
          extend = this.source.toolbox.getOptionJSON(toolboxData.option);
        }
      }

      this.openByToolboxData(toolboxData, extend, fromTab);
    },
    toCopyTab(tab) {
      let extend = tab.extend;
      this.openByToolboxData(
        tab.toolboxData,
        extend,
        tab,
        tab.openData.createTime
      );
    },
    toolboxDataOpenSfpt(toolboxData) {
      this.tool.stopEvent();
      this.source.toolbox.hideToolboxType();

      this.openByToolboxData(toolboxData, { isFTP: true });
    },
    async openByToolboxData(toolboxData, extend, fromTab, createTime) {
      let openData = await this.source.toolbox.open(
        toolboxData.toolboxId,
        extend,
        createTime
      );
      if (openData == null) {
        return;
      }
      let tab = await this.openByOpenData(openData, fromTab);
      if (tab != null) {
        this.doActiveTab(tab);
      }
    },
    async openByOpenData(openData, fromTab) {
      let toolboxData = this.source.toolbox.getToolboxData(openData.toolboxId);
      if (toolboxData == null) {
        await this.source.toolbox.closeOpen(openData.openId);
        return;
      }
      let toolboxType = this.source.toolbox.getToolboxType(
        toolboxData.toolboxType
      );
      if (toolboxType == null) {
        await this.source.toolbox.closeOpen(openData.openId);
        return;
      }
      openData.toolboxData = toolboxData;
      openData.toolboxType = toolboxType;
      if (this.tool.isNotEmpty(openData.extend)) {
        openData.extend = JSON.parse(openData.extend);
      } else {
        openData.extend = {};
      }
      let key = this.getTabKeyByData(openData);

      let tab = this.getTab(key);
      if (tab == null) {
        tab = this.source.toolbox.createToolboxDataTab(openData);
        tab.key = key;
      }

      this.addTab(tab, fromTab);
      return tab;
    },
    async initOpens() {
      let opens = await this.source.toolbox.loadOpens();

      await opens.forEach(async (openData) => {
        await this.openByOpenData(openData);
      });

      // 激活最后
      let activeOpenData = null;
      opens.forEach(async (openData) => {
        if (
          this.hasOpenNewWindow &&
          openData.extend &&
          openData.extend.openNewWindow
        ) {
          let key = this.getTabKeyByData(openData);
          let tab = this.getTab(key);
          if (tab != null) {
            this.openNewWindow(tab);
            return;
          }
        }
        if (activeOpenData == null) {
          activeOpenData = openData;
        } else {
          if (
            new Date(openData.openTime).getTime() >
            new Date(activeOpenData.openTime).getTime()
          ) {
            activeOpenData = openData;
          }
        }
      });
      if (activeOpenData != null) {
        this.doActiveTab(activeOpenData.openId);
      } else {
        this.source.toolbox.showToolboxType();
      }
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.init();
    this.source.toolbox.toolboxDataOpenSfpt = this.toolboxDataOpenSfpt;
    this.source.toolbox.toolboxDataOpen = this.toolboxDataOpen;
    this.source.toolbox.getTabByData = this.getTabByData;
  },
};
</script>

<style>
.toolbox-main {
  width: 100%;
  height: 100%;
  position: relative;
}
.toolbox-tab-title {
  font-size: 13px;
  line-height: 24px;
}
.toolbox-tab-title .icon {
  margin-right: 5px;
}
</style>

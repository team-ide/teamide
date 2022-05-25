<template>
  <div class="toolbox-main">
    <TabEditor
      ref="TabEditor"
      :source="source"
      :onRemoveTab="onRemoveTab"
      :onActiveTab="onActiveTab"
      :slotTab="true"
      :copyTab="toCopyTab"
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
    return {};
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
    getTabs() {
      return this.$refs.TabEditor && this.$refs.TabEditor.getTabs();
    },
    getTab(tab) {
      return this.$refs.TabEditor && this.$refs.TabEditor.getTab(tab);
    },
    onRemoveTab(tab) {
      this.toolbox.closeOpen(tab.openId);
      if (this.getTabs().length == 0) {
        this.toolbox.showToolboxType();
      }
    },
    onActiveTab(tab) {
      this.toolbox.activeOpen(tab.openId);
      this.toolbox.hideToolboxType();
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
    createTabByOpenData(openData) {
      let key = this.getTabKeyByData(openData);

      let tab = this.getTab(key);
      if (tab == null) {
        tab = this.toolbox.createToolboxDataTab(openData);
        tab.key = key;
      }
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
      this.toolbox.hideToolboxType();

      let extend = {};
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
      this.toolbox.hideToolboxType();

      this.openByToolboxData(toolboxData, { isFTP: true });
    },
    async openByToolboxData(toolboxData, extend, fromTab, createTime) {
      let openData = await this.toolbox.open(
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
      let toolboxData = this.getToolboxData(openData.toolboxId);
      if (toolboxData == null) {
        await this.toolbox.closeOpen(openData.openId);
        return;
      }
      let toolboxType = this.getToolboxType(toolboxData.toolboxType);
      if (toolboxType == null) {
        await this.toolbox.closeOpen(openData.openId);
      }
      openData.toolboxData = toolboxData;
      openData.toolboxType = toolboxType;
      if (this.tool.isNotEmpty(openData.extend)) {
        openData.extend = JSON.parse(openData.extend);
      } else {
        openData.extend = {};
      }
      let tab = this.createTabByOpenData(openData);
      this.addTab(tab, fromTab);
      return tab;
    },
    getToolboxType(type) {
      let res = null;
      this.toolbox.types.forEach((one) => {
        if (one == type || one.name == type || one.name == type.name) {
          res = one;
        }
      });
      return res;
    },
    getToolboxData(toolboxData) {
      let res = null;
      for (let type in this.context) {
        if (this.context[type] == null) {
          continue;
        }
        this.context[type].forEach((one) => {
          if (
            one == toolboxData ||
            one.toolboxId == toolboxData ||
            one.toolboxId == toolboxData.toolboxId
          ) {
            res = one;
          }
        });
      }
      return res;
    },
    async initOpens() {
      let opens = await this.toolbox.loadOpens();

      await opens.forEach(async (openData) => {
        await this.openByOpenData(openData);
      });

      // 激活最后
      let activeOpenData = null;
      opens.forEach(async (openData) => {
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
        this.toolbox.showToolboxType();
      }
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.initOpens();
    this.toolbox.toolboxDataOpenSfpt = this.toolboxDataOpenSfpt;
    this.toolbox.toolboxDataOpen = this.toolboxDataOpen;
    this.toolbox.getTabByData = this.getTabByData;
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

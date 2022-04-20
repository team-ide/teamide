<template>
  <div class="toolbox-main">
    <TabEditor
      ref="TabEditor"
      :source="source"
      :onRemoveTab="onRemoveTab"
      :onActiveTab="onActiveTab"
      :onOffsetRightDistance="onOffsetRightDistance"
      :slotTab="true"
      :copyTab="toCopyTab"
    >
      <template v-slot:tab="{ tab }">
        <span>
          <template v-if="tab.toolboxType.name == 'database'">
            <i class="mgr-3">
              <IconDatabase></IconDatabase>
            </i>
          </template>
          <template v-else-if="tab.toolboxType.name == 'redis'">
            <i class="mgr-3">
              <IconRedis></IconRedis>
            </i>
          </template>
          <template v-else-if="tab.toolboxType.name == 'elasticsearch'">
            <i class="mgr-3">
              <IconElasticsearch></IconElasticsearch>
            </i>
          </template>
          <template v-else-if="tab.toolboxType.name == 'kafka'">
            <i class="mgr-3">
              <IconKafka></IconKafka>
            </i>
          </template>
          <template v-else-if="tab.toolboxType.name == 'zookeeper'">
            <i class="mgr-3">
              <IconZookeeper></IconZookeeper>
            </i>
          </template>
          <template
            v-else-if="tab.toolboxType.name == 'ssh' && tab.extend.isFTP"
          >
            <i class="mgr-3">
              <IconFtp></IconFtp>
            </i>
          </template>
          <template v-else-if="tab.toolboxType.name == 'ssh'">
            <i class="mgr-3">
              <IconSsh></IconSsh>
            </i>
          </template>
          {{ tab.name }}
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
        >
        </ToolboxEditor>
      </template>
      <div slot="extend" class="tab-header-extend">
        <el-dropdown
          trigger="click"
          size="mini"
          :placement="dropdownPlacement"
          :visible-change="dropdownVisible"
          ref="dropdown"
        >
          <div class="tab-header-nav tm-pointer">
            <i class="mdi mdi-plus"></i>
          </div>
          <el-dropdown-menu slot="dropdown" class="pd-0">
            <MenuBox class="bd-0" :subLeft="dropdownMenuSubLeft" size="mini">
              <template v-for="toolboxType in source.toolboxTypes">
                <MenuItem :key="toolboxType.name">
                  <template v-if="toolboxType.name == 'database'">
                    <i class="mgr-3">
                      <IconDatabase></IconDatabase>
                    </i>
                  </template>
                  <template v-else-if="toolboxType.name == 'redis'">
                    <i class="mgr-3">
                      <IconRedis></IconRedis>
                    </i>
                  </template>
                  <template v-else-if="toolboxType.name == 'elasticsearch'">
                    <i class="mgr-3">
                      <IconElasticsearch></IconElasticsearch>
                    </i>
                  </template>
                  <template v-else-if="toolboxType.name == 'kafka'">
                    <i class="mgr-3">
                      <IconKafka></IconKafka>
                    </i>
                  </template>
                  <template v-else-if="toolboxType.name == 'zookeeper'">
                    <i class="mgr-3">
                      <IconZookeeper></IconZookeeper>
                    </i>
                  </template>
                  <template v-else-if="toolboxType.name == 'ssh'">
                    <i class="mgr-3">
                      <IconSsh></IconSsh>
                    </i>
                    <i class="mgr-3">
                      <IconFtp></IconFtp>
                    </i>
                  </template>
                  {{ toolboxType.text || toolboxType.name }}
                  <span
                    class="tm-link color-green mgl-10"
                    title="新增"
                    @click="toInsert(toolboxType)"
                  >
                    <i class="mdi mdi-plus ft-14"></i>
                  </span>
                  <MenuSubBox slot="MenuSubBox">
                    <MenuItem
                      class="tm-pointer"
                      @click="toInsert(toolboxType)"
                      title="新增"
                    >
                      <span class="tm-link color-green">
                        <i class="mdi mdi-plus ft-16"></i>
                      </span>
                    </MenuItem>
                    <template v-if="context[toolboxType.name] != null">
                      <template
                        v-for="toolboxData in context[toolboxType.name]"
                      >
                        <MenuItem :key="toolboxData.toolboxId"
                          ><span
                            class="tm-link color-green mgr-10"
                            title="打开"
                            @click="toolboxDataOpen(toolboxData)"
                          >
                            {{ toolboxData.name }}
                          </span>
                          <span
                            title="打开FTP"
                            v-if="toolboxType.name == 'ssh'"
                            class="tm-link color-orange mgr-10"
                            @click="toolboxDataOpenSfpt(toolboxData)"
                          >
                            <i class="mdi mdi-folder-outline ft-13"></i>
                          </span>
                          <span
                            title="编辑"
                            class="tm-link color-blue mgr-10"
                            @click="toUpdate(toolboxType, toolboxData)"
                          >
                            <i class="mdi mdi-folder-edit-outline ft-13"></i>
                          </span>
                          <span
                            title="复制"
                            class="tm-link color-green mgr-10"
                            @click="toCopy(toolboxType, toolboxData)"
                          >
                            <i class="mdi mdi-content-copy ft-12"></i>
                          </span>
                          <span
                            title="删除"
                            class="tm-link color-orange mgr-10"
                            @click="toDelete(toolboxType, toolboxData)"
                          >
                            <i class="mdi mdi-delete-outline ft-14"></i>
                          </span>
                        </MenuItem>
                      </template>
                    </template>
                  </MenuSubBox>
                </MenuItem>
              </template>
            </MenuBox>
          </el-dropdown-menu>
        </el-dropdown>
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
      dropdownPlacement: "bottom-start",
      dropdownVisible: false,
      dropdownMenuSubLeft: false,
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
    onOffsetRightDistance(offsetRightDistance) {
      offsetRightDistance = parseInt(offsetRightDistance);
      if (offsetRightDistance < 500) {
        this.dropdownPlacement = "bottom-end";
        this.dropdownMenuSubLeft = true;
      } else {
        this.dropdownPlacement = "bottom-start";
        this.dropdownMenuSubLeft = false;
      }
    },
    getTab(tab) {
      return this.$refs.TabEditor.getTab(tab);
    },
    onRemoveTab(tab) {
      this.toolbox.closeOpen(tab.openId);
    },
    onActiveTab(tab) {
      this.toolbox.activeOpen(tab.openId);
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
    async updateOpenExtend(openId, keys, value) {
      let tab = this.getTab("" + openId);
      if (tab == null) {
        return;
      }
      let obj = tab.extend;
      keys.forEach((key, index) => {
        if (index < keys.length - 1) {
          obj[key] = obj[key] || {};
          obj = obj[key];
        } else {
          obj[key] = value;
        }
      });
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
      this.$refs.dropdown.hide();

      let extend = {};
      this.openByToolboxData(toolboxData, extend, fromTab);
    },
    toCopyTab(tab) {
      let extend = tab.extend;
      this.openByToolboxData(tab.toolboxData, extend, tab);
    },
    toolboxDataOpenSfpt(toolboxData) {
      this.tool.stopEvent();
      this.$refs.dropdown.hide();

      this.openByToolboxData(toolboxData, { isFTP: true });
    },
    async openByToolboxData(toolboxData, extend, fromTab) {
      let openData = await this.toolbox.open(toolboxData.toolboxId, extend);
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
      this.source.toolboxTypes.forEach((one) => {
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
    toInsert(toolboxType) {
      this.tool.stopEvent();
      this.$refs.dropdown.hide();
      let toolboxData = {};
      this.toolbox.showToolboxForm(toolboxType, toolboxData, (g, m) => {
        let flag = this.doInsert(g, m);
        return flag;
      });
    },
    toCopy(toolboxType, copy) {
      this.tool.stopEvent();
      this.$refs.dropdown.hide();
      let toolboxData = {};
      Object.assign(toolboxData, copy);
      delete toolboxData.toolboxId;
      toolboxData.name = toolboxData.name + " Copy";
      this.toolbox.showToolboxForm(toolboxType, toolboxData, (g, m) => {
        let flag = this.doInsert(g, m);
        return flag;
      });
    },
    toUpdate(toolboxType, toolboxData) {
      this.tool.stopEvent();
      this.$refs.dropdown.hide();
      this.updateData = toolboxData;
      this.toolbox.showToolboxForm(toolboxType, toolboxData, (g, m) => {
        let flag = this.doUpdate(g, m);
        return flag;
      });
    },
    toDelete(toolboxType, toolboxData) {
      this.tool.stopEvent();
      this.$refs.dropdown.hide();
      this.tool
        .confirm(
          "删除[" +
            toolboxType.text +
            "]工具[" +
            toolboxData.name +
            "]将无法回复，确定删除？"
        )
        .then(async () => {
          return this.doDelete(toolboxType, toolboxData);
        })
        .catch((e) => {});
    },
    async doDelete(toolboxType, toolboxData) {
      let res = await this.server.toolbox.delete(toolboxData);
      if (res.code == 0) {
        this.toolbox.initContext();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doUpdate(toolboxType, toolboxData) {
      toolboxData.toolboxType = toolboxType.name;
      toolboxData.toolboxId = this.updateData.toolboxId;
      let res = await this.server.toolbox.update(toolboxData);
      if (res.code == 0) {
        this.toolbox.initContext();
        let tab = this.getTabByData(toolboxData);
        if (tab != null) {
          Object.assign(tab.toolboxData, toolboxData);
        }
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doInsert(toolboxType, toolboxData) {
      toolboxData.toolboxType = toolboxType.name;
      let res = await this.server.toolbox.insert(toolboxData);
      if (res.code == 0) {
        this.toolbox.initContext();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
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
      }
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.initOpens();
  },
};
</script>

<style>
.toolbox-main {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>

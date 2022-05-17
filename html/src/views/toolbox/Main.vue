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
        <el-dropdown
          trigger="click"
          size="mini"
          :placement="dropdownPlacement"
          @visible-change="dropdownVisible"
          ref="dropdown"
        >
          <div class="tab-header-nav tm-pointer">
            <i class="mdi mdi-plus"></i>
          </div>
          <el-dropdown-menu slot="dropdown" class="toolbox-dropdown-box-menu">
            <div class="toolbox-dropdown-box">
              <WaterfallLayout
                class="toolbox-type-box"
                :columns="3"
                :gap="10"
                :startLeft="10"
                :startTop="10"
                ref="WaterfallLayout"
              >
                <template v-for="toolboxType in toolbox.types">
                  <div
                    :key="toolboxType.name"
                    class="toolbox-type-one waterfall-layout-item"
                  >
                    <div class="toolbox-type-title">
                      <template v-if="toolboxType.name == 'database'">
                        <IconFont class="teamide-database"> </IconFont>
                      </template>
                      <template v-else-if="toolboxType.name == 'redis'">
                        <IconFont class="teamide-redis"> </IconFont>
                      </template>
                      <template v-else-if="toolboxType.name == 'elasticsearch'">
                        <IconFont class="teamide-elasticsearch"> </IconFont>
                      </template>
                      <template v-else-if="toolboxType.name == 'kafka'">
                        <IconFont class="teamide-kafka"> </IconFont>
                      </template>
                      <template v-else-if="toolboxType.name == 'zookeeper'">
                        <IconFont class="teamide-zookeeper"> </IconFont>
                      </template>
                      <template v-else-if="toolboxType.name == 'ssh'">
                        <IconFont class="teamide-ssh"> </IconFont>
                        <IconFont class="teamide-ftp"> </IconFont>
                      </template>
                      <span class="toolbox-type-text">
                        {{ toolboxType.text || toolboxType.name }}
                      </span>
                      <span
                        class="tm-link color-green mgl-10"
                        title="新增"
                        @click="toInsert(toolboxType)"
                      >
                        <i class="mdi mdi-plus ft-14"></i>
                      </span>
                    </div>
                    <div class="toolbox-type-data-box">
                      <template
                        v-if="
                          context[toolboxType.name] == null ||
                          context[toolboxType.name].length == 0
                        "
                      >
                        <span
                          class="tm-link color-green"
                          title="新增"
                          @click="toInsert(toolboxType)"
                        >
                          新增
                        </span>
                      </template>
                      <template v-else>
                        <template
                          v-for="toolboxData in context[toolboxType.name]"
                        >
                          <div
                            :key="toolboxData.toolboxId"
                            class="toolbox-type-data"
                          >
                            <span
                              class="
                                toolbox-type-data-text
                                tm-link
                                color-grey
                                mgr-10
                              "
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
                              class="tm-link color-grey mgr-10"
                              @click="toCopy(toolboxType, toolboxData)"
                            >
                              <i class="mdi mdi-content-copy ft-12"></i>
                            </span>
                            <span
                              title="删除"
                              class="tm-link color-red mgr-10"
                              @click="toDelete(toolboxType, toolboxData)"
                            >
                              <i class="mdi mdi-delete-outline ft-14"></i>
                            </span>
                          </div>
                        </template>
                      </template>
                    </div>
                  </div>
                </template>
              </WaterfallLayout>
            </div>
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </TabEditor>
    <FormDialog
      ref="InsertToolbox"
      :source="source"
      title="新增Toolbox"
      :onSave="doInsert"
    ></FormDialog>
    <FormDialog
      ref="UpdateToolbox"
      :source="source"
      title="编辑Toolbox"
      :onSave="doUpdate"
    ></FormDialog>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "context"],
  data() {
    return {
      dropdownPlacement: "bottom-start",
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
    dropdownVisible(dropdownVisible) {
      if (dropdownVisible) {
        this.$nextTick(() => {
          this.$refs.WaterfallLayout && this.$refs.WaterfallLayout.doLayout();
        });
      }
    },
    onOffsetRightDistance(offsetRightDistance) {
      // offsetRightDistance = parseInt(offsetRightDistance);
      // if (offsetRightDistance < 500) {
      //   this.dropdownPlacement = "bottom-end";
      //   this.dropdownMenuSubLeft = true;
      // } else {
      //   this.dropdownPlacement = "bottom-start";
      //   this.dropdownMenuSubLeft = false;
      // }
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
      this.$refs.dropdown.hide();

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
      this.$refs.dropdown.hide();

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
    toInsert(toolboxType) {
      this.tool.stopEvent();
      this.$refs.dropdown.hide();
      let toolboxData = {};
      let optionsJSON = {};

      this.$refs.InsertToolbox.show({
        title: `新增[${toolboxType.text}]工具`,
        form: [this.form.toolbox, toolboxType.configForm],
        data: [toolboxData, optionsJSON],
        toolboxType,
      });
    },
    onInsertSuccess() {},
    toCopy(toolboxType, copy) {
      this.tool.stopEvent();
      this.$refs.dropdown.hide();
      let toolboxData = {};
      Object.assign(toolboxData, copy);
      delete toolboxData.toolboxId;
      toolboxData.name = toolboxData.name + " Copy";

      let optionsJSON = this.getOptionJSON(toolboxData.option);

      this.$refs.InsertToolbox.show({
        title: `新增[${toolboxType.text}]工具`,
        form: [this.form.toolbox, toolboxType.configForm],
        data: [toolboxData, optionsJSON],
        toolboxType,
      });
    },
    toUpdate(toolboxType, toolboxData) {
      this.tool.stopEvent();
      this.$refs.dropdown.hide();
      this.updateData = toolboxData;

      let optionsJSON = this.getOptionJSON(toolboxData.option);

      this.$refs.UpdateToolbox.show({
        title: `编辑[${toolboxType.text}][${toolboxData.name}]工具`,
        form: [this.form.toolbox, toolboxType.configForm],
        data: [toolboxData, optionsJSON],
        toolboxType,
      });
    },
    getOptionJSON(option) {
      let json = {};
      if (this.tool.isNotEmpty(option)) {
        json = JSON.parse(option);
      }
      return json;
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
        this.tool.success("删除成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doUpdate(dataList, config) {
      let toolboxData = dataList[0];
      let optionJSON = dataList[1];
      let toolboxType = config.toolboxType;
      toolboxData.toolboxType = toolboxType.name;
      toolboxData.toolboxId = this.updateData.toolboxId;
      toolboxData.option = JSON.stringify(optionJSON);
      let res = await this.server.toolbox.update(toolboxData);
      if (res.code == 0) {
        this.toolbox.initContext();
        let tab = this.getTabByData(toolboxData);
        if (tab != null) {
          Object.assign(tab.toolboxData, toolboxData);
        }
        this.tool.success("修改成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doInsert(dataList, config) {
      let toolboxData = dataList[0];
      let optionJSON = dataList[1];
      let toolboxType = config.toolboxType;
      toolboxData.toolboxType = toolboxType.name;
      toolboxData.option = JSON.stringify(optionJSON);
      let res = await this.server.toolbox.insert(toolboxData);
      if (res.code == 0) {
        this.toolbox.initContext();
        this.tool.success("新增成功");
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
.toolbox-tab-title {
  font-size: 13px;
  line-height: 24px;
}
.toolbox-tab-title .icon {
  margin-right: 5px;
}
.el-dropdown-menu.toolbox-dropdown-box-menu {
  padding: 0px;
  margin: 0px;
  border: 0px;
  box-shadow: none;
  background: #2a3036;
  padding: 0px 0px 10px;
}
.el-dropdown-menu.toolbox-dropdown-box-menu.el-popper[x-placement^="bottom"]
  .popper__arrow {
  border-bottom-color: #2a3036;
}
.el-dropdown-menu.toolbox-dropdown-box-menu.el-popper[x-placement^="bottom"]
  .popper__arrow::after {
  border-bottom-color: #2a3036;
}
.toolbox-dropdown-box {
  width: 940px;
  /* height: 600px; */
  /* background: #2a4a67; */
}

.toolbox-type-box {
  width: 100%;
  font-size: 12px;
  padding: 10px;
}

.toolbox-type-one {
  width: 300px;
}
.toolbox-type-title {
  padding: 0px 10px;
  background: #2b3f51;
  color: #ffffff;
  line-height: 23px;
}
.toolbox-type-title .icon {
  margin-right: 5px;
}
.toolbox-type-title .tm-link {
  padding: 0px;
}
.toolbox-type-data-box {
  background: #1b2a38;
  padding: 5px 10px;
}
.toolbox-type-data {
  display: flex;
  overflow: hidden;
  padding: 2px 0px;
}
.toolbox-type-data-text {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  text-align: left;
}
.toolbox-type-data .tm-link {
  padding: 0px;
}
</style>

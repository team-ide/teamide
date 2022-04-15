<template>
  <div class="toolbox-main">
    <TabEditor
      ref="TabEditor"
      :source="source"
      :onRemoveTab="onRemoveTab"
      :onActiveTab="onActiveTab"
      :onOffsetRightDistance="onOffsetRightDistance"
    >
      <template v-slot:body="{ tab }">
        <ToolboxEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="tab.toolboxType"
          :data="tab.data"
          :extend="tab.extend"
          :active="tab.active"
        ></ToolboxEditor>
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
          <el-dropdown-menu slot="dropdown">
            <MenuBox class="bd-0" :subLeft="dropdownMenuSubLeft">
              <template v-for="toolboxType in source.toolboxTypes">
                <MenuItem :key="toolboxType.name">
                  {{ toolboxType.text || toolboxType.name }}
                  <MenuSubBox slot="MenuSubBox">
                    <MenuItem @click="toInsert(toolboxType)">
                      <span class="tm-pointer color-green" title="新增">
                        <i class="mdi mdi-plus ft-16"></i>
                      </span>
                    </MenuItem>
                    <template v-if="context[toolboxType.name] != null">
                      <template v-for="data in context[toolboxType.name]">
                        <MenuItem
                          :key="data.toolboxId"
                          @click="dataOpen(toolboxType, data)"
                        >
                          {{ data.name }}

                          <span
                            title="打开FTP"
                            v-if="toolboxType.name == 'ssh'"
                            class="tm-pointer color-orange mgl-5"
                            @click="dataOpenSfpt(toolboxType, data)"
                          >
                            <i class="mdi mdi-folder-outline ft-16"></i>
                          </span>
                          <span
                            title="编辑"
                            class="tm-pointer color-blue mgl-5"
                            @click="toUpdate(toolboxType, data)"
                          >
                            <i class="mdi mdi-folder-edit-outline ft-16"></i>
                          </span>
                          <span
                            title="复制"
                            class="tm-pointer color-green mgl-5"
                            @click="toCopy(toolboxType, data)"
                          >
                            <i class="mdi mdi-content-copy ft-15"></i>
                          </span>
                          <span
                            title="删除"
                            class="tm-pointer color-orange mgl-5"
                            @click="toDelete(toolboxType, data)"
                          >
                            <i class="mdi mdi-delete-outline ft-15"></i>
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
      dropdownPlacement: "right",
      dropdownVisible: false,
      dropdownMenuSubLeft: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    onOffsetRightDistance(offsetRightDistance) {
      offsetRightDistance = parseInt(offsetRightDistance);
      if (offsetRightDistance < 500) {
        this.dropdownPlacement = "left";
        this.dropdownMenuSubLeft = true;
      } else {
        this.dropdownPlacement = "right";
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
    addTab(tab) {
      return this.$refs.TabEditor.addTab(tab);
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
    createTabByData(openData) {
      let key = this.getTabKeyByData(openData);

      let tab = this.getTab(key);
      if (tab == null) {
        let toolboxType = openData.toolboxType;
        let data = openData.data;
        let extend = openData.extend;
        let title = toolboxType.text + " : " + data.name;
        let name = data.name;

        extend = extend || {};
        this.toolbox.formatExtend(toolboxType, data, extend);
        if (extend.isFTP) {
          title = "FTP : " + data.name;
          name = "FTP : " + data.name;
        } else if (toolboxType.name == "ssh") {
          title = "SSH : " + data.name;
          name = "SSH : " + data.name;
        }
        tab = {
          key,
          data,
          title,
          name,
          toolboxType,
          extend,
          openId: openData.openId,
        };
        tab.active = false;
      }
      return tab;
    },

    dataOpen(toolboxType, data) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
        window.event.preventDefault && window.event.preventDefault();
      }
      this.$refs.dropdown.hide();
      this.open(data);
    },
    dataOpenSfpt(toolboxType, data) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
        window.event.preventDefault && window.event.preventDefault();
      }
      this.$refs.dropdown.hide();
      this.open(data, {
        isFTP: true,
      });
    },
    async open(data, extend) {
      let extendStr = null;
      if (extend != null) {
        extendStr = JSON.stringify(extend);
      }
      let param = {
        toolboxId: data.toolboxId,
        extend: extendStr,
      };
      let res = await this.server.toolbox.open(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      let openData = res.data.open;
      let tab = await this.openByOpenData(openData);
      if (tab != null) {
        this.toolbox.doActiveTab(tab);
      }
    },
    async openByOpenData(openData) {
      let data = this.getToolboxData(openData.toolboxId);
      if (data == null) {
        await this.closeOpen(openData.openId);
        return;
      }
      let toolboxType = this.getToolboxType(data.toolboxType);
      if (toolboxType == null) {
        await this.closeOpen(openData.openId);
      }
      openData.data = data;
      openData.toolboxType = toolboxType;
      if (this.tool.isNotEmpty(openData.extend)) {
        openData.extend = JSON.parse(openData.extend);
      } else {
        openData.extend = null;
      }
      let tab = this.toolbox.createTabByData(openData);
      this.toolbox.addTab(tab);
      return tab;
    },
    async activeOpen(openId) {
      let param = {
        openId: openId,
      };
      let res = await this.server.toolbox.open(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
    },
    async closeOpen(openId) {
      let param = {
        openId: openId,
      };
      let res = await this.server.toolbox.close(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
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
    getToolboxData(data) {
      let res = null;
      for (let type in this.context) {
        if (this.context[type] == null) {
          continue;
        }
        this.context[type].forEach((one) => {
          if (
            one == data ||
            one.toolboxId == data ||
            one.toolboxId == data.toolboxId
          ) {
            res = one;
          }
        });
      }
      return res;
    },
    toInsert(toolboxType) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
        window.event.preventDefault && window.event.preventDefault();
      }
      this.$refs.dropdown.hide();
      let data = {};
      this.toolbox.showToolboxForm(toolboxType, data, (g, m) => {
        let flag = this.doInsert(g, m);
        return flag;
      });
    },
    toCopy(toolboxType, copy) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
        window.event.preventDefault && window.event.preventDefault();
      }
      this.$refs.dropdown.hide();
      let data = {};
      Object.assign(data, copy);
      delete data.toolboxId;
      data.name = data.name + " Copy";
      this.toolbox.showToolboxForm(toolboxType, data, (g, m) => {
        let flag = this.doInsert(g, m);
        return flag;
      });
    },
    toUpdate(toolboxType, data) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
        window.event.preventDefault && window.event.preventDefault();
      }
      this.$refs.dropdown.hide();
      this.updateData = data;
      this.toolbox.showToolboxForm(toolboxType, data, (g, m) => {
        let flag = this.doUpdate(g, m);
        return flag;
      });
    },
    toDelete(toolboxType, data) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
        window.event.preventDefault && window.event.preventDefault();
      }
      this.$refs.dropdown.hide();
      this.tool
        .confirm(
          "删除[" +
            toolboxType.text +
            "]工具[" +
            data.name +
            "]将无法回复，确定删除？"
        )
        .then(async () => {
          return this.doDelete(toolboxType, data);
        })
        .catch((e) => {});
    },
    async doDelete(toolboxType, data) {
      let res = await this.server.toolbox.delete(data);
      if (res.code == 0) {
        this.toolbox.initContext();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doUpdate(toolboxType, data) {
      data.toolboxType = toolboxType.name;
      data.toolboxId = this.updateData.toolboxId;
      let res = await this.server.toolbox.update(data);
      if (res.code == 0) {
        this.toolbox.initContext();
        let tab = this.toolbox.getTabByData(data);
        if (tab != null) {
          Object.assign(tab.data, data);
        }
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doInsert(toolboxType, data) {
      data.toolboxType = toolboxType.name;
      let res = await this.server.toolbox.insert(data);
      if (res.code == 0) {
        this.toolbox.initContext();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async initOpens() {
      let opens = await this.loadOpens();

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
        this.toolbox.doActiveTab(activeOpenData.openId);
      }
    },
    async loadOpens() {
      let param = {};
      let res = await this.server.toolbox.queryOpens(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      let opens = res.data.opens || [];
      return opens;
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.initOpens();
    this.toolbox.closeOpen = this.closeOpen;
    this.toolbox.activeOpen = this.activeOpen;
    this.toolbox.saveToolbox = this.saveToolbox;

    this.toolbox.doActiveTab = this.doActiveTab;
    this.toolbox.createTab = this.createTab;
    this.toolbox.createTabByData = this.createTabByData;
    this.toolbox.getTabByData = this.getTabByData;
    this.toolbox.addTab = this.addTab;
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

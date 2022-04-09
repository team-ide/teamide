<template>
  <div class="toolbox-main">
    <div class="toolbox-main-header">
      <div class="toolbox-tab-box">
        <template v-for="(one, index) in toolbox.tabs">
          <div
            :key="'tab-' + index"
            class="toolbox-tab"
            :title="one.title"
            :class="{ active: one.active }"
          >
            <span class="text" @click="toSelectTab(one)">{{ one.name }}</span>
            <span
              class="delete-btn tm-pointer color-orange"
              @click="toDeleteTab(one)"
              title="关闭"
            >
              <b-icon icon="x"></b-icon>
            </span>
          </div>
        </template>
      </div>
    </div>
    <div class="toolbox-main-body">
      <div class="toolbox-tab-span-box">
        <template v-for="one in toolbox.tabs">
          <div
            :key="one.key"
            class="toolbox-tab-span"
            :class="{ active: one.active }"
          >
            <ToolboxEditor
              :source="source"
              :toolbox="source.toolbox"
              :toolboxType="one.toolboxType"
              :data="one.data"
              :extend="one.extend"
              :active="one.active"
            ></ToolboxEditor>
          </div>
        </template>
      </div>
    </div>
    <div class="toolbox-main-footer">
      <template v-if="toolbox.activeTab != null">
        <span class="ft-12 pdlr-10">{{ toolbox.activeTab.title }}</span>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "app"],
  data() {
    return {};
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    toSelectTab(tab) {
      this.doActiveTab(tab);
    },
    getTab(tab) {
      let res = null;
      this.toolbox.tabs.forEach((one) => {
        if (one == tab || one.key == tab || one.key == tab.key) {
          res = one;
        }
      });
      return res;
    },
    toDeleteTab(tab) {
      this.removeTab(tab);
    },
    addTab(tab) {
      let find = this.getTab(tab);
      if (find != null) {
        return;
      }
      this.toolbox.tabs.push(tab);
    },
    removeTab(tab) {
      let find = this.getTab(tab);
      if (find == null) {
        return;
      }
      let tabIndex = this.toolbox.tabs.indexOf(find);
      this.toolbox.tabs.splice(tabIndex, 1);
      if (find.active) {
        let nextTabIndex = tabIndex - 1;
        if (nextTabIndex < 0) {
          nextTabIndex = 0;
        }
        this.doActiveTab(this.toolbox.tabs[nextTabIndex]);
      }
      this.toolbox.closeOpen(find.openId);
    },
    doActiveTab(tab) {
      this.$nextTick(() => {
        tab = this.getTab(tab);
        this.toolbox.tabs.forEach((one) => {
          if (one != tab) {
            one.active = false;
          }
        });
        this.toolbox.tabs.forEach((one) => {
          if (one == tab) {
            one.active = true;
          }
        });
        this.toolbox.activeTab = tab;
        if (tab != null) {
          this.toolbox.activeOpen(tab.openId);
        }
      });
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
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
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
.toolbox-main-header {
  width: 100%;
  height: 25px;
  line-height: 25px;
  font-size: 14px;
  position: relative;
}
.toolbox-main-body {
  width: 100%;
  height: calc(100% - 25px - 20px);
  border-bottom: 1px solid #4e4e4e;
  position: relative;
}
.toolbox-main-footer {
  width: 100%;
  height: 20px;
  line-height: 20px;
  position: relative;
  display: flex;
}
.toolbox-tab-box {
  display: flex;
  position: relative;
  background-color: #383838;
}
.toolbox-tab {
  display: flex;
  border-right: 1px solid #4e4e4e;
  position: relative;
  border-top-left-radius: 0px;
  border-top-right-radius: 10px;
}
.toolbox-tab.active {
  background-color: #2d2d2d;
}
.toolbox-tab .text {
  padding: 0px 5px 0px 10px;
  cursor: pointer;
}
.toolbox-tab .delete-btn {
  transition: all 0.1s;
  transform: scale(0);
  margin: 0px 5px;
}
.toolbox-tab:hover .delete-btn {
  transform: scale(1);
}
.toolbox-tab-span-box {
  width: 100%;
  height: 100%;
  position: relative;
}
.toolbox-tab-span {
  width: 100%;
  height: 100%;
  position: absolute;
  left: 0px;
  right: 0px;
  transition: all 0s;
  transform: scale(0);
  background-color: #2d2d2d;
}
.toolbox-tab-span.active {
  transform: scale(1);
}
</style>

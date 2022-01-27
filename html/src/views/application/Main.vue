<template>
  <div class="application-main">
    <div class="application-main-header">
      <div class="application-main-tab-box">
        <template v-for="(one, index) in application.tabs">
          <div
            :key="'tab-' + index"
            class="application-main-tab"
            :title="one.title"
            :class="{ active: one.active }"
          >
            <span class="text" @click="toSelectTab(one)">{{ one.text }}</span>
            <span
              v-if="one.changed"
              class="color-orange mgr-2 ft-22"
              title="有更改"
            >
              <b-icon icon="dot"></b-icon>
            </span>
            <span
              v-if="!one.changed"
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
    <div class="application-main-body">
      <div class="application-main-span-box">
        <template v-for="(one, index) in application.tabs">
          <div
            :key="'span-' + index"
            class="application-main-span"
            :class="{ active: one.active }"
          >
            <ModelEditor
              :source="source"
              :group="one.group"
              :model="one.model"
              :context="application.context"
              @change="onModelChange"
              @save="onModelSave"
            ></ModelEditor>
          </div>
        </template>
      </div>
    </div>
    <div class="application-main-footer">
      <template v-if="application.activeTab != null">
        <span class="ft-12 pdlr-10">{{ application.activeTab.title }}</span>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "application", "app", "context"],
  data() {
    return {};
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    async onModelChange(group, model) {
      let key = this.getTabKeyByModel(group, model);
      let tab = this.getTab(key);
      if (tab == null) {
        return;
      }
      tab.changed = true;
      tab.last_model = model;
    },
    async onModelSave(group, model) {
      let key = this.getTabKeyByModel(group, model);
      let tab = this.getTab(key);
      if (tab == null) {
        return;
      }
      let flag = await this.application.saveModel(group, model);
      if (flag) {
        tab.changed = false;
        tab.last_model = null;
        tab.model = model;
      }
      return flag;
    },
    toSelectTab(tab) {
      this.doActiveTab(tab);
    },
    getTab(tab) {
      let res = null;
      this.application.tabs.forEach((one) => {
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
      this.application.tabs.push(tab);
    },
    removeTab(tab) {
      let find = this.getTab(tab);
      if (find == null) {
        return;
      }
      let tabIndex = this.application.tabs.indexOf(find);
      this.application.tabs.splice(tabIndex, 1);
      if (find.active) {
        let nextTabIndex = tabIndex - 1;
        if (nextTabIndex < 0) {
          nextTabIndex = 0;
        }
        this.doActiveTab(this.application.tabs[nextTabIndex]);
      }
    },
    doActiveTab(tab) {
      this.$nextTick(() => {
        tab = this.getTab(tab);
        this.application.tabs.forEach((one) => {
          if (one != tab) {
            one.active = false;
          }
        });
        this.application.tabs.forEach((one) => {
          if (one == tab) {
            one.active = true;
          }
        });
        this.application.activeTab = tab;
      });
    },
    getTabKeyByModel(group, model) {
      let key = this.app.name + ":" + group.name + ":" + model.name;
      return key;
    },
    createTabByModel(group, model) {
      let key = this.getTabKeyByModel(group, model);

      let tab = this.getTab(key);
      if (tab == null) {
        let title =
          "应用：" + this.app.name + " > " + group.text + " > " + model.name;
        if (this.tool.isNotEmpty(model.comment)) {
          title += "(" + model.comment + ")";
        }
        let text = model.name;
        tab = {
          key,
          text,
          title,
          model,
          group,
        };
        tab.active = false;
        tab.changed = false;
      }
      return tab;
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.application.doActiveTab = this.doActiveTab;
    this.application.createTab = this.createTab;
    this.application.createTabByModel = this.createTabByModel;
    this.application.addTab = this.addTab;
  },
};
</script>

<style>
.application-main {
  width: 100%;
  height: 100%;
  position: relative;
}
.application-main-header {
  width: 100%;
  height: 25px;
  line-height: 25px;
  font-size: 14px;
  position: relative;
}
.application-main-body {
  width: 100%;
  height: calc(100% - 25px - 20px);
  border-bottom: 1px solid #4e4e4e;
  position: relative;
}
.application-main-footer {
  width: 100%;
  height: 20px;
  line-height: 20px;
  position: relative;
  display: flex;
}
.application-main-tab-box {
  display: flex;
  position: relative;
  background-color: #383838;
}
.application-main-tab {
  display: flex;
  border-right: 1px solid #4e4e4e;
  position: relative;
  border-top-left-radius: 0px;
  border-top-right-radius: 10px;
}
.application-main-tab.active {
  background-color: #2d2d2d;
}
.application-main-tab .text {
  padding: 0px 5px 0px 10px;
  cursor: pointer;
}
.application-main-tab .delete-btn {
  transition: all 0.1s;
  transform: scale(0);
  margin: 0px 5px;
}
.application-main-tab:hover .delete-btn {
  transform: scale(1);
}
.application-main-span-box {
  width: 100%;
  height: 100%;
  position: relative;
}
.application-main-span {
  width: 100%;
  height: 100%;
  position: absolute;
  left: 0px;
  right: 0px;
  transition: all 0s;
  transform: scale(0);
  background-color: #2d2d2d;
}
.application-main-span.active {
  transform: scale(1);
}
</style>

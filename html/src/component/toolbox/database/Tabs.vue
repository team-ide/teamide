<template>
  <div class="toolbox-database-tabs">
    <template v-if="ready">
      <div class="toolbox-database-tabs-header">
        <div class="toolbox-tab-box">
          <template v-for="(one, index) in tabs">
            <div
              :key="'tab-' + index"
              class="toolbox-tab"
              :title="one.title"
              :class="{ active: one.active }"
            >
              <span class="text" @click="toSelectTab(one)">
                {{ one.name }}
              </span>
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
      <div class="toolbox-database-tabs-body">
        <div class="toolbox-tab-span-box">
          <template v-for="(one, index) in tabs">
            <div
              :key="'span-' + index"
              class="toolbox-tab-span"
              :class="{ active: one.active }"
            >
              <template v-if="one.type == 'data'">
                <ToolboxDatabaseTableData
                  :source="source"
                  :toolbox="toolbox"
                  :toolboxType="toolboxType"
                  :data="data"
                  :wrap="wrap"
                  :database="one.data.database"
                  :table="one.data"
                >
                </ToolboxDatabaseTableData>
              </template>
              <template v-else-if="one.type == 'sql'">
                <ToolboxDatabaseSql
                  :source="source"
                  :toolbox="toolbox"
                  :toolboxType="toolboxType"
                  :data="data"
                  :wrap="wrap"
                  :sqlData="one.data"
                >
                </ToolboxDatabaseSql>
              </template>
            </div>
          </template>
        </div>
      </div>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "data", "toolboxType", "toolbox", "option", "wrap"],
  data() {
    return {
      ready: false,
      tabs: [],
      activeTab: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.wrap.doActiveTab = this.doActiveTab;
      this.wrap.createTab = this.createTab;
      this.wrap.createTabByData = this.createTabByData;
      this.wrap.getTabByData = this.getTabByData;
      this.wrap.addTab = this.addTab;
      this.ready = true;
    },
    toSelectTab(tab) {
      this.doActiveTab(tab);
    },
    getTab(tab) {
      let res = null;
      this.tabs.forEach((one) => {
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
      this.tabs.push(tab);
    },
    removeTab(tab) {
      let find = this.getTab(tab);
      if (find == null) {
        return;
      }
      let tabIndex = this.tabs.indexOf(find);
      this.tabs.splice(tabIndex, 1);
      if (find.active) {
        let nextTabIndex = tabIndex - 1;
        if (nextTabIndex < 0) {
          nextTabIndex = 0;
        }
        this.doActiveTab(this.tabs[nextTabIndex]);
      }
    },
    doActiveTab(tab) {
      this.$nextTick(() => {
        tab = this.getTab(tab);
        this.tabs.forEach((one) => {
          if (one != tab) {
            one.active = false;
          }
        });
        this.tabs.forEach((one) => {
          if (one == tab) {
            one.active = true;
          }
        });
        this.activeTab = tab;
      });
    },
    getTabKeyByData(type, data) {
      let key;
      if (type == "data") {
        key = "" + data.database.name + ":" + data.name;
      } else {
        key = data.key;
      }

      return key;
    },
    getTabByData(type, data) {
      let key = this.getTabKeyByData(type, data);
      let tab = this.getTab(key);
      return tab;
    },
    createTabByData(type, data) {
      let key = this.getTabKeyByData(type, data);

      let tab = this.getTab(key);
      if (tab == null) {
        let title = "";
        let name = "";
        if (type == "data") {
          title = data.database.name + "." + data.name;
          name = data.database.name + "." + data.name;
        } else {
          title = data.title || data.name;
          name = data.name || data.title;
        }
        tab = {
          key,
          data,
          type,
          title,
          name,
        };
        tab.active = false;
        tab.changed = false;
      }
      return tab;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-tabs {
  width: 100%;
  height: 100%;
}
.toolbox-database-tabs-header {
  width: 100%;
  height: 25px;
  line-height: 25px;
  font-size: 14px;
  position: relative;
}
.toolbox-database-tabs-body {
  width: 100%;
  height: calc(100% - 25px);
  border-bottom: 1px solid #4e4e4e;
  position: relative;
}
</style>

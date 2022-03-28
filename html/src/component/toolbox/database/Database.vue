<template>
  <div class="toolbox-database-database">
    <template v-if="ready">
      <div class="pd-10">
        <el-tree
          ref="tree"
          :load="loadNode"
          lazy
          :props="defaultProps"
          :default-expanded-keys="expands"
          node-key="key"
          :expand-on-click-node="false"
          @node-click="nodeClick"
        >
          <span class="toolbox-editor-tree-span" slot-scope="{ node, data }">
            <span>{{ node.label }}</span>
            <div class="toolbox-editor-tree-btn-group">
              <a
                v-if="data.isDatabase || data.isDatabaseTables"
                class="tm-link color-grey ft-14 mgr-4"
                @click="toReloadChildren(data)"
              >
                <i class="mdi mdi-reload"></i>
              </a>
              <a
                v-if="data.isDatabase || data.isTable"
                class="tm-link color-grey ft-12 mgr-4"
                title="表数据"
                @click="toShowCreate(data)"
              >
                <i class="">DDL</i>
              </a>
              <a
                v-if="data.isDatabase || data.isTable"
                class="tm-link color-orange ft-15 mgr-4"
                @click="toDelete(data)"
              >
                <i class="mdi mdi-delete-outline"></i>
              </a>
            </div>
          </span>
        </el-tree>
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
      expands: [],
      opens: [],
      defaultProps: {
        children: "children",
        label: "name",
        isLeaf: "leaf",
      },
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.ready = true;
    },
    toReloadChildren(data) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
      }
      this.reloadChildren(data);
    },
    reloadChildren(key) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
      }
      let node = this.$refs.tree.getNode(key);
      if (node) {
        node.loaded = false;
        node.expand();
      }
    },
    async toShowCreate(data) {
      if (data.isDatabase) {
        this.wrap.showDatabaseCreate(data);
      } else if (data.isTable) {
        this.wrap.showTableCreate(data);
      }
    },
    nodeClick(data, node, nodeView) {
      let nowTime = new Date().getTime();
      let clickTime = node.clickTime;
      node.clickTime = nowTime;
      if (clickTime) {
        let timeout = nowTime - clickTime;
        if (timeout < 300) {
          node.clickTime = null;
          this.nodeDbClick(data, node, nodeView);
        }
      }
    },
    nodeDbClick(data, node, nodeView) {
      if (data.isDatabase) {
        if (node.expanded) {
          node.collapse();
        } else {
          node.expand();
        }
      } else if (data.isDatabaseTables) {
        if (node.expanded) {
          node.collapse();
        } else {
          node.expand();
        }
      } else if (data.isTable) {
        this.toOpenTable(data);
      }
    },
    toOpenTable(data) {
      let tab = this.wrap.createTabByData("data", data);
      this.wrap.addTab(tab);
      this.wrap.doActiveTab(tab);
    },
    toDelete(data) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
      }
      let msg = "确认删除";
      if (data.isDatabase) {
        msg += "库[" + data.name + "]";
      } else if (data.isTable) {
        msg += "表[" + data.name + "]";
      }
      msg += "?";
      this.tool
        .confirm(msg)
        .then(async () => {
          if (data.isDatabase) {
            this.doDeleteDatabase(data.name);
          } else if (data.isTable) {
            this.doDeleteTable(data.database.name, data.name);
          }
        })
        .catch((e) => {});
    },
    async loadNode(node, resolve) {
      if (node.level === 0) {
        let databases = await this.loadDatabases();

        let list = [];
        databases.forEach((one) => {
          let database = {};
          database.name = one.name;
          database.isDatabase = true;
          database.key = "database:" + database.name;
          database.leaf = false;

          list.push(database);
        });

        resolve(list);
        this.initTreeWidth();
        return;
      }
      if (node.data.isDatabase) {
        let database = node.data;
        resolve([
          {
            name: "Tables",
            isDatabaseTables: true,
            key: "database:tables:" + database.name,
            leaf: false,
            database: database,
          },
        ]);
        this.initTreeWidth();
        return;
      }
      if (node.data.isDatabaseTables) {
        let database = node.data.database;
        let tables = await this.loadTables(database.name);
        let list = [];
        tables.forEach((one) => {
          let table = {};
          table.name = one.name;
          table.database = database;
          table.isTable = true;
          table.key = "database:" + database.name + ":" + table.name;
          table.leaf = true;

          list.push(table);
        });
        resolve(list);
        this.initTreeWidth();
      }
    },
    initTreeWidth() {
      // setTimeout(() => {
      //   this.$nextTick(() => {
      //     this.tool.initTreeWidth(this.$refs.tree, this.$refs.treeBox);
      //   });
      // }, 100);
    },
    async loadDatabases() {
      let param = {};
      let res = await this.wrap.work("databases", param);
      res.data = res.data || {};
      return res.data.databases || [];
    },
    async loadTables(database) {
      let param = {
        database: database,
      };
      let res = await this.wrap.work("tables", param);
      res.data = res.data || {};
      return res.data.tables || [];
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-database {
  width: 100%;
  height: 100%;
}
</style>

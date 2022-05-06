<template>
  <div class="toolbox-database-database">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="50px">
          <div class="pdlr-10 pdt-10">
            <div class="tm-btn tm-btn-sm bg-grey-6 ft-13" @click="refresh">
              刷新
            </div>
            <div class="tm-btn tm-btn-sm bg-teal-8 ft-13" @click="toInsert">
              新建库
            </div>
            <div class="tm-btn tm-btn-sm bg-green ft-13" @click="toOpenSql">
              新建查询
            </div>
          </div>
        </tm-layout>
        <tm-layout height="auto" class="scrollbar">
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
              <span
                class="toolbox-editor-tree-span"
                slot-scope="{ node, data }"
              >
                <span>{{ node.label }}</span>
                <div class="toolbox-editor-tree-btn-group">
                  <div
                    v-if="data.isDatabase || data.isDatabaseTables"
                    class="tm-link color-grey ft-14 mgr-4"
                    @click="toReloadChildren(data)"
                  >
                    <i class="mdi mdi-reload"></i>
                  </div>
                  <div
                    v-if="data.isTable"
                    class="tm-link color-grey ft-14 mgr-4"
                    title="表数据"
                    @click="toOpenTable(data)"
                  >
                    <i class="mdi mdi-database-outline"></i>
                  </div>
                  <div
                    v-if="data.isDatabase || data.isTable"
                    class="tm-link color-grey ft-13 mgr-4"
                    title="DDL"
                    @click="toShowDDL(data)"
                  >
                    <IconFont
                      class="teamide-suffix-sql"
                      style="vertical-align: -1px"
                    ></IconFont>
                  </div>
                  <div
                    v-if="data.isDatabase || data.isTable"
                    class="tm-link color-orange ft-15 mgr-4"
                    @click="toDelete(data)"
                  >
                    <i class="mdi mdi-delete-outline"></i>
                  </div>
                </div>
              </span>
            </el-tree>
          </div>
        </tm-layout>
      </tm-layout>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxType", "toolbox", "option", "wrap"],
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
    refresh() {
      this.reloadChildren(this.$refs.tree.root);
    },
    toReloadChildren(data) {
      this.tool.stopEvent();
      this.reloadChildren(data);
    },
    reloadChildren(key) {
      this.tool.stopEvent();
      let node = this.$refs.tree.getNode(key);
      if (node) {
        node.loaded = false;
        node.expand();
      }
    },
    async toShowDDL(data) {
      if (data.isDatabase) {
        let extend = {
          name: data.name + ">DDL",
          title: data.name + ">DDL",
          type: "ddl",
          database: data.name,
        };
        this.wrap.openTabByExtend(extend);
      } else if (data.isTable) {
        let extend = {
          name: data.database.name + "." + data.name + ">DDL",
          title: data.database.name + "." + data.name + ">DDL",
          type: "ddl",
          database: data.database.name,
          table: data.name,
        };
        this.wrap.openTabByExtend(extend);
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
    toInsert() {},
    toOpenSql() {
      let extend = {
        name: "新建SQL",
        title: "新建SQL",
        type: "sql",
      };
      this.wrap.openTabByExtend(extend);
    },
    toOpenTable(data) {
      let extend = {
        name: data.database.name + "." + data.name,
        title: data.database.name + "." + data.name,
        type: "data",
        database: data.database.name,
        table: data.name,
      };
      this.wrap.openTabByExtend(extend);
    },
    toDelete(data) {
      this.tool.stopEvent();
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
    async getTableDetail(database, table) {
      let key = database + "-" + table;
      this.tableCache = this.tableCache || {};
      let res = this.tableCache[key];
      if (res == null) {
        res = await this.loadTableDetail(database, table);
        if (res != null) {
          this.tableCache[key] = res;
        }
      }
      return res;
    },
    async loadTableDetail(database, table) {
      let param = {
        database: database,
        table: table,
      };
      let res = await this.wrap.work("tableDetail", param);
      if (res.code != 0) {
        this.tool.error(res.msg);
        return null;
      }
      res.data = res.data || {};
      return res.data.table;
    },
  },
  created() {},
  mounted() {
    this.wrap.getTableDetail = this.getTableDetail;
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

<template>
  <div class="toolbox-database-database">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="50px">
          <div class="pdlr-10 pdt-10">
            <div class="tm-btn tm-btn-xs bg-grey-6" @click="refresh">刷新</div>
            <div class="tm-btn tm-btn-xs bg-teal-8" @click="toCreateDatabase">
              新建库
            </div>
            <div class="tm-btn tm-btn-xs bg-green" @click="toOpenSql">
              新建查询
            </div>
            <div
              class="tm-btn tm-btn-xs bg-grey"
              @click="toolboxWorker.showInfo()"
            >
              信息
            </div>
          </div>
        </tm-layout>
        <tm-layout height="auto" class="app-scroll-bar">
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
              @node-contextmenu="nodeContextmenu"
              @node-expand="nodeExpand"
              @node-collapse="nodeCollapse"
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
  props: ["source", "toolboxWorker", "extend", "databasesChange"],
  data() {
    return {
      ready: false,
      expands: [],
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
      if (this.extend && this.extend.expands) {
        this.expands = this.extend.expands;
      }
      this.ready = true;
    },
    nodeExpand(data) {
      let index = this.expands.indexOf(data.key);
      if (index < 0) {
        this.expands.push(data.key);
        this.toolboxWorker.updateExtend({
          expands: this.expands,
        });
      }
    },
    nodeCollapse(data) {
      let index = this.expands.indexOf(data.key);
      if (index >= 0) {
        this.expands.splice(index, 1);
        this.toolboxWorker.updateExtend({
          expands: this.expands,
        });
      }
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
        if (
          node.data &&
          node.data.isDatabase &&
          node.loaded &&
          node.childNodes
        ) {
          node.childNodes.forEach((one) => {
            one.loaded = false;
            one.expand();
          });
          return;
        }
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
        this.toolboxWorker.openTabByExtend(extend);
      } else if (data.isTable) {
        let extend = {
          name: data.database.name + "." + data.name + ">DDL",
          title: data.database.name + "." + data.name + ">DDL",
          type: "ddl",
          database: data.database.name,
          table: data.name,
        };
        this.toolboxWorker.openTabByExtend(extend);
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
    nodeContextmenu(event, data, node, nodeView) {
      let menus = [];
      if (data.isDatabase || data.isDatabaseTables) {
        menus.push({
          text: "刷新",
          onClick: () => {
            this.toReloadChildren(data);
          },
        });
        menus.push({
          text: "新增表",
          onClick: () => {
            this.toCreateTable(data);
          },
        });
        menus.push({
          text: "新建SQL查询",
          onClick: () => {
            let extend = {
              name: "查询[" + data.name + "]库SQL",
              title: "查询[" + data.name + "]库SQL",
              type: "sql",
              database: data.name,
              executeSQL: "SHOW TABLES;",
            };
            this.toolboxWorker.openTabByExtend(extend);
          },
        });
      }
      if (data.isTable) {
        menus.push({
          text: "查看数据",
          onClick: () => {
            this.toOpenTable(data);
          },
        });
        menus.push({
          text: "新建SQL查询",
          onClick: () => {
            let extend = {
              name:
                "查询[" + data.database.name + "]库[" + data.name + "]表SQL",
              title:
                "查询[" + data.database.name + "]库[" + data.name + "]表SQL",
              type: "sql",
              database: data.database.name,
              executeSQL:
                "SELECT * FROM `" +
                data.database.name +
                "`.`" +
                data.name +
                "`;",
            };
            this.toolboxWorker.openTabByExtend(extend);
          },
        });
        menus.push({
          text: "编辑表",
          onClick: () => {
            this.toUpdateTable(data);
          },
        });
        menus.push({
          text: "导出数据（SQL、Excel等）",
          onClick: () => {
            this.toExport(data);
          },
        });
        menus.push({
          text: "导入数据（策略、SQL、Excel等）",
          onClick: () => {
            this.toImport(data);
          },
        });
      }
      if (data.isDatabase || data.isTable) {
        menus.push({
          text: "查看DDL",
          onClick: () => {
            this.toShowDDL(data);
          },
        });
        menus.push({
          text: "复制名称",
          onClick: async () => {
            let res = await this.tool.clipboardWrite(data.name);
            if (res.success) {
              this.tool.success("复制成功");
            } else {
              this.tool.warn("复制失败，请允许访问剪贴板！");
            }
          },
        });
        menus.push({
          text: "删除",
          onClick: () => {
            this.toDelete(data);
          },
        });
      }

      if (menus.length > 0) {
        this.tool.showContextmenu(menus);
      }
    },
    toOpenSql() {
      let extend = {
        name: "新建SQL",
        title: "新建SQL",
        type: "sql",
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    toOpenTable(data) {
      let extend = {
        name: data.database.name + "." + data.name,
        title: data.database.name + "." + data.name,
        type: "data",
        database: data.database.name,
        table: data.name,
      };
      this.toolboxWorker.openTabByExtend(extend);
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
            await this.doDeleteDatabase(data.name);
            this.refresh();
          } else if (data.isTable) {
            await this.doDeleteTable(data.database.name, data.name);
            this.reloadChildren(data.database);
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
        this.databasesChange(list);
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
    toCreateDatabase() {
      this.toolboxWorker.showCreateDatabase(() => {
        this.refresh();
      });
    },
    toCreateTable(database) {
      if (database.isDatabaseTables) {
        database = database.database;
      }
      let extend = {
        name: "新建[" + database.name + "]库表",
        title: "新建[" + database.name + "]库表",
        type: "table",
        database: database.name,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    async toUpdateTable(table) {
      let database = table.database.name;
      let extend = {
        name: "编辑[" + database + "]库表[" + table.name + "]",
        title: "编辑[" + database + "]库表[" + table.name + "]",
        type: "table",
        database: database,
        table: table.name,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    async toExport(table) {
      let database = table.database.name;
      let extend = {
        name: "导出[" + database + "]库表[" + table.name + "]数据",
        title: "导出[" + database + "]库表[" + table.name + "]数据",
        type: "export",
        database: database,
        table: table.name,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    async toImport(table) {
      let database = table.database.name;
      let extend = {
        name: "导入[" + database + "]库表[" + table.name + "]数据",
        title: "导入[" + database + "]库表[" + table.name + "]数据",
        type: "import",
        database: database,
        table: table.name,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    async loadDatabases() {
      let param = {};
      let res = await this.toolboxWorker.work("databases", param);
      res.data = res.data || {};
      return res.data.databases || [];
    },
    async loadTables(database) {
      let param = {
        database: database,
      };
      let res = await this.toolboxWorker.work("tables", param);
      res.data = res.data || {};
      return res.data.tables || [];
    },
    async getTableDetail(database, table) {
      let res = await this.loadTableDetail(database, table);
      return res;
    },
    async doDeleteDatabase(database) {
      let param = {
        database: database,
      };
      let res = await this.toolboxWorker.work("deleteDatabase", param);
      if (res.code != 0) {
        return false;
      }
      this.tool.success("删除成功");
      return true;
    },
    async doDeleteTable(database, table) {
      let param = {
        database: database,
        table: table,
      };
      let res = await this.toolboxWorker.work("deleteTable", param);
      if (res.code != 0) {
        return false;
      }
      this.tool.success("删除成功");
      return true;
    },
    async loadTableDetail(database, table) {
      let param = {
        database: database,
        table: table,
      };
      let res = await this.toolboxWorker.work("tableDetail", param);
      if (res.code != 0) {
        return null;
      }
      res.data = res.data || {};
      let tableDetail = res.data.table;
      if (tableDetail) {
        tableDetail.columnList = tableDetail.columnList || [];
        tableDetail.indexList = tableDetail.indexList || [];
      }
      return tableDetail;
    },
  },
  created() {},
  mounted() {
    this.toolboxWorker.getTableDetail = this.getTableDetail;
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

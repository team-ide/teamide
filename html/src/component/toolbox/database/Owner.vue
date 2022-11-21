<template>
  <div class="toolbox-database-owner">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="50px">
          <div class="pdlr-10 pdt-10">
            <div class="tm-btn tm-btn-xs bg-grey-6" @click="refresh">刷新</div>
            <div class="tm-btn tm-btn-xs bg-teal-8" @click="toOwnerCreate">
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
            <el-input placeholder="输入关键字进行过滤" v-model="filterText">
            </el-input>
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
              :filter-node-method="filterNode"
            >
              <span
                class="toolbox-editor-tree-span"
                slot-scope="{ node, data }"
              >
                <span>{{ node.label }}</span>
                <div class="toolbox-editor-tree-btn-group">
                  <div
                    v-if="data.isOwner || data.isOwnerTables"
                    class="tm-link color-grey ft-14 mgr-4"
                    @click="toReloadChildren(data)"
                  >
                    <i class="mdi mdi-reload"></i>
                  </div>
                  <div
                    v-if="data.isTable"
                    class="tm-link color-grey ft-14 mgr-4"
                    title="表数据"
                    @click="toTableOpen(data)"
                  >
                    <i class="mdi mdi-database-outline"></i>
                  </div>
                  <div
                    v-if="data.isOwner || data.isTable"
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
                    v-if="data.isOwner || data.isTable"
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
  props: ["source", "toolboxWorker", "extend", "ownersChange"],
  data() {
    return {
      ready: false,
      expands: [],
      defaultProps: {
        children: "children",
        label: "text",
        isLeaf: "leaf",
      },
      filterText: "",
    };
  },
  computed: {},
  watch: {
    filterText(val) {
      this.$refs.tree.filter(val);
    },
  },
  methods: {
    init() {
      if (this.extend && this.extend.expands) {
        this.expands = this.extend.expands;
      }
      this.ready = true;
    },
    filterNode(value, data) {
      if (!value) return true;
      return (
        data.text && data.text.toLowerCase().indexOf(value.toLowerCase()) !== -1
      );
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
      let needDeletes = [];
      needDeletes.push(data.key);
      if (data.isOwner) {
        this.expands.forEach((one) => {
          if (one == "owner:tables:" + data.ownerName) {
            needDeletes.push(one);
          } else if (("" + one).startsWith("owner:" + data.ownerName + ":")) {
            needDeletes.push(one);
          }
        });
      } else if (data.isOwnerTables) {
        this.expands.forEach((one) => {
          if (one == "owner:tables:" + data.owner.ownerName) {
            needDeletes.push(one);
          } else if (
            ("" + one).startsWith("owner:" + data.owner.ownerName + ":")
          ) {
            needDeletes.push(one);
          }
        });
      }
      if (needDeletes.length > 0) {
        needDeletes.forEach((one) => {
          let index = this.expands.indexOf(one);
          if (index >= 0) {
            this.expands.splice(index, 1);
          }
        });
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
        if (node.data && node.data.isOwner && node.loaded && node.childNodes) {
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
      if (data.isOwner) {
        let extend = {
          name: data.ownerName + ">DDL",
          title: data.ownerName + ">DDL",
          type: "ddl",
          ownerName: data.ownerName,
        };
        this.toolboxWorker.openTabByExtend(extend);
      } else if (data.isTable) {
        let extend = {
          name: data.owner.ownerName + "." + data.tableName + ">DDL",
          title: data.owner.ownerName + "." + data.tableName + ">DDL",
          type: "ddl",
          ownerName: data.owner.ownerName,
          tableName: data.tableName,
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
      if (data.isOwner) {
        if (node.expanded) {
          node.collapse();
        } else {
          node.expand();
        }
      } else if (data.isOwnerTables) {
        if (node.expanded) {
          node.collapse();
        } else {
          node.expand();
        }
      } else if (data.isTable) {
        this.toTableOpen(data);
      }
    },
    nodeContextmenu(event, data, node, nodeView) {
      let menus = [];
      if (data.isOwner || data.isOwnerTables) {
        menus.push({
          text: "刷新",
          onClick: () => {
            this.toReloadChildren(data);
          },
        });
        menus.push({
          text: "新增表",
          onClick: () => {
            this.toTableCreate(data);
          },
        });
        menus.push({
          text: "新建SQL查询",
          onClick: () => {
            let extend = {
              name: "查询[" + data.ownerName + "]库SQL",
              title: "查询[" + data.ownerName + "]库SQL",
              type: "sql",
              ownerName: data.ownerName,
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
            this.toTableOpen(data);
          },
        });
        menus.push({
          text: "新建SQL查询",
          onClick: () => {
            let extend = {
              name:
                "查询[" +
                data.owner.ownerName +
                "]库[" +
                data.tableName +
                "]表SQL",
              title:
                "查询[" +
                data.owner.ownerName +
                "]库[" +
                data.tableName +
                "]表SQL",
              type: "sql",
              ownerName: data.owner.ownerName,
              executeSQL:
                "SELECT * FROM " +
                data.owner.ownerName +
                "." +
                data.tableName +
                ";",
            };
            this.toolboxWorker.openTabByExtend(extend);
          },
        });
        menus.push({
          text: "编辑表",
          onClick: () => {
            this.toTableUpdate(data);
          },
        });
      }
      if (data.isOwner || data.isTable) {
        menus.push({
          text: "查看DDL",
          onClick: () => {
            this.toShowDDL(data);
          },
        });
        menus.push({
          text: "导出",
          onClick: () => {
            this.toExport(data);
          },
        });
        menus.push({
          text: "导入",
          onClick: () => {
            this.toImport(data);
          },
        });
        menus.push({
          text: "同步",
          onClick: () => {
            this.toSync(data);
          },
        });
        menus.push({
          text: "复制名称",
          onClick: async () => {
            let res = await this.tool.clipboardWrite(
              data.tableName || data.ownerName
            );
            if (res.success) {
              this.tool.success("复制成功");
            } else {
              this.tool.warn("复制失败，请允许访问剪贴板！");
            }
          },
        });
        menus.push({
          text: "清空数据",
          onClick: () => {
            this.toDataTrim(data);
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
    toTableOpen(data) {
      let extend = {
        name: data.owner.ownerName + "." + data.tableName,
        title: data.owner.ownerName + "." + data.tableName,
        type: "data",
        ownerName: data.owner.ownerName,
        tableName: data.tableName,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    toDelete(data) {
      this.tool.stopEvent();
      let msg = "确认删除";
      if (data.isOwner) {
        msg += "库[" + data.ownerName + "]";
      } else if (data.isTable) {
        msg +=
          "库[" + data.owner.ownerName + "]" + "表[" + data.tableName + "]";
      }
      msg += "?";
      this.tool
        .confirm(msg)
        .then(async () => {
          if (data.isOwner) {
            await this.doOwnerDelete(data.ownerName);
            this.refresh();
          } else if (data.isTable) {
            await this.doTableDelete(data.owner.ownerName, data.tableName);
            this.reloadChildren(data.owner);
          }
        })
        .catch((e) => {});
    },
    toDataTrim(data) {
      this.tool.stopEvent();
      let msg = "清空";
      if (data.isOwner) {
        msg += "库[" + data.ownerName + "]";
      } else if (data.isTable) {
        msg +=
          "库[" + data.owner.ownerName + "]" + "表[" + data.tableName + "]";
      }
      msg += "数据，将无法恢复，确认清空?";
      this.tool
        .confirm(msg)
        .then(async () => {
          if (data.isOwner) {
            await this.doOwnerDataTrim(data.ownerName);
          } else if (data.isTable) {
            await this.doTableDataTrim(data.owner.ownerName, data.tableName);
          }
        })
        .catch((e) => {});
    },
    async loadNode(node, resolve) {
      if (node.level === 0) {
        let owners = await this.loadOwners();

        let list = [];
        owners.forEach((one) => {
          let owner = {};
          owner.ownerName = one.ownerName;
          owner.text = one.ownerName;
          owner.isOwner = true;
          owner.key = "owner:" + owner.ownerName;
          owner.leaf = false;

          list.push(owner);
        });
        this.ownersChange(list);
        resolve(list);
        this.initTreeWidth();
        return;
      }
      if (node.data.isOwner) {
        let owner = node.data;
        resolve([
          {
            text: "Tables",
            isOwnerTables: true,
            key: "owner:tables:" + owner.ownerName,
            leaf: false,
            owner: owner,
          },
        ]);
        this.initTreeWidth();
        return;
      }
      if (node.data.isOwnerTables) {
        let owner = node.data.owner;
        let tables = await this.loadTables(owner.ownerName);
        let list = [];
        tables.forEach((one) => {
          let table = {};
          table.tableName = one.tableName;
          table.text = one.tableName;
          table.owner = owner;
          table.isTable = true;
          table.key = "owner:" + owner.ownerName + ":" + table.tableName;
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
    toOwnerCreate() {
      this.toolboxWorker.showOwnerCreate(() => {
        this.refresh();
      });
    },
    toTableCreate(owner) {
      if (owner.isOwnerTables) {
        owner = owner.owner;
      }
      let extend = {
        name: "新建[" + owner.ownerName + "]库表",
        title: "新建[" + owner.ownerName + "]库表",
        type: "table",
        ownerName: owner.ownerName,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    async toTableUpdate(table) {
      let ownerName = table.owner.ownerName;
      let extend = {
        name: "编辑[" + ownerName + "]库表[" + table.tableName + "]",
        title: "编辑[" + ownerName + "]库表[" + table.tableName + "]",
        type: "table",
        ownerName: ownerName,
        tableName: table.tableName,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    async toExport(data) {
      let ownerName = null;
      let tableName = null;
      if (data.isOwner) {
        ownerName = data.ownerName;
      } else if (data.isTable) {
        ownerName = data.owner.ownerName;
        tableName = data.tableName;
      }
      let name = "导出[" + ownerName + "]库";
      if (tableName) {
        name += "[" + tableName + "]表";
      }
      let extend = {
        name: name,
        title: name,
        type: "export",
        ownerName: ownerName,
        tableName: tableName,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    async toImport(data) {
      let ownerName = null;
      let tableName = null;
      if (data.isOwner) {
        ownerName = data.ownerName;
      } else if (data.isTable) {
        ownerName = data.owner.ownerName;
        tableName = data.tableName;
      }
      let name = "导入[" + ownerName + "]库";
      if (tableName) {
        name += "[" + tableName + "]表";
      }
      let extend = {
        name: name,
        title: name,
        type: "import",
        ownerName: ownerName,
        tableName: tableName,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    async toSync(data) {
      let ownerName = null;
      let tableName = null;
      if (data.isOwner) {
        ownerName = data.ownerName;
      } else if (data.isTable) {
        ownerName = data.owner.ownerName;
        tableName = data.tableName;
      }
      let name = "同步[" + ownerName + "]库";
      if (tableName) {
        name += "[" + tableName + "]表";
      }
      let extend = {
        name: name,
        title: name,
        type: "sycn",
        ownerName: ownerName,
        tableName: tableName,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    async loadOwners() {
      let param = {};
      let res = await this.toolboxWorker.work("owners", param);
      res.data = res.data || {};
      return res.data.owners || [];
    },
    async loadTables(ownerName) {
      let param = {
        ownerName: ownerName,
      };
      let res = await this.toolboxWorker.work("tables", param);
      res.data = res.data || {};
      return res.data.tables || [];
    },
    async getTableDetail(ownerName, tableName) {
      let res = await this.loadTableDetail(ownerName, tableName);
      return res;
    },
    async doOwnerDelete(ownerName) {
      let param = {
        ownerName: ownerName,
      };
      let res = await this.toolboxWorker.work("ownerDelete", param);
      if (res.code != 0) {
        return false;
      }
      this.tool.success("删除成功");
      return true;
    },
    async doTableDelete(ownerName, tableName) {
      let param = {
        ownerName: ownerName,
        tableName: tableName,
      };
      let res = await this.toolboxWorker.work("tableDelete", param);
      if (res.code != 0) {
        return false;
      }
      this.tool.success("删除成功");
      return true;
    },
    async doOwnerDataTrim(ownerName) {
      let param = {
        ownerName: ownerName,
      };
      let res = await this.toolboxWorker.work("ownerDataTrim", param);
      if (res.code != 0) {
        return false;
      }
      this.tool.success("清空成功");
      return true;
    },
    async doTableDataTrim(ownerName, tableName) {
      let param = {
        ownerName: ownerName,
        tableName: tableName,
      };
      let res = await this.toolboxWorker.work("tableDataTrim", param);
      if (res.code != 0) {
        return false;
      }
      this.tool.success("清空成功");
      return true;
    },
    async loadTableDetail(ownerName, tableName) {
      let param = {
        ownerName: ownerName,
        tableName: tableName,
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
    this.toolboxWorker.loadOwners = this.loadOwners;
    this.toolboxWorker.loadTables = this.loadTables;
    this.init();
  },
};
</script>

<style>
.toolbox-database-owner {
  width: 100%;
  height: 100%;
}
</style>

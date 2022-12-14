<template>
  <div class="toolbox-redis-topic">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="90px">
          <el-form class="pdt-10 pdlr-10" size="mini" inline>
            <el-form-item label="Database" class="mgb-5">
              <el-input v-model="searchForm.database" style="width: 80px" />
            </el-form-item>
            <el-form-item label="搜索" class="mgb-5">
              <el-input v-model="searchForm.pattern" style="width: 160px" />
            </el-form-item>
            <el-form-item label="数量" class="mgb-5">
              <el-input v-model="searchForm.size" style="width: 55px" />
            </el-form-item>
            <el-form-item label="" class="mgb-5">
              <div class="tm-btn tm-btn-xs bg-teal-8" @click="toSearch">
                搜索
              </div>
              <div class="tm-btn tm-btn-xs bg-green" @click="toInsert">
                新增
              </div>
              <div class="tm-btn tm-btn-xs bg-grey" @click="toImport()">
                导入
              </div>
              <div class="tm-btn tm-btn-xs bg-grey" @click="toExport()">
                导出
              </div>
              <div
                class="tm-btn tm-btn-xs bg-grey"
                @click="toolboxWorker.showInfo()"
              >
                信息
              </div>
              <div
                class="tm-btn tm-btn-xs bg-orange"
                @click="
                  toDeletePattern(searchForm.database, searchForm.pattern)
                "
              >
                删除
              </div>
            </el-form-item>
          </el-form>
        </tm-layout>
        <tm-layout height="auto">
          <template v-if="searchResult == null">
            <div class="text-center ft-13 pdtb-10">数据加载中，请稍后!</div>
          </template>
          <template v-else>
            <div class="text-center ft-12 pdtb-10" style="height: 40px">
              共
              <span class="color-green-2 pdlr-3">
                {{ searchResult.count }}</span
              >
              Keys
              <span class="pdlr-2"></span>
              加载
              <span class="color-green pdlr-3">
                {{ searchResult.keyList.length }}
              </span>
              个
              <el-radio-group v-model="viewModel">
                <el-radio label="list" class="mglr-0 mgl-10">列表</el-radio>
                <el-radio label="tree" class="mglr-0 mgl-10"
                  >树形(分割符号'{{ splitChars.join("','") }}')</el-radio
                >
              </el-radio-group>
            </div>
            <div>
              <el-input placeholder="输入关键字进行过滤" v-model="filterText">
              </el-input>
            </div>
            <div
              class="app-scroll-bar"
              style="height: calc(100% - 80px); user-select: text"
            >
              <div class="pd-10">
                <el-tree
                  ref="tree"
                  :props="defaultProps"
                  node-key="key"
                  :expand-on-click-node="false"
                  :data="
                    viewModel == 'list'
                      ? searchResult.keyList
                      : searchResult.treeDatas
                  "
                  :filter-node-method="filterNode"
                  @node-click="nodeClick"
                  @node-contextmenu="nodeContextmenu"
                  empty-text="暂无匹配数据"
                >
                  <span
                    class="toolbox-editor-tree-span"
                    slot-scope="{ node, data }"
                  >
                    <span>{{ node.label }}</span>
                    <template v-if="data.isData">
                      <div class="toolbox-editor-tree-btn-group">
                        <div
                          class="tm-link color-green ft-15 mgr-2"
                          @click="toUpdate(data)"
                        >
                          <i class="mdi mdi-text-box-edit-outline"></i>
                        </div>
                        <div
                          class="tm-link color-orange ft-15 mgr-2"
                          @click="toDelete(data)"
                        >
                          <i class="mdi mdi-delete-outline"></i>
                        </div>
                      </div>
                    </template>
                  </span>
                </el-tree>
              </div>
            </div>
          </template>
        </tm-layout>
      </tm-layout>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      ready: false,
      searchForm: {
        database: 0,
        pattern: "xx*",
        size: 200,
      },
      splitChars: [":", "-", "/"],
      searchResult: null,
      viewModel: "list",
      filterText: "",
      defaultProps: {
        children: "children",
        label: "name",
        isLeaf: "leaf",
      },
    };
  },
  computed: {},
  watch: {
    filterText(val) {
      this.$refs.tree && this.$refs.tree.filter(val);
    },
    viewModel() {
      this.$nextTick(() => {
        this.$refs.tree && this.$refs.tree.filter(this.filterText);
      });
    },
  },
  methods: {
    filterNode(value, data) {
      if (!value) return true;
      return (
        data.key && data.key.toLowerCase().indexOf(value.toLowerCase()) !== -1
      );
    },
    init() {
      this.ready = true;
      if (this.extend && this.extend.search) {
        if (this.extend.search.pattern) {
          this.searchForm.pattern = this.extend.search.pattern;
        }
        if (this.extend.search.database >= 0) {
          this.searchForm.database = this.extend.search.database;
        }
      }
      this.toolboxWorker.loadKeys = this.loadKeys;
      this.loadKeys();
    },
    refresh() {
      this.toSearch();
    },
    toSearch() {
      this.loadKeys();
      this.toolboxWorker.updateExtend({
        search: {
          pattern: this.searchForm.pattern,
          database: Number(this.searchForm.database),
        },
      });
    },
    async loadKeys() {
      this.searchResult = null;
      this.searchForm.database = Number(this.searchForm.database);

      let param = this.toolboxWorker.getWorkParam(
        Object.assign({}, this.searchForm)
      );
      Object.assign(param, this.searchForm);
      if (this.tool.isEmpty(param.size)) {
        param.size = 50;
      }
      param.size = Number(param.size);
      if (this.tool.isEmpty(param.pattern)) {
        this.tool.warn("请输入“*”或“user*”等关键字模糊搜索");
      }
      let res = await this.server.redis.keys(param);
      if (res.code == 0) {
        let keysData = res.data || {};
        this.formatData(keysData);
        this.searchResult = keysData;
      } else {
        this.tool.error(res.msg);
      }
    },
    formatData(keysData) {
      keysData = keysData || {};
      keysData.keyList = keysData.keyList || [];
      keysData.treeDatas = [];
      var treeDataCache = {};
      keysData.keyList.forEach((data) => {
        data.isData = true;
        data.name = data.key;
        let treeData = {
          database: data.database,
          key: data.key,
          name: data.key,
          isData: true,
          children: [],
        };
        let lastFind = null;
        let splitChar = null;
        this.splitChars.forEach((one) => {
          if (splitChar == null) {
            if (data.key.indexOf(one) >= 0) {
              splitChar = one;
            }
          }
        });
        if (data.key.indexOf(splitChar) >= 0) {
          let ss = data.key.split(splitChar);
          let lastK = "";
          ss.forEach((s, i) => {
            treeData.name = s;
            if (i >= ss.length - 1) {
              return;
            }
            if (i > 0) {
              lastK += splitChar;
            }
            lastK += s;
            let find = treeDataCache[lastK];
            if (find == null) {
              find = {
                database: data.database,
                key: lastK,
                name: s,
                isData: false,
                children: [],
              };
              treeDataCache[lastK] = find;
              if (lastFind != null) {
                lastFind.children.push(find);
              } else {
                keysData.treeDatas.push(find);
              }
            }
            lastFind = find;
          });
        }
        if (lastFind != null) {
          lastFind.children.push(treeData);
        } else {
          keysData.treeDatas.push(treeData);
        }
      });
    },
    rowClick(data) {
      this.rowClickTimeCache = this.rowClickTimeCache || {};
      let nowTime = new Date().getTime();
      let clickTime = this.rowClickTimeCache[data];
      this.rowClickTimeCache[data] = nowTime;
      if (clickTime) {
        let timeout = nowTime - clickTime;
        if (timeout < 300) {
          delete this.rowClickTimeCache[data];
          this.rowDbClick(data);
        }
      }
    },
    rowDbClick(data) {
      this.toUpdate(data);
    },

    nodeClick(data, node) {
      this.rowClickTimeCache = this.rowClickTimeCache || {};
      let nowTime = new Date().getTime();
      let clickTime = this.rowClickTimeCache[node];
      this.rowClickTimeCache[node] = nowTime;
      if (clickTime) {
        let timeout = nowTime - clickTime;
        if (timeout < 300) {
          delete this.rowClickTimeCache[node];
          this.nodeDbClick(node);
        }
      }
    },
    nodeDbClick(node) {
      if (node.data.isData) {
        this.toUpdate(node.data);
        return;
      }
      if (node.expanded) {
        node.expanded = false;
      } else {
        node.loaded = false;
        node.expand();
      }
    },
    nodeContextmenu(event, data, node, nodeView) {
      if (data.isData) {
        this.dataContextmenu(data);
      }
    },
    dataContextmenu(data) {
      let menus = [];
      menus.push({
        header: data.key,
      });
      menus.push({
        text: "修改",
        onClick: () => {
          this.toUpdate(data);
        },
      });
      menus.push({
        text: "新增",
        onClick: () => {
          this.toInsert();
        },
      });
      menus.push({
        text: "导入",
        onClick: () => {
          this.toImport(data);
        },
      });
      menus.push({
        text: "导出",
        onClick: () => {
          this.toExport(data);
        },
      });
      menus.push({
        text: "删除",
        onClick: () => {
          this.toDelete(data);
        },
      });

      if (menus.length > 0) {
        this.tool.showContextmenu(menus);
      }
    },
    async toExport(data) {
      this.tool.warn("功能还未完善，敬请期待！");
      return;
    },
    async toImport(data) {
      this.tool.warn("功能还未完善，敬请期待！");
      return;
      // data = data || {};
      // let extend = {
      //   name: "导入",
      //   title: "导入",
      //   type: "import",
      //   database: data.database || this.searchForm.database,
      // };
      // this.toolboxWorker.openTabByExtend(extend);
    },
    toInsert() {
      let extend = {
        name: "新增数据",
        title: "新增数据",
        type: "data",
        database: this.searchForm.database,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    toUpdate(data) {
      let extend = {
        name: `编辑[${data.key}]数据`,
        title: `编辑[${data.key}]数据`,
        type: "data",
        key: data.key,
        database: data.database,
        onlyOpenOneKey: "redis:data:key" + data.key,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    toDelete(data) {
      this.tool
        .confirm("确认删除[" + data.key + "]？")
        .then(async () => {
          this.doDelete(data.database, data.key);
        })
        .catch((e) => {});
    },
    toDeletePattern(database, pattern) {
      this.tool
        .confirm("将删除所有匹配[" + pattern + "]的Key，确定删除？")
        .then(async () => {
          this.doDeletePattern(database, pattern);
        })
        .catch((e) => {});
    },
    async doDelete(database, key) {
      let param = this.toolboxWorker.getWorkParam({
        database: Number(database),
        key: key,
      });
      let res = await this.server.redis.delete(param);
      if (res.code == 0) {
        this.tool.success("删除成功!");
        this.toSearch();
      } else {
        this.tool.error(res.msg);
      }
    },
    async doDeletePattern(database, pattern) {
      let param = this.toolboxWorker.getWorkParam({
        database: Number(database),
        pattern: pattern,
      });
      let res = await this.server.redis.deletePattern(param);
      if (res.code == 0) {
        this.tool.success("删除成功!");
        this.toSearch();
      } else {
        this.tool.error(res.msg);
      }
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-redis-topic {
  width: 100%;
  height: 100%;
}
</style>

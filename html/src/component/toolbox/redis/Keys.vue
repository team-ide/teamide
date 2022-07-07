<template>
  <div class="toolbox-redis-topic">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="90px">
          <el-form class="pdt-10 pdlr-10" size="mini" inline>
            <el-form-item label="Database" label-width="70px" class="mgb-5">
              <el-input v-model="searchForm.database" style="width: 80px" />
            </el-form-item>
            <el-form-item label="Key(支持*模糊搜索)" class="mgb-5">
              <el-input v-model="searchForm.pattern" style="width: 150px" />
            </el-form-item>
            <el-form-item label="数量" label-width="70px" class="mgb-5">
              <el-input v-model="searchForm.size" style="width: 80px" />
            </el-form-item>
            <el-form-item label="" class="mgb-5">
              <div class="tm-btn tm-btn-sm bg-teal-8 ft-13" @click="toSearch">
                搜索
              </div>
              <div class="tm-btn tm-btn-sm bg-green ft-13" @click="toInsert">
                新增
              </div>
              <div class="tm-btn tm-btn-sm bg-grey ft-13" @click="toImport()">
                导入
              </div>
              <div class="tm-btn tm-btn-sm bg-grey ft-13" @click="toExport()">
                导出
              </div>
              <div
                class="tm-btn tm-btn-sm bg-orange ft-13"
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
          <template
            v-else-if="
              searchResult.dataList == null || searchResult.dataList.length == 0
            "
          >
            <div class="text-center ft-13 pdtb-10">暂无匹配数据!</div>
          </template>
          <template v-else>
            <div class="text-center ft-13 pdtb-10" style="height: 40px">
              Keys （{{ searchResult.count }}）
            </div>
            <div
              class="data-list-box scrollbar"
              style="height: calc(100% - 40px); user-select: text"
            >
              <template v-for="(one, index) in searchResult.dataList">
                <div
                  :key="index"
                  class="data-list-one"
                  @click="rowClick(one)"
                  @contextmenu="dataContextmenu(one)"
                >
                  <div class="data-list-one-text">
                    {{ one.key }}
                  </div>
                </div>
              </template>
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
        size: 50,
      },
      searchResult: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
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
      let param = {};
      this.searchForm.database = Number(this.searchForm.database);
      Object.assign(param, this.searchForm);
      if (this.tool.isEmpty(param.size)) {
        param.size = 50;
      }
      param.size = Number(param.size);
      let res = await this.toolboxWorker.work("keys", param);
      this.searchResult = res.data;
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
      data = data || {};
      let extend = {
        name: "导入",
        title: "导入",
        type: "import",
        database: data.database || this.searchForm.database,
      };
      this.toolboxWorker.openTabByExtend(extend);
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
      let param = {
        database: Number(database),
        key: key,
      };
      let res = await this.toolboxWorker.work("delete", param);
      if (res.code == 0) {
        this.tool.success("删除成功!");
        this.toSearch();
      }
    },
    async doDeletePattern(database, pattern) {
      let param = {
        database: Number(database),
        pattern: pattern,
      };
      let res = await this.toolboxWorker.work("deletePattern", param);
      if (res.code == 0) {
        this.tool.success("删除成功!");
        this.toSearch();
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

<template>
  <div class="toolbox-elasticsearch-indexName">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="80px">
          <el-form class="pdt-10 pdlr-10" size="mini" inline>
            <el-form-item label="搜索" class="mgb-5">
              <el-input v-model="searchForm.pattern" style="width: 300px" />
            </el-form-item>
            <el-form-item label="" class="mgb-5">
              <div class="tm-btn tm-btn-xs bg-grey-6" @click="loadIndexNames">
                刷新
              </div>
              <div class="tm-btn tm-btn-xs bg-teal-8" @click="toInsert">
                新建索引
              </div>
              <div
                class="tm-btn tm-btn-xs bg-grey"
                @click="toolboxWorker.showInfo()"
              >
                信息
              </div>
            </el-form-item>
          </el-form>
        </tm-layout>
        <tm-layout height="auto">
          <template v-if="indexNames == null">
            <div class="text-center ft-13 pdtb-10">数据加载中，请稍后!</div>
          </template>
          <template v-else-if="indexNames.length == 0">
            <div class="text-center ft-13 pdtb-10">暂无匹配数据!</div>
          </template>
          <template v-else>
            <div class="text-center ft-13 pdtb-10" style="height: 40px">
              IndexNames （{{ indexNames.length }}）
            </div>
            <div
              class="data-list-box scrollbar"
              style="height: calc(100% - 40px); user-select: text"
            >
              <template v-for="(one, index) in indexNames">
                <div
                  :key="index"
                  v-if="
                    tool.isEmpty(searchForm.pattern) ||
                    one.name
                      .toLowerCase()
                      .indexOf(searchForm.pattern.toLowerCase()) >= 0
                  "
                  class="data-list-one"
                  @click="rowClick(one)"
                  @contextmenu="dataContextmenu(one)"
                >
                  <div class="data-list-one-text">
                    {{ one.name }}
                  </div>
                </div>
              </template>
            </div>
          </template>
        </tm-layout>
      </tm-layout>
    </template>
    <FormDialog
      ref="InsertIndexName"
      :source="source"
      title="新增索引"
      :onSave="doInsert"
    ></FormDialog>
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
        pattern: null,
      },
      indexNames: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.toolboxWorker.getMapping = this.getMapping;
      this.ready = true;
      this.loadIndexNames();
    },
    refresh() {
      this.loadIndexNames();
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
      this.toOpenIndexName(data);
    },
    toOpenIndexName(data) {
      let extend = {
        name: data.name,
        title: data.name,
        type: "data",
        indexName: data.name,
      };
      this.toolboxWorker.openTabByExtend(extend);
    },
    dataContextmenu(data) {
      let menus = [];
      menus.push({
        header: data.name,
      });
      menus.push({
        text: "数据",
        onClick: () => {
          this.toOpenIndexName(data);
        },
      });
      menus.push({
        text: "结构",
        onClick: () => {
          this.toUpdateMapping(data);
        },
      });
      menus.push({
        text: "新增索引",
        onClick: () => {
          this.toInsert();
        },
      });
      menus.push({
        text: "迁移",
        onClick: () => {
          this.toReindex(data);
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
    toInsert() {
      let data = {};

      this.$refs.InsertIndexName.show({
        title: `新增索引`,
        form: [this.form.toolbox.elasticsearch.index],
        data: [data],
      });
    },
    async doInsert(dataList) {
      let data = dataList[0];
      let param = {
        indexName: data.indexName,
        mapping: data.mapping,
      };
      let res = await this.toolboxWorker.work("createIndex", param);
      if (res.code == 0) {
        await this.loadIndexNames();
        return true;
      } else {
        return false;
      }
    },
    toReindex(data) {
      this.toolboxWorker.showReindexForm(
        {
          indexName: data.name,
        },
        async (m) => {
          let flag = await this.doReindex(m);
          return flag;
        }
      );
    },
    async doReindex(data) {
      let param = {
        sourceIndexName: data.sourceIndexName,
        destIndexName: data.destIndexName,
      };
      let res = await this.toolboxWorker.work("reindex", param);
      if (res.code == 0) {
        await this.loadIndexNames();
        return true;
      } else {
        return false;
      }
    },
    toDelete(data) {
      let msg = "确认删除";
      msg += "索引[" + data.name + "]";
      msg += "?";
      this.tool
        .confirm(msg)
        .then(async () => {
          this.doDelete(data.name);
        })
        .catch((e) => {});
    },
    async toUpdateMapping(data) {
      let indexName = data.name;
      let mapping = await this.getMapping(indexName);
      this.toolboxWorker.showMappingForm(
        {
          indexName: indexName,
          mapping: mapping,
        },
        async (m) => {
          let flag = await this.putMapping(m.indexName, m.mapping);
          return flag;
        }
      );
    },
    async doDelete(indexName) {
      let param = {
        indexName: indexName,
      };
      let res = await this.toolboxWorker.work("deleteIndex", param);
      if (res.code == 0) {
        this.tool.success("删除成功!");
        this.loadIndexNames();
      }
    },
    async loadIndexNames() {
      this.indexNames = null;
      let param = {};
      let res = await this.toolboxWorker.work("indexNames", param);
      res.data = res.data || {};
      res.data.indexNames = res.data.indexNames || [];
      let indexNames = [];
      res.data.indexNames.forEach((one) => {
        let indexName = {};
        indexName.name = one;

        indexNames.push(indexName);
      });
      this.indexNames = indexNames;
    },
    async getMapping(indexName) {
      let param = {
        indexName: indexName,
      };
      let res = await this.toolboxWorker.work("getMapping", param);
      res.data = res.data || {};
      res.data.mapping = res.data.mapping || {};
      return res.data.mapping;
    },
    async putMapping(indexName, mapping) {
      let param = {
        indexName: indexName,
        mapping: mapping,
      };
      let res = await this.toolboxWorker.work("putMapping", param);
      if (res.code != 0) {
        return false;
      }
      return true;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-elasticsearch-indexName {
  width: 100%;
  height: 100%;
}
</style>

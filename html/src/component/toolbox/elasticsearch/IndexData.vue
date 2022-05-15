<template>
  <div class="toolbox-elasticsearch-indexName-data">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="140px">
          <el-form
            class="pdt-20 mglr-10"
            ref="form"
            :model="form"
            label-width="90px"
            size="mini"
            :inline="true"
          >
            <el-form-item label="IndexName">
              <el-input v-model="searchForm.indexName" />
            </el-form-item>
            <el-form-item label="PageIndex">
              <el-input v-model="searchForm.pageIndex" />
            </el-form-item>
            <el-form-item label="PageSize">
              <el-input v-model="searchForm.pageSize" />
            </el-form-item>
            <div class="pdt-25">
              <div class="tm-btn tm-btn-sm bg-teal-8 ft-13" @click="toSearch">
                搜索
              </div>
              <div class="tm-btn tm-btn-sm bg-green ft-13" @click="toInsert">
                新增
              </div>
            </div>
          </el-form>
        </tm-layout>
        <tm-layout height="auto" class="scrollbar">
          <div class="pd-10" style="o"></div>
        </tm-layout>
      </tm-layout>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxType", "toolbox", "option", "indexName", "wrap"],
  data() {
    return {
      ready: false,
      searchForm: {
        indexName: this.indexName,
        pageIndex: 1,
        pageSize: 10,
      },
      msgs: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      this.ready = true;
      // await this.toSearch();
    },
    async toSearch() {
      await this.doSearch();
    },
    toIndex() {},
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
      this.wrap.showData(data);
    },
    toDelete(data) {
      let indexName = data.indexName;
      let _id = data._id;

      let msg = "确认删除索引[" + indexName + "]数据[" + _id + "]";
      msg += "?";
      this.tool
        .confirm(msg)
        .then(async () => {
          this.doDelete({ indexName, _id });
        })
        .catch((e) => {});
    },
    toInsert() {
      this.tool.warn("暂不支持ES数据新增，敬请期待！");
    },
    async doDelete(data) {
      let param = {};
      Object.assign(param, data);
      let res = await this.wrap.work("deleteData", param);
      if (res.code == 0) {
        this.doSearch();
        return true;
      } else {
        return false;
      }
    },
    async doSearch() {
      this.tool.warn("暂不支持ES查询，敬请期待！");
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-elasticsearch-indexName-data {
  width: 100%;
  height: 100%;
}
</style>

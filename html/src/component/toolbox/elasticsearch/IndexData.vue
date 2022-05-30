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
          <div class="pd-10" style="o">
            <table>
              <thead>
                <tr>
                  <th width="100">_id</th>
                  <th width="">_source</th>
                  <th width="150">操作</th>
                </tr>
              </thead>
              <tbody>
                <template v-if="dataList == null">
                  <tr>
                    <td colspan="3">
                      <div class="text-center ft-13 pdtb-10">查询中...</div>
                    </td>
                  </tr>
                </template>
                <template v-else-if="dataList.length == 0">
                  <tr>
                    <td colspan="3">
                      <div class="text-center ft-13 pdtb-10">暂无匹配数据!</div>
                    </td>
                  </tr>
                </template>
                <template v-else>
                  <template v-for="(one, index) in dataList">
                    <tr :key="index" @click="rowClick(one)">
                      <td>{{ one._id }}</td>
                      <td>{{ one._source }}</td>
                      <td>
                        <div style="width: 150px">
                          <div
                            class="tm-btn color-grey tm-btn-xs"
                            @click="wrap.showData(one)"
                            title="查看"
                          >
                            <i class="mdi mdi-eye-outline"></i>
                          </div>
                          <div
                            class="tm-btn color-blue tm-btn-xs"
                            @click="toUpdate(one)"
                          >
                            修改
                          </div>
                          <div
                            class="tm-btn color-grey tm-btn-xs"
                            @click="toCopy(one)"
                          >
                            复制
                          </div>
                          <div
                            class="tm-btn color-orange tm-btn-xs"
                            @click="toDelete(one)"
                            title="删除"
                          >
                            <i class="mdi mdi-delete-outline"></i>
                          </div>
                        </div>
                      </td>
                    </tr>
                  </template>
                </template>
              </tbody>
            </table>
          </div>
        </tm-layout>
      </tm-layout>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: [
    "source",
    "toolboxType",
    "toolbox",
    "tab",
    "option",
    "indexName",
    "wrap",
  ],
  data() {
    return {
      ready: false,
      searchForm: {
        indexName: this.indexName,
        pageIndex: 1,
        pageSize: 10,
      },
      dataList: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      this.ready = true;
      await this.toSearch();
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
    async toInsert() {
      let indexName = this.indexName;
      let data = {
        indexName: indexName,
      };
      let mapping = await this.wrap.getMapping(indexName);
      this.wrap.showDataForm(data, mapping, async (m) => {
        let flag = await this.doInsert(m);
        return flag;
      });
    },
    async doInsert(data) {
      let param = {
        indexName: data.indexName,
        doc: data.doc,
        id: data.id,
      };
      let res = await this.wrap.work("insertData", param);
      if (res.code == 0) {
        await this.toSearch();
        return true;
      } else {
        return false;
      }
    },
    async toCopy(data) {
      let indexName = this.indexName;
      let param = {
        indexName: indexName,
        doc: data._source,
        id: data._id + "xxx",
      };
      let mapping = await this.wrap.getMapping(indexName);
      this.wrap.showDataForm(param, mapping, async (m) => {
        let flag = await this.doInsert(m);
        return flag;
      });
    },
    async toUpdate(data) {
      let indexName = this.indexName;
      let param = {
        indexName: indexName,
        doc: data._source,
        id: data._id,
      };
      let mapping = await this.wrap.getMapping(indexName);
      this.wrap.showDataForm(param, mapping, async (m) => {
        let flag = await this.doUpdate(m);
        return flag;
      });
    },
    async doUpdate(data) {
      let param = {
        indexName: data.indexName,
        doc: data.doc,
        id: data.id,
      };
      let res = await this.wrap.work("updateData", param);
      if (res.code == 0) {
        await this.toSearch();
        return true;
      } else {
        return false;
      }
    },
    async doDelete(data) {
      let indexName = this.indexName;
      let param = {
        indexName: indexName,
        id: data._id,
      };
      let res = await this.wrap.work("deleteData", param);
      if (res.code == 0) {
        this.doSearch();
        return true;
      } else {
        return false;
      }
    },
    async doSearch() {
      let param = {};
      this.searchForm.pageIndex = Number(this.searchForm.pageIndex);
      this.searchForm.pageSize = Number(this.searchForm.pageSize);
      Object.assign(param, this.searchForm);
      let res = await this.wrap.work("search", param);
      res.data = res.data || {};
      let result = res.data.result || {};
      let hits = result.hits || {};
      this.dataList = hits.hits || [];
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

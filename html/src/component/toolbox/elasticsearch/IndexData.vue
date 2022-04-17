<template>
  <div class="toolbox-elasticsearch-indexName-data">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="140px">
          <b-form inline class="pdt-20 mglr-10">
            <b-form-group label="Topic" label-size="sm" class="pdr-10">
              <b-form-input size="sm" v-model="searchForm.indexName">
              </b-form-input>
            </b-form-group>
            <b-form-group label="PageIndex" label-size="sm" class="pdr-10">
              <b-form-input size="sm" v-model="searchForm.pageIndex">
              </b-form-input>
            </b-form-group>
            <b-form-group label="PageSize" label-size="sm" class="pdr-10">
              <b-form-input size="sm" v-model="searchForm.pageSize">
              </b-form-input>
            </b-form-group>
            <b-form-group label="">
              <div class="pdt-25">
                <div class="tm-btn tm-btn-sm bg-teal-8 ft-13" @click="toSearch">
                  搜索
                </div>
                <div class="tm-btn tm-btn-sm bg-green ft-13" @click="toIndex">
                  新增
                </div>
              </div>
            </b-form-group>
          </b-form>
        </tm-layout>
        <tm-layout height="auto" class="scrollbar">
          <div class="pd-10" style="o">
            <table>
              <thead>
                <tr>
                  <th width="100">Partition</th>
                  <th width="80">Offset</th>
                  <th>Key</th>
                  <th>Value</th>
                  <th width="150"></th>
                </tr>
              </thead>
              <tbody>
                <template v-if="msgs == null">
                  <tr>
                    <td colspan="5">
                      <div class="text-center ft-13 pdtb-10">拉取中...</div>
                    </td>
                  </tr>
                </template>
                <template v-else-if="msgs.length == 0">
                  <tr>
                    <td colspan="5">
                      <div class="text-center ft-13 pdtb-10">暂无匹配数据!</div>
                    </td>
                  </tr>
                </template>
                <template v-else>
                  <template v-for="(one, index) in msgs">
                    <tr :key="index" @click="rowClick(one)">
                      <td>{{ one.partition }}</td>
                      <td>{{ one.offset }}</td>
                      <td>{{ one.key }}</td>
                      <td>{{ one.value }}</td>
                      <td>
                        <div style="width: 140px">
                          <div
                            class="tm-btn color-grey tm-btn-xs"
                            @click="wrap.showData(one)"
                          >
                            查看
                          </div>
                          <div
                            class="tm-btn color-orange tm-btn-xs"
                            @click="toDelete(one)"
                          >
                            删除
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
    init() {
      this.ready = true;
      this.toSearch();
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
      this.msgs = null;
      let param = {};
      Object.assign(param, this.searchForm);
      let res = await this.wrap.work("search", param);
      res.data = res.data || {};
      let result = res.data.result || {};
      this.result = result;
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

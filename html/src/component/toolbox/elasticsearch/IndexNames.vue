<template>
  <div class="toolbox-elasticsearch-indexName">
    <template v-if="ready">
      <div class="pd-10">
        <table>
          <thead>
            <tr>
              <th>IndexName</th>
              <th>
                <div style="width: 120px">
                  <div
                    class="tm-link color-grey-3 ft-14 mglr-2"
                    @click="loadIndexNames()"
                  >
                    <i class="mdi mdi-reload"></i>
                  </div>
                  <div
                    class="tm-link color-green-3 ft-14 mglr-2"
                    @click="toInsert()"
                  >
                    <i class="mdi mdi-plus"></i>
                  </div>
                </div>
              </th>
            </tr>
          </thead>
          <tbody>
            <template v-if="indexNames == null">
              <tr>
                <td colspan="2">
                  <div class="text-center ft-13 pdtb-10">加载中...</div>
                </td>
              </tr>
            </template>
            <template v-else-if="indexNames.length == 0">
              <tr>
                <td colspan="2">
                  <div class="text-center ft-13 pdtb-10">暂无匹配数据!</div>
                </td>
              </tr>
            </template>
            <template v-else>
              <template v-for="(one, index) in indexNames">
                <tr :key="index" @click="rowClick(one)">
                  <td>{{ one.name }}</td>
                  <td>
                    <div
                      class="tm-btn color-blue tm-btn-xs"
                      @click="toOpenIndexName(one)"
                    >
                      数据
                    </div>
                    <div
                      class="tm-btn color-orange tm-btn-xs"
                      @click="toDelete(one)"
                    >
                      删除
                    </div>
                  </td>
                </tr>
              </template>
            </template>
          </tbody>
        </table>
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
      indexNames: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.ready = true;
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
      let tab = this.wrap.createTabByData(data);
      this.wrap.addTab(tab);
      this.wrap.doActiveTab(tab);
    },
    toInsert() {
      let data = {};
      this.wrap.showIndexForm(data, (m) => {
        let flag = this.doInsert(m);
        return flag;
      });
    },
    async doInsert(data) {
      let param = {
        indexName: data.indexName,
      };
      let res = await this.wrap.work("createIndex", param);
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
    async doDelete(indexName) {
      let param = {
        indexName: indexName,
      };
      let res = await this.wrap.work("deleteIndex", param);
      if (res.code == 0) {
        this.tool.info("删除成功!");
        this.loadIndexNames();
      }
    },
    async loadIndexNames() {
      this.indexNames = null;
      let param = {};
      let res = await this.wrap.work("indexNames", param);
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

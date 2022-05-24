<template>
  <div class="toolbox-elasticsearch-indexName">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="50px">
          <div class="pdlr-10 pdt-10">
            <div
              class="tm-btn tm-btn-sm bg-grey-6 ft-13"
              @click="loadIndexNames"
            >
              刷新
            </div>
            <div class="tm-btn tm-btn-sm bg-teal-8 ft-13" @click="toInsert">
              新建索引
            </div>
          </div>
        </tm-layout>
        <tm-layout height="auto" class="scrollbar">
          <div class="">
            <table>
              <thead>
                <tr>
                  <th>IndexName</th>
                  <th width="175px">
                    <div style="width: 175px">
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
                          class="tm-btn color-grey tm-btn-xs"
                          @click="toUpdateMapping(one)"
                        >
                          结构
                        </div>
                        <div
                          class="tm-btn color-green tm-btn-xs"
                          @click="toReindex(one)"
                        >
                          迁移
                        </div>
                        <div
                          class="tm-btn color-orange tm-btn-xs"
                          @click="toDelete(one)"
                        >
                          <i class="mdi mdi-delete-outline"></i>
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
  props: ["source", "toolboxType", "toolbox", "option", "wrap"],
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
      this.wrap.getMapping = this.getMapping;
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
      this.wrap.openTabByExtend(extend);
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
      let res = await this.wrap.work("createIndex", param);
      if (res.code == 0) {
        await this.loadIndexNames();
        return true;
      } else {
        return false;
      }
    },
    toReindex(data) {
      this.wrap.showReindexForm(
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
      let res = await this.wrap.work("reindex", param);
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
      this.wrap.showMappingForm(
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
      let res = await this.wrap.work("deleteIndex", param);
      if (res.code == 0) {
        this.tool.success("删除成功!");
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
    async getMapping(indexName) {
      let param = {
        indexName: indexName,
      };
      let res = await this.wrap.work("getMapping", param);
      res.data = res.data || {};
      res.data.mapping = res.data.mapping || {};
      return res.data.mapping;
    },
    async putMapping(indexName, mapping) {
      let param = {
        indexName: indexName,
        mapping: mapping,
      };
      let res = await this.wrap.work("putMapping", param);
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

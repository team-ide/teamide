<template>
  <div class="toolbox-elasticsearch-indexName-data">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="120px" style="overflow: hidden">
          <tm-layout width="400px">
            <ul class="part-box scrollbar mg-0">
              <template v-for="(one, index) in searchForm.whereList">
                <li :key="index">
                  <input v-model="one.checked" type="checkbox" />
                  <template v-if="one.sqlConditionalOperation == 'custom'">
                  </template>
                  <template v-else>
                    <select
                      v-model="one.name"
                      @change="initInputWidth"
                      class="part-form-input"
                    >
                      <option :value="null" text="请选择">请选择</option>
                      <template v-for="(one, index) in pointColumnList">
                        <option
                          :key="index"
                          :value="one.name"
                          :text="one.name + '&nbsp;'"
                        >
                          {{ one.name }}
                          <template v-if="tool.isNotEmpty(one.comment)">
                            （{{ one.comment }}）
                          </template>
                        </option>
                      </template>
                    </select>
                  </template>
                  <select
                    v-model="one.sqlConditionalOperation"
                    @change="initInputWidth"
                    class="part-form-input"
                  >
                    <option value="=">等于</option>
                    <option value="like">包含</option>
                    <option value="not like">不包含</option>
                    <option value="like start">开始以</option>
                    <option value="not like start">开始不是以</option>
                    <option value="like end">结束以</option>
                    <option value="not like end">结束不是以</option>
                    <option value="between">介于</option>
                    <option value="not between">不介于</option>
                    <option value="in">在列表</option>
                    <option value="not in">不在列表</option>
                  </select>
                  <template
                    v-if="
                      one.sqlConditionalOperation == 'is null' ||
                      one.sqlConditionalOperation == 'is not null' ||
                      one.sqlConditionalOperation == 'is empty' ||
                      one.sqlConditionalOperation == 'is not empty'
                    "
                  >
                  </template>
                  <template
                    v-else-if="
                      one.sqlConditionalOperation == 'between' ||
                      one.sqlConditionalOperation == 'not between'
                    "
                  >
                    <input
                      v-model="one.before"
                      type="text"
                      @input="initInputWidth"
                      @change="initInputWidth"
                      class="part-form-input"
                    />
                    <span class="mglr-5">,</span>
                    <input
                      v-model="one.after"
                      type="text"
                      @input="initInputWidth"
                      @change="initInputWidth"
                      class="part-form-input"
                    />
                  </template>
                  <template v-else>
                    <input
                      v-model="one.value"
                      type="text"
                      @input="initInputWidth"
                      @change="initInputWidth"
                      class="part-form-input"
                    />
                  </template>
                  <!-- <select
                    v-model="one.andOr"
                    @change="initInputWidth"
                    class="part-form-input"
                  >
                    <option value="AND">AND</option>
                    <option value="OR">OR</option>
                  </select> -->
                  <div
                    @click="removeWhere(one)"
                    class="color-grey tm-link mgl-10"
                  >
                    删除
                  </div>
                </li>
              </template>
              <li class="pdl-5">
                <div @click="addWhere" class="color-green tm-link mgr-10">
                  添加条件
                </div>
              </li>
            </ul>
          </tm-layout>
          <!-- <tm-layout-bar right></tm-layout-bar> -->
          <tm-layout width="400px">
            <ul class="part-box scrollbar mg-0">
              <template v-for="(one, index) in searchForm.orderList">
                <li :key="index">
                  <input v-model="one.checked" type="checkbox" />
                  <select
                    v-model="one.name"
                    @change="initInputWidth"
                    class="part-form-input"
                  >
                    <option :value="null" text="请选择">请选择</option>
                    <template v-for="(one, index) in pointColumnList">
                      <option
                        :key="index"
                        :value="one.name"
                        :text="one.name + '&nbsp;'"
                      >
                        {{ one.name }}
                        <template v-if="tool.isNotEmpty(one.comment)">
                          （{{ one.comment }}）
                        </template>
                      </option>
                    </template>
                  </select>
                  <select
                    v-model="one.ascDesc"
                    @change="initInputWidth"
                    class="part-form-input"
                  >
                    <option value="DESC">DESC</option>
                    <option value="ASC">ASC</option>
                  </select>
                  <div
                    @click="removeOrder(one)"
                    class="color-grey tm-link mgl-10"
                  >
                    删除
                  </div>
                </li>
              </template>
              <li class="pdl-5">
                <div @click="addOrder" class="color-green tm-link mgr-10">
                  添加排序
                </div>
              </li>
            </ul>
          </tm-layout>
          <!-- <tm-layout-bar right></tm-layout-bar> -->
          <tm-layout>
            <ul class="part-box scrollbar mg-0">
              <li></li>
            </ul>
          </tm-layout>
        </tm-layout>
        <tm-layout height="30px">
          <div class="pdl-10">
            <div
              class="tm-btn tm-btn-sm bg-teal-8 ft-13"
              @click="toSearch"
              :class="{ 'tm-disabled': dataListLoading }"
            >
              搜索
            </div>
            <div class="tm-btn tm-btn-sm bg-green ft-13" @click="toInsert">
              新增
            </div>
          </div>
        </tm-layout>
        <tm-layout height="auto">
          <div style="height: 100%">
            <el-table
              :data="dataList"
              :border="true"
              height="100%"
              style="width: 100%"
              size="mini"
              v-loading="dataListLoading"
            >
              <!-- @row-dblclick="rowDblClick" -->
              <el-table-column label="_id" fixed>
                <template slot-scope="scope">
                  <span>{{ scope.row._id }}</span>
                </template>
              </el-table-column>
              <template v-for="(column, index) in columnList">
                <template v-if="column.checked">
                  <el-table-column
                    :key="index"
                    :prop="column.name"
                    :label="column.name"
                  >
                    <template slot-scope="scope">
                      <span class="">
                        {{ scope.row._source[column.name] }}
                      </span>
                    </template>
                  </el-table-column>
                </template>
              </template>
              <el-table-column width="180" label="操作" fixed="right">
                <template slot-scope="scope">
                  <div
                    class="tm-btn color-grey tm-btn-xs"
                    @click="toolboxWorker.showData(scope.row)"
                  >
                    查看
                  </div>
                  <div
                    class="tm-btn color-blue tm-btn-xs"
                    @click="toUpdate(scope.row)"
                  >
                    修改
                  </div>
                  <div
                    class="tm-btn color-grey tm-btn-xs"
                    @click="toCopy(scope.row)"
                  >
                    复制
                  </div>
                  <div
                    class="tm-btn color-orange tm-btn-xs"
                    @click="toDelete(scope.row)"
                    title="删除"
                  >
                    <i class="mdi mdi-delete-outline"></i>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </tm-layout>
        <tm-layout height="30px">
          <div class="ft-12 pdt-2 text-center">
            <el-pagination
              small
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              :current-page="searchForm.pageIndex"
              :page-sizes="[10, 50, 100, 200, 500]"
              :page-size="searchForm.pageSize"
              layout="total, sizes, prev, pager, next, jumper"
              :total="total"
              :disabled="total <= 0"
            >
            </el-pagination>
          </div>
        </tm-layout>
      </tm-layout>
    </template>
  </div>
</template>


<script>
var JSONbig = require("json-bigint");
var JSONbigString = JSONbig({ storeAsString: true });

export default {
  components: {},
  props: ["source", "indexName", "toolboxWorker"],
  data() {
    return {
      ready: false,
      searchForm: {
        indexName: this.indexName,
        whereList: [],
        orderList: [],
      },
      pageIndex: 1,
      pageSize: 10,
      total: 0,
      dataList: null,
      mapping: null,
      columnList: [],
      pointColumnList: [],
      dataListLoading: false,
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
    handleSizeChange(pageSize) {
      this.pageSize = pageSize;
      this.doSearch();
    },
    handleCurrentChange(pageIndex) {
      this.pageIndex = pageIndex;
      this.doSearch();
    },
    async initMapping() {
      let indexName = this.indexName;
      this.mapping = await this.toolboxWorker.getMapping(indexName);
      let columnList = [];
      let pointColumnList = [];
      if (
        this.mapping &&
        this.mapping.mappings &&
        this.mapping.mappings.properties
      ) {
        this.appendColumnList(
          columnList,
          pointColumnList,
          this.mapping.mappings.properties
        );
      }
      this.columnList = columnList;
      this.pointColumnList = pointColumnList;
    },
    appendColumnList(columnList, pointColumnList, properties, parentName) {
      if (properties == null) {
        return;
      }
      for (let key in properties) {
        let option = properties[key];
        let name = key;
        if (this.tool.isNotEmpty(parentName)) {
          name = parentName + "." + key;
        }
        let isStruct = option.properties != null;
        if (isStruct) {
          this.appendColumnList(
            columnList,
            pointColumnList,
            option.properties,
            name
          );
        }
        let column = {
          name: name,
          type: option.type,
          checked: true,
          isStruct: isStruct,
        };
        pointColumnList.push(column);
        if (this.tool.isEmpty(parentName)) {
          columnList.push(column);
        }
      }
    },
    inputValueChange() {},
    initInputWidth() {
      this.$nextTick(() => {
        if (this.initInputWidthIng) {
          return;
        }
        this.initInputWidthIng = true;
        let es = this.$el.getElementsByClassName("part-form-input");
        if (es) {
          Array.prototype.forEach.call(es, (one) => {
            this.tool.initInputWidth(one);
          });
        }
        this.initInputWidthIng = false;
      });
    },
    addOrder() {
      let order = {
        checked: true,
        name: null,
        ascDesc: "ASC",
      };
      let column = null;
      if (this.pointColumnList.length > 0) {
        this.pointColumnList.forEach((one) => {
          if (column != null) {
            return;
          }
          let find = false;
          this.searchForm.orderList.forEach((w) => {
            if (w.name == one.name) {
              find = true;
            }
          });
          if (find) {
            return;
          }
          column = one;
        });
      }
      if (column != null) {
        order.name = column.name;
      }

      this.searchForm.orderList.push(order);
      this.initInputWidth();
      return order;
    },
    removeOrder(order) {
      let orderList = this.searchForm.orderList;
      if (orderList.indexOf(order) >= 0) {
        orderList.splice(orderList.indexOf(order), 1);
      }
    },
    removeWhere(where) {
      let whereList = this.searchForm.whereList;
      if (whereList.indexOf(where) >= 0) {
        whereList.splice(whereList.indexOf(where), 1);
      }
    },
    addWhere() {
      let where = {
        checked: true,
        name: null,
        value: null,
        before: null,
        after: null,
        sqlConditionalOperation: "=",
        andOr: "AND",
      };
      let column = null;
      if (this.pointColumnList.length > 0) {
        this.pointColumnList.forEach((one) => {
          if (column != null) {
            return;
          }
          let find = false;
          this.searchForm.whereList.forEach((w) => {
            if (w.name == one.name) {
              find = true;
            }
          });
          if (find) {
            return;
          }
          column = one;
        });
      }
      if (column != null) {
        where.name = column.name;
      }

      this.searchForm.whereList.push(where);
      this.initInputWidth();
      return where;
    },
    async saveExtend() {
      let keyValueMap = {};
      keyValueMap.orderList = this.searchForm.orderList;
      keyValueMap.whereList = this.searchForm.whereList;
      keyValueMap.pageSize = this.pageSize;
      keyValueMap.pageIndex = this.pageIndex;
      await this.toolboxWorker.updateOpenTabExtend(this.tabId, keyValueMap);
    },
    toIndex() {},
    rowDblClick(row, column, event) {
      this.toolboxWorker.showData(row);
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
      this.toolboxWorker.showDataForm(data, this.mapping, async (m) => {
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
      let res = await this.toolboxWorker.work("insertData", param);
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
      this.toolboxWorker.showDataForm(param, this.mapping, async (m) => {
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
      this.toolboxWorker.showDataForm(param, this.mapping, async (m) => {
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
      let res = await this.toolboxWorker.work("updateData", param);
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
      let res = await this.toolboxWorker.work("deleteData", param);
      if (res.code == 0) {
        this.doSearch();
        return true;
      } else {
        return false;
      }
    },
    formatSourceJSON(_source) {
      let sourceJSON = {};
      try {
        sourceJSON = JSONbigString.parse(_source);
      } catch (error) {}
      return sourceJSON;
    },
    async doSearch() {
      this.dataListLoading = true;
      try {
        await this.initMapping();
        let param = {};
        Object.assign(param, this.searchForm);
        let whereList = [];
        let orderList = [];
        this.searchForm.whereList.forEach((one) => {
          if (one.checked) {
            whereList.push(one);
          }
        });
        this.searchForm.orderList.forEach((one) => {
          if (one.checked) {
            orderList.push(one);
          }
        });
        param.whereList = whereList;
        param.orderList = orderList;

        param.pageIndex = Number(this.pageIndex);
        param.pageSize = Number(this.pageSize);
        let res = await this.toolboxWorker.work("search", param);
        res.data = res.data || {};
        let result = res.data.result || {};
        let hits = result.hits || [];
        hits.forEach((one) => {
          one._source = this.formatSourceJSON(one._source);
        });
        this.dataList = hits;
        this.total = 0;
        if (result.total) {
          this.total = result.total.value;
        }
      } catch (error) {}
      this.dataListLoading = false;
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

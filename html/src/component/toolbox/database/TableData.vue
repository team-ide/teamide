<template>
  <div class="toolbox-database-table-data">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="120px" style="overflow: hidden">
          <tm-layout width="auto">
            <ul class="part-box app-scroll-bar mg-0" v-if="tableDetail != null">
              <template v-for="(one, index) in form.wheres">
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
                      <template v-for="(one, index) in columnList">
                        <option
                          :key="index"
                          :value="one.columnName"
                          :text="one.columnName + '&nbsp;'"
                        >
                          {{ one.columnName }}
                          <template v-if="tool.isNotEmpty(one.columnComment)">
                            （{{ one.columnComment }}）
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
                    <template
                      v-for="(one, index) in source.sqlConditionalOperations"
                    >
                      <option :key="index" :value="one.value" :text="one.text">
                        {{ one.text }}
                        <template v-if="tool.isNotEmpty(one.comment)">
                          （{{ one.comment }}）
                        </template>
                      </option>
                    </template>
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
                  <template v-else-if="one.sqlConditionalOperation == 'custom'">
                    <input
                      v-model="one.customSql"
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

                  <select
                    v-model="one.andOr"
                    @change="initInputWidth"
                    class="part-form-input"
                  >
                    <option value="AND">AND</option>
                    <option value="OR">OR</option>
                  </select>
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
          <tm-layout width="300px">
            <ul class="part-box app-scroll-bar mg-0">
              <template v-for="(one, index) in form.orders">
                <li :key="index">
                  <input v-model="one.checked" type="checkbox" />
                  <select
                    v-model="one.name"
                    @change="initInputWidth"
                    class="part-form-input"
                  >
                    <option :value="null" text="请选择">请选择</option>
                    <template v-for="(one, index) in columnList">
                      <option
                        :key="index"
                        :value="one.columnName"
                        :text="one.columnName + '&nbsp;'"
                      >
                        {{ one.columnName }}
                        <template v-if="tool.isNotEmpty(one.columnComment)">
                          （{{ one.columnComment }}）
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
        </tm-layout>
        <!-- <tm-layout-bar bottom></tm-layout-bar> -->
        <tm-layout height="20px">
          <div class="ft-12 pdl-10" v-if="tableDetail != null">
            <div class="color-grey tm-link mgr-10" @click="toSelectAll">
              全选
            </div>
            <div class="color-grey tm-link mgr-10" @click="toUnselectAll">
              取消全选
            </div>
            <div @click="toInsert({})" class="color-blue tm-link mgr-10">
              新增
            </div>
            <div
              class="color-red tm-link mgr-10"
              @click="toDeleteSelect"
              :class="{ 'tm-disabled': selects.length == 0 }"
            >
              删除选中({{ selects.length }})
            </div>
            <div
              class="color-green tm-link mgr-10"
              @click="toSaveSelect"
              :class="{
                'tm-disabled': updates.length == 0 && inserts.length == 0,
              }"
            >
              保存修改(编辑:{{ updates.length }}/新增:{{ inserts.length }})
            </div>
            <div
              class="color-grey tm-link mgr-10"
              @click="showTableDataListExecSql"
              :class="{
                'tm-disabled':
                  updates.length == 0 &&
                  inserts.length == 0 &&
                  selects.length == 0,
              }"
            >
              查看SQL(编辑:{{ updates.length }}/新增:{{
                inserts.length
              }}/删除:{{ selects.length }})
            </div>
            <div
              @click="showTableDataListSql"
              class="color-grey tm-link mgr-10"
              :class="{ 'tm-disabled': selects.length == 0 }"
            >
              导出选中(SQL)({{ selects.length }})
            </div>
            <div @click="doSearch" class="color-green tm-link mgr-10">查询</div>
          </div>
        </tm-layout>
        <tm-layout height="auto" v-loading="dataList_loading">
          <DataTable
            ref="DataTable"
            :source="source"
            :selects="selects"
            :inserts="inserts"
            :updates="updates"
            @keyup="tableKeyUp"
          >
          </DataTable>
        </tm-layout>
        <!-- <tm-layout-bar top></tm-layout-bar> -->
        <tm-layout height="50px" class="app-scroll-bar">
          <div class="ft-12 pdlr-10" v-if="tableDetail != null && sql != null">
            <div style="line-height: 20px">
              <span style="word-break: break-all; user-select: text">{{
                executeSql
              }}</span>
            </div>
          </div>
        </tm-layout>
        <tm-layout height="30px">
          <div
            class="ft-12 pdt-2 text-center"
            v-if="tableDetail != null && sql != null"
          >
            <el-pagination
              small
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              :current-page="pageNo"
              :page-sizes="[10, 50, 100, 200, 500]"
              :page-size="pageSize"
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
export default {
  components: {},
  props: [
    "source",
    "toolboxWorker",
    "actived",
    "ownerName",
    "tableName",
    "extend",
    "tabId",
  ],
  data() {
    return {
      ready: false,
      tableDetail: null,
      dataList_loading: false,
      dataList: [],
      sql: null,
      args: null,
      executeSql: null,
      inserts: [],
      updates: [],
      selects: [],
      keys: [],
      columnList: [],
      pageSize: 50,
      pageNo: 1,
      total: 0,
      form: {
        wheres: [],
        orders: [],
        columnList: [],
      },
    };
  },
  computed: {},
  watch: {},
  methods: {
    onFocus() {
      if (this.inited) {
        return;
      }
      this.$nextTick(async () => {
        this.init();
      });
    },
    async init() {
      this.inited = true;
      await this.initTable();
      if (this.tableDetail == null) {
        return;
      }
      this.result = null;
      this.keys = [];
      this.form.wheres = [];
      this.form.orders = [];
      this.form.columnList = [];

      if (this.tableDetail && this.tableDetail.columnList) {
        this.tableDetail.columnList.forEach((one, index) => {
          if (one.primaryKey) {
            this.keys.push(one.columnName);
          }
          this.columnList.push(one);
          let column = Object.assign({}, one);
          column.checked = true;
          this.form.columnList.push(column);
        });
      }

      if (this.extend.wheres) {
        this.extend.wheres.forEach((one) => {
          this.form.wheres.push(one);
        });
      }
      if (this.extend.orders) {
        this.extend.orders.forEach((one) => {
          this.form.orders.push(one);
        });
      }
      if (this.extend.pageSize) {
        this.pageSize = this.extend.pageSize;
      }
      if (this.extend.pageNo) {
        this.pageNo = this.extend.pageNo;
      }
      this.ready = true;
      this.$nextTick(() => {
        if (this.tableDetail && this.tableDetail.columnList) {
          if (this.form.wheres.length == 0) {
            this.tableDetail.columnList.forEach((one, index) => {
              if (index < 3) {
                let where = this.addWhere();
                where.checked = false;
              }
            });
          }
        }
      });
      this.initInputWidth();
      this.isInitSearch = true;
      await this.doSearch();
      this.isInitSearch = false;
    },
    async initTable() {
      this.tableDetail = await this.toolboxWorker.getTableDetail(
        this.ownerName,
        this.tableName
      );
    },
    handleSizeChange(pageSize) {
      this.pageSize = pageSize;
      this.doSearch();
    },
    handleCurrentChange(pageNo) {
      this.pageNo = pageNo;
      this.doSearch();
    },
    initInputWidth() {
      this.$nextTick(() => {
        if (this.initInputWidthIng) {
          return;
        }
        this.initInputWidthIng = true;
        let es = this.$el.getElementsByClassName("part-form-input");
        if (es) {
          Array.prototype.forEach.call(es, (one) => {
            this.initWidth(one);
          });
        }
        this.initInputWidthIng = false;
      });
    },
    initWidth(event) {
      let target = event;
      if (event.target) {
        target = event.target;
      }
      let value = target.value;
      this.tool.initInputWidth(event);
    },
    addOrder() {
      let order = {
        checked: true,
        name: null,
        ascDesc: "ASC",
      };
      let column = null;
      if (this.tableDetail && this.tableDetail.columnList) {
        if (this.tableDetail.columnList.length > 0) {
          this.tableDetail.columnList.forEach((one) => {
            if (column != null) {
              return;
            }
            let find = false;
            this.form.orders.forEach((w) => {
              if (w.name == one.columnName) {
                find = true;
              }
            });
            if (find) {
              return;
            }
            column = one;
          });
        }
      }
      if (column != null) {
        order.name = column.columnName;
      }

      this.form.orders.push(order);
      this.initInputWidth();
      return order;
    },
    removeOrder(order) {
      let orders = this.form.orders;
      if (orders.indexOf(order) >= 0) {
        orders.splice(orders.indexOf(order), 1);
      }
    },
    removeWhere(where) {
      let wheres = this.form.wheres;
      if (wheres.indexOf(where) >= 0) {
        wheres.splice(wheres.indexOf(where), 1);
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
      if (this.tableDetail && this.tableDetail.columnList) {
        if (this.tableDetail.columnList.length > 0) {
          this.tableDetail.columnList.forEach((one) => {
            if (column != null) {
              return;
            }
            let find = false;
            this.form.wheres.forEach((w) => {
              if (w.name == one.columnName) {
                find = true;
              }
            });
            if (find) {
              return;
            }
            column = one;
          });
        }
      }
      if (column != null) {
        where.name = column.columnName;
      }

      this.form.wheres.push(where);
      this.initInputWidth();
      return where;
    },
    async saveExtend() {
      let keyValueMap = {};
      keyValueMap.orders = this.form.orders;
      keyValueMap.wheres = this.form.wheres;
      keyValueMap.pageSize = this.pageSize;
      keyValueMap.pageNo = this.pageNo;
      await this.toolboxWorker.updateOpenTabExtend(this.tabId, keyValueMap);
    },
    async doSearch() {
      if (!this.isInitSearch) {
        await this.saveExtend();
      }
      let wheres = [];
      let orders = [];
      let columnList = [];
      this.form.wheres.forEach((where) => {
        if (!where.checked) {
          return;
        }
        if (this.tool.isEmpty(where.name)) {
          return;
        }
        wheres.push(where);
      });
      this.form.orders.forEach((order) => {
        if (!order.checked) {
          return;
        }
        orders.push(order);
      });
      this.form.columnList.forEach((column) => {
        if (!column.checked) {
          return;
        }
        columnList.push(column);
      });
      let data = {};
      data.ownerName = this.ownerName;
      data.tableName = this.tableName;
      data.wheres = wheres;
      data.orders = orders;
      data.columnList = columnList;
      data.pageNo = this.pageNo;
      data.pageSize = this.pageSize;
      this.dataList_loading = true;

      this.dataList = [];
      this.total = 0;
      this.sql = null;
      this.args = null;
      this.executeSql = null;
      this.updates = [];
      this.inserts = [];
      this.selects = [];

      let param = this.toolboxWorker.getWorkParam(data);
      let res = await this.server.database.tableData(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      res.data = res.data || {};

      let dataList = res.data.dataList || [];
      this.dataListCache = [];
      this.dataListCacheForIndex = [];
      dataList.forEach((data) => {
        this.tableDetail.columnList.forEach((column) => {
          if (data[column.columnName] != null) {
            data[column.columnName] = this.toolboxWorker.formatDateColumn(
              column,
              data[column.columnName]
            );
          } else {
            data[column.columnName] = null;
          }
        });
        this.dataListCache.push(Object.assign({}, data));
        this.dataListCacheForIndex.push(data);
      });

      let options = {};
      options.dataList = dataList;

      options.columnList = [];
      columnList.forEach((one) => {
        let column = Object.assign({}, one);
        column.name = column.columnName;
        options.columnList.push(column);
      });
      options.getRowMenus = this.getRowMenus;
      options.inputValueChange = this.inputValueChange;
      options.checkboxChange = this.checkboxChange;
      this.$nextTick(() => {
        this.$refs.DataTable.build(options);
      });
      this.dataList = dataList;

      this.sql = res.data.sql;
      this.total = Number(res.data.total || 0);
      this.args = res.data.args || [];
      this.dataList_loading = false;
      let executeSql = this.sql || "";
      this.executeSql = executeSql;
    },
    toSelectAll() {
      if (this.dataList.length == this.selects.length) {
        return;
      } else {
        this.dataList.forEach((one) => {
          if (this.selects.indexOf(one) < 0) {
            this.selects.push(one);
          }
        });
      }
    },
    checkboxChange(data) {
      let index = this.selects.indexOf(data);
      if (index >= 0) {
        this.selects.splice(index, 1);
      } else {
        this.selects.push(data);
      }
    },
    inputValueChange(data, column) {
      if (this.inserts.indexOf(data) >= 0) {
        return;
      }
      let dataUpdated = false;
      let dataCache = this.getCacheData(data);
      this.columnList.forEach((column) => {
        let name = column.columnName;
        if (data[name] !== dataCache[name]) {
          // if (typeof dataCache[name] == "number") {
          //   if (this.tool.isNotEmpty(data[name])) {
          //     data[name] = Number(data[name]);
          //   }
          // }
          if (data[name] != dataCache[name]) {
            dataUpdated = true;
          }
        }
      });
      let dataUpdateIndex = this.updates.indexOf(data);
      if (dataUpdated) {
        if (dataUpdateIndex < 0) {
          this.updates.push(data);
        }
      } else {
        if (dataUpdateIndex >= 0) {
          this.updates.splice(dataUpdateIndex, 1);
        }
      }
    },
    toUnselectAll() {
      this.selects.splice(0, this.selects.length);
    },
    getRowMenus(row, column, event) {
      let menus = [];
      if (event.target.tagName == "INPUT") {
        let input = this.tool.jQuery(event.target);
        let name = input.attr("name");
        let startIndex = event.target.selectionStart || 0;
        let endIndex = event.target.selectionEnd || 0;
        if (endIndex > startIndex) {
          let selectText = event.target.value.substring(startIndex, endIndex);
          menus.push({
            text: "复制选中文案",
            onClick: async () => {
              let res = await this.tool.clipboardWrite(selectText);
              if (res.success) {
                this.tool.success("复制成功");
              } else {
                this.tool.warn("复制失败，请允许访问剪贴板！");
              }
            },
          });
        }
        let menu_w_1 = {
          text: "追加粘贴",
          disabled: true,
        };
        menus.push(menu_w_1);
        let menu_w_2 = {
          text: "覆盖粘贴",
          disabled: true,
        };
        menus.push(menu_w_2);
        setTimeout(async () => {
          if (navigator.clipboard == null || navigator.permissions == null) {
            return;
          }
          let readResult = await this.tool.readClipboardText();
          if (!readResult.success) {
            menu_w_1.text += "(请允许访问剪贴板)";
            menu_w_2.text += "(请允许访问剪贴板)";
          }
          menu_w_1.disabled = this.tool.isEmpty(readResult.text);
          menu_w_1.onClick = () => {
            input.val(input.val() + readResult.text);
            input.change();
          };
          menu_w_2.disabled = this.tool.isEmpty(readResult.text);
          menu_w_2.onClick = () => {
            input.val(readResult.text);
            input.change();
          };
        }, 200);

        menus.push({
          text: "设置为空字符串",
          onClick: () => {
            row[name] = "";
            this.inputValueChange(row);
          },
        });
        menus.push({
          text: "设置为NULL",
          onClick: () => {
            row[name] = null;
            this.inputValueChange(row);
          },
        });
        menus.push({
          text: "查看列数据",
          onClick: () => {
            this.tool.showJSONData(row[name]);
          },
        });
      }
      menus.push({
        text: "查看行数据",
        onClick: () => {
          this.tool.showJSONData(row);
        },
      });
      let dataCache = this.getCacheData(row);
      let insertData = Object.assign({}, row);
      let updateData = {};
      let updateWhereData = {};
      let deleteData = {};
      if (this.keys.length > 0) {
        for (let key in row) {
          if (this.keys.indexOf(key) < 0) {
            updateData[key] = row[key];
          }
        }
        this.keys.forEach((key) => {
          if (dataCache != null) {
            deleteData[key] = dataCache[key];
            updateWhereData[key] = dataCache[key];
          } else {
            deleteData[key] = row[key];
            updateWhereData[key] = row[key];
          }
        });
      } else {
        updateData = Object.assign({}, row);
        if (dataCache != null) {
          updateWhereData = Object.assign({}, dataCache);
          deleteData = Object.assign({}, dataCache);
        } else {
          updateWhereData = Object.assign({}, row);
          deleteData = Object.assign({}, row);
        }
      }

      let insertList = [insertData];
      let updateList = [updateData];
      let updateWhereList = [updateWhereData];
      let deleteList = [deleteData];

      menus.push({
        text: "查看SQL",
        onClick: () => {
          this.toolboxWorker.showTableDataListSql(
            this.ownerName,
            this.tableDetail,
            [row]
          );
        },
      });
      if (this.updates.indexOf(row) >= 0) {
        menus.push({
          text: "保存该记录",
          onClick: () => {
            this.doSave(null, updateList, updateWhereList);
          },
        });
      }
      if (this.inserts.indexOf(row) >= 0) {
        menus.push({
          text: "保存该记录",
          onClick: () => {
            this.doSave(insertList, null, null);
          },
        });
      }
      menus.push({
        text: "追加记录",
        onClick: () => {
          this.toInsert({}, row);
        },
      });
      menus.push({
        text: "复制追加记录",
        onClick: () => {
          this.toInsert(Object.assign({}, row), row);
        },
      });
      menus.push({
        text: "删除记录",
        onClick: () => {
          if (this.inserts.indexOf(row) >= 0) {
            this.removeInsert(row);
          } else {
            if (dataCache != null) {
              let msg = "删除该记录将无法恢复，确认删除？";
              this.tool
                .confirm(msg)
                .then(async () => {
                  await this.doDelete([dataCache]);
                  this.removeDatas([row]);
                })
                .catch((e) => {});
            }
          }
        },
      });

      return menus;
    },
    tableKeyUp(event) {
      event = event || window.event;
      if (this.tool.keyIsCtrlS(event)) {
        this.tool.stopEvent(event);
        this.toSaveSelect();
      }
    },
    getCacheData(data) {
      let index = this.dataListCacheForIndex.indexOf(data);
      if (index < 0) {
        return null;
      }
      let dataCache = this.dataListCache[index];
      return dataCache;
    },
    toDeleteSelect() {
      let deleteList = this.getDeleteList();
      if (deleteList.length == 0) {
        this.tool.warn("暂无需要删除的数据");
        return;
      }
      let msg = "此次将删除（" + deleteList.length + "）条数据";
      msg += "，确认删除？";
      this.tool
        .confirm(msg)
        .then(async () => {
          await this.doDelete(deleteList);
        })
        .catch((e) => {});
    },
    removeDatas(datas) {
      datas = datas || [];
      let list = [];
      datas.forEach((data) => {
        list.push(data);
      });
      list.forEach((data) => {
        if (this.inserts.indexOf(data) >= 0) {
          this.inserts.splice(this.inserts.indexOf(data), 1);
        }
        if (this.updates.indexOf(data) >= 0) {
          this.updates.splice(this.updates.indexOf(data), 1);
        }
        let dataCache = this.getCacheData(data);
        if (dataCache != null) {
          if (this.dataListCache.indexOf(dataCache) >= 0) {
            this.dataListCache.splice(this.dataListCache.indexOf(dataCache), 1);
          }
        }
        if (this.dataListCacheForIndex.indexOf(data) >= 0) {
          this.dataListCacheForIndex.splice(
            this.dataListCacheForIndex.indexOf(data),
            1
          );
        }
        if (this.dataList.indexOf(data) >= 0) {
          this.dataList.splice(this.dataList.indexOf(data), 1);
        }
        if (this.selects.indexOf(data) >= 0) {
          this.selects.splice(this.selects.indexOf(data), 1);
        }
      });

      this.$refs.DataTable.refreshData();
    },
    getDeleteList() {
      let deleteList = this.getDeleteListByDatas(this.selects);
      return deleteList;
    },
    getDeleteListByDatas(datas) {
      datas = datas || [];
      let deleteList = [];
      datas.forEach((data) => {
        let dataCache = this.getCacheData(data);
        if (dataCache == null) {
          return;
        }
        if (this.keys.length > 0) {
          let deleteData = {};
          this.keys.forEach((key) => {
            deleteData[key] = dataCache[key];
          });
          deleteList.push(deleteData);
        } else {
          deleteList.push(dataCache);
        }
      });
      return deleteList;
    },
    getInsertList() {
      let insertList = this.getInsertListByDatas(this.inserts);
      return insertList;
    },
    getInsertListByDatas(datas) {
      datas = datas || [];
      let insertList = [];
      datas.forEach((data) => {
        insertList.push(data);
      });
      return insertList;
    },
    getUpdateList() {
      let updateList = this.getUpdateListByDatas(this.updates);
      return updateList;
    },
    getUpdateListByDatas(datas) {
      datas = datas || [];
      let updateList = [];
      datas.forEach((data) => {
        let dataCache = this.getCacheData(data);
        if (dataCache == null) {
          return;
        }
        let updateData = {};
        for (let key in dataCache) {
          if (data[key] != dataCache[key]) {
            updateData[key] = data[key];
          }
        }
        updateList.push(updateData);
      });
      return updateList;
    },
    getUpdateWhereList() {
      let updateWhereList = this.getUpdateWhereListByDatas(this.updates);
      return updateWhereList;
    },
    getUpdateWhereListByDatas(datas) {
      datas = datas || [];
      let updateWhereList = [];
      datas.forEach((data) => {
        let dataCache = this.getCacheData(data);
        if (dataCache == null) {
          return;
        }
        if (this.keys.length > 0) {
          let whereData = {};
          this.keys.forEach((key) => {
            whereData[key] = dataCache[key];
          });
          updateWhereList.push(whereData);
        } else {
          updateWhereList.push(dataCache);
        }
      });
      return updateWhereList;
    },
    toSaveSelect() {
      let insertList = this.getInsertList();
      let updateList = this.getUpdateList();
      let updateWhereList = this.getUpdateWhereList();

      if (insertList.length == 0 && updateList.length == 0) {
        this.tool.warn("暂无需要保存的数据");
        return;
      }
      let msg = "此次将，";
      if (insertList.length > 0) {
        msg += "新增（" + insertList.length + "）条记录，";
      }
      if (updateList.length > 0) {
        msg += "更新（" + updateList.length + "）条记录，";
      }
      msg += "确认保存？";
      this.tool
        .confirm(msg)
        .then(async () => {
          await this.doSave(insertList, updateList, updateWhereList);
          this.savedInserts(this.inserts);
          this.savedUpdates(this.updates);
        })
        .catch((e) => {});
    },
    async doSave(insertList, updateList, updateWhereList) {
      let data = {};
      data.ownerName = this.ownerName;
      data.tableName = this.tableName;
      data.columnList = this.tableDetail.columnList;
      data.insertList = insertList;
      data.updateList = updateList;
      data.updateWhereList = updateWhereList;

      let param = this.toolboxWorker.getWorkParam(data);
      let res = await this.server.database.dataListExec(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
        return;
      }
      let task = res.data || {};
      let info = "保存成功，";
      info +=
        "成功记录数（" +
        task.saveSuccess +
        "）条，耗时（" +
        task.useTime +
        "）毫秒！";
      this.tool.success(info);
      this.doSearch();
    },
    async doDelete(deleteList) {
      let data = {};
      data.ownerName = this.ownerName;
      data.tableName = this.tableName;
      data.columnList = this.tableDetail.columnList;
      data.deleteList = deleteList;

      let param = this.toolboxWorker.getWorkParam(data);
      let res = await this.server.database.dataListExec(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
        return;
      }
      let task = res.data || {};
      let info = "删除成功，";
      info +=
        "成功记录数（" +
        task.saveSuccess +
        "）条，耗时（" +
        task.useTime +
        "）毫秒！";
      this.tool.success(info);
      this.doSearch();
    },
    removeInsert(data) {
      if (this.inserts.indexOf(data) >= 0) {
        this.inserts.splice(this.inserts.indexOf(data), 1);
      }
      if (this.dataList.indexOf(data) >= 0) {
        this.dataList.splice(this.dataList.indexOf(data), 1);
      }
      if (this.selects.indexOf(data) >= 0) {
        this.selects.splice(this.selects.indexOf(data), 1);
      }

      this.$refs.DataTable.refreshData();
    },
    toInsert(data, row) {
      data = data || {};
      this.tableDetail.columnList.forEach((column) => {
        if (data[column.columnName] == null) {
          data[column.columnName] = null;
        }
      });
      this.inserts.push(data);
      if (row && this.dataList.indexOf(row) >= 0) {
        this.dataList.splice(this.dataList.indexOf(row) + 1, 0, data);
      } else {
        this.dataList.push(data);
      }

      this.$refs.DataTable.refreshData();
    },
    showTableDataListSql() {
      this.toolboxWorker.showTableDataListSql(
        this.ownerName,
        this.tableDetail,
        this.selects
      );
    },

    showTableDataListExecSql() {
      let insertList = this.getInsertList();
      let updateList = this.getUpdateList();
      let updateWhereList = this.getUpdateWhereList();
      let deleteList = this.getDeleteList();

      this.toolboxWorker.showTableDataListExecSql(
        this.ownerName,
        this.tableDetail,
        {
          insertList,
          updateList,
          updateWhereList,
          deleteList,
        }
      );
    },
  },
  created() {},
  mounted() {
    if (this.actived) {
      this.init();
    }
  },
};
</script>

<style>
.toolbox-database-table-data {
  width: 100%;
  height: 100%;
  user-select: none;
}
</style>

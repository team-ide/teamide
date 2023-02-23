<template>
  <div @keyup="tableKeyUp" style="height: 100%">
            <el-table
              :data="dataList"
              :border="true"
              height="100%"
              style="width: 100%"
              size="mini"
              @row-contextmenu="rowContextmenu"
            >
              <el-table-column width="70" label="序号" fixed>
                <template slot-scope="scope">
                  <!-- :checked="selects.indexOf(scope.row) >= 0" -->
                  <input
                    type="checkbox"
                    class="mgl-5"
                    style="width: auto"
                    v-model="selects"
                    :value="scope.row"
                  />
                  <span class="mgl-5">{{ scope.$index + 1 }}</span>
                  <template v-if="updates.indexOf(scope.row) >= 0">
                    <i
                      class="mgl-5 mdi mdi-database-edit-outline"
                      style="vertical-align: 0px"
                    ></i>
                  </template>
                  <template v-if="inserts.indexOf(scope.row) >= 0">
                    <i
                      class="mgl-5 mdi mdi-database-plus-outline"
                      style="vertical-align: 0px"
                    ></i>
                  </template>
                </template>
              </el-table-column>
              <template v-for="(column, index) in form.columnList">
                <template v-if="column.checked">
                  <el-table-column
                    :key="index"
                    :prop="column.columnName"
                    :label="column.columnName"
                  >
                    <template slot-scope="scope">
                      <div class="">
                        <input
                          v-model="scope.row[column.columnName]"
                          :name="column.columnName"
                          @change="inputValueChange(scope.row, column, $event)"
                          @input="inputValueChange(scope.row, column, $event)"
                          :placeholder="
                            scope.row[column.columnName] == null ? 'NULL' : ''
                          "
                          type="text"
                        />
                      </div>
                    </template>
                  </el-table-column>
                </template>
              </template>
            </el-table>
          </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "item"],
  data() {
    return {
      dataList: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      let dataList = [];
      for (let i = 0; i < this.item.dataList.length; i++) {
        if (i < 50) {
          dataList.push(this.item.dataList[i]);
        } else {
          break;
        }
      }

      this.dataList = dataList;
      this.bind();
    },
    bind() {
      let that = this;
      const selectWrap = this.$el.querySelector(".el-table__body-wrapper");
      selectWrap.addEventListener("scroll", function () {
        const scrollDistance =
          this.scrollHeight - this.scrollTop - this.clientHeight;
        // 判断是否到底,scrollTop为已滚动到元素上边界的像素数，scrollHeight为元素完整的高度，clientHeight为内容的可见宽度
        if (scrollDistance <= 10) {
          that.loadmore();
        }
      });
    },
    rowContextmenu(row, column, event) {
      let menus = [];

      menus.push({
        text: "查看行数据",
        onClick: () => {
          this.tool.showJSONData(row);
        },
      });

      if (this.item && this.item.columnList) {
        let insertData = Object.assign({}, row);

        let ownerName = "ownerName";
        let insertList = [insertData];
        let tableDetail = {
          tableName: "tableName",
          columnList: [],
        };
        this.item.columnList.forEach((one) => {
          tableDetail.columnList.push({
            columnName: one.name,
          });
        });
        menus.push({
          text: "查看新增记录SQL",
          onClick: () => {
            this.toolboxWorker.showTableDataListExecSql(
              ownerName,
              tableDetail,
              {
                insertList,
              }
            );
          },
        });
      }

      if (menus.length > 0) {
        this.tool.showContextmenu(menus);
      }
    },
    loadmore() {
      if (this.loadmore_ing) {
        return;
      }

      this.loadmore_ing = true;
      let allDataList = this.item.dataList;
      let dataList = this.dataList;
      if (allDataList.length <= dataList.length) {
        return;
      }

      let start = dataList.length;
      for (let i = start; i < allDataList.length; i++) {
        if (i - start < 50) {
          dataList.push(allDataList[i]);
        } else {
          break;
        }
      }
      this.loadmore_ing = false;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
</style>

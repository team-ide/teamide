<template>
  <div class="toolbox-database-sql-data-list">
    <DataTable ref="DataTable" :source="source"></DataTable>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "item"],
  data() {
    return {
      updates: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      let options = {};
      options.showSize = 50;
      options.dataList = this.item.dataList;
      options.columnList = this.item.columnList;
      options.getRowMenus = this.getRowMenus;
      this.$nextTick(() => {
        this.$refs.DataTable.build(options);
      });
    },
    getRowMenus(row, column, event) {
      let menus = [];

      if (column) {
        menus.push({
          text: "查看列数据",
          onClick: () => {
            this.tool.showJSONData(row[column.name]);
          },
        });
      }
      menus.push({
        text: "查看行数据",
        onClick: () => {
          this.tool.showJSONData(row);
        },
      });

      if (this.item && this.item.columnList) {
        let ownerName = "ownerName";
        let dataList = [Object.assign({}, row)];
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
          text: "查看SQL",
          onClick: () => {
            this.toolboxWorker.showTableDataListSql(
              ownerName,
              tableDetail,
              dataList
            );
          },
        });
      }

      return menus;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-sql-data-list {
  width: 100%;
  height: 100%;
}
</style>

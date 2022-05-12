<template>
  <div class="toolbox-database-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout width="300px" class="">
          <ToolboxDatabaseDatabase
            ref="ToolboxDatabaseDatabase"
            :source="source"
            :toolbox="toolbox"
            :toolboxType="toolboxType"
            :wrap="wrap"
          >
          </ToolboxDatabaseDatabase>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <ToolboxDatabaseTabs
            :source="source"
            :toolbox="toolbox"
            :toolboxType="toolboxType"
            :wrap="wrap"
          >
          </ToolboxDatabaseTabs>
        </tm-layout>
      </tm-layout>
      <ShowDatabaseCreate :source="source" :wrap="wrap"> </ShowDatabaseCreate>
      <ShowTableCreate :source="source" :wrap="wrap"> </ShowTableCreate>
      <ShowExportSql :source="source" :wrap="wrap"> </ShowExportSql>
      <ShowSaveSql :source="source" :wrap="wrap"> </ShowSaveSql>
      <ShowImportDataForStrategy :source="source" :wrap="wrap">
      </ShowImportDataForStrategy>
    </template>
  </div>
</template>


<script>
import ShowDatabaseCreate from "./ShowDatabaseCreate";
import ShowTableCreate from "./ShowTableCreate";
import ShowExportSql from "./ShowExportSql";
import ShowSaveSql from "./ShowSaveSql";
import ShowImportDataForStrategy from "./ShowImportDataForStrategy";

export default {
  components: {
    ShowDatabaseCreate,
    ShowTableCreate,
    ShowExportSql,
    ShowSaveSql,
    ShowImportDataForStrategy,
  },
  props: ["source", "toolboxType", "toolbox", "option", "wrap"],
  data() {
    return {
      ready: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.wrap.columnIsNumber = this.columnIsNumber;
      this.wrap.columnIsDate = this.columnIsDate;
      this.wrap.formatDateColumn = this.formatDateColumn;
      this.ready = true;
    },
    columnIsNumber(column) {
      let type = ("" + column.type).toLowerCase();
      if (
        type == "int" ||
        type == "bigint" ||
        type == "bit" ||
        type == "number" ||
        type == "tinyint" ||
        type == "smallint" ||
        type == "integer" ||
        type == "float" ||
        type == "double" ||
        type == "dec" ||
        type == "decimal"
      ) {
        return true;
      }
      return false;
    },
    columnIsDate(column) {
      let type = ("" + column.type).toLowerCase();
      if (
        type == "year" ||
        type == "date" ||
        type == "time" ||
        type == "datetime" ||
        type == "timestamp"
      ) {
        return true;
      }
      return false;
    },
    formatDateColumn(column, value) {
      if (value == null) {
        return;
      }
      let type = ("" + column.type).toLowerCase();

      try {
        if (type == "year") {
          value = this.tool.formatDate(new Date(value), "yyyy");
        } else if (type == "date") {
          value = this.tool.formatDate(new Date(value), "yyyy-MM-dd");
        } else if (type == "time") {
          value = this.tool.formatDate(new Date(value), "hh:mm:ss");
        } else if (type == "datetime" || type == "timestamp") {
          value = this.tool.formatDate(new Date(value), "yyyy-MM-dd hh:mm:ss");
        }
      } catch (e) {
        this.tool.error(e);
      }
      return value;
    },
    refresh() {
      this.$refs.ToolboxDatabaseDatabase.refresh();
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-editor {
  width: 100%;
  height: 100%;
}
</style>

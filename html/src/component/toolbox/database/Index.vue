<template>
  <div class="toolbox-database-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout width="300px" class="">
          <Database
            ref="ToolboxDatabaseDatabase"
            :source="source"
            :toolboxWorker="toolboxWorker"
            :extend="extend"
            :databasesChange="databasesChange"
          >
          </Database>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <Tabs
            :source="source"
            :toolboxWorker="toolboxWorker"
            :databases="databases"
          >
          </Tabs>
        </tm-layout>
      </tm-layout>
      <ShowExportSql :source="source" :toolboxWorker="toolboxWorker">
      </ShowExportSql>
      <ShowSaveSql :source="source" :toolboxWorker="toolboxWorker">
      </ShowSaveSql>
      <CreateDatabase :source="source" :toolboxWorker="toolboxWorker">
      </CreateDatabase>
      <Table :source="source" :toolboxWorker="toolboxWorker"> </Table>
    </template>
    <ShowInfo :source="source" :toolboxWorker="toolboxWorker"> </ShowInfo>
  </div>
</template>


<script>
import Database from "./Database";
import Tabs from "./Tabs";
import ShowExportSql from "./ShowExportSql";
import ShowSaveSql from "./ShowSaveSql";
import CreateDatabase from "./CreateDatabase";
import ShowInfo from "./ShowInfo";

export default {
  components: {
    Database,
    Tabs,
    ShowExportSql,
    ShowSaveSql,
    CreateDatabase,
    ShowInfo,
  },
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      ready: false,
      databases: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.toolboxWorker.columnIsNumber = this.columnIsNumber;
      this.toolboxWorker.columnIsDate = this.columnIsDate;
      this.toolboxWorker.formatDateColumn = this.formatDateColumn;
      this.ready = true;
    },
    databasesChange(databases) {
      this.databases = databases;
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

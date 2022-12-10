<template>
  <div class="toolbox-database-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout width="300px" class="">
          <Owner
            ref="Owner"
            :source="source"
            :toolboxWorker="toolboxWorker"
            :extend="extend"
            :ownersChange="ownersChange"
          >
          </Owner>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <Tabs
            :source="source"
            :toolboxWorker="toolboxWorker"
            :owners="owners"
            :columnTypeInfoList="columnTypeInfoList"
            :indexTypeInfoList="indexTypeInfoList"
          >
          </Tabs>
        </tm-layout>
      </tm-layout>
      <TableDataListExecSql :source="source" :toolboxWorker="toolboxWorker">
      </TableDataListExecSql>
      <TableDataListSql :source="source" :toolboxWorker="toolboxWorker">
      </TableDataListSql>
      <OwnerCreate :source="source" :toolboxWorker="toolboxWorker">
      </OwnerCreate>
      <Table :source="source" :toolboxWorker="toolboxWorker"> </Table>
    </template>
    <ShowInfo :source="source" :toolboxWorker="toolboxWorker"> </ShowInfo>
  </div>
</template>


<script>
import Owner from "./Owner";
import Tabs from "./Tabs";
import TableDataListExecSql from "./TableDataListExecSql";
import TableDataListSql from "./TableDataListSql";
import OwnerCreate from "./OwnerCreate";
import ShowInfo from "./ShowInfo";

export default {
  components: {
    Owner,
    Tabs,
    TableDataListExecSql,
    TableDataListSql,
    OwnerCreate,
    ShowInfo,
  },
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      ready: false,
      owners: [],
      columnTypeInfoList: [],
      indexTypeInfoList: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      this.toolboxWorker.columnIsNumber = this.columnIsNumber;
      this.toolboxWorker.columnIsDate = this.columnIsDate;
      this.toolboxWorker.formatDateColumn = this.formatDateColumn;
      this.toolboxWorker.formatParam = this.formatParam;
      await this.loadData();
      this.ready = true;
    },
    async loadData() {
      let param = this.toolboxWorker.getWorkParam({});
      let res = await this.server.database.data(param);
      let data = res.data || {};
      this.columnTypeInfoList = data.columnTypeInfoList || [];
      this.indexTypeInfoList = data.indexTypeInfoList || [];
    },
    formatParam(param) {
      param = param || {};
      if (this.tool.isEmpty(param.ownerNamePackChar)) {
        delete param.ownerNamePack;
        delete param.ownerNamePackChar;
      } else {
        if (param.ownerNamePackChar == "-") {
          param.ownerNamePack = false;
        } else {
          param.ownerNamePack = true;
        }
      }
      if (this.tool.isEmpty(param.tableNamePackChar)) {
        delete param.tableNamePack;
        delete param.tableNamePackChar;
      } else {
        if (param.tableNamePackChar == "-") {
          param.tableNamePack = false;
        } else {
          param.tableNamePack = true;
        }
      }
      if (this.tool.isEmpty(param.columnNamePackChar)) {
        delete param.columnNamePack;
        delete param.columnNamePackChar;
      } else {
        if (param.columnNamePackChar == "-") {
          param.columnNamePack = false;
        } else {
          param.columnNamePack = true;
        }
      }
      if (this.tool.isEmpty(param.sqlValuePackChar)) {
        delete param.sqlValuePackChar;
      }
    },
    ownersChange(owners) {
      this.owners = owners;
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
        this.tool.error(e.toString());
      }
      return value;
    },
    refresh() {
      this.$refs.Owner.refresh();
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    let param = this.toolboxWorker.getWorkParam({});
    this.server.database.close(param);
  },
};
</script>

<style>
.toolbox-database-editor {
  width: 100%;
  height: 100%;
}
</style>

<template>
  <div class="toolbox-database-tabs">
    <div class="default-tabs-container">
      <WorkspaceTabs :source="source" :itemsWorker="toolboxWorker.itemsWorker">
      </WorkspaceTabs>
    </div>
    <div class="default-spans-container">
      <WorkspaceSpans :source="source" :itemsWorker="toolboxWorker.itemsWorker">
        <template v-slot:span="{ item }">
          <template v-if="item.extend.type == 'data'">
            <TableData
              :source="source"
              :toolboxWorker="toolboxWorker"
              :ownerName="item.extend.ownerName"
              :tableName="item.extend.tableName"
              :extend="item.extend"
              :tabId="item.tabId"
              :actived="
                toolboxWorker.itemsWorker.activeItem &&
                item.key == toolboxWorker.itemsWorker.activeItem.key
              "
            >
            </TableData>
          </template>
          <template v-else-if="item.extend.type == 'sql'">
            <Sql
              :source="source"
              :toolboxWorker="toolboxWorker"
              :extend="item.extend"
              :owners="owners"
              :tabId="item.tabId"
              :actived="
                toolboxWorker.itemsWorker.activeItem &&
                item.key == toolboxWorker.itemsWorker.activeItem.key
              "
            >
            </Sql>
          </template>
          <template v-else-if="item.extend.type == 'ddl'">
            <DDL
              :source="source"
              :toolboxWorker="toolboxWorker"
              :ownerName="item.extend.ownerName"
              :tableName="item.extend.tableName"
              :tabId="item.tabId"
              :actived="
                toolboxWorker.itemsWorker.activeItem &&
                item.key == toolboxWorker.itemsWorker.activeItem.key
              "
            >
            </DDL>
          </template>
          <template v-else-if="item.extend.type == 'model'">
            <Model
              :source="source"
              :toolboxWorker="toolboxWorker"
              :ownerName="item.extend.ownerName"
              :tableName="item.extend.tableName"
              :tabId="item.tabId"
              :actived="
                toolboxWorker.itemsWorker.activeItem &&
                item.key == toolboxWorker.itemsWorker.activeItem.key
              "
            >
            </Model>
          </template>
          <template v-else-if="item.extend.type == 'table'">
            <Table
              :source="source"
              :toolboxWorker="toolboxWorker"
              :ownerName="item.extend.ownerName"
              :tableName="item.extend.tableName"
              :columnTypeInfoList="columnTypeInfoList"
              :indexTypeInfoList="indexTypeInfoList"
              :extend="item.extend"
              :tabId="item.tabId"
              :actived="
                toolboxWorker.itemsWorker.activeItem &&
                item.key == toolboxWorker.itemsWorker.activeItem.key
              "
            >
            </Table>
          </template>
          <template v-else-if="item.extend.type == 'export'">
            <Export
              :source="source"
              :toolboxWorker="toolboxWorker"
              :ownerName="item.extend.ownerName"
              :tableName="item.extend.tableName"
              :columnTypeInfoList="columnTypeInfoList"
              :indexTypeInfoList="indexTypeInfoList"
              :extend="item.extend"
              :tabId="item.tabId"
              :actived="
                toolboxWorker.itemsWorker.activeItem &&
                item.key == toolboxWorker.itemsWorker.activeItem.key
              "
            >
            </Export>
          </template>
          <template v-else-if="item.extend.type == 'import'">
            <Import
              :source="source"
              :toolboxWorker="toolboxWorker"
              :ownerName="item.extend.ownerName"
              :tableName="item.extend.tableName"
              :columnTypeInfoList="columnTypeInfoList"
              :indexTypeInfoList="indexTypeInfoList"
              :extend="item.extend"
              :tabId="item.tabId"
              :actived="
                toolboxWorker.itemsWorker.activeItem &&
                item.key == toolboxWorker.itemsWorker.activeItem.key
              "
            >
            </Import>
          </template>
          <template v-else-if="item.extend.type == 'sync'">
            <Sync
              :source="source"
              :toolboxWorker="toolboxWorker"
              :ownerName="item.extend.ownerName"
              :tableName="item.extend.tableName"
              :columnTypeInfoList="columnTypeInfoList"
              :indexTypeInfoList="indexTypeInfoList"
              :extend="item.extend"
              :tabId="item.tabId"
              :actived="
                toolboxWorker.itemsWorker.activeItem &&
                item.key == toolboxWorker.itemsWorker.activeItem.key
              "
            >
            </Sync>
          </template>
        </template>
      </WorkspaceSpans>
    </div>
  </div>
</template>


<script>
import DDL from "./DDL";
import Model from "./Model";
import Sql from "./Sql";
import Table from "./Table";
import TableData from "./TableData";
import Export from "./Export";
import Import from "./Import";
import Sync from "./Sync";

export default {
  components: { DDL, Model, Sql, Table, TableData, Export, Import, Sync },
  props: [
    "source",
    "toolboxWorker",
    "owners",
    "columnTypeInfoList",
    "indexTypeInfoList",
  ],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    init() {},
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-tabs {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>

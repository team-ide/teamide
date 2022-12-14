<template>
  <div class="toolbox-kafka-topic-data">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="110px">
          <el-form
            class="pdt-10 pdlr-10"
            size="mini"
            :inline="true"
            @submit.native.prevent
          >
            <el-form-item label="Topic">
              <el-input v-model="pullForm.topic" />
            </el-form-item>
            <el-form-item label="GroupId">
              <el-input v-model="pullForm.groupId" />
            </el-form-item>
            <el-form-item label="KeyType">
              <el-select
                placeholder="请选择类型"
                v-model="pullForm.keyType"
                style="width: 120px"
              >
                <el-option label="String" value="string"> </el-option>
                <el-option label="Long（int64）" value="long"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="ValueType">
              <el-select
                placeholder="请选择类型"
                v-model="pullForm.valueType"
                style="width: 120px"
              >
                <el-option label="String" value="string"> </el-option>
                <el-option label="Long（int64）" value="long"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="">
              <div class="">
                <div
                  class="tm-btn tm-btn-sm bg-teal-8 ft-13"
                  @click="toPull"
                  :class="{ 'tm-disabled': dataListLoading }"
                >
                  拉取
                </div>
                <div class="tm-btn tm-btn-sm bg-green ft-13" @click="toPush">
                  新增
                </div>
              </div>
            </el-form-item>
          </el-form>
        </tm-layout>
        <tm-layout height="auto" class="">
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
              <el-table-column width="60" label="分区">
                <template slot-scope="scope">
                  <span>{{ scope.row.partition }}</span>
                </template>
              </el-table-column>
              <el-table-column width="80" label="偏移量">
                <template slot-scope="scope">
                  <span>{{ scope.row.offset }}</span>
                </template>
              </el-table-column>
              <template v-for="(column, index) in headerColumnList">
                <template v-if="column.checked">
                  <el-table-column
                    :key="index"
                    :prop="column.name"
                    :label="`header:${column.name}`"
                  >
                    <template slot-scope="scope">
                      <div class="">
                        {{ scope.row.headerValue[column.name] }}
                      </div>
                    </template>
                  </el-table-column>
                </template>
              </template>
              <el-table-column width="120" label="Key">
                <template slot-scope="scope">
                  <span>{{ scope.row.key }}</span>
                </template>
              </el-table-column>
              <el-table-column label="Value">
                <template slot-scope="scope">
                  <span>{{ scope.row.value }}</span>
                </template>
              </el-table-column>
              <template v-for="(column, index) in columnList">
                <template v-if="column.checked">
                  <el-table-column
                    :key="index"
                    :prop="column.name"
                    :label="`json:${column.name}`"
                  >
                    <template slot-scope="scope">
                      <div class="">
                        {{ scope.row.jsonValue[column.name] }}
                      </div>
                    </template>
                  </el-table-column>
                </template>
              </template>
              <el-table-column width="150" label="操作" fixed="right">
                <template slot-scope="scope">
                  <div
                    class="tm-btn color-grey tm-btn-xs"
                    @click="tool.showJSONData(scope.row)"
                  >
                    查看
                  </div>
                  <div
                    class="tm-btn color-blue tm-btn-xs"
                    @click="toCommit(scope.row)"
                  >
                    消费
                  </div>
                  <div
                    class="tm-btn color-orange tm-btn-xs"
                    @click="toDelete(scope.row)"
                  >
                    删除
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </tm-layout>
      </tm-layout>
    </template>
    <FormDialog
      ref="PushForm"
      :source="source"
      title="推送数据"
      :onSave="doPush"
    ></FormDialog>
  </div>
</template>


<script>
var JSONbig = require("json-bigint");
var JSONbigString = JSONbig({});

export default {
  components: {},
  props: ["source", "topic", "toolboxWorker"],
  data() {
    return {
      ready: false,
      pullForm: {
        groupId: "test-group",
        topic: this.topic,
        keyType: "string",
        valueType: "string",
        pullSize: 10,
        pullTimeout: 1000,
      },
      dataList: [],
      headerColumnList: [],
      columnList: [],
      dataListLoading: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.ready = true;
      this.toPull();
    },
    async toPull() {
      await this.doPull();
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
      this.tool.showJSONData(data);
    },
    toPush() {
      let data = {
        topic: this.pullForm.topic,
        keyType: this.pullForm.keyType,
        valueType: this.pullForm.valueType,
      };

      this.$refs.PushForm.show({
        title: `推送数据至:${data.topic}`,
        form: [this.form.toolbox.kafka.push],
        data: [data],
      });
    },
    async doPush(dataList) {
      let data = dataList[0];
      let param = this.toolboxWorker.getWorkParam(Object.assign({}, data));
      let res = await this.server.kafka.push(param);
      if (res.code == 0) {
        this.doPull();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    toCommit(data) {
      let groupId = this.pullForm.groupId;
      let topic = data.topic;
      let partition = data.partition;
      let offset = data.offset + 1;

      let msg =
        "使用groupId[" +
        groupId +
        "]消费主题[" +
        topic +
        "]分区[" +
        partition +
        "]位置[" +
        offset +
        "]";
      msg += "?";
      this.tool
        .confirm(msg)
        .then(async () => {
          this.doCommit({ groupId, topic, partition, offset });
        })
        .catch((e) => {});
    },
    async doCommit(data) {
      let param = this.toolboxWorker.getWorkParam(Object.assign({}, data));
      let res = await this.server.kafka.commit(param);
      if (res.code == 0) {
        this.doPull();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    toDelete(data) {
      let topic = data.topic;
      let partition = data.partition;
      let offset = data.offset + 1;

      let msg =
        "确认删除主题[" +
        topic +
        "]分区[" +
        partition +
        "]位置[" +
        offset +
        "]包含之前所有数据";
      msg += "?";
      this.tool
        .confirm(msg)
        .then(async () => {
          this.doDelete({ topic, partition, offset });
        })
        .catch((e) => {});
    },
    async doDelete(data) {
      let param = this.toolboxWorker.getWorkParam(Object.assign({}, data));
      let res = await this.server.kafka.deleteRecords(param);
      if (res.code == 0) {
        this.doPull();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    initDataList(dataList) {
      dataList = dataList || [];
      let columnList = [];
      let columnNameList = [];
      let headerColumnList = [];
      let headerColumnNameList = [];
      dataList.forEach((one) => {
        let headerValue = {};
        let jsonValue = {};
        let value = one.value;
        if (this.tool.isJSONString(value)) {
          try {
            let data = JSONbigString.parse(value);
            let jsonStr = JSON.stringify(data, null, "    ");
            jsonValue = JSON.parse(jsonStr);
          } catch (e) {
            this.form.valueJson = e.toString();
          }
        }
        for (let key in jsonValue) {
          if (columnNameList.indexOf(key) >= 0) {
            continue;
          }
          columnNameList.push(key);
          columnList.push({
            name: key,
            checked: true,
          });
        }
        if (one.headers && one.headers.length > 0) {
          one.headers.forEach((data) => {
            if (headerColumnNameList.indexOf(data.key) >= 0) {
              return;
            }
            headerColumnNameList.push(data.key);
            headerValue[data.key] = data.value;
            headerColumnList.push({
              name: data.key,
              checked: true,
            });
          });
        }
        one.headerValue = headerValue;
        one.jsonValue = jsonValue;
      });
      this.headerColumnList = headerColumnList;
      this.columnList = columnList;
    },
    async doPull() {
      this.dataListLoading = true;
      try {
        let param = this.toolboxWorker.getWorkParam(
          Object.assign({}, this.pullForm)
        );
        let res = await this.server.kafka.pull(param);
        if (res.code != 0) {
          this.tool.error(res.msg);
        }
        let dataList = res.data || [];
        this.initDataList(dataList);
        this.dataList = dataList;
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
.toolbox-kafka-topic-data {
  width: 100%;
  height: 100%;
}
</style>

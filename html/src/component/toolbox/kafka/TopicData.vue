<template>
  <div class="toolbox-kafka-topic-data">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="140px">
          <el-form
            class="pdt-10"
            label-width="90px"
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
              <el-input v-model="pullForm.keyType" />
            </el-form-item>
            <el-form-item label="ValueType">
              <el-input v-model="pullForm.valueType" />
            </el-form-item>
            <el-form-item label="">
              <div class="pd">
                <div class="tm-btn tm-btn-sm bg-teal-8 ft-13" @click="toPull">
                  拉取
                </div>
                <div class="tm-btn tm-btn-sm bg-green ft-13" @click="toPush">
                  新增
                </div>
              </div>
            </el-form-item>
          </el-form>
        </tm-layout>
        <tm-layout height="auto" class="scrollbar">
          <div class="pd-10" style="o">
            <table>
              <thead>
                <tr>
                  <th width="100">Partition</th>
                  <th width="80">Offset</th>
                  <th>Key</th>
                  <th>Value</th>
                  <th width="150">
                    <div
                      class="tm-link color-green-3 ft-14 mglr-2"
                      @click="toPush()"
                    >
                      <i class="mdi mdi-plus"></i>
                    </div>
                  </th>
                </tr>
              </thead>
              <tbody>
                <template v-if="msgList == null">
                  <tr>
                    <td colspan="5">
                      <div class="text-center ft-13 pdtb-10">拉取中...</div>
                    </td>
                  </tr>
                </template>
                <template v-else-if="msgList.length == 0">
                  <tr>
                    <td colspan="5">
                      <div class="text-center ft-13 pdtb-10">暂无匹配数据!</div>
                    </td>
                  </tr>
                </template>
                <template v-else>
                  <template v-for="(one, index) in msgList">
                    <tr :key="index" @click="rowClick(one)">
                      <td>{{ one.partition }}</td>
                      <td>{{ one.offset }}</td>
                      <td>{{ one.key }}</td>
                      <td>{{ one.value }}</td>
                      <td>
                        <div style="width: 140px">
                          <div
                            class="tm-btn color-grey tm-btn-xs"
                            @click="wrap.showData(one)"
                          >
                            查看
                          </div>
                          <div
                            class="tm-btn color-blue tm-btn-xs"
                            @click="toCommit(one)"
                          >
                            消费
                          </div>
                          <div
                            class="tm-btn color-orange tm-btn-xs"
                            @click="toDelete(one)"
                          >
                            删除
                          </div>
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
      ref="PushForm"
      :source="source"
      title="推送数据"
      :onSave="doPush"
    ></FormDialog>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxType", "toolbox", "option", "topic", "wrap"],
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
      msgList: null,
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
      this.wrap.showData(data);
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
      let param = {};
      Object.assign(param, data);
      let res = await this.wrap.work("push", param);
      if (res.code == 0) {
        this.doPull();
        return true;
      } else {
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
      let param = {};
      Object.assign(param, data);
      let res = await this.wrap.work("commit", param);
      if (res.code == 0) {
        this.doPull();
        return true;
      } else {
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
      let param = {};
      Object.assign(param, data);
      let res = await this.wrap.work("deleteRecords", param);
      if (res.code == 0) {
        this.doPull();
        return true;
      } else {
        return false;
      }
    },
    async doPull() {
      this.msgList = null;
      let param = {};
      Object.assign(param, this.pullForm);
      let res = await this.wrap.work("pull", param);
      res.data = res.data || {};
      let msgList = res.data.msgList;
      msgList.forEach((one) => {
        if (this.tool.isNotEmpty(one.value)) {
          try {
            if (
              (one.value.startsWith("{") && one.value.endsWith("}")) ||
              (one.value.startsWith("[") && one.value.endsWith("]"))
            ) {
              let data = JSON.parse(one.value);
              one.valueJson = JSON.stringify(data, null, "    ");
            }
          } catch (e) {
            one.valueJsonError = e;
          }
        }
      });
      this.msgList = res.data.msgList || [];
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

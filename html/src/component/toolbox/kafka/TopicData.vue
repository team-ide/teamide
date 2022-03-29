<template>
  <div class="toolbox-kafka-topic-data">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="120px">
          <b-form inline class="pdt-20 pdlr-10">
            <b-form-group label="Topic" label-size="sm" class="mgr-10">
              <b-form-input size="sm" v-model="pullForm.topic"> </b-form-input>
            </b-form-group>
            <b-form-group label="GroupId" label-size="sm" class="mgr-10">
              <b-form-input size="sm" v-model="pullForm.groupId">
              </b-form-input>
            </b-form-group>
            <b-form-group label="KeyType" label-size="sm" class="mgr-10">
              <b-form-input size="sm" v-model="pullForm.keyType">
              </b-form-input>
            </b-form-group>
            <b-form-group label="ValueType" label-size="sm" class="mgr-10">
              <b-form-input size="sm" v-model="pullForm.valueType">
              </b-form-input>
            </b-form-group>
          </b-form>
          <div class="pdlr-10 pdt-10">
            <div class="tm-btn tm-btn-sm bg-teal-8 ft-13" @click="toPull">
              拉取
            </div>
            <div class="tm-btn tm-btn-sm bg-green ft-13" @click="toPush">
              新增
            </div>
          </div>
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
                <template v-if="msgs == null">
                  <tr>
                    <td colspan="5">
                      <div class="text-center ft-13 pdtb-10">拉取中...</div>
                    </td>
                  </tr>
                </template>
                <template v-else-if="msgs.length == 0">
                  <tr>
                    <td colspan="5">
                      <div class="text-center ft-13 pdtb-10">暂无匹配数据!</div>
                    </td>
                  </tr>
                </template>
                <template v-else>
                  <template v-for="(one, index) in msgs">
                    <tr :key="index" @click="rowClick(one)" class="tm-pointer">
                      <td>{{ one.partition }}</td>
                      <td>{{ one.offset }}</td>
                      <td>{{ one.key }}</td>
                      <td>{{ one.value }}</td>
                      <td>
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
  </div>
</template>


<script>
export default {
  components: {},
  props: [
    "source",
    "data",
    "toolboxType",
    "toolbox",
    "option",
    "topic",
    "wrap",
  ],
  data() {
    return {
      ready: false,
      pullForm: {
        groupId: "test-group",
        topic: this.topic.name,
        keyType: "string",
        valueType: "string",
      },
      msgs: null,
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
      let nowTime = new Date().getTime();
      let clickTime = data.clickTime;
      data.clickTime = nowTime;
      if (clickTime) {
        let timeout = nowTime - clickTime;
        if (timeout < 300) {
          data.clickTime = null;
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
      this.wrap.showPushForm(data, (m) => {
        let flag = this.doPush(m);
        return flag;
      });
    },
    async doPush(data) {
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
      this.msgs = null;
      let param = {};
      Object.assign(param, this.pullForm);
      let res = await this.wrap.work("pull", param);
      res.data = res.data || {};
      let msgs = res.data.msgs;
      msgs.forEach((one) => {
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
      this.msgs = res.data.msgs || [];
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

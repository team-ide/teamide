<template>
  <div class="toolbox-kafka-tabs">
    <div class="default-tabs-container">
      <WorkspaceTabs :source="source" :itemsWorker="toolboxWorker.itemsWorker">
      </WorkspaceTabs>
    </div>
    <div class="default-spans-container">
      <WorkspaceSpans :source="source" :itemsWorker="toolboxWorker.itemsWorker">
        <template v-slot:span="{ item }">
          <TopicData
            :source="source"
            :toolboxWorker="toolboxWorker"
            :topic="item.extend.topic"
          >
          </TopicData>
        </template>
      </WorkspaceSpans>
    </div>
  </div>
</template>


<script>
import TopicData from "./TopicData";
export default {
  components: { TopicData },
  props: ["source", "toolboxWorker"],
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
.toolbox-kafka-tabs {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>

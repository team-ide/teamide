<template>
  <div class="toolbox-kafka-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout width="400px" class="">
          <Topic
            ref="Topic"
            :source="source"
            :toolboxWorker="toolboxWorker"
            :extend="extend"
          >
          </Topic>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <Tabs :source="source" :toolboxWorker="toolboxWorker"> </Tabs>
        </tm-layout>
      </tm-layout>
      <TopicForm :source="source" :toolboxWorker="toolboxWorker"> </TopicForm>
      <ShowData :source="source" :toolboxWorker="toolboxWorker"> </ShowData>
    </template>
  </div>
</template>


<script>
import Topic from "./Topic";
import Tabs from "./Tabs";
import TopicForm from "./TopicForm";
import ShowData from "./ShowData";
export default {
  components: { Topic, Tabs, TopicForm, ShowData },
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      ready: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.ready = true;
    },
    refresh() {
      this.$refs.Topic.refresh();
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-kafka-editor {
  width: 100%;
  height: 100%;
  user-select: text;
}
</style>

<template>
  <div class="toolbox-kafka-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout width="400px" class="scrollbar">
          <Topic
            ref="ToolboxKafkaTopic"
            :source="source"
            :toolbox="toolbox"
            :toolboxType="toolboxType"
            :wrap="wrap"
          >
          </Topic>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <Tabs
            :source="source"
            :toolbox="toolbox"
            :toolboxType="toolboxType"
            :wrap="wrap"
          >
          </Tabs>
        </tm-layout>
      </tm-layout>
      <TopicForm :source="source" :toolbox="toolbox" :wrap="wrap"> </TopicForm>
      <PushForm :source="source" :toolbox="toolbox" :wrap="wrap"> </PushForm>
      <ShowData :source="source" :toolbox="toolbox" :wrap="wrap"> </ShowData>
    </template>
  </div>
</template>


<script>
import Topic from "./Topic";
import Tabs from "./Tabs";
import TopicForm from "./TopicForm";
import PushForm from "./PushForm";
import ShowData from "./ShowData";
export default {
  components: { Topic, Tabs, TopicForm, PushForm, ShowData },
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
      this.ready = true;
    },
    refresh() {
      this.$refs.ToolboxKafkaTopic.refresh();
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

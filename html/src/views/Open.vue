<template>
  <div class="open-page-box">
    <template v-if="$route.params.type == 'toolbox'">
      <ToolboxOpen :source="source" :data="$route.params.data"></ToolboxOpen>
    </template>
  </div>
</template>

<script>
import ToolboxOpen from "@/views/toolbox/Open.vue";

export default {
  components: { ToolboxOpen },
  props: ["source"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {},
  created() {},
  mounted() {},
};
</script>

<style>
.open-page-box {
  width: 100%;
  height: 100%;
}
</style>

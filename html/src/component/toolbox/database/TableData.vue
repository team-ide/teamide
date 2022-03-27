<template>
  <div class="toolbox-database-table-data">
    <template v-if="ready">
      <div>
        {{ database }}
      </div>
      <div>
        {{ table }}
      </div>
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
    "database",
    "table",
    "wrap",
  ],
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
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-table-data {
  width: 100%;
  height: 100%;
}
</style>

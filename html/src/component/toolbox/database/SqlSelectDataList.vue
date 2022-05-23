<template>
  <div class="">
    <el-table
      :data="dataList"
      :border="true"
      height="100%"
      style="width: 100%"
      size="mini"
    >
      <el-table-column width="70" label="序号">
        <template slot-scope="scope">
          <span class="mgl-5">{{ scope.$index + 1 }}</span>
        </template>
      </el-table-column>
      <template v-for="(column, index) in tab.columnList">
        <el-table-column
          :key="index"
          :prop="column.name"
          :label="column.name"
          width="120"
        >
          <template slot-scope="scope">
            <div class="">
              <input
                v-model="scope.row[column.name]"
                :placeholder="scope.row[column.name] == null ? 'NULL' : ''"
                type="text"
              />
            </div>
          </template>
        </el-table-column>
      </template>
    </el-table>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "wrap", "tab"],
  data() {
    return {
      dataList: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      let dataList = [];
      for (let i = 0; i < this.tab.dataList.length; i++) {
        if (i < 50) {
          dataList.push(this.tab.dataList[i]);
        } else {
          break;
        }
      }

      this.dataList = dataList;
      this.bind();
    },
    bind() {
      let that = this;
      const selectWrap = this.$el.querySelector(".el-table__body-wrapper");
      selectWrap.addEventListener("scroll", function () {
        const scrollDistance =
          this.scrollHeight - this.scrollTop - this.clientHeight;
        // 判断是否到底,scrollTop为已滚动到元素上边界的像素数，scrollHeight为元素完整的高度，clientHeight为内容的可见宽度
        if (scrollDistance <= 10) {
          that.loadmore();
        }
      });
    },
    loadmore() {
      if (this.loadmore_ing) {
        return;
      }

      this.loadmore_ing = true;
      let allDataList = this.tab.dataList;
      let dataList = this.dataList;
      if (allDataList.length <= dataList.length) {
        return;
      }

      let start = dataList.length;
      for (let i = start; i < allDataList.length; i++) {
        if (i - start < 50) {
          dataList.push(allDataList[i]);
        } else {
          break;
        }
      }
      this.loadmore_ing = false;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
</style>

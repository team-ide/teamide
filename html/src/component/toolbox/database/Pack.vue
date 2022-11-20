<template>
  <span>
    <el-form-item label="数据库类型">
      <el-select
        placeholder="当前库类型"
        v-model="form.targetDatabaseType"
        style="width: 120px"
      >
        <el-option value="" label="当前库类型"> </el-option>
        <template v-for="(one, index) in source.databaseTypes">
          <el-option
            :key="index"
            :label="one.dialectName"
            :value="one.dialectName"
          >
          </el-option>
        </template>
      </el-select>
    </el-form-item>
    <el-form-item label="追加库名">
      <el-switch v-model="form.appendOwnerName" @change="change"> </el-switch>
    </el-form-item>
    <template v-if="form.appendOwnerCreateSql || form.appendOwnerName">
      <el-form-item label="库名包装">
        <el-select
          placeholder="默认"
          v-model="form.ownerNamePackChar"
          @change="change"
          style="width: 90px"
        >
          <el-option
            v-for="(one, index) in source.namePackChars"
            :key="index"
            :value="one.value"
            :label="one.text"
          >
          </el-option>
        </el-select>
      </el-form-item>
    </template>
    <el-form-item label="表名包装">
      <el-select
        placeholder="默认"
        v-model="form.tableNamePackChar"
        @change="change"
        style="width: 90px"
      >
        <el-option
          v-for="(one, index) in source.namePackChars"
          :key="index"
          :value="one.value"
          :label="one.text"
        >
        </el-option>
      </el-select>
    </el-form-item>
    <el-form-item label="字段包装">
      <el-select
        placeholder="默认"
        v-model="form.columnNamePackChar"
        @change="change"
        style="width: 90px"
      >
        <el-option
          v-for="(one, index) in source.namePackChars"
          :key="index"
          :value="one.value"
          :label="one.text"
        >
        </el-option>
      </el-select>
    </el-form-item>
    <el-form-item label="字符值包装">
      <el-select
        v-model="form.sqlValuePackChar"
        @change="change"
        style="width: 90px"
      >
        <el-option
          v-for="(one, index) in source.sqlValuePackChars"
          :key="index"
          :value="one.value"
          :label="one.text"
        >
        </el-option>
      </el-select>
    </el-form-item>
  </span>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "form", "change"],
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
</style>

<template>
  <div class="pdlr-10 color-grey ft-12" style="user-select: text">
    <div>
      <span class="color-orange pdr-10">表达式</span>
      <span>表达式，如：'aa' + 'c'，返回“aac”；1 + 2，返回“3”</span>
    </div>
    <div class="pdl-10">
      <div>
        <span class="color-orange pdr-10">_$index</span>
        <span>索引，每个策略数据从0开始，最大为当前策略数据数量-1</span>
      </div>
      <div>
        <span class="color-orange pdr-10">_$value_index</span>
        <span>值索引，从0开始，最大为值数量-1</span>
      </div>
      <div>
        <span class="color-orange pdr-10">_$now()</span>
        <span>当前时间对象</span>
      </div>
      <div>
        <span class="color-orange pdr-10">_$nowTime()</span>
        <span>当前时间戳</span>
      </div>
      <div>
        <span class="color-orange pdr-10">_$uuid()</span>
        <span>生成UUID</span>
      </div>
      <div>
        <span class="color-orange pdr-10">
          _$randomString(minLength, maxLength)
        </span>
        <span>随机字符串</span>
      </div>
      <div>
        <span class="color-orange pdr-10">_$randomInt(min, max)</span>
        <span>随机数字</span>
      </div>
      <div>
        <span class="color-orange pdr-10">
          _$randomUserName(minLength, maxLength)
        </span>
        <span>随机用户姓名</span>
      </div>
      <div>
        <span class="color-orange pdr-10"> _$toPinYin(str) </span>
        <span>转为拼音</span>
      </div>
    </div>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source"],
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
</style>

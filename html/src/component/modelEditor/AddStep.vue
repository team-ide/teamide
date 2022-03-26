<template>
  <b-dropdown class="trim-dropdown-style" no-caret ref="dropdown">
    <template #button-content>
      <template v-if="isChild">
        <span class="tm-pointer color-green mgl-5"> 添加子步骤 </span>
      </template>
      <template v-else>
        <span class="tm-pointer color-green mgl-5"> 添加步骤 </span>
      </template>
    </template>
    <MenuBox class="bd-0">
      <template v-for="(one, index) in stepTypes">
        <template v-if="one.subs == null">
          <MenuItem :key="'stepType-' + index" @click="addStep(one)">
            {{ one.text }}
          </MenuItem>
        </template>
        <template v-else>
          <MenuItem :key="'stepType-' + index">
            {{ one.text }}
            <MenuSubBox slot="MenuSubBox">
              <template v-for="(sub, subIndex) in one.subs">
                <MenuItem :key="'stepType-' + subIndex" @click="addStep(sub)">
                  {{ sub.text }}
                </MenuItem>
              </template>
            </MenuSubBox>
          </MenuItem>
        </template>
      </template>
    </MenuBox>
  </b-dropdown>
</template>


<script>
export default {
  components: {},
  props: ["source", "context", "wrap", "bean", "isChild"],
  data() {
    return {
      ready: false,
      stepTypes: [
        { value: "variables", text: "定义变量", isList: true },
        { value: "validates", text: "变量验证", isList: true },
        { value: "error", text: "抛出异常" },
        {
          text: "锁操作",
          subs: [
            { value: "lock", text: "锁定" },
            { value: "unlock", text: "解锁" },
          ],
        },
        {
          text: "Sql操作",
          subs: [
            { value: "sqlSelect", text: "Sql Select" },
            { value: "sqlInsert", text: "Sql Insert" },
            { value: "sqlUpdate", text: "Sql Update" },
            { value: "sqlDelete", text: "Sql Delete" },
          ],
        },
        {
          text: "Redis操作",
          subs: [
            { value: "redisGet", text: "Redis Get" },
            { value: "redisSet", text: "Redis Set" },
            { value: "redisDel", text: "Redis Del" },
          ],
        },
        {
          text: "文件操作",
          subs: [
            { value: "fileSave", text: "文件保存" },
            { value: "fileGet", text: "文件信息获取" },
            { value: "fileRead", text: "文件读取" },
            { value: "fileDelete", text: "文件删除" },
          ],
        },
      ],
    };
  },
  computed: {},
  watch: {
    bean() {
      this.init();
    },
  },
  methods: {
    addStep(stepType) {
      let step = {};
      if (stepType.isList) {
        step[stepType.value] = [];
      } else {
        step[stepType.value] = {};
      }
      this.$refs.dropdown.hide();
      this.wrap.push(this.bean, "steps", step);
    },
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
.trim-dropdown-style {
  vertical-align: unset !important;
  display: unset !important;
}
.trim-dropdown-style .btn {
  padding: 0px !important;
  margin: 0px !important;
  font-weight: unset !important;
  font-size: unset !important;
  color: unset !important;
  background: unset !important;
  border: unset !important;
  line-height: unset !important;
  border-radius: unset !important;
  transition: unset !important;
  user-select: none !important;
  box-shadow: none !important;
  vertical-align: unset !important;
}
.trim-dropdown-style .btn .b-icon {
  font-size: unset !important;
}
.trim-dropdown-style li {
  line-height: unset !important;
  margin: 0px !important;
}
.trim-dropdown-style .dropdown-item {
  padding: 2px 10px;
}
</style>

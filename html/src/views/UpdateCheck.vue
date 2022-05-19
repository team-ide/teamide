<template>
  <el-dialog
    ref="modal"
    :title="`检测新版本`"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="900px"
  >
    <div class="mgt--20" style="user-select: text">
      <div class="ft-14 mgb-10">
        当前版本
        <span class="color-green mgl-10">{{ currentVersion }}</span>
      </div>
      <div class="ft-14 mgb-10">
        新版本请查看
        <a
          class="tm-link color-green mgl-10"
          target="_blank"
          :href="githubReleasesURL"
          >{{ githubReleasesURL }}</a
        >
      </div>
      <div v-html="releaseHistory"></div>
    </div>
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {
      showDialog: false,
      currentVersion: `检测中...`,
      githubReleasesURL: `https://github.com/team-ide/teamide/releases`,
      releaseHistory: ``,
    };
  },
  computed: {},
  watch: {},
  methods: {
    show() {
      this.showDialog = true;
    },
    hide() {
      this.showDialog = false;
    },
    getHTML(releaseHistory) {
      try {
        let div = this.tool.jQuery("<div>" + releaseHistory + "</div>");
        let find = false;
        div.find("h2").each((index, one) => {
          one = this.tool.jQuery(one);
          one.find("a").remove();
          if (
            one
              .text()
              .toLowerCase()
              .indexOf(("v" + this.currentVersion).toLowerCase()) >= 0
          ) {
            find = true;
            one.addClass("color-orange");
            one.append('<span class="mgl-5 ft-12">当前版本</span>');
          }
          if (!find) {
            this.source.hasNewVersion = true;
            one.addClass("color-green");
            one.append('<span class="mgl-5 ft-12">新版本</span>');
          }
        });
        releaseHistory = div.html();
      } catch (error) {
        if (error != null) {
          return releaseHistory;
        }
      }
      return releaseHistory;
    },
    updateCheck() {
      this.server
        .updateCheck({})
        .then((res) => {
          if (res.code == 0 && res.data != null) {
            if (this.tool.isNotEmpty(res.data.currentVersion)) {
              this.currentVersion = res.data.currentVersion;
            }
            if (this.tool.isNotEmpty(res.data.githubReleasesURL)) {
              this.githubReleasesURL = res.data.githubReleasesURL;
            }
            if (this.tool.isNotEmpty(res.data.releaseHistory)) {
              this.releaseHistory = this.getHTML(res.data.releaseHistory);
            }
          }
          setTimeout(() => {
            this.updateCheck();
          }, 1000 * 60 * 10);
        })
        .catch(() => {
          setTimeout(() => {
            this.updateCheck();
          }, 1000 * 60 * 10);
        });
    },
  },
  created() {},
  mounted() {
    this.tool.showUpdateCheck = this.show;
    this.updateCheck();
  },
};
</script>

<style>
</style>

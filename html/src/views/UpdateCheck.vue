<template>
  <el-dialog
    ref="modal"
    :title="`检测新版本`"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    top="20px"
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
      currentVersionNumber: 0,
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
    getVersionNumberByVersion(version) {
      let number = 0;
      if (this.tool.isEmpty(version)) {
        return number;
      }
      version = ("" + version).trim();
      let rule = /([1-9]\d|[0-9])(.([1-9]\d|\d)){1,}/;
      let res = version.match(rule);
      if (res && res[0]) {
        let str = res[0];
        let ss = str.split(".");
        ss.forEach((s, i) => {
          let sN = Number(s);
          for (let n = ss.length - i - 1; n > 0; n--) {
            sN = sN * 100;
          }
          number = Number(number) + Number(sN);
        });
      }
      return number;
    },
    getHTML(releaseHistory) {
      try {
        let div = this.tool.jQuery("<div>" + releaseHistory + "</div>");
        div.find("h2").each((index, one) => {
          one = this.tool.jQuery(one);
          one.find("a").remove();
          let version = one.text().toLowerCase();
          let versionNumber = this.getVersionNumberByVersion(version);
          if (versionNumber == this.currentVersionNumber) {
            one.addClass("color-green");
            one.append('<span class="mgl-5 ft-12">当前版本</span>');
          } else if (versionNumber > this.currentVersionNumber) {
            this.source.hasNewVersion = true;
            one.addClass("color-orange");
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
              this.currentVersionNumber = this.getVersionNumberByVersion(
                this.currentVersion
              );
            } else {
              this.currentVersion = "未检测到版本";
              this.currentVersionNumber = this.getVersionNumberByVersion("");
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
          }, 1000 * 60 * 30);
        })
        .catch(() => {
          setTimeout(() => {
            this.updateCheck();
          }, 1000 * 60 * 30);
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

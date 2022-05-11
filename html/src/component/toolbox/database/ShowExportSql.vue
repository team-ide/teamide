<template>
  <el-dialog
    ref="modal"
    :title="'导出：' + (tableDetail == null ? '' : tableDetail.name)"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="1200px"
  >
    <div class="mgt--20 toolbox-database-export-sql">
      <el-form
        class="mgt-10"
        ref="form"
        :model="form"
        label-width="90px"
        size="mini"
        :inline="true"
      >
        <el-form-item label="SQL类型">
          <el-select v-model="form.sqlType" @change="toLoad">
            <el-option
              v-for="(one, index) in sqlTypes"
              :key="index"
              :value="one.value"
            >
              {{ one.text }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="追加库名">
          <el-switch v-model="form.appendDatabase" @change="toLoad">
          </el-switch>
        </el-form-item>
        <template v-if="form.appendDatabase">
          <el-form-item label="库名包装">
            <el-select
              placeholder="不包装"
              v-model="form.databasePackingCharacter"
              @change="toLoad"
              style="width: 100px"
            >
              <el-option
                v-for="(one, index) in packingCharacters"
                :key="index"
                :value="one.value"
              >
                {{ one.text }}
              </el-option>
            </el-select>
          </el-form-item>
        </template>
        <el-form-item label="表名包装">
          <el-select
            placeholder="不包装"
            v-model="form.tablePackingCharacter"
            @change="toLoad"
            style="width: 100px"
          >
            <el-option
              v-for="(one, index) in packingCharacters"
              :key="index"
              :value="one.value"
            >
              {{ one.text }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="字段包装">
          <el-select
            placeholder="不包装"
            v-model="form.columnPackingCharacter"
            @change="toLoad"
            style="width: 100px"
          >
            <el-option
              v-for="(one, index) in packingCharacters"
              :key="index"
              :value="one.value"
            >
              {{ one.text }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="字符值包装">
          <el-select
            placeholder="不包装"
            v-model="form.stringPackingCharacter"
            @change="toLoad"
            style="width: 100px"
          >
            <el-option
              v-for="(one, index) in stringPackingCharacters"
              :key="index"
              :value="one.value"
            >
              {{ one.text }}
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div>
        <textarea v-model="showSQL"> </textarea>
      </div>
    </div>
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "wrap"],
  data() {
    return {
      showDialog: false,
      showSQL: null,
      tableDetail: null,
      sqlTypes: [
        { value: "insert", text: "Insert" },
        { value: "update", text: "Update" },
        { value: "delete", text: "Delete" },
      ],
      packingCharacters: [
        { value: "", text: "不包装" },
        { value: "'", text: "'" },
        { value: '"', text: '"' },
        { value: "`", text: "`" },
      ],
      stringPackingCharacters: [
        { value: "'", text: "'" },
        { value: '"', text: '"' },
      ],
      form: {
        sqlType: "insert",
        appendDatabase: false,
        databasePackingCharacter: null,
        tablePackingCharacter: null,
        columnPackingCharacter: null,
        stringPackingCharacter: "'",
      },
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    async show(database, tableDetail, dataList) {
      this.database = database;
      this.dataList = dataList || [];
      this.tableDetail = tableDetail;
      await this.toLoad();
      this.showDialog = true;
    },
    hide() {
      this.showDialog = false;
    },
    async toLoad() {
      this.showSQL = "";
      let sqls = await this.initSqls();
      sqls.forEach((sql) => {
        this.showSQL += sql + ";\n";
      });
    },
    packingCharacterDatabase(value) {
      if (
        this.tool.isEmpty(this.form.databasePackingCharacter) ||
        this.tool.isEmpty(value)
      ) {
        return value;
      }
      return (
        this.form.databasePackingCharacter +
        value +
        this.form.databasePackingCharacter
      );
    },
    packingCharacterTable(value) {
      if (
        this.tool.isEmpty(this.form.tablePackingCharacter) ||
        this.tool.isEmpty(value)
      ) {
        return value;
      }
      return (
        this.form.tablePackingCharacter +
        value +
        this.form.tablePackingCharacter
      );
    },
    packingCharacterColumn(value) {
      if (
        this.tool.isEmpty(this.form.columnPackingCharacter) ||
        this.tool.isEmpty(value)
      ) {
        return value;
      }
      return (
        this.form.columnPackingCharacter +
        value +
        this.form.columnPackingCharacter
      );
    },
    packingCharacterString(value) {
      if (
        this.tool.isEmpty(this.form.stringPackingCharacter) ||
        this.tool.isEmpty(value)
      ) {
        return value;
      }
      if (typeof value != "string") {
        return value;
      }
      if (value.indexOf(this.form.stringPackingCharacter) < 0) {
        return (
          this.form.stringPackingCharacter +
          value +
          this.form.stringPackingCharacter
        );
      }
      return (
        this.form.stringPackingCharacter +
        this.packingCharacterEscape(this.form.stringPackingCharacter, value) +
        this.form.stringPackingCharacter
      );
    },
    packingCharacterEscape(packing, value) {
      let res = "";
      var arr = value.split("");
      for (var i = 0; i < arr.length; i++) {
        let c = arr[i];
        if (packing == c) {
          res += "\\" + c;
        } else if (c == "\\") {
          res += "\\" + c;
        } else {
          res += c;
        }
      }

      return res;
    },
    async initSqls() {
      let sqls = [];
      if (this.tableDetail == null) {
        return sqls;
      }

      let keys = [];
      this.tableDetail.columnList.forEach((column) => {
        if (column.primaryKey) {
          keys.push(column.name);
        }
      });
      this.dataList.forEach((one) => {
        let insertSql = "INSERT INTO ";

        if (this.form.appendDatabase) {
          insertSql += this.packingCharacterDatabase(this.database) + ".";
        }
        insertSql += this.packingCharacterTable(this.tableDetail.name);

        let updateSql = "UPDATE ";

        if (this.form.appendDatabase) {
          updateSql += this.packingCharacterDatabase(this.database) + ".";
        }
        updateSql +=
          this.packingCharacterTable(this.tableDetail.name) + " SET ";

        let deleteSql = "DELETE FROM ";

        if (this.form.appendDatabase) {
          deleteSql += this.packingCharacterDatabase(this.database) + ".";
        }
        deleteSql += this.packingCharacterTable(this.tableDetail.name);

        insertSql += " (";
        this.tableDetail.columnList.forEach((column) => {
          insertSql += "" + this.packingCharacterColumn(column.name) + ", ";
        });
        if (insertSql.endsWith(", ")) {
          insertSql = insertSql.substring(0, insertSql.length - 2);
        }
        insertSql += ")";

        insertSql += " VALUES (";

        let whereSql = "WHERE ";

        this.tableDetail.columnList.forEach((column) => {
          let value = one[column.name];
          let valueSql = value;
          if (value == null) {
            valueSql = "NULL";
          } else {
            if (this.wrap.columnIsNumber(column)) {
              valueSql = value;
            } else {
              value = this.wrap.formatDateColumn(column, value);
              valueSql = this.packingCharacterString(value);
            }
          }

          insertSql += valueSql + ", ";

          if (keys.length > 0) {
            if (keys.indexOf(column.name) >= 0) {
              whereSql +=
                this.packingCharacterColumn(column.name) +
                "=" +
                valueSql +
                " AND ";
            } else {
              updateSql +=
                this.packingCharacterColumn(column.name) +
                "=" +
                valueSql +
                ", ";
            }
          } else {
            updateSql +=
              this.packingCharacterColumn(column.name) + "=" + valueSql + ", ";
            whereSql +=
              this.packingCharacterColumn(column.name) +
              "=" +
              valueSql +
              " AND ";
          }
        });
        if (insertSql.endsWith(", ")) {
          insertSql = insertSql.substring(0, insertSql.length - 2);
        }
        insertSql += ")";

        if (updateSql.endsWith(", ")) {
          updateSql = updateSql.substring(0, updateSql.length - 2);
        }
        if (whereSql.endsWith("AND ")) {
          whereSql = whereSql.substring(0, whereSql.length - "AND ".length);
        }
        updateSql += " " + whereSql;
        deleteSql += " " + whereSql;
        switch (this.form.sqlType) {
          case "insert":
            sqls.push(insertSql);
            break;
          case "update":
            sqls.push(updateSql);
            break;
          case "delete":
            sqls.push(deleteSql);
            break;
        }
      });
      return sqls;
    },
    init() {},
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.wrap.showExportSql = this.show;
    this.init();
  },
};
</script>

<style>
.toolbox-database-export-sql textarea {
  width: 100%;
  height: 400px;
  letter-spacing: 1px;
  word-spacing: 5px;
  word-break: break-all;
  font-size: 12px;
  border: 1px solid #ddd;
  padding: 0px 5px;
  outline: none;
  user-select: none;
}
</style>

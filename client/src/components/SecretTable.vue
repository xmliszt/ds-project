<template>
  <div>
    <SecretCreation :role="role" @refreshTable="refreshTable" />
    <div class="secret-table">
      <div class="secret-table-wrapper">
        <el-table
          class="secret-data-table"
          :data="filteredData"
          style="width: 100%"
          height="38rem"
          v-loading="loading"
          @selection-change="handleSelectChange"
        >
          <el-table-column type="index" width="80"> </el-table-column>
          <el-table-column
            prop="alias"
            label="Alias / Description"
            width="500"
            show-overflow-tooltip
            sortable
          >
          </el-table-column>
          <el-table-column
            label="Secret Value"
            width="500"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              <span v-show="scope.row.show">{{ scope.row.value }}</span>
              <span v-show="scope.row.edit">
                <el-input
                  v-model="scope.row.value"
                  show-password
                  clearable
                  size="mini"
                ></el-input>
              </span>
              <span v-show="!scope.row.show && !scope.row.edit">{{
                "*".repeat(scope.row.value.length)
              }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="role" sortable label="Role"></el-table-column>
          <el-table-column align="right">
            <!-- eslint-disable-next-line vue/no-unused-vars -->
            <template slot="header" slot-scope="scope">
              <el-input
                v-model="search"
                size="small"
                placeholder="Search alias/description..."
                clearable
              />
            </template>
            <template slot-scope="scope">
              <el-tooltip v-if="scope.row.show" content="Hide this secret">
                <el-button size="mini" @click="handleHide(scope.row)" plain
                  >Hide</el-button
                >
              </el-tooltip>
              <el-tooltip v-else content="Show this secret">
                <el-button size="mini" @click="handleShow(scope.row)" plain
                  >Show</el-button
                >
              </el-tooltip>
              <el-tooltip
                v-if="scope.row.edit"
                content="Finish edit this secret"
              >
                <el-button size="mini" @click="handleEditDone(scope.row)" plain
                  >Done</el-button
                >
              </el-tooltip>
              <el-tooltip v-else content="Edit this secret">
                <el-button size="mini" @click="handleEdit(scope.row)" plain
                  >Edit</el-button
                >
              </el-tooltip>

              <el-tooltip content="Delete this secret">
                <el-button
                  size="mini"
                  type="danger"
                  @click="handleDelete(scope.row)"
                  >Delete</el-button
                >
              </el-tooltip>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script>
import SecretCreation from "../components/SecretCreation";

import { getSecrets, putSecret, deleteSecret } from "../service/secret";

import { parseJwt } from "../util/util";
export default {
  components: {
    SecretCreation,
  },
  data() {
    return {
      role: 0,
      loading: false,
      search: "",
      secretsData: [],
      selectedSecrets: [],
    };
  },
  methods: {
    handleSelectChange(val) {
      this.selectedSecrets = val;
    },
    handleShow(row) {
      row.show = true;
      row.edit = false;
    },
    handleHide(row) {
      row.show = false;
    },
    handleEdit(row) {
      row.edit = true;
      row.show = false;
    },
    handleEditDone(row) {
      row.edit = false;
      putSecret(row.alias, row.value, this.role).then((result) => {
        if (result.success) {
          this.$message.success("Secret has been updated!");
        } else {
          this.$message.error(
            "Failed to update the secret: " + result.error.response.data.Error
          );
        }
      });
    },
    handleDelete(row) {
      this.$confirm("Are you sure to delete this secret?").then(() => {
        deleteSecret(row.alias).then((result) => {
          if (result.success) {
            this.$message.success("Secret has been deleted!");
            this.refreshTable();
          } else {
            this.$message.error(
              "Failed to delete the secret: " + result.error.response.data.Error
            );
          }
        });
      });
    },
    refreshTable() {
      this.loading = true;
      getSecrets().then((result) => {
        if (result.success) {
          this.secretsData = result.data;
          if (result.deadNodes.length > 0) {
            result.deadNodes.forEach((nid) => {
              this.$message.warning(`Physical Node ${nid} is down!`);
            });
          }
        } else {
          this.$message.error(
            "Failed to fetch secrets: " + result.error.response.data.Error
          );
        }
      });
      this.loading = false;
    },
  },
  created() {
    var jwt = localStorage.getItem("token");
    if (jwt === null) {
      this.$message.error("You are not logged in!");
      this.$router.push("/");
    }
    var json = parseJwt(jwt);
    this.role = json.role;
    this.refreshTable();
  },
  computed: {
    filteredData() {
      var data = this.secretsData.filter(
        (data) =>
          !this.search ||
          data.alias.toLowerCase().includes(this.search.toLowerCase())
      );
      return data;
    },
  },
};
</script>

<style scoped>
.secret-table {
  margin-top: 2rem;
}

.el-table__body-wrapper::-webkit-scrollbar-track {
  background: #dba95f; /* color of the tracking area */
}

.el-table__body-wrapper::-webkit-scrollbar-thumb {
  background-color: rgb(253, 253, 253); /* color of the scroll thumb */
  border-radius: 20px; /* roundness of the scroll thumb */
}

* {
  scrollbar-color: #dba95f rgb(253, 253, 253);
}
</style>

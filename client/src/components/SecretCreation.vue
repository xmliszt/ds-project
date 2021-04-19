<template>
  <div class="secret-creation">
    <el-tooltip content="Create a new secret">
      <el-button
        type="success"
        icon="el-icon-plus"
        @click="openCreateSecretDialog"
        >Create Secret</el-button
      >
    </el-tooltip>
    <el-dialog
      :title="dialogTitle"
      :visible.sync="showCreateSecretDialog"
      width="50%"
      :before-close="closeCreateSecretDialog"
    >
      <el-form
        ref="secretForm"
        :rules="rules"
        :model="form"
        label-width="150px"
        label-position="left"
      >
        <el-form-item label="Alias / Description" prop="alias" required>
          <el-input
            size="medium"
            placeholder="Alias/Description..."
            v-model="form.alias"
            clearable
          ></el-input>
        </el-form-item>
        <el-form-item label="Secret Value" prop="value" required>
          <el-input
            size="medium"
            placeholder="Secret value..."
            v-model="form.value"
            show-password
            clearable
          ></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button type="success" @click="submitForm">Create</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { putSecret, getSecrets } from "../service/secret";
import { itemExistsInArray } from "../util/util";

export default {
  props: ["role"],
  data() {
    return {
      showCreateSecretDialog: false,
      dialogTitle: `Create a new secret for role: ${this.role}`,
      form: {
        alias: "",
        value: "",
      },
      rules: {
        alias: [
          {
            required: true,
            message: "Alias / Description cannot be empty!",
            trigger: "blur",
          },
        ],
        value: [
          {
            required: true,
            message: "Secret value cannot be empty!",
            trigger: "blur",
          },
        ],
      },
    };
  },
  methods: {
    openCreateSecretDialog() {
      this.showCreateSecretDialog = true;
    },
    closeCreateSecretDialog() {
      this.$confirm("Are you sure to abort creation?").then(() => {
        this.$message.error("You cancelled secret creation.");
        this.showCreateSecretDialog = false;
      });
    },
    submitForm() {
      this.$refs["secretForm"].validate((valid) => {
        if (valid) {
          getSecrets().then((result) => {
            if (result.success) {
              var secrets = result.data;
              var exist = itemExistsInArray(secrets, this.form.alias, "alias");
              if (exist) {
                this.$message.error(
                  "The alias / description already exists! Please change your alias / description!"
                );
              } else {
                putSecret(this.form.alias, this.form.value, this.role).then(
                  (result) => {
                    if (result.success) {
                      this.$message.success("New secret has been created!");
                      this.showCreateSecretDialog = false;
                      this.$emit("refreshTable");
                    } else {
                      this.$message.error(
                        "Failed to create new secret: " +
                          result.error.response.data.Error
                      );
                    }
                  }
                );
              }
            } else {
              this.$message.error(
                "Failed to check secrets: " + result.error.response.data.Error
              );
            }
          });
        } else {
          return false;
        }
      });
    },
  },
};
</script>

<style scoped>
.secret-creation {
  margin-top: 2rem;
  display: flex;
  justify-content: space-between;
}

.dialog-wrapper {
  z-index: 999;
  position: fixed;
}
</style>

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
        <el-button type="primary" @click="submitForm">Confirm</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
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
      this.$confirm("Are you sure to abort creation?")
        .then(() => {
          this.showCreateSecretDialog = false;
        })
        .catch(() => {
          this.$message.error("You cancelled secret creation.");
          this.showCreateSecretDialog = false;
        });
    },
    submitForm() {
      this.$refs["secretForm"].validate((valid) => {
        if (valid) {
          //TODO: create a secret
          console.log("Create a secret");
          this.showCreateSecretDialog = false;
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

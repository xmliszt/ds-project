<template>
  <div class="login">
    <el-form
      ref="loginForm"
      :rules="rules"
      :model="form"
      style="width: 30rem"
      label-width="120px"
      label-position="left"
    >
      <el-form-item label="Username" prop="username" required>
        <el-input
          size="medium"
          placeholder="Username"
          v-model="form.username"
          clearable
        ></el-input>
      </el-form-item>
      <el-form-item label="Password" prop="password" required>
        <el-input
          size="medium"
          placeholder="Password"
          v-model="form.password"
          show-password
          clearable
        ></el-input>
      </el-form-item>
      <el-form-item style="margin-left: -120px">
        <el-button
          style="margin-top: 20px"
          size="medium"
          type="danger"
          plain
          @click="submitForm"
          >LOGIN</el-button
        >
      </el-form-item>
      <el-link href="/register"
        >Don't have an account? Click here to sign up!</el-link
      >
    </el-form>
  </div>
</template>

<script>
import { loginUser } from "../service/user.js";

export default {
  data() {
    return {
      form: {
        username: "",
        password: "",
      },
      rules: {
        username: [
          {
            required: true,
            message: "Username cannot be empty!",
            trigger: "blur",
          },
        ],
        password: [
          {
            required: true,
            message: "Password cannot be empty!",
            trigger: "blur",
          },
        ],
      },
    };
  },
  methods: {
    submitForm() {
      this.$refs["loginForm"].validate((valid) => {
        if (valid) {
          // register
          loginUser(this.form.username, this.form.password).then((result) => {
            if (result.success) {
              this.$message.success("You are logged in!");
              this.$router.push("/dashboard");
            } else {
              this.$message.error(
                "Login failed: " + result.error.response.data.Error
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
.not-register {
  font-size: 1rem;
  font-weight: 100;
}

.login {
  display: flex;
  justify-content: center;
  align-items: center;
}

.el-input--medium .el-input__inner {
  height: 80px !important;
  font-weight: 700;
}

.el-input.is-active .el-input__inner,
.el-input__inner:focus {
  border-color: #f56c6c !important;
}

.el-button--medium {
  font-weight: 700 !important;
}

.el-slider__button {
  border: 2px solid #f56c6c !important;
}
</style>

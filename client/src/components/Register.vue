<template>
  <div class="register-container">
    <el-form
      ref="registerForm"
      :model="form"
      style="width: 40rem"
      label-width="120px"
      label-position="left"
      :rules="rules"
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
      <el-form-item prop="role" label="Role" required>
        <el-slider
          v-model="form.role"
          :step="1"
          show-stops
          :min="roleStepper.min"
          :max="roleStepper.max"
        ></el-slider>
      </el-form-item>
      <el-form-item style="margin-left: -120px">
        <el-button
          style="margin-top: 20px"
          size="medium"
          type="danger"
          plain
          @click="submitForm"
          >Register</el-button
        >
        <el-button
          style="margin-top: 20px"
          size="medium"
          type="primary"
          plain
          @click="$router.push('/')"
        >
          Back To Login
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { registerUser } from "../service/user.js";
export default {
  data() {
    return {
      form: {
        username: "",
        password: "",
        role: 1,
      },
      roleStepper: {
        min: 1,
        max: 5,
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
        role: [
          {
            required: true,
            message: "Role cannot be empty!",
            trigger: "blur",
          },
        ],
      },
    };
  },
  methods: {
    submitForm() {
      this.$refs["registerForm"].validate((valid) => {
        if (valid) {
          // register
          registerUser(
            this.form.username,
            this.form.password,
            this.form.role
          ).then((result) => {
            if (result.success) {
              this.$message.success(
                "Your registration is successful! Please login"
              );
              this.$router.push("/");
            } else {
              this.$message.error(
                "Registration failed: " + result.error.response.data.Error
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

<style scoped></style>

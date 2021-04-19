<template>
  <div class="logout">
    <el-tooltip content="Logout">
      <el-button
        type="danger"
        icon="el-icon-switch-button"
        @click="logoutUser"
        plain
      ></el-button>
    </el-tooltip>
  </div>
</template>

<script>
export default {
  methods: {
    logoutUser() {
      window.localStorage.removeItem("token");
      this.$message.warning("You have been logged out!");
      this.$router.push("/");
    },
  },
};
</script>

<style scoped>
.logout {
  position: fixed;
  top: 5rem;
  right: 5rem;
}
</style>

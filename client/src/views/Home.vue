<template>
  <div class="home">
    <div class="home-wrapper">
      <div class="home-title">
        <div class="title-text">{{ title }}</div>
        <el-tag type="info" class="subtitle">{{ subtitle }}</el-tag>
      </div>
      <div class="login-container">
        <Login />
      </div>
    </div>
  </div>
</template>

<script>
import Login from "../components/Login";
export default {
  components: { Login },
  name: "Home",
  data() {
    return {
      title: "E-SAFE: Login",
      subtitle: "IAM Secret Management",
    };
  },
};
</script>

<style scoped>
.home {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}
.home-title {
  padding: 0rem 10rem 3rem 10rem;
}

.title-text {
  color: #f56c6c;
  font-weight: 100;
  font-size: 4rem;
}

.subtitle {
  color: #606266;
  font-weight: 100;
  font-size: 1rem;
}
</style>

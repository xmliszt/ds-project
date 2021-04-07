<template>
  <div class="register">
    <div class="register-parent-wrapper">
      <div class="home-title">
        <div class="title-text">{{ title }}</div>
        <el-tag type="info" class="subtitle">{{ subtitle }}</el-tag>
      </div>
      <div class="register-container">
        <Register />
      </div>
    </div>
  </div>
</template>

<script>
import Register from "../components/Register.vue";
export default {
  components: {
    Register,
  },
  data() {
    return {
      title: "E-SAFE: Register An Account",
      subtitle: "IAM Secret Management",
    };
  },
};
</script>

<style scoped>
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
.register {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>

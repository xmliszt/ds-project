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
  </div>
</template>

<script>
export default {
  data() {
    return {
      showCreateSecretDialog: false,
    };
  },
  methods: {
    openCreateSecretDialog() {
      this.showCreateSecretDialog = true;
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
</style>

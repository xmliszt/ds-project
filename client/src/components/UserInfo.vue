<template>
  <div class="user-info">
    <h1>Welcome! {{ username }}</h1>
    <el-tooltip effect="dark" :content="role_tooltip_text">
      <el-rate
        v-model="computedRole"
        disabled
        :icon-classes="icon_classes"
        disabled-void-icon-class="el-icon-user"
      ></el-rate>
    </el-tooltip>
  </div>
</template>

<script>
export default {
  props: ["username", "role"],
  data() {
    return {
      role_tooltip_text: `Your IAM accessibility is: ${this.role}`,
      icon_classes: { 5: "el-icon-user-solid" },
    };
  },
  computed: {
    computedRole() {
      return 6 - Number(this.role);
    },
  },
};
</script>

<style scoped></style>

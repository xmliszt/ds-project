<template>
  <div class="monitor-switch">
    <el-radio-group
      v-model="onMonitor"
      size="small"
      fill="#67C23A"
      @change="handleSwitchMonitor"
    >
      <el-radio-button label="Secrets"></el-radio-button>
      <el-radio-button label="Monitor"></el-radio-button>
    </el-radio-group>
  </div>
</template>

<script>
export default {
  data() {
    return {
      onMonitor: "Secrets",
    };
  },
  methods: {
    handleSwitchMonitor(val) {
      if (val === "Secrets") {
        // show secret table
        this.$router.push("/dashboard");
      } else {
        // show monitor
        this.$router.push("/dashboard/monitor");
      }
    },
  },
  created() {
    if (this.onMonitor === "Secrets") {
      // show secret table
      this.$router.push("/dashboard");
    } else {
      // show monitor
      this.$router.push("/dashboard/monitor");
    }
  },
};
</script>

<style scoped>
.monitor-switch {
  position: fixed;
  top: 5rem;
  left: 5rem;
}
</style>

<template>
  <div class="dashboard">
    <Logout />
    <MonitorSwitch />
    <div class="dashboard-wrapper drop-shadow-box">
      <UserInfo :username="user.username" :role="user.role" />
      <router-view />
    </div>
  </div>
</template>

<script>
import UserInfo from "../components/UserInfo";
import Logout from "../components/Logout";
import MonitorSwitch from "../components/MonitorSwitch";

import { parseJwt } from "../util/util";
export default {
  components: {
    UserInfo,
    Logout,
    MonitorSwitch,
  },
  data() {
    return {
      user: null,
    };
  },
  created() {
    var jwt = localStorage.getItem("token");
    if (jwt === null) {
      this.$message.error("You are not logged in!");
      this.$router.push("/");
    }
    var json = parseJwt(jwt);
    this.user = json;
  },
};
</script>

<style scoped>
.dashboard {
  height: 95vh;
  max-height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.dashboard-wrapper {
  width: 90vw;
  height: 80vh;
  border: 1px solid #dcdfe6;
  border-radius: 20px;
  padding: 2vw;
}
.drop-shadow-box {
  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;
}
</style>

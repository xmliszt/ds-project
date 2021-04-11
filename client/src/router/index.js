import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";
import Register from "../views/Register.vue";
import Dashboard from "../views/Dashboard.vue";
import SecretTable from "../components/SecretTable.vue";
import Monitor from "../components/Monitor.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    component: Home,
  },
  {
    path: "/register",
    component: Register,
  },
  {
    path: "/dashboard",
    component: Dashboard,
    children: [
      {
        path: "",
        component: SecretTable,
      },
      {
        path: "monitor",
        component: Monitor,
      },
    ],
  },
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

export default router;

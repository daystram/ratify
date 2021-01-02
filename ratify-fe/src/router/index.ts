import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import {
  ApplicationList,
  ApplicationDetail,
  Authorize,
  Home,
  Manage,
  Placeholder,
  Profile,
  Signup,
  User
} from "@/views";
import {
  authenticatedOnly,
  callback,
  login,
  logout,
  unAuthenticatedOnly
} from "@/auth";
import Dashboard from "@/views/manage/Dashboard.vue";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "home",
    component: Home
  },
  {
    path: "/authorize",
    name: "authorize",
    component: Authorize
  },
  {
    path: "/manage",
    beforeEnter: authenticatedOnly,
    component: Manage,
    redirect: "/manage/dashboard",
    children: [
      {
        path: "dashboard",
        name: "manage:dashboard",
        component: Dashboard
      },
      {
        path: "profile",
        name: "manage:profile",
        component: Profile
      },
      {
        path: "incident",
        name: "manage:incident",
        component: Placeholder
      },
      {
        path: "session",
        name: "manage:session",
        component: Placeholder
      },
      {
        path: "setting",
        name: "manage:setting",
        component: Placeholder
      },
      {
        path: "user",
        name: "manage:user",
        component: User
      },
      {
        path: "application",
        name: "manage:application",
        component: ApplicationList
      },
      {
        path: "application/:clientId",
        name: "manage:application-detail",
        component: ApplicationDetail
      },
      {
        path: "log",
        name: "manage:log",
        component: Placeholder
      },
      { path: "*", redirect: "/" }
    ]
  },
  {
    path: "/login",
    name: "login",
    beforeEnter: unAuthenticatedOnly,
    component: login
  },
  {
    path: "/logout",
    name: "logout",
    component: logout
  },
  {
    path: "/signup",
    name: "signup",
    beforeEnter: unAuthenticatedOnly,
    component: Signup
  },
  {
    path: "/callback",
    name: "callback",
    beforeEnter: unAuthenticatedOnly,
    component: callback
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;

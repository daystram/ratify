import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import { Authorize, Home, Manage, User } from "@/views";
import {
  authenticatedOnly,
  callback,
  login,
  logout,
  unAuthenticatedOnly
} from "@/auth";

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
    name: "manage",
    beforeEnter: authenticatedOnly,
    component: Manage,
    redirect: "/manage/profile",
    children: [
      {
        path: "profile",
        name: "manage:user",
        component: User
      }
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
    path: "/callback",
    name: "callback",
    component: callback
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;

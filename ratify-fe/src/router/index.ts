import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import { Authorize, Home, Manage, Profile, Signup, User } from "@/views";
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
    beforeEnter: authenticatedOnly,
    component: Manage,
    redirect: "/manage/dashboard",
    children: [
      {
        path: "dashboard",
        name: "manage:dashboard",
        component: Home
      },
      {
        path: "profile",
        name: "manage:profile",
        component: Profile
      },
      {
        path: "incident",
        name: "manage:incident",
        component: User
      },
      {
        path: "session",
        name: "manage:session",
        component: User
      },
      {
        path: "setting",
        name: "manage:setting",
        component: User
      },
      {
        path: "user",
        name: "manage:user",
        component: User
      },
      {
        path: "application",
        name: "manage:application",
        component: User
      },
      {
        path: "log",
        name: "manage:log",
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
    path: "/signup",
    name: "signup",
    component: Signup
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

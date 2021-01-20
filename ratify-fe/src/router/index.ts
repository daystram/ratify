import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import {
  Activity,
  ApplicationList,
  ApplicationDetail,
  Authorize,
  Home,
  Log,
  Manage,
  Placeholder,
  Profile,
  Signup,
  User,
  Verify
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
    component: Home,
    meta: {
      title: "Ratify"
    }
  },
  {
    path: "/authorize",
    name: "authorize",
    component: Authorize,
    meta: {
      title: "Authorize | Ratify"
    }
  },
  {
    path: "/manage",
    component: Manage,
    redirect: "/manage/dashboard",
    children: [
      {
        path: "dashboard",
        name: "manage:dashboard",
        beforeEnter: authenticatedOnly,
        component: Dashboard,
        meta: {
          title: "Dashboard | Ratify"
        }
      },
      {
        path: "profile",
        name: "manage:profile",
        beforeEnter: authenticatedOnly,
        component: Profile,
        meta: {
          title: "Profile | Ratify"
        }
      },
      {
        path: "activity",
        name: "manage:activity",
        beforeEnter: authenticatedOnly,
        component: Activity,
        meta: {
          title: "Activity | Ratify"
        }
      },
      {
        path: "session",
        name: "manage:session",
        beforeEnter: authenticatedOnly,
        component: Placeholder,
        meta: {
          title: "WIP:Sessions | Ratify"
        }
      },
      {
        path: "setting",
        name: "manage:setting",
        beforeEnter: authenticatedOnly,
        component: Placeholder,
        meta: {
          title: "WIP:Settings | Ratify"
        }
      },
      {
        path: "user",
        name: "manage:user",
        beforeEnter: authenticatedOnly,
        component: User,
        meta: {
          title: "WIP:Users | Ratify"
        }
      },
      {
        path: "application",
        name: "manage:application",
        beforeEnter: authenticatedOnly,
        component: ApplicationList,
        meta: {
          title: "Applications | Ratify"
        }
      },
      {
        path: "application/:clientId",
        name: "manage:application-detail",
        beforeEnter: authenticatedOnly,
        component: ApplicationDetail,
        meta: {
          title: "Application Detail | Ratify"
        }
      },
      {
        path: "log",
        name: "manage:log",
        beforeEnter: authenticatedOnly,
        component: Log,
        meta: {
          title: "Logs | Ratify"
        }
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
    component: logout(false)
  },
  {
    path: "/logout-global",
    name: "logout:global",
    component: logout(true)
  },
  {
    path: "/signup",
    name: "signup",
    beforeEnter: unAuthenticatedOnly,
    component: Signup,
    meta: {
      title: "Signup | Ratify"
    }
  },
  {
    path: "/verify",
    name: "verify",
    beforeEnter: unAuthenticatedOnly,
    component: Verify,
    meta: {
      title: "Verify | Ratify"
    }
  },
  {
    path: "/callback",
    name: "callback",
    beforeEnter: unAuthenticatedOnly,
    component: callback
  },
  { path: "*", redirect: { name: "home", query: {} } }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

router.beforeEach((to, from, next) => {
  document.title = to.meta.title || "Ratify";
  next();
});

export default router;

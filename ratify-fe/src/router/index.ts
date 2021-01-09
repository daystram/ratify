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
    beforeEnter: authenticatedOnly,
    component: Manage,
    redirect: "/manage/dashboard",
    children: [
      {
        path: "dashboard",
        name: "manage:dashboard",
        component: Dashboard,
        meta: {
          title: "Dashboard | Ratify"
        }
      },
      {
        path: "profile",
        name: "manage:profile",
        component: Profile,
        meta: {
          title: "Profile | Ratify"
        }
      },
      {
        path: "incident",
        name: "manage:incident",
        component: Placeholder,
        meta: {
          title: "WIP:Incidents | Ratify"
        }
      },
      {
        path: "session",
        name: "manage:session",
        component: Placeholder,
        meta: {
          title: "WIP:Sessions | Ratify"
        }
      },
      {
        path: "setting",
        name: "manage:setting",
        component: Placeholder,
        meta: {
          title: "WIP:Settings | Ratify"
        }
      },
      {
        path: "user",
        name: "manage:user",
        component: User,
        meta: {
          title: "WIP:Users | Ratify"
        }
      },
      {
        path: "application",
        name: "manage:application",
        component: ApplicationList,
        meta: {
          title: "Applications | Ratify"
        }
      },
      {
        path: "application/:clientId",
        name: "manage:application-detail",
        component: ApplicationDetail,
        meta: {
          title: "Application Detail | Ratify"
        }
      },
      {
        path: "log",
        name: "manage:log",
        component: Placeholder,
        meta: {
          title: "WIP:Logs | Ratify"
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
    path: "/callback",
    name: "callback",
    beforeEnter: unAuthenticatedOnly,
    component: callback
  },
  { path: "*", redirect: "/" }
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

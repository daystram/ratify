import { ACCESS_TOKEN, RatifyClient } from "@daystram/ratify-client";
import router from "@/router";
import { Route } from "vue-router";

const CLIENT_ID = process.env.VUE_APP_CLIENT_ID;
const ISSUER = process.env.VUE_APP_OAUTH_ISSUER;
const REDIRECT_URI = `${location.origin}/callback`;

const authManager = new RatifyClient({
  clientId: CLIENT_ID,
  redirectUri: REDIRECT_URI,
  issuer: ISSUER,
  storage: localStorage
});

const login = function() {
  authManager.authorize();
};

const logout = function(global: boolean) {
  return function() {
    authManager.logout(global).finally(() => {
      router.replace({ name: "home" });
    });
  };
};

/*
// logout implementation for clients:
const logout = function() {
  authManager.logout().then(() => {
    router.replace({ name: "home" });
  });
};
 */

const callback = function() {
  const params = new URLSearchParams(document.location.search);
  const code = params.get("code");
  const state = params.get("state");
  if (!code || !state || !authManager.checkState(state)) {
    router.replace("/");
    return;
  }
  authManager
    .redeemToken(code)
    .then(() => {
      const lastRoute = sessionStorage.getItem("last_route");
      if (lastRoute) {
        sessionStorage.removeItem("last_route");
        router.replace({
          path: lastRoute?.toString()
        });
      } else {
        router.replace({
          name: "manage:dashboard"
        });
      }
    })
    .catch(() => {
      router.replace({ name: "home" });
    });
};

const refreshAuth = function(destinationPath: string) {
  sessionStorage.setItem("last_route", destinationPath);
  authManager.reset();
  authManager.authorize(true);
};

const authenticatedOnly = function(to: Route, from: Route, next: () => void) {
  if (authManager.getToken(ACCESS_TOKEN)) {
    next();
  } else {
    refreshAuth(to.fullPath);
  }
};

const unAuthenticatedOnly = function(
  to: object,
  from: object,
  next: () => void
) {
  if (!authManager.getToken(ACCESS_TOKEN)) {
    next();
  } else {
    router.push({ name: "manage:dashboard" });
  }
};

export {
  authManager,
  login,
  logout,
  callback,
  refreshAuth,
  authenticatedOnly,
  unAuthenticatedOnly
};

import router from "@/router";
import { AuthManager, ACCESS_TOKEN } from "@/auth/AuthManager";

const CLIENT_ID = process.env.VUE_APP_CLIENT_ID;
const ISSUER = process.env.VUE_APP_OAUTH_ISSUER;
const REDIRECT_URI = `${location.origin}/callback`;

const authManager = new AuthManager({
  clientId: CLIENT_ID,
  redirectUri: REDIRECT_URI,
  issuer: ISSUER
});

const login = function() {
  authManager.authorize();
};

const logout = function() {
  authManager.revokeToken();
  authManager.reset();
  router.replace({ name: "home" });
};

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
    .then(response => {
      router.replace({ name: "manage:dashboard" });
      console.log(response.data);
    })
    .catch(error => {
      console.error(error.response.data);
      router.replace({ name: "home" });
    });
};

const isAuthenticated = function(): boolean {
  return authManager.getToken(ACCESS_TOKEN) !== "";
};

const user = function() {
  return authManager.getUser();
};

const authenticatedOnly = function(to: any, from: any, next: () => void) {
  if (authManager.getToken(ACCESS_TOKEN)) {
    next();
  } else {
    authManager.reset();
    // TODO: deauth link followup, store @ localStorage
    router.push({ name: "login" });
  }
};

const unAuthenticatedOnly = function(to: any, from: any, next: () => void) {
  if (!authManager.getToken(ACCESS_TOKEN)) {
    next();
  } else {
    router.push({ name: "manage:dashboard" });
  }
};

export {
  isAuthenticated,
  user,
  login,
  logout,
  callback,
  authenticatedOnly,
  unAuthenticatedOnly
};

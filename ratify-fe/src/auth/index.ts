import router from "@/router";
import { AuthManager, ACCESS_TOKEN, MemoryStorage } from "@/auth/AuthManager";

const CLIENT_ID = process.env.VUE_APP_CLIENT_ID;
const ISSUER = process.env.VUE_APP_OAUTH_ISSUER;
const REDIRECT_URI = `${location.origin}/callback`;

const authManager = new AuthManager({
  clientId: CLIENT_ID,
  redirectUri: REDIRECT_URI,
  issuer: ISSUER,
  storage: new MemoryStorage()
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
      router.replace({ name: "manage:dashboard" });
    })
    .catch(() => {
      router.replace({ name: "home" });
    });
};

const authenticatedOnly = function(to: object, from: object, next: () => void) {
  if (authManager.getToken(ACCESS_TOKEN)) {
    next();
  } else {
    authManager.reset();
    authManager.authorize(true);
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
  authenticatedOnly,
  unAuthenticatedOnly
};

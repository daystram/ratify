import qs from "qs";
import oauth from "@/apis/oauth";
import pkceChallenge from "pkce-challenge";
import { v4 as uuidv4 } from "uuid";
import { AxiosResponse } from "axios";
import jwtDecode from "jwt-decode";

export const KEY_STATE = "state";
export const KEY_CODE = "code";
export const KEY_TOKEN = "token";
export const ACCESS_TOKEN = "access_token";
export const ID_TOKEN = "id_token";

interface AuthManagerOptions {
  clientId: string;
  redirectUri: string;
  issuer: string;
  storage: TokenStorage;
}

interface User {
  subject: string;
  given_name: string;
  family_name: string;
  preferred_username: string;
  is_superuser: boolean;
}

interface TokenStorage {
  getItem: (key: string) => string | null;
  setItem: (key: string, value: string) => void;
  removeItem: (key: string) => void;
}

export class MemoryStorage implements TokenStorage {
  private tokens: { [key: string]: string };

  constructor() {
    this.tokens = {};
  }

  getItem(key: string): string | null {
    return this.tokens[key];
  }

  setItem(key: string, value: string): void {
    this.tokens[key] = value;
  }

  removeItem(key: string): void {
    delete this.tokens[key];
  }
}

export class AuthManager {
  private options: AuthManagerOptions;
  private storageManager: TokenStorage;

  constructor(opts: AuthManagerOptions) {
    this.options = opts;
    this.storageManager = opts.storage;
    // code and state will still use sessionStorage (need to persist during page reloads)
  }

  isAuthenticated(): boolean {
    return this.getToken(ACCESS_TOKEN) !== "";
  }

  getToken(tokenKey: string): string {
    return (
      JSON.parse(this.storageManager.getItem(KEY_TOKEN) || "{}")[tokenKey] || ""
    );
  }

  getUser(): User {
    return jwtDecode(this.getToken(ID_TOKEN));
  }

  reset() {
    this.storageManager.removeItem(KEY_TOKEN);
  }

  logout(global?: boolean) {
    return oauth
      .logout({
        /* eslint-disable @typescript-eslint/camelcase */
        client_id: this.options.clientId,
        global: global || false
        /* eslint-enable @typescript-eslint/camelcase */
      })
      .then(() => {
        this.reset();
      })
      .catch(() => {
        this.reset();
      });
  }

  authorize(immediate?: boolean, scopes?: string[]): void {
    window.location.href = `${this.options.issuer}/authorize?${qs.stringify({
      /* eslint-disable @typescript-eslint/camelcase */
      client_id: this.options.clientId,
      response_type: "code",
      redirect_uri: this.options.redirectUri,
      scope:
        "openid profile" + (scopes || []).map(scope => " " + scope).join(""),
      state: this.getState(),
      code_challenge: this.getCodeChallenge(),
      code_challenge_method: "S256",
      immediate: immediate || false
      /* eslint-enable @typescript-eslint/camelcase */
    })}`;
  }

  redeemToken(authorizationCode: string): Promise<AxiosResponse> {
    return oauth
      .token({
        /* eslint-disable @typescript-eslint/camelcase */
        client_id: this.options.clientId,
        grant_type: "authorization_code",
        code: authorizationCode,
        code_verifier: this.getCodeVerifier()
        /* eslint-enable @typescript-eslint/camelcase */
      })
      .then(response => {
        this.storageManager.setItem(KEY_TOKEN, JSON.stringify(response.data));
        return response;
      });
  }

  getState(): string {
    const state = uuidv4();
    sessionStorage.setItem(KEY_STATE, state);
    return state;
  }

  checkState(state: string): boolean {
    const temp = sessionStorage.getItem(KEY_STATE);
    sessionStorage.removeItem(KEY_STATE);
    return temp === state;
  }

  getCodeChallenge() {
    const pkce = pkceChallenge();
    sessionStorage.setItem(KEY_CODE, JSON.stringify(pkce));
    return pkce.code_challenge;
  }

  getCodeVerifier() {
    const pkce = JSON.parse(sessionStorage.getItem(KEY_CODE) || "");
    sessionStorage.removeItem(KEY_CODE);
    return pkce.code_verifier;
  }
}

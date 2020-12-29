import qs from "qs";
import oauth from "@/apis/oauth";
import pkceChallenge from "pkce-challenge";
import { v4 as uuidv4 } from "uuid";
import { AxiosResponse } from "axios";

export const KEY_STATE = "state";
export const KEY_CODE = "code";
export const KEY_ACCESS_TOKEN = "access_token";

interface AuthManagerOptions {
  clientId: string;
  redirectUri: string;
  issuer: string;
}

export class AuthManager {
  private options: AuthManagerOptions;
  private storageManager: Storage;

  constructor(opts: AuthManagerOptions) {
    this.options = opts;
    this.storageManager = localStorage;
  }

  getToken(tokenKey: string): string {
    return this.storageManager.getItem(tokenKey) || "";
  }

  reset() {
    this.storageManager.removeItem(KEY_ACCESS_TOKEN);
    this.storageManager.removeItem(KEY_CODE);
    this.storageManager.removeItem(KEY_STATE);
  }

  authorize(): void {
    window.location.href = `${this.options.issuer}/authorize?${qs.stringify({
      /* eslint-disable @typescript-eslint/camelcase */
      client_id: this.options.clientId,
      response_type: "code",
      redirect_uri: this.options.redirectUri,
      state: this.getState(),
      code_challenge: this.getCodeChallenge(),
      code_challenge_method: "S256"
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
        this.storageManager.setItem(
          KEY_ACCESS_TOKEN,
          response.data[KEY_ACCESS_TOKEN]
        );
        return response;
      });
  }

  revokeToken(): Promise<AxiosResponse> {
    return oauth.revoke({
      /* eslint-disable @typescript-eslint/camelcase */
      client_id: this.options.clientId,
      token: this.getToken(KEY_ACCESS_TOKEN)
      /* eslint-enable @typescript-eslint/camelcase */
    });
  }

  getState(): string {
    const state = uuidv4();
    this.storageManager.setItem(KEY_STATE, state);
    return state;
  }

  checkState(state: string): boolean {
    const temp = localStorage.getItem(KEY_STATE);
    this.storageManager.removeItem(KEY_STATE);
    return temp === state;
  }

  getCodeChallenge() {
    const pkce = pkceChallenge();
    this.storageManager.setItem(KEY_CODE, JSON.stringify(pkce));
    return pkce.code_challenge;
  }

  getCodeVerifier() {
    const pkce = JSON.parse(this.storageManager.getItem(KEY_CODE) || "");
    this.storageManager.removeItem(KEY_CODE);
    return pkce.code_verifier;
  }
}

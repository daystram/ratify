import axios, { AxiosInstance, AxiosResponse } from "axios";
import qs from "qs";

const oauthClient: AxiosInstance = axios.create({
  baseURL: "/oauth/"
});

export default {
  authorize: function(authRequest: object): Promise<AxiosResponse> {
    return oauthClient.post(`authorize`, authRequest);
  },
  token: function(tokenRequest: object): Promise<AxiosResponse> {
    return oauthClient.post(`token`, qs.stringify(tokenRequest), {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    });
  },
  revoke: function(revokeRequest: object): Promise<AxiosResponse> {
    return oauthClient.post(`revoke`, revokeRequest);
  }
};

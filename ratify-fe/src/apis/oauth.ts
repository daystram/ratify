import axios, { AxiosInstance, AxiosResponse } from "axios";

const oauthClient: AxiosInstance = axios.create({
  baseURL: "/oauth/"
});

export default {
  authorize: function(authRequest: object): Promise<AxiosResponse> {
    return oauthClient.post(`authorize`, authRequest);
  }
};

import axios, { AxiosInstance, AxiosResponse } from "axios";
import { authManager } from "@/auth";
import { ACCESS_TOKEN } from "@/auth/AuthManager";

const apiClient: AxiosInstance = axios.create({
  baseURL: "/api/v1/"
});

export default {
  application: {
    getOne: function(clientId: string): Promise<AxiosResponse> {
      return apiClient.get(`application/${clientId}`);
    },
    getAll: function(): Promise<AxiosResponse> {
      return apiClient.get(`application`, {
        headers: {
          Authorization: `Bearer ${authManager.getToken(ACCESS_TOKEN)}`
        }
      });
    }
  },
  form: {
    checkUnique: function(uniqueRequest: object): Promise<AxiosResponse> {
      return apiClient.post("form/unique", uniqueRequest);
    }
  },
  user: {
    detail: function(subject?: string): Promise<AxiosResponse> {
      return apiClient.get(`user/${subject || ""}`, {
        headers: {
          Authorization: `Bearer ${authManager.getToken(ACCESS_TOKEN)}`
        }
      });
    },
    update: function(user: object): Promise<AxiosResponse> {
      return apiClient.put(`user`, user, {
        headers: {
          Authorization: `Bearer ${authManager.getToken(ACCESS_TOKEN)}`
        }
      });
    },
    signup: function(userSignup: object): Promise<AxiosResponse> {
      return apiClient.post("user", userSignup);
    }
  }
};

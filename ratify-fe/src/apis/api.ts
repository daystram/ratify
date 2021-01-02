import axios, { AxiosInstance, AxiosResponse } from "axios";
import { authManager } from "@/auth";
import { ACCESS_TOKEN } from "@/auth/AuthManager";

const apiClient: AxiosInstance = axios.create({
  baseURL: "/api/v1/"
});

const withAuth = {
  headers: {
    Authorization: `Bearer ${authManager.getToken(ACCESS_TOKEN)}`
  }
};

export default {
  application: {
    detail: function(
      clientId: string,
      complete?: boolean
    ): Promise<AxiosResponse> {
      return apiClient.get(`application/${clientId}`, complete ? withAuth : {});
    },
    update: function(
      clientId: string,
      application: object
    ): Promise<AxiosResponse> {
      return apiClient.put(`application/${clientId}`, application, withAuth);
    },
    delete: function(clientId: string): Promise<AxiosResponse> {
      return apiClient.delete(`application/${clientId}`, withAuth);
    },
    list: function(): Promise<AxiosResponse> {
      return apiClient.get(`application/`, withAuth);
    },
    register: function(application: object): Promise<AxiosResponse> {
      return apiClient.post(`application`, application, withAuth);
    }
  },
  form: {
    checkUnique: function(uniqueRequest: object): Promise<AxiosResponse> {
      return apiClient.post(`form/unique`, uniqueRequest);
    }
  },
  user: {
    detail: function(subject?: string): Promise<AxiosResponse> {
      return apiClient.get(`user/${subject || ""}`, withAuth);
    },
    update: function(user: object): Promise<AxiosResponse> {
      return apiClient.put(`user`, user, withAuth);
    },
    signup: function(userSignup: object): Promise<AxiosResponse> {
      return apiClient.post("user", userSignup);
    }
  }
};

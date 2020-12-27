import axios, { AxiosInstance, AxiosResponse } from "axios";

const apiClient: AxiosInstance = axios.create({
  baseURL: "api/v1/"
});

export default {
  application: {
    getOne: function(clientId: string): Promise<AxiosResponse> {
      return apiClient.get(`application/${clientId}`);
    }
  }
};

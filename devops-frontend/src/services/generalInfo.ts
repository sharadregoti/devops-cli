import api from "./index";
import { AxiosError, AxiosResponse } from "axios";
import { InfoCardTypes } from "../types/InfoCardTypes";

export const getGeneralInfoService = (pluginName: string, authId: string, contextId: string) => {
  return new Promise((resolve, reject) => {
    api.get<InfoCardTypes>(`/connect/${pluginName}/${authId}/${contextId}`).then((res: AxiosResponse) => {
      if(res.statusText === 'OK'){
        resolve(res.data);
      }
    }).catch((ex: AxiosError) => reject(ex))
  })
}

export const getAppConfigService = () => {
  return new Promise((resolve, reject) => {
    api.get<InfoCardTypes>('/config').then((res: AxiosResponse) => {
      if(res.statusText === 'OK'){
        resolve(res.data);
      }
    }).catch((ex: AxiosError) => reject(ex))
  })
}

export const getPluginAuthService = (pluginName: string) => {
  return new Promise((resolve, reject) => {
    api.get<InfoCardTypes>(`/auth/${pluginName}`).then((res: AxiosResponse) => {
      if(res.statusText === 'OK'){
        resolve(res.data);
      }
    }).catch((ex: AxiosError) => reject(ex))
  })
}

export const sendEvent = (id: string, data: any) => {
  return new Promise((resolve, reject) => {
    api.post<InfoCardTypes>(`/events/${id}`, data).then((res: AxiosResponse) => {
      if(res.statusText === 'OK'){
        resolve(res.data);
      }
    }).catch((ex: AxiosError) => reject(ex))
  })
}
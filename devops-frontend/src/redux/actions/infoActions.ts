import { createAsyncThunk } from '@reduxjs/toolkit';
import { getGeneralInfoService, getAppConfigService, getPluginAuthService } from '../../services/generalInfo';
import { notify } from '../../utils/utils';

export const getGeneralInfoAction = createAsyncThunk(
  'generalInfo/getGeneralInfoAction',
  async () => {
    return getGeneralInfoService().then(res => {
      return res;
    }).catch((ex: string) => notify('error', 'Error in getting general info', ex));
  }
);

export const getAppConfig = createAsyncThunk(
  'generalInfo/getAppConfig',
  async () => {
    return getAppConfigService().then(res => {
      return res;
    }).catch((ex: string) => notify('error', 'Error in getting app config', ex));
  }
);

export const getPluginAuth = createAsyncThunk(
  'generalInfo/getPluginAuth',
  async (pluginName: string) => {
    return getPluginAuthService(pluginName).then(res => {
      return res;
    }).catch((ex: string) => notify('error', 'Error in getting plugin auth', ex));
  }
);
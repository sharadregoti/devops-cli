import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { InfoCardTypes } from '../../types/InfoCardTypes';
import { failed, loading, success } from '../../utils/constants';
import { getGeneralInfoAction } from '../actions/infoActions';

interface InitialStateTypes {
  info: InfoCardTypes | null,
  status: string | null
}
const initialState= {
  info: null,
  status: null
} as InitialStateTypes

export const generalInfoSlice = createSlice({
  name: 'generalInfo',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(getGeneralInfoAction.pending, (state) => {
      state.status = loading;
    })
    .addCase(getGeneralInfoAction.fulfilled, (state, action: PayloadAction<any>) => {
      state.status = success;
      state.info = action.payload;
    })
    .addCase(getGeneralInfoAction.rejected, (state) => {
      state.status = failed;
    });
  }
});

export default generalInfoSlice.reducer;

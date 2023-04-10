import { configureStore } from "@reduxjs/toolkit";
import generalInfoReducer from './reducers/infoReducer';

const store = configureStore({
  reducer: {
    generalInfo: generalInfoReducer 
  },
  devTools: true
});

export default store;

export type StoreType = ReturnType<typeof store.getState>;

export type DispatchType = typeof store.dispatch;
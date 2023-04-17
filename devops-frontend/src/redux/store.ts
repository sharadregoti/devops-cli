import { configureStore } from "@reduxjs/toolkit";
import { combineReducers } from 'redux';
import generalInfoReducer from './reducers/infoReducer';
import { ModelConfig } from "../generated-sources/openapi";
import { navBarReducer, pluginSelectorReducer } from "./reducers/PluginSelectorReducer";
import { homeReducer } from "./reducers/Home";

const rootReducer = combineReducers({
  pluginSelector: pluginSelectorReducer,
  navBar: navBarReducer,
  home: homeReducer,
});

const store = configureStore({
  reducer: rootReducer,
  devTools: true
});

export default store;

export type StoreType = ReturnType<typeof store.getState>;

export type DispatchType = typeof store.dispatch;
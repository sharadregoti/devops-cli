export interface InfoCardTypes {
  id: string;
  general: General;
  plugins: Plugins;
  actions: Action[];
  resourceTypes: string[];
  defaultIsolator: string;
  isolatorType: string;
}

export interface Action {
  type: string;
  name: string;
  key_binding: string;
  output_type: string;
  schema?: any;
}

export interface Plugins {
  'alt-0': string;
}

export interface General {
  Cluster: string;
  Context: string;
  'Server Version': string;
  User: string;
}
import { ModelConfig, ModelInfoResponse } from "../../generated-sources/openapi";


export type PluginSelectorState = {
    pluginSelector: {
        pluginName: string;
        serverConfig: ModelConfig,
        pluginAuthData: TableDataType[];
    };
};

export interface TableDataType {
    key: string;
    name: string;
    context: string;
}

export const pluginSelectorReducer = (state = { pluginName: "", serverConfig: {}, pluginAuthData: [] }, action) => {
    switch (action.type) {
        case 'SET_PLUGIN_SELECTOR_STATE':
            return { ...state, ...action.payload };
        default:
            return state;
    }
};

export type NavBarState = {
    navBar: {
        items: NavBarItem[];
    };
};

export type NavBarItem = {
    pluginName: string;
    authId: string;
    contextId: string;
    sessionId: string;
    generalInfo: ModelInfoResponse
}

export const navBarReducer = (state = { items: [{ pluginName: "", sessionId: "0", authId: "", contextId: "Plugins", generalInfo: {} } as NavBarItem] }, action) => {
    switch (action.type) {
        case 'SET_NAV_BAR_STATE':
            return { ...state, ...action.payload };
        default:
            return state;
    }
};


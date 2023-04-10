// Constants
export const settingPlugin = "@use.plugin";
export const settingAuthentication = "@use.auth";

// A function which returns a string for use.plugin
export function getPluginSetting(name: string) {
    return `${settingPlugin}.${name}`;
}

// A function to parse the plugin setting and extract the plugin name
export function parsePluginSetting(setting: string) {
    return setting.split(".")[2];
}

// A function to parse the authentication setting and extract identifyingName and name
export function parseAuthenticationSetting(setting: string) {
    const arr = setting.split(".");
    return {
        "identifyingName": arr[2],
        "name": arr[3]
    }
}

// A function which returns a string for authentication setting
export function getAuthenticationSetting(identifyingName: string, name: string) {
    return `${settingAuthentication}.${identifyingName}.${name}`;
}

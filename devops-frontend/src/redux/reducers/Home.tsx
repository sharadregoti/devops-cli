import { DrawerPropsTypes } from "../../components/drawer/Drawer";
import { InfoCardPropsTypes } from "../../components/isolator/Isolator";
import { SpecificActionFormProps } from "../../components/specificActionForm/SpecificActionForm";
import { ModelConfig } from "../../generated-sources/openapi";
import { AppState } from "../../types/Event";


export type HomeState = {
    home: Record<string, {
        recentlyUsedItems: string[],
        isolatorsList: string[],
        appConfig: AppState,
        drawerState: DrawerPropsTypes,
        specificActionFormState: SpecificActionFormProps,
        isolatorCardState: InfoCardPropsTypes
    }>;
};

export const homeReducer = (
    state: HomeState['home'] = {},
    action: { type: string; payload: any; key: string }
) => {
    switch (action.type) {
        case 'ADD_ISOLATOR': {
            const { key, isolatorName } = action;
            const currentIsolatorsList = state[key]?.isolatorsList || [];

            // Check if isolatorName already exists in the currentIsolatorsList
            if (currentIsolatorsList.includes(isolatorName)) {
                return state; // Return the unchanged state if the isolator already exists
            }

            return {
                ...state,
                [key]: {
                    ...state[key],
                    isolatorsList: [...currentIsolatorsList, isolatorName],
                },
            };
        }
        case 'SET_HOME_STATE':
            return {
                ...state,
                [action.key]: {
                    ...(state[action.key] || {
                        recentlyUsedItems: [],
                        isolatorsList: [],
                        drawerState: {},
                        specificActionFormState: {
                            formItems: {},
                        } as SpecificActionFormProps,
                        isolatorCardState: {} as InfoCardPropsTypes,
                    }),
                    ...action.payload,
                },
            };
        default:
            return state;
    }
};


// export type NavBarState = {
//     navBar: {
//         items: NavBarItem[];
//     };
// };

// export type NavBarItem = {
//     pluginName: string;
//     authId: string;
//     contextId: string;
//     sessionId: string;
// }

// export const navBarReducer = (state = { items: [{ pluginName: "", sessionId: "0", authId: "", contextId: "Plugins" } as NavBarItem] }, action) => {
//     switch (action.type) {
//         case 'SET_NAV_BAR_STATE':
//             return { ...state, ...action.payload };
//         default:
//             return state;
//     }
// };


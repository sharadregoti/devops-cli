/* tslint:disable */
/* eslint-disable */
/**
 * NEW Devops API
 * Devops API Sec
 *
 * The version of the OpenAPI document: v0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
import type { ProtoServerInput } from './ProtoServerInput';
import {
    ProtoServerInputFromJSON,
    ProtoServerInputFromJSONTyped,
    ProtoServerInputToJSON,
} from './ProtoServerInput';
import type { ProtoUserInput } from './ProtoUserInput';
import {
    ProtoUserInputFromJSON,
    ProtoUserInputFromJSONTyped,
    ProtoUserInputToJSON,
} from './ProtoUserInput';

/**
 * 
 * @export
 * @interface ProtoExecution
 */
export interface ProtoExecution {
    /**
     * 
     * @type {string}
     * @memberof ProtoExecution
     */
    cmd?: string;
    /**
     * 
     * @type {boolean}
     * @memberof ProtoExecution
     */
    isLongRunning?: boolean;
    /**
     * 
     * @type {ProtoServerInput}
     * @memberof ProtoExecution
     */
    serverInput?: ProtoServerInput;
    /**
     * 
     * @type {ProtoUserInput}
     * @memberof ProtoExecution
     */
    userInput?: ProtoUserInput;
}

/**
 * Check if a given object implements the ProtoExecution interface.
 */
export function instanceOfProtoExecution(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ProtoExecutionFromJSON(json: any): ProtoExecution {
    return ProtoExecutionFromJSONTyped(json, false);
}

export function ProtoExecutionFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProtoExecution {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'cmd': !exists(json, 'cmd') ? undefined : json['cmd'],
        'isLongRunning': !exists(json, 'is_long_running') ? undefined : json['is_long_running'],
        'serverInput': !exists(json, 'server_input') ? undefined : ProtoServerInputFromJSON(json['server_input']),
        'userInput': !exists(json, 'user_input') ? undefined : ProtoUserInputFromJSON(json['user_input']),
    };
}

export function ProtoExecutionToJSON(value?: ProtoExecution | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'cmd': value.cmd,
        'is_long_running': value.isLongRunning,
        'server_input': ProtoServerInputToJSON(value.serverInput),
        'user_input': ProtoUserInputToJSON(value.userInput),
    };
}

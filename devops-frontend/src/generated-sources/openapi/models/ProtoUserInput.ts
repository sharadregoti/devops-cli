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
import type { StructpbValue } from './StructpbValue';
import {
    StructpbValueFromJSON,
    StructpbValueFromJSONTyped,
    StructpbValueToJSON,
} from './StructpbValue';

/**
 * 
 * @export
 * @interface ProtoUserInput
 */
export interface ProtoUserInput {
    /**
     * 
     * @type {{ [key: string]: StructpbValue; }}
     * @memberof ProtoUserInput
     */
    args?: { [key: string]: StructpbValue; };
    /**
     * 
     * @type {boolean}
     * @memberof ProtoUserInput
     */
    required?: boolean;
}

/**
 * Check if a given object implements the ProtoUserInput interface.
 */
export function instanceOfProtoUserInput(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ProtoUserInputFromJSON(json: any): ProtoUserInput {
    return ProtoUserInputFromJSONTyped(json, false);
}

export function ProtoUserInputFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProtoUserInput {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'args': !exists(json, 'args') ? undefined : (mapValues(json['args'], StructpbValueFromJSON)),
        'required': !exists(json, 'required') ? undefined : json['required'],
    };
}

export function ProtoUserInputToJSON(value?: ProtoUserInput | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'args': value.args === undefined ? undefined : (mapValues(value.args, StructpbValueToJSON)),
        'required': value.required,
    };
}

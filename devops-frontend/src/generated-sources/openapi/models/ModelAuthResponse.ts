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
import type { ProtoAuthInfo } from './ProtoAuthInfo';
import {
    ProtoAuthInfoFromJSON,
    ProtoAuthInfoFromJSONTyped,
    ProtoAuthInfoToJSON,
} from './ProtoAuthInfo';

/**
 * 
 * @export
 * @interface ModelAuthResponse
 */
export interface ModelAuthResponse {
    /**
     * 
     * @type {Array<ProtoAuthInfo>}
     * @memberof ModelAuthResponse
     */
    auths?: Array<ProtoAuthInfo>;
}

/**
 * Check if a given object implements the ModelAuthResponse interface.
 */
export function instanceOfModelAuthResponse(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ModelAuthResponseFromJSON(json: any): ModelAuthResponse {
    return ModelAuthResponseFromJSONTyped(json, false);
}

export function ModelAuthResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): ModelAuthResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'auths': !exists(json, 'auths') ? undefined : ((json['auths'] as Array<any>).map(ProtoAuthInfoFromJSON)),
    };
}

export function ModelAuthResponseToJSON(value?: ModelAuthResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'auths': value.auths === undefined ? undefined : ((value.auths as Array<any>).map(ProtoAuthInfoToJSON)),
    };
}


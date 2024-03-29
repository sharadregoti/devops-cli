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
/**
 * 
 * @export
 * @interface ProtoAuthInfo
 */
export interface ProtoAuthInfo {
    /**
     * 
     * @type {Array<string>}
     * @memberof ProtoAuthInfo
     */
    defaultIsolators?: Array<string>;
    /**
     * 
     * @type {string}
     * @memberof ProtoAuthInfo
     */
    identifyingName?: string;
    /**
     * 
     * @type {{ [key: string]: string; }}
     * @memberof ProtoAuthInfo
     */
    info?: { [key: string]: string; };
    /**
     * 
     * @type {boolean}
     * @memberof ProtoAuthInfo
     */
    isDefault?: boolean;
    /**
     * 
     * @type {string}
     * @memberof ProtoAuthInfo
     */
    name?: string;
    /**
     * 
     * @type {string}
     * @memberof ProtoAuthInfo
     */
    path?: string;
}

/**
 * Check if a given object implements the ProtoAuthInfo interface.
 */
export function instanceOfProtoAuthInfo(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ProtoAuthInfoFromJSON(json: any): ProtoAuthInfo {
    return ProtoAuthInfoFromJSONTyped(json, false);
}

export function ProtoAuthInfoFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProtoAuthInfo {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'defaultIsolators': !exists(json, 'default_isolators') ? undefined : json['default_isolators'],
        'identifyingName': !exists(json, 'identifying_name') ? undefined : json['identifying_name'],
        'info': !exists(json, 'info') ? undefined : json['info'],
        'isDefault': !exists(json, 'is_default') ? undefined : json['is_default'],
        'name': !exists(json, 'name') ? undefined : json['name'],
        'path': !exists(json, 'path') ? undefined : json['path'],
    };
}

export function ProtoAuthInfoToJSON(value?: ProtoAuthInfo | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'default_isolators': value.defaultIsolators,
        'identifying_name': value.identifyingName,
        'info': value.info,
        'is_default': value.isDefault,
        'name': value.name,
        'path': value.path,
    };
}


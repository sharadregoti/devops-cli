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
 * @interface ModelErrorResponse
 */
export interface ModelErrorResponse {
    /**
     * 
     * @type {string}
     * @memberof ModelErrorResponse
     */
    message?: string;
}

/**
 * Check if a given object implements the ModelErrorResponse interface.
 */
export function instanceOfModelErrorResponse(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ModelErrorResponseFromJSON(json: any): ModelErrorResponse {
    return ModelErrorResponseFromJSONTyped(json, false);
}

export function ModelErrorResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): ModelErrorResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'message': !exists(json, 'message') ? undefined : json['message'],
    };
}

export function ModelErrorResponseToJSON(value?: ModelErrorResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'message': value.message,
    };
}


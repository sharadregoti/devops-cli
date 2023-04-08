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
import type { ProtoExecution } from './ProtoExecution';
import {
    ProtoExecutionFromJSON,
    ProtoExecutionFromJSONTyped,
    ProtoExecutionToJSON,
} from './ProtoExecution';
import type { StructpbValue } from './StructpbValue';
import {
    StructpbValueFromJSON,
    StructpbValueFromJSONTyped,
    StructpbValueToJSON,
} from './StructpbValue';

/**
 * 
 * @export
 * @interface ProtoAction
 */
export interface ProtoAction {
    /**
     * 
     * @type {{ [key: string]: StructpbValue; }}
     * @memberof ProtoAction
     */
    args?: { [key: string]: StructpbValue; };
    /**
     * 
     * @type {ProtoExecution}
     * @memberof ProtoAction
     */
    execution?: ProtoExecution;
    /**
     * 
     * @type {string}
     * @memberof ProtoAction
     */
    keyBinding?: string;
    /**
     * 
     * @type {string}
     * @memberof ProtoAction
     */
    name?: string;
    /**
     * 
     * @type {string}
     * @memberof ProtoAction
     */
    outputType?: string;
    /**
     * 
     * @type {{ [key: string]: StructpbValue; }}
     * @memberof ProtoAction
     */
    schema?: { [key: string]: StructpbValue; };
}

/**
 * Check if a given object implements the ProtoAction interface.
 */
export function instanceOfProtoAction(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ProtoActionFromJSON(json: any): ProtoAction {
    return ProtoActionFromJSONTyped(json, false);
}

export function ProtoActionFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProtoAction {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'args': !exists(json, 'args') ? undefined : (mapValues(json['args'], StructpbValueFromJSON)),
        'execution': !exists(json, 'execution') ? undefined : ProtoExecutionFromJSON(json['execution']),
        'keyBinding': !exists(json, 'key_binding') ? undefined : json['key_binding'],
        'name': !exists(json, 'name') ? undefined : json['name'],
        'outputType': !exists(json, 'output_type') ? undefined : json['output_type'],
        'schema': !exists(json, 'schema') ? undefined : (mapValues(json['schema'], StructpbValueFromJSON)),
    };
}

export function ProtoActionToJSON(value?: ProtoAction | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'args': value.args === undefined ? undefined : (mapValues(value.args, StructpbValueToJSON)),
        'execution': ProtoExecutionToJSON(value.execution),
        'key_binding': value.keyBinding,
        'name': value.name,
        'output_type': value.outputType,
        'schema': value.schema === undefined ? undefined : (mapValues(value.schema, StructpbValueToJSON)),
    };
}


import { ModelAuthResponse, ModelConfig, ModelInfoResponse, ProtoAction } from "../generated-sources/openapi";

export interface TableEvent {
    eventType: string;
    record: any;
}

export interface AppState {
    serverConfig: ModelConfig,
    pluginAuth: ModelAuthResponse,
    generalInfo: ModelInfoResponse,
    currentIsolator: string,
    currentResourceType: string,
    currentPluginName: string
}

export interface WebsocketData {
    id: string;
    name: string;
    eventType: string;
    data: Datum[];
    specificActions: Array<SpecificAction>;
}

export interface Datum {
    data: string[];
    color: string;
}

export interface SpecificAction {
    name: string;
    key_binding: string;
    output_type: string;
    execution: Execution
}

export interface Execution {
    cmd: string;
    user_input: UserInput;
}

export interface UserInput {
    required: boolean;
    args: object;
}
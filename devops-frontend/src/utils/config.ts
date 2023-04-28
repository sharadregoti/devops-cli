import { Configuration, DefaultApi, ModelConfig, HandleAuthRequest, HandleInfoRequest, HandleEventRequest, ModelFrontendEvent, ModelFrontendEventNameEnum, ProtoAction } from "../generated-sources/openapi";

// http://localhost:9753
// For Prod
export const apiHost = window.location.host;
// export const apiHost = "localhost:9753";
export const httpAPI = `http://${apiHost}`;

const configuration = new Configuration({
    basePath: `${httpAPI}`,
});

export const api = new DefaultApi(configuration);
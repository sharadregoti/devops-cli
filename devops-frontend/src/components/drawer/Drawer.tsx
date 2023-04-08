import React, { useEffect, useRef, useState } from 'react';
import { Drawer, Spin, Button } from 'antd';
import XTermComponent from '../xTerm/XTerm';
import Editor from '@monaco-editor/react';
import { showNotification } from '../../utils/notification';
import { api } from '../../utils/config';
import { AppState } from '../../types/Event'
import { ModelFrontendEvent, ModelFrontendEventNameEnum, HandleEventRequest } from '../../generated-sources/openapi';
import yaml from "yamljs";

type DrawerBodyType = "editor" | "xterm" | "table";

export type DrawerPropsTypes = {
    socketUrl: string
    drawerBodyType: DrawerBodyType
    isDrawerOpen: boolean
    resourceName: string
    appConfig: AppState
    editorOptions: EditorOptions
    onDrawerClose: () => void
}

type EditorOptions = {
    isReadOnly: boolean
    defaultText: string
}

const SideDrawer: React.FC<DrawerPropsTypes> = ({ isDrawerOpen, socketUrl, drawerBodyType, resourceName, appConfig, editorOptions, onDrawerClose }) => {
    console.log('Rendering SideDrawer...');

    const terminalRef = useRef(null);
    const editorRef = useRef(null);

    // const [drawerLoading, setDrawerLoading] = useState(false);

    const handleEditorSaveButton = () => {
        // TODO: Handle spinner
        // setDrawerLoading(true);

        // Get the current content of the editor
        const content = editorRef.current.getValue();

        let yamlobj = {}
        try {
            yamlobj = yaml.parse(content);
        } catch (error) {
            // TODO: Handle spinner
            // setDrawerLoading(false);
            showNotification('error', 'Invalid YAML', error.message)
            return
        }

        const e: ModelFrontendEvent = {
            eventType: "normal-action",
            name: ModelFrontendEventNameEnum.Edit,
            isolatorName: appConfig?.currentIsolator,
            pluginName: appConfig?.currentPluginName,
            resourceName: resourceName,
            resourceType: appConfig?.currentResourceType,
            args: yamlobj,
        }

        let params: HandleEventRequest = {
            id: appConfig.generalInfo.id,
            modelFrontendEvent: e
        }

        api.handleEvent(params)
            .then(res => {
                // setDrawerLoading(false);
                // setOpen(false);
                onDrawerClose()
                showNotification('success', `Successfully updated ${e.resourceType}`, '')
            })
            .catch(err => {
                // setDrawerLoading(false);
                showNotification('error', 'Event invocation failed', err)
            })
    }

    return (
        <Drawer
            keyboard={true}
            width={"80%"}
            title="Basic Drawer"
            placement="right"
            onClose={onDrawerClose}
            open={isDrawerOpen}>
            {/* TODO: Handle spinner */}
            {/* <Spin spinning={drawerLoading}> */}
            {drawerBodyType == 'editor' && <>
                <Editor
                    language="yaml"
                    value={editorOptions.defaultText}
                    options={{ readOnly: editorOptions.isReadOnly }}
                    height={"700px"}
                    onMount={(editor, monaco) => {
                        // Store a reference to the editor instance
                        editorRef.current = editor;
                    }}
                />
                {!editorOptions.isReadOnly &&
                    <Button type="primary" onClick={handleEditorSaveButton} >Save</Button>
                }
            </>}
            {drawerBodyType == 'xterm' &&
                <XTermComponent ref={terminalRef} socketUrl={socketUrl} isDrawerOpen={isDrawerOpen} />
            }
            {/* </Spin> */}
        </Drawer>
    );
};

export default SideDrawer;
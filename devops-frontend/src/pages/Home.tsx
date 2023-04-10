import { AutoComplete, Col, Modal, Row } from 'antd';
import Fuse from 'fuse.js';
import React, { useEffect, useRef, useState } from 'react';
import '../../node_modules/xterm/css/xterm.css';
import InfoCard from '../components/infoCard/InfoCard';
import IsolatorCard, { InfoCardPropsTypes } from '../components/isolator/Isolator';
import ResourceTable, { CustomTable, InternalTable } from '../components/resourceTable/ResourceTable';
import { HandleEventRequest, HandleInfoRequest, ModelConfig, ModelFrontendEvent, ModelFrontendEventNameEnum, ModelPlugin, ProtoAuthInfo } from "../generated-sources/openapi";
import { AppState, SpecificAction } from '../types/Event';
import { api, apiHost } from '../utils/config';
import { showNotification } from '../utils/notification';
import { getAuthenticationSetting, getPluginSetting, parseAuthenticationSetting, parsePluginSetting, settingAuthentication, settingPlugin } from '../utils/settings';
import './Home.css';
import yaml from "js-yaml";
import SideDrawer, { DrawerPropsTypes } from '../components/drawer/Drawer';
import SpecificActionForm, { SpecificActionFormProps } from '../components/specificActionForm/SpecificActionForm';

const Home: React.FC = () => {

  const [drawerState, setDrawerState] = useState<DrawerPropsTypes>({} as DrawerPropsTypes)

  const [specificActionFormState, setSpecificActionFormState] = useState<SpecificActionFormProps>({
    formItems: {},
  } as SpecificActionFormProps)

  const [searchOptions, setSearchOptions] = useState([]);

  const [currentSettings, setCurrentSettings] = useState<Array<string>>([])
  const [appConfig, setAppConfig] = useState<AppState>({} as AppState);


  // Resource table
  const [currentResourceSpecificActions, setCurrentResourceSpecificActions] = useState<Array<SpecificAction>>([]);
  const [customTable, setCustomTable] = useState<CustomTable>({} as CustomTable);
  const [websocketURL, setwebsocketURL] = useState("");
  const [internalTable, setInternalTable] = useState({} as InternalTable);

  const [isolatorsList, setIsolatorsList] = useState<Array<string>>([])

  // Search bar
  const [hideSearchBar, setHideSearchBar] = useState(false);

  const authTablecolumns = ["ID", "NAME"].map((title, index) => ({
    title,
    dataIndex: index,
    key: index,
  }));

  // Modal
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [deleteModalMessage, setDeleteModalMessage] = useState("");

  const handleOnDeleteOkButtonClick = () => {
    const e: ModelFrontendEvent = {
      eventType: "normal-action",
      name: ModelFrontendEventNameEnum.Delete,
      isolatorName: appConfig?.currentIsolator,
      pluginName: appConfig?.currentPluginName,
      resourceName: tableRow["1"],
      resourceType: appConfig?.currentResourceType,
      args: {},
    }

    let params: HandleEventRequest = {
      id: appConfig.generalInfo.id,
      modelFrontendEvent: e
    }

    api.handleEvent(params)
      .then(res => {
        showNotification('success', `Successfully deleted ${e.resourceType}`, '')
      })
      .catch(err => {
        showNotification('error', 'Event invocation failed', err)
      })

    setIsModalOpen(false);
  };

  const handleOnDeleteCancelButtonClick = () => {
    setIsModalOpen(false);
  };

  // Main use effect
  useEffect(() => {
    fetchData()
  }, []);


  const fetchData = async () => {

    let authDataRows: Array<object> = []
    let serverConfig: ModelConfig = {} as ModelConfig
    let pluginName: string = ""
    let appC = {} as AppState

    try {
      api
      serverConfig = await api.handleConfig()

      if (serverConfig.plugins.length === 0) {
        showNotification('error', 'No plugins found', '')
        throw new Error("no plugins found")
      }

      pluginName = serverConfig.plugins[0].name
      const p: ModelPlugin | undefined = serverConfig.plugins.find((plugin) => {
        if (plugin.isDefault) {
          return true
        }
      })

      if (p !== undefined) {
        pluginName = p?.name
      }

      const obj = await loadPlugin(pluginName, serverConfig)
      console.log("return", obj);
      authDataRows = obj.authDataRows
      if (obj.error !== "") {
        throw new Error(obj.error)
      }

      //   // Close socket on unmount:
      // return () => ws.close();
    } catch (error) {
      showNotification('error', 'Failed to load plugin or connection failure', error.message)
      setCustomTable({
        dataRows: authDataRows,
        headerRow: authTablecolumns,
        tableName: "authentication",
      })
      // setHeaderRow(authTablecolumns)
      // setDataRows(authDataRows)

      // setTableName("Authentication")
      console.log("Hi", authDataRows);

      setAppConfig({
        ...appConfig,
        serverConfig: serverConfig,
        currentPluginName: pluginName,
      })

    }
  }

  const loadPlugin = async (pluginName: string, serverConfig: ModelConfig) => {
    const pluginAuths = await api.handleAuth({ pluginName: pluginName })

    if (pluginAuths?.auths?.length === 0) {
      showNotification('error', 'No auths found', `no authentication found for default plugin ${pluginName}`)
      // TODO: We should not throw an error here
      throw new Error(`no authentication found for default plugin ${pluginName}`)
    }

    let authDataRows: Array<object> = []

    let pluginAuth: ProtoAuthInfo = {}

    let settings = []

    for (let i = 0; i < pluginAuths?.auths?.length; i++) {

      // TODO: context creation
      authDataRows.push({ "0": pluginAuths.auths[i].identifyingName, "1": pluginAuths.auths[i].name, "key": pluginAuths.auths[i].name } as object)

      if (pluginAuths.auths[i].isDefault) {
        pluginAuth = pluginAuths.auths[i]
      }

      settings.push(getAuthenticationSetting(pluginAuths.auths[i].identifyingName, pluginAuths.auths[i].name))
    }

    for (let i = 0; i < serverConfig.plugins.length; i++) {
      settings.push(getPluginSetting(serverConfig.plugins[i].name))
    }

    setCurrentSettings(settings)

    // If default auth is not found, ask the user to select one
    if (pluginAuth?.identifyingName === "") {
      // Show auth selection
      return {
        "error": `Default authentiation not found ${pluginName}`,
        "authDataRows": authDataRows
      }
    }

    setAppConfig({ ...appConfig, serverConfig: serverConfig, currentPluginName: pluginName } as AppState)

    try {
      await connectAndLoadData(pluginName, pluginAuth);
    } catch (e) {
      console.log("Im here", authDataRows);
      return {
        "error": e.message,
        "authDataRows": authDataRows
      };
    }

    return {
      "error": "",
      "authDataRows": authDataRows
    }
  }


  const connectAndLoadData = async (pluginName: string, pluginAuth: ProtoAuthInfo) => {
    try {
      // TODO: Fix error, genrated by the server not getting catched
      const generalInfo = await api.handleInfo({
        pluginName: pluginName, authId: pluginAuth.identifyingName, contextId: pluginAuth.name
      } as HandleInfoRequest)

      // setDataRows(dataRows => []);
      setCustomTable({ dataRows: [], headerRow: [], tableName: "Clearing Table..." })
      setIsolatorsList(isolatorsList => [])

      setSearchOptions(s => {
        return [...generalInfo.resourceTypes.map(type => {
          return { value: type }
        }), ...currentSettings.map(type => {
          return { value: type }
        })]
      })

      setwebsocketURL(`ws://${apiHost}/v1/ws/${generalInfo.id}`)

      setIsolatorCardState({
        ...isolatorCardState,
        defaultIsolator: generalInfo.defaultIsolator[0],
        currentIsolator: generalInfo.defaultIsolator[0],
        isolators: [],
        frequentlyUsed: generalInfo.defaultIsolator,
      })

      setAppConfig(a => ({
        ...a,
        generalInfo: generalInfo,
        currentIsolator: generalInfo.defaultIsolator[0],
        currentPluginName: pluginName,
        currentResourceType: generalInfo.isolatorType,
      }))

      console.log("Connect & load successfully");

    } catch (e) {
      throw e;
    }
  }

  const handleOnResourceEvent = (e: ModelFrontendEvent, skipResolveArgs: boolean = false) => {


    console.log("App config", appConfig);

    if (e.name === "connect") {
      // This is id
      // e.isolatorName = ""
      // This is context name
      // e.resourceName = ""

      const pluginAuth: ProtoAuthInfo = {
        identifyingName: e.isolatorName,
        name: e.resourceName,
      }

      console.log("Connect called with plugin auth", pluginAuth, appConfig.currentPluginName);
      connectAndLoadData(appConfig.currentPluginName, pluginAuth).catch((e) => { showNotification("error", e.message, "") })
      return
    }


    let params: HandleEventRequest = {
      id: appConfig.generalInfo.id,
      modelFrontendEvent: e
    }


    console.log(`${e.name} event is triggered`);

    if (e.eventType === "specfic-action") {
      // Find this event in our app config
      const sa: SpecificAction | undefined = currentResourceSpecificActions.find((action) => {
        if (action.name === e.name) {
          return true
        }
      })

      if (e.name !== "delete-long-running" && sa === undefined) {
        showNotification("error", "Specific action not found", "")
        return
      }

      // debugger
      if (sa?.execution?.user_input?.required && !skipResolveArgs) {
        const newEvent = { ...e, args: sa.execution?.user_input?.args, name: "specific-action-resolve-args" as ModelFrontendEventNameEnum }
        // params.modelFrontendEvent.args = sa.execution?.user_input?.args
        // params.modelFrontendEvent.name = "specific-action-resolve-args" as ModelFrontendEventNameEnum
        params.modelFrontendEvent = newEvent

        api.handleEvent(params)
          .then(res => {
            // TODO: Show form
            setSpecificActionFormState({ ...specificActionFormState, event: e, open: true, formItems: res.result })
            // setstate, pass res.result to form
          })
        return
      }

      api.handleEvent(params)
        .then(res => {
          if (e.name === "delete-long-running") {
            return
          }

          if (sa.output_type === "event") {
            const tablecolumns = ["ID", "Name", "Status", "Message"].map((title, index) => ({
              title,
              dataIndex: title.toLowerCase(),
              key: title.toLowerCase(),
            }));

            // res.result is an object whose value is also an object
            // Iterater over res.result
            const dataRows = Object.values(res.result).map((val, index) => {
              return {
                key: index,
                ...val
              }
            })

            console.log("Data rows", dataRows);
            console.log("Table columns", tablecolumns);
            setInternalTable({ dataRows: dataRows, headerRow: tablecolumns, key: e.resourceName })
            return
          }

          if (sa.output_type === "string") {
            setDrawerState({ ...drawerState, drawerBodyType: 'editor', isDrawerOpen: true, editorOptions: { defaultText: res.result, isReadOnly: true } });
            return
          }

          if (sa.output_type === "bidirectional" || sa.output_type === "stream") {
            setDrawerState({
              ...drawerState,
              drawerBodyType: 'xterm',
              isDrawerOpen: true,
              socketUrl: `ws://${apiHost}/v1/ws/action/${appConfig?.generalInfo.id}/${res.id}`
            });
            return
          }

        })
        .catch(err => showNotification('error', 'Failed to perform specific action', err.message))

      return
    }

    switch (e.name) {
      case ModelFrontendEventNameEnum.Read:
        api.handleEvent(params)
          .then(res => {
            const prettyYaml = yaml.dump(res.result, {
              indent: 2,
              noArrayIndent: true,
              lineWidth: -1, // Disables line wrapping
              quotingType: '"', // Use double quotes for strings
            });
            setDrawerState({ ...drawerState, drawerBodyType: 'editor', isDrawerOpen: true, editorOptions: { defaultText: prettyYaml, isReadOnly: true } });
          })
          .catch(err => showNotification('error', 'Failed to perform read operation', err.message))
        break;

      case ModelFrontendEventNameEnum.Edit:
        // Read first, save button will handle the rest
        params.modelFrontendEvent.name = ModelFrontendEventNameEnum.Read
        api.handleEvent(params)
          .then(res => {
            const prettyYaml = yaml.dump(res.result, {
              indent: 2,
              noArrayIndent: true,
              lineWidth: -1, // Disables line wrapping
              quotingType: '"', // Use double quotes for strings
            });
            // console.log("here 0", tableRow);
            setDrawerState({
              ...drawerState, resourceName: e.resourceName, drawerBodyType: 'editor', isDrawerOpen: true, editorOptions: { defaultText: prettyYaml, isReadOnly: false }, appConfig: appConfig
            });
          })
          .catch(err => showNotification('error', 'Failed to perfrom edit operation', err.message))
        break;

      case ModelFrontendEventNameEnum.Delete:
        setDeleteModalMessage(`Are you sure you want to delete ${e.resourceType} ${e.resourceName}?`)
        setIsModalOpen(true);
        break;

      case "use" as ModelFrontendEventNameEnum:
        console.log("Frequency used", isolatorCardState.frequentlyUsed);
        setIsolatorCardState((prevState) => {


          const resultSet = new Set(appConfig.generalInfo.defaultIsolator);

          // Add newIsolator to resultSet
          resultSet.add(e.resourceName)

          // Add frequentlyUsedIsolators to resultSet
          prevState.frequentlyUsed.forEach((isolator) => resultSet.add(isolator));

          // return Array.from(resultSet);

          // const isResourceAlreadyInFrequentlyUsed = prevState.isolators.find((item) => item === e.resourceName);

          return {
            ...prevState,
            frequentlyUsed: Array.from(resultSet)
          };
        });

        break;

      case ModelFrontendEventNameEnum.ResourceTypeChange:
        api.handleEvent(params)
          .then(res => showNotification('success', 'Successfully updated', `${ModelFrontendEventNameEnum.ResourceTypeChange} success`))
          .catch(err => showNotification('error', 'Failed to change resource type', err))
        break;

      case ModelFrontendEventNameEnum.IsolatorChange:
        api.handleEvent(params)
          .then(res => showNotification('success', 'Successfully updated', `${ModelFrontendEventNameEnum.IsolatorChange} success`))
          .catch(err => showNotification('error', 'Failed to change isolator', err))
        break;

      case "close":
        api.handleEvent(params)
          .then(res => showNotification('success', 'Successfully closed previous plugin', ""))
          .catch(err => showNotification('error', 'Failed to close previous plugin', err))
        break;

      default:
        break;
    }
  }

  const handleOnKeyBoardPress = (event: React.KeyboardEvent<HTMLElement>) => {
    console.log("Somethis is pressed", event.key);

    if (event.ctrlKey && event.key.toLowerCase() === 'a') {
      console.log("Ctrl + A is pressed");
      event.preventDefault();
      setHideSearchBar(!hideSearchBar);
    }

    if (event.key.toLowerCase() === 'd') {
      const e: ModelFrontendEvent = {
        eventType: "specfic-action",
        name: "describe",
        isolatorName: appConfig?.currentIsolator,
        pluginName: appConfig?.currentPluginName,
        resourceName: tableRow["1"],
        resourceType: appConfig?.currentResourceType,
        args: {},
      }

      let params: HandleEventRequest = {
        id: appConfig.generalInfo.id,
        modelFrontendEvent: e
      }

      // d.handleEvent(params)
      //   .then(res => {
      //     setOpen(true);
      //     // const prettyYaml = yaml.stringify(res.result, 4);
      //     setDrawerMessage(res.result);
      //   })
      //   .catch(err => openNotificationWithIcon('error', err))
    }
  };

  const onSearch = async (value: string) => {
    console.log("Search result is", value);


    if (value.startsWith(settingPlugin)) {
      const pluginName = parsePluginSetting(value)

      handleOnResourceEvent({ eventType: "internal-action", name: "close", pluginName: appConfig.currentPluginName } as ModelFrontendEvent)
      const obj = await loadPlugin(pluginName, appConfig.serverConfig)
      if (obj.error !== "") {
        showNotification('error', `Failed to load ${pluginName} plugin or connection failure`, obj.error.message)
        setCustomTable({
          dataRows: obj.authDataRows,
          headerRow: authTablecolumns,
          tableName: "authentication",
        })
        // setHeaderRow(authTablecolumns)
        // setDataRows(obj.authDataRows)
        // setTableName("Authentication")
      }
      return
    }

    if (value.startsWith(settingAuthentication)) {
      const parseResult = parseAuthenticationSetting(value)
      console.log("parseResult", parseResult);


      handleOnResourceEvent({ eventType: "internal-action", name: "close", pluginName: appConfig.currentPluginName } as ModelFrontendEvent)

      connectAndLoadData(appConfig.currentPluginName, { identifyingName: parseResult.identifyingName, name: parseResult.name } as ProtoAuthInfo)
      return
    }

    const event: ModelFrontendEvent = {
      eventType: "internal-action",
      name: "resource-type-change",
      isolatorName: appConfig?.currentIsolator,
      pluginName: appConfig?.currentPluginName,
      resourceName: "",
      resourceType: value.toLowerCase(),
      args: {},
    }

    handleOnResourceEvent(event)
    setCustomTable({ dataRows: [], headerRow: [], tableName: value.toLowerCase() })
    // setDataRows(() => [])
    // setTableName(value.toLowerCase())
    setAppConfig({ ...appConfig, currentResourceType: value.toLowerCase() })
    // appConfig.currentResourceType = value.toLowerCase()
    // setHideSearchBar(true);
    focusOnTable()
    // Clearing search bar
    // setSearchValue("")
    // setSelectedValue('');
    // TODO: Clear table
  }

  const handleNamespaceChange = (value: string) => {

    const event: ModelFrontendEvent = {
      eventType: "normal-action",
      name: "isolator-change",
      isolatorName: value,
      pluginName: appConfig?.currentPluginName,
      resourceName: "",
      resourceType: appConfig?.currentResourceType,
      args: {},
    }

    // setDataRows(dataRows => []);
    setCustomTable({ dataRows: [], headerRow: [], tableName: appConfig?.currentResourceType })
    handleOnResourceEvent(event)
    setIsolatorCardState({
      ...isolatorCardState,
      currentIsolator: value,
    })
    appConfig.currentIsolator = value
  };

  const onNamespaceChange = (value) => {
    console.log("App config 2 is", appConfig);
    handleNamespaceChange(value)
  }

  const [isolatorCardState, setIsolatorCardState] = useState<InfoCardPropsTypes>({} as InfoCardPropsTypes)

  const myInputRef = useRef(null);
  const focusOnTable = () => {
    myInputRef.current.focus();
  };

  const [tableRow, setTableRow] = useState()

  const handleTableRowClick = (record: any) => {
    setTableRow(record)
  }

  const handleOnSpecificActionOKButtonClick = (event: ModelFrontendEvent) => {
    handleOnResourceEvent(event, true)
    setSpecificActionFormState(p => {
      return {
        ...p,
        formItems: {},
        open: false,
      }
    })
  }

  const fuzzySearch = (searchText) => {
    // Combine the resourceTypes and currentSettings into a single array
    const list = [
      ...appConfig.generalInfo.resourceTypes,
      ...currentSettings,
    ];

    // Configure Fuse.js options
    const options = {
      keys: ['value'],
      includeScore: true,
      threshold: 0.9,
    };

    // Initialize Fuse.js
    const fuse = new Fuse(list.map((type) => ({ value: type })), options);

    // Perform fuzzy search
    const results = fuse.search(searchText);

    // Extract and return the matched items
    return results.map((result) => result.item);
  };

  const handleSearch = (searchText) => {
    const results = fuzzySearch(searchText);
    setSearchOptions(results);
    // Perform any additional actions with the search results
  };


  return (
    <>
      <SideDrawer {...drawerState} onDrawerClose={() => {
        console.log("onDrawerClose close called 2");
        setDrawerState({
          ...drawerState,
          isDrawerOpen: false,
          appConfig: appConfig
        })
      }} ></SideDrawer>
      <SpecificActionForm
        {
        ...specificActionFormState
        }
        onCancel={() => {
          setSpecificActionFormState(p => {
            return {
              ...p,
              formItems: {},
              open: false,
            }
          })
        }}
        onSubmit={handleOnSpecificActionOKButtonClick}
      ></SpecificActionForm >
      <Modal title="Confirmation" open={isModalOpen} onOk={handleOnDeleteOkButtonClick} onCancel={handleOnDeleteCancelButtonClick}>
        <p>{deleteModalMessage}</p>
      </Modal>
      {/* <Row style={{ "margin": "8px" }} align={"top"}> */}
      {/* Show Info */}
      <Row style={{ "margin": "8px" }}>
        {appConfig && appConfig.generalInfo && <>
          <Col>
            <InfoCard
              title='General Info'
              content={appConfig.generalInfo.general}
            />
          </Col>
          <Col push={1} style={{ "width": "25%" }} >
            <IsolatorCard {...isolatorCardState} isolators={isolatorsList} onNamespaceChange={onNamespaceChange}
            />
          </Col>
          <Col push={2}>
            <InfoCard
              title='Actions'
              content={appConfig.generalInfo.actions.reduce((acc, curVal) => {
                return { ...acc, [curVal.keyBinding]: curVal.name }
              }, {})}
            />
          </Col>
          <Col push={3}>
            <InfoCard
              title='Specific Actions'
              content={currentResourceSpecificActions?.reduce((acc, curVal) => {
                return { ...acc, [curVal.key_binding]: curVal.name }
              }, {})}
            />
          </Col>
        </>
        }
      </Row>
      {/* Resource Table */}
      {/* <Row style={{ "margin": "8px" }} className="row-flex-height" onKeyDown={handleOnKeyBoardPress} tabIndex={0} ref={myInputRef}> */}
      <Row style={{ "margin": "8px" }} className="row-flex-height" tabIndex={0} ref={myInputRef}>
        <Col flex={10} >
          {/* {appConfig && appConfig.generalInfo && */}
          {appConfig &&
            <ResourceTable
              handleTableRowClick={handleTableRowClick}
              onEvent={handleOnResourceEvent}
              handleResourceSpecificAction={(sas: Array<SpecificAction>) => setCurrentResourceSpecificActions(sas)}
              handleIsolator={(isolatorName) => setIsolatorsList(isolatorsList => [...isolatorsList, isolatorName])}
              isCurrentResourceAnIsolator={appConfig.currentResourceType === appConfig?.generalInfo?.isolatorType}
              pluginName={appConfig.currentPluginName}
              websocketURL={websocketURL}
              customTable={customTable}
              internalTable={internalTable}
            />}
        </Col>
      </Row>
      {/*  */}
      {
        !hideSearchBar &&
        <Row style={{ "margin": "8px", position: 'fixed', bottom: 0, left: 0, right: 0 }} >
          <Col span={24}>
            {appConfig && appConfig.generalInfo &&
              <AutoComplete
                options={searchOptions}
                filterOption={false}
                style={{ width: '100%' }}
                autoFocus={true}
                onSelect={onSearch}
                onSearch={handleSearch}
                backfill={true}
                placeholder="Search Resource Types"
              >
              </AutoComplete>
            }
          </Col>
        </Row>
      }
    </>
  );
}

export default Home;
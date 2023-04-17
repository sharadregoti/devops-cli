import {
    AppstoreFilled, RightCircleFilled
} from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { Breadcrumb, Empty, Layout, Select, Space, Spin, Table, theme } from 'antd';
import { ColumnsType } from 'antd/es/table';
import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import SideNav from '../components/sideNav/SideNav';
import { HandleInfoRequest, ModelAuthResponse, ModelConfig, ModelInfoResponse } from '../generated-sources/openapi';
import { NavBarItem, PluginSelectorState, TableDataType } from '../redux/reducers/PluginSelectorReducer';
import { api } from '../utils/config';
import { showNotification } from '../utils/notification';

const { Header, Content, Footer, Sider } = Layout;

type MenuItem = Required<MenuProps>['items'][number];

const gridStyle: React.CSSProperties = {
    width: '25%',
    textAlign: 'center',
};

const PluginSelector: React.FC = () => {
    const navigate = useNavigate();
    const dispatch = useDispatch();
    const { pluginName, serverConfig, pluginAuthData } = useSelector((state: PluginSelectorState) => state.pluginSelector);
    const [navItem, setNavItem] = useState<NavBarItem>();

    const [drawerLoading, setDrawerLoading] = useState(false);
    const {
        token: { colorBgContainer },
    } = theme.useToken();

    const columns: ColumnsType<TableDataType> = [
        {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
            // render: (text) => <a>{text}</a>,
        },
        {
            title: 'Context',
            dataIndex: 'context',
            key: 'context',
            render: (text) => <a>{text}</a>,
        },
        {
            title: 'Action',
            key: 'action',
            render: (_, record) => (
                <Space size="middle">
                    <a
                        onClick={() => {
                            const encodedPluginName = encodeURIComponent(pluginName);
                            const encodedAuthId = encodeURIComponent(record.name);
                            const encodedContextId = encodeURIComponent(record.context);
                            async function connectWithPlugin() {
                                try {
                                    setDrawerLoading(true);
                                    const generalInfo: ModelInfoResponse = await api.handleInfo({
                                        pluginName: pluginName, authId: record.name, contextId: record.context
                                    } as HandleInfoRequest)

                                    console.log("Id is ", generalInfo.id);
                                    setNavItem({ pluginName: pluginName, sessionId: generalInfo.id, authId: record.name, contextId: record.context, generalInfo: generalInfo } as NavBarItem)
                                    // dispatch({ type: 'SET_NAV_BAR_STATE', payload: { items: [...items, { pluginName: pluginName, sessionId: generalInfo.id, authId: record.name, contextId: record.context } as NavBarItem] } })
                                    setDrawerLoading(false);
                                    // navigate(`/plugin/${encodedPluginName}/${encodedAuthId}/${encodedContextId}`);
                                } catch (error) {
                                    setDrawerLoading(false);
                                    showNotification('error', 'Unable to connect', error.message);
                                }
                            }

                            connectWithPlugin();
                        }}
                    >Connect</a>
                </Space>
            ),
        },
    ];

    useEffect(() => {
        async function fetchPlugins() {
            try {
                const serverConfig: ModelConfig = await api.handleConfig()

                if (serverConfig.plugins.length === 0) {
                    showNotification('error', 'No plugins found', '')
                    throw new Error("no plugins found")
                }
                // setServerConfig(serverConfig)

                let pluginName: string = serverConfig.plugins[0].name
                serverConfig.plugins.find((plugin) => {
                    if (plugin.isDefault) {
                        pluginName = plugin.name
                        return true
                    }
                })
                // setCurrentPluginName(pluginName)
                dispatch({ type: 'SET_PLUGIN_SELECTOR_STATE', payload: { pluginName: pluginName, serverConfig: serverConfig } })
            }
            catch (error) {
                // Handle network error, non-2xx status codes, or other errors here
                showNotification('error', 'Unable to fetch plugins due to an error', error.message);
            }
        }

        fetchPlugins()
    }, [])

    const [selectWidth, setSelectWidth] = useState()

    function getTextWidth(text, font) {
        const canvas = document.createElement('canvas');
        const context = canvas.getContext('2d');
        context.font = font;
        const metrics = context.measureText(text);
        return metrics.width;
    }

    useEffect(() => {
        if (pluginName === "") {
            return
        }

        const longestPluginName = serverConfig.plugins
            ? serverConfig.plugins.reduce((longest, plugin) => {
                return plugin.name.length > longest.length ? plugin.name : longest;
            }, '')
            : '';

        const selectWidth = getTextWidth(longestPluginName, '16px Arial') + 50; // 50px extra for padding and dropdown arrow
        setSelectWidth(selectWidth);

        async function fetchPluginAuthData() {
            const pluginAuthData: ModelAuthResponse = await api.handleAuth({ pluginName: pluginName })

            if (pluginAuthData === undefined || pluginAuthData.auths === undefined) {
                showNotification('error', 'No auths found', '')
                throw new Error("no auths found")
            }

            const pluginAuthDataList: TableDataType[] = pluginAuthData.auths.map((auth) => {
                // TODO: Backend should always send this values. So frontend doesn't need to check for undefined
                return {
                    key: auth.name || "",
                    name: auth.identifyingName || "",
                    context: auth.name || "",
                }
            })

            // setPluginAuthData(pluginAuthDataList)
            dispatch({ type: 'SET_PLUGIN_SELECTOR_STATE', payload: { pluginName: pluginName, serverConfig: serverConfig, pluginAuthData: pluginAuthDataList } })
        }

        fetchPluginAuthData()
    }, [pluginName])

    const onChange = (value: string) => {
        dispatch({ type: 'SET_PLUGIN_SELECTOR_STATE', payload: { pluginName: value, serverConfig: serverConfig, pluginAuthData: pluginAuthData } })
    };

    const onSearch = (value: string) => {
        console.log('search:', value);
    };

    return (
        <Spin size='large' tip="Connecting..." spinning={drawerLoading}>
            <Layout style={{ minHeight: '100vh' }}>
                <SideNav selectedItem='0' newNavItem={navItem}></SideNav>
                <Layout className="site-layout">
                    <Content style={{ margin: '0 16px' }}>
                        <Breadcrumb style={{ margin: '16px 0' }}>
                            <Breadcrumb.Item>Plugins</Breadcrumb.Item>
                        </Breadcrumb>
                        {pluginName !== "" ?
                            (<>
                                <Select
                                    style={{ width: selectWidth }}
                                    showSearch
                                    placeholder="Select a person"
                                    optionFilterProp="children"
                                    onChange={onChange}
                                    onSearch={onSearch}
                                    defaultValue={pluginName}
                                    filterOption={(input, option) =>
                                        (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                                    }
                                    options={serverConfig.plugins !== undefined ? serverConfig.plugins.map((plugin) => {
                                        return { label: plugin.name, value: plugin.name }
                                    }) : []}
                                />
                                <Table columns={columns} dataSource={pluginAuthData} />
                            </>) : (<>
                                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100%' }}>
                                    <Empty
                                        description="Failed to load plugins"
                                        imageStyle={{ height: 60 }}
                                    >
                                        {/* <Button type="primary" onClick={handleRefresh}>Refresh</Button> */}
                                    </Empty>
                                </div>
                            </>)
                        }
                    </Content>
                    <Footer style={{ textAlign: 'center' }}>Ant Design ©2023 Created by Ant UED</Footer>
                </Layout>
            </Layout>
        </Spin>
    );
};

export default PluginSelector;
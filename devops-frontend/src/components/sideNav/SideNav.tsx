import { Menu, MenuProps } from "antd"
import Sider from "antd/es/layout/Sider"
import { ReactNode, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import Icon, { HomeOutlined } from '@ant-design/icons';
import { NavBarItem, NavBarState } from "../../redux/reducers/PluginSelectorReducer";
import {
    PieChartOutlined, AppstoreFilled, RightCircleFilled
} from '@ant-design/icons';
import { useNavigate } from "react-router-dom";
import { showNotification } from "../../utils/notification";

type MenuItem = Required<MenuProps>['items'][number];

export type SideNavProps = {
    selectedItem: string
    newNavItem?: NavBarItem
}


import kubernetesImage from '../../assets/kubernetes32.png';
import helmImage from '../../assets/helmicon-32.png';

const SideNav: React.FC<SideNavProps> = ({ selectedItem, newNavItem }) => {
    const navigate = useNavigate();
    const dispatch = useDispatch();
    const { items } = useSelector((state: NavBarState) => state.navBar);
    const [collapsed, setCollapsed] = useState(true);
    const [navBarItems, setNavBarItems] = useState<MenuItem[]>([])


    function getNavBarItem(
        auhtId: string,
        pluginName: string,
        label: React.ReactNode,
        key: React.Key,
        children?: MenuItem[],
        icon?: ReactNode,
    ): MenuItem {
        return {
            // <img src={pluginName === "kubernetes" ? kubernetesImage : helmImage} />
            style: { "display": "flex", "align-items": "center" },
            key,
            icon: label === "Plugins" ? <AppstoreFilled /> : <img src={pluginName === "kubernetes" ? kubernetesImage : helmImage} />,
            children,
            label: label === "Plugins" ? label : `${pluginName} :: ${auhtId} :: ${label}`,
            onClick: () => {
                if (label === "Plugins") {
                    navigate(`/`);
                    return;
                }
                const encodedPluginName = encodeURIComponent(pluginName);
                const encodedAuthId = encodeURIComponent(auhtId);
                const encodedContextId = encodeURIComponent(label as string);
                console.log('Clicked on item: ', key, label);
                navigate(`/plugin/${encodedPluginName}/session/${key as string}/${encodedAuthId}/${encodedContextId}`);
            }
        } as MenuItem;
    }

    useEffect(() => {
        if (items === undefined) {
            return
        }
        const navItem = items.map((item: NavBarItem) => getNavBarItem(item.authId, item.pluginName, item.contextId, item.sessionId))
        setNavBarItems(navItem)
    }, [items])

    useEffect(() => {
        if (newNavItem === undefined) {
            return
        }

        // Handle same sessionID
        const existingItem = items.find((item: NavBarItem) => item.sessionId === newNavItem.sessionId)
        if (existingItem !== undefined) {
            showNotification('error', 'Try again', 'Session ID already exists.')
            return
        }

        dispatch({ type: 'SET_NAV_BAR_STATE', payload: { items: [...items, newNavItem] } })
    }, [newNavItem])


    return (
        <Sider collapsible collapsed={collapsed} onCollapse={(value) => setCollapsed(value)}>
            {/* <div style={{ height: 32, margin: 16, background: 'rgba(255, 255, 255, 0.2)' }} /> */}
            <Menu style={{ paddingTop: '8px' }} theme="dark" defaultSelectedKeys={[selectedItem]} mode="inline" items={navBarItems} />
        </Sider>
    )
}

export default SideNav;
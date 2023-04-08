import { notification } from 'antd';


type NotificationType = 'success' | 'info' | 'warning' | 'error';

export const showNotification = (nType: NotificationType, title: string, description: string) => {

    switch (nType) {
        case 'success':
            notification.success({
                message: title,
                description: description,
            });
            break;
        case 'info':
            notification.info({
                message: title,
                description: description,
            });
            break;
        case 'warning':
            notification.warning({
                message: title,
                description: description,
            });
            break;
        case 'error':
            notification.error({
                message: title,
                description: description,
            });
            break;
    }
};
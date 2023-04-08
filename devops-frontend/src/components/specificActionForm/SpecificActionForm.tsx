import React, { useEffect, useState } from 'react';
import { Button, Form, Input, Modal, Radio } from 'antd';
import { ModelFrontendEvent, ModelFrontendEventNameEnum, HandleEventRequest } from '../../generated-sources/openapi';
import { showNotification } from '../../utils/notification';

export interface Values {
    title: string;
    description: string;
    modifier: string;
}

export type SpecificActionFormProps = {
    open: boolean;
    formItems: object;
    event: ModelFrontendEvent;
    onSubmit: (event: ModelFrontendEvent) => void;
    onCancel: () => void;
}

const SpecificActionForm: React.FC<SpecificActionFormProps> = ({
    open,
    event,
    formItems,
    onSubmit,
    onCancel,
}) => {


    // const [localFormItem, setlocalFormItem] = useState(formItems)
    const [form] = Form.useForm();

    useEffect(() => {
        form.setFieldsValue(formItems);
    }, [formItems])


    return (
        <Modal
            open={open}
            title="Form"
            okText="Ok"
            cancelText="Cancel"
            onCancel={onCancel}
            onOk={() => {
                form
                    .validateFields()
                    .then((values) => {
                        form.resetFields();
                        // console.log('Received values of form: ', values);
                        // return
                        onSubmit({ ...event, args: values });
                    })
                    .catch((info) => {
                        console.log('Validate Failed: 2', info);
                        showNotification('error', 'Event invocation failed', info)
                    });
            }}
        >
            <Form
                form={form}
                layout="vertical"
                name="form_in_modal"
                initialValues={{ modifier: 'public' }}
            >
                {Object.keys(formItems).map((key, index) => {
                    return (
                        <Form.Item
                            name={key}
                            label={key}
                            initialValue={formItems[key]}
                            rules={[{ required: true, message: 'Please input this key!' }]}
                        >
                            <Input />
                        </Form.Item>
                    )
                })
                }
            </Form>
        </Modal>
    );
};

export default SpecificActionForm;
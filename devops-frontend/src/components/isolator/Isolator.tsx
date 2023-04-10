import React, { useEffect, useState } from "react";
import { Card, Typography, Select, Radio, Space, Row, Col } from "antd";
import type { RadioChangeEvent } from 'antd';
import { AppState } from '../../types/Event'

// import './InfoCard.css'
import Paragraph from "antd/es/skeleton/Paragraph";

export type InfoCardPropsTypes = {
    currentIsolator: string,
    defaultIsolator: string,
    isolators: string[],
    frequentlyUsed: string[],
    appConfig: AppState
    onNamespaceChange: (isolatorName: string) => void
}

const IsolatorCard: React.FC<InfoCardPropsTypes> = ({ currentIsolator, defaultIsolator, isolators, frequentlyUsed, appConfig, onNamespaceChange }) => {
    const onChange = (e: RadioChangeEvent) => {
        console.log("Namespace change 1", e.target.value);
        onNamespaceChange(e.target.value)
    };

    return (
        <>
            <Row>
                <Col>
                    <Typography.Title style={{ "margin": "8px" }} level={5}>Isolator</Typography.Title>
                </Col>
                <Col>
                    <Select
                        defaultValue={defaultIsolator}
                        showSearch
                        style={{ width: "150px", margin: "5px" }}
                        onChange={(value: string) => onNamespaceChange(value)}
                        options={isolators.map((val) => ({ "value": val, "label": val }))}
                    />
                </Col>
            </Row>
            <Row>
                <Col flex={""} >
                    <Radio.Group onChange={onChange} value={currentIsolator}>
                        {/* <Space direction="vertical"> */}
                        <Row justify={"space-evenly"}>
                            {frequentlyUsed.map((val, i) => {
                                if (i < 5) {
                                    return (<Col>
                                        <Radio defaultChecked={val === defaultIsolator} key={val} value={val}>({i}): {val}</Radio>
                                    </Col>)
                                }
                            })}
                        </Row>
                        {/* </Space>/ */}
                    </Radio.Group>
                </Col>
            </Row>

        </>
    );
}

export default IsolatorCard;
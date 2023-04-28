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
        onNamespaceChange(e.target.value)
    };

    const [selectWidth, setSelectWidth] = useState()

    function getTextWidth(text, font) {
        const canvas = document.createElement('canvas');
        const context = canvas.getContext('2d');
        context.font = font;
        const metrics = context.measureText(text);
        return metrics.width;
    }

    useEffect(() => {
        if (isolators.length === 0) {
            return
        }

        const longestIsolatorName = isolators
            ? isolators.reduce((longest, isolator) => {
                return isolator.length > longest.length ? isolator : longest;
            }, '')
            : '';

        const selectWidth = getTextWidth(longestIsolatorName, '16px Arial'); // 50px extra for padding and dropdown arrow
        setSelectWidth(selectWidth);
    }, [isolators])

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
                        value={currentIsolator}
                        style={{ width: selectWidth, margin: "5px" }}
                        onChange={(value: string) => onNamespaceChange(value)}
                        filterSort={(optionA, optionB) => optionA.value.localeCompare(optionB.value)}
                        options={isolators.map((val) => {
                            return { "value": val, "label": val }
                        })}
                    />
                </Col>
            </Row>
            <Row style={{ width: selectWidth }}>
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
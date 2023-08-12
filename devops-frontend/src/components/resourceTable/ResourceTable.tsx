import { DeleteOutlined, DownOutlined, EditOutlined, FileAddOutlined, MoreOutlined, ReadOutlined } from '@ant-design/icons';
import type { TableColumnsType } from 'antd';
import { Badge, Col, Dropdown, Menu, Row, Space, Table, Tooltip, Typography } from 'antd';
import React, { useEffect, useState } from 'react';
import { ModelFrontendEvent, ProtoAction } from "../../generated-sources/openapi";
import './ResourceTable.css';
import { SpecificAction, WebsocketData } from '../../types/Event';

export type CustomTable = {
  headerRow: [],
  dataRows: [],
  tableName: string,
}

export type InternalTable = {
  key: string,
  headerRow: [],
  dataRows: [],
}

export type ResourceTablePropsType = {
  // headerRow: [],
  // dataRows: [],
  // tableName: string,
  // specificActions: Array<string>,
  pluginName: string
  websocketURL: string,
  defaultIsolatorResourceType: string,
  customTable: CustomTable
  internalTable: InternalTable,
  handleCloseGlobalLoading: () => void,
  onEvent: (e: ModelFrontendEvent) => void,
  handleIsolator: (isolatorName: string) => void,
  handleTableRowClick: (record: any) => void,
  handleResourceSpecificAction: (sas: Array<SpecificAction>) => void,
}

const ResourceTable: React.FC<ResourceTablePropsType> = ({ onEvent, handleCloseGlobalLoading, handleIsolator, handleTableRowClick, handleResourceSpecificAction, defaultIsolatorResourceType, pluginName, websocketURL, customTable, internalTable }) => {

  let count = 0;
  let wsTableName = "";
  const [itableName, setTableName] = useState("");
  const [idataRows, setDataRows] = useState<any>([]);
  const [iheaderRow, setHeaderRow] = useState<any>([]);
  const [ispecificActions, setSpecificActions] = useState<Array<ProtoAction>>([])
  const [currentInternalTable, setCurrentInternalTable] = useState(internalTable);

  // Add state to manage expanded row keys
  const [expandedRowKeys, setExpandedRowKeys] = useState<React.Key[]>([]);

  useEffect(() => {
    console.log("internal table called with key", internalTable.key);
    if (internalTable.key === "") {
      return
    }
    const record = idataRows.find((row) => row[1] === internalTable.key)
    console.log("internal table record object is", record);
    if (record === undefined) {
      return
    }
    console.log("setting expand value");
    setExpandedRowKeys([record["key"]])
    setCurrentInternalTable(internalTable)
  }, [internalTable])


  useEffect(() => {
    if (customTable.tableName !== "") {
      setTableName(() => customTable.tableName)
      setDataRows(() => customTable.dataRows)
      setHeaderRow(() => customTable.headerRow)
    }
  }, [customTable])

  const BATCH_SIZE = 500;
  const PROCESSING_INTERVAL = 500; // 500ms
  let messageBatch = [];
  let messageBuffer = [];
  let batchProcessingTimeout;
  let isProcessingBatch = false;

  const processMessageBatch = () => {
    if (messageBatch.length === 0) return;

    isProcessingBatch = true;
    const batchToProcess = messageBatch;
    messageBatch = messageBuffer;
    messageBuffer = [];

    setDataRows(currentDataRows => {
      if (currentDataRows === undefined) {
        console.log("Current data rows is undefined", idataRows.length);
        return idataRows
      }

      let newRows = [...currentDataRows];

      batchToProcess.forEach((message, index) => {
        count++;
        const resultObj = calculateDataSource(message, newRows, count);
        // if (message.name === generalInfo.isolatorType) {
        //   setIsolatorsList(isolatorsList => [...isolatorsList, resultObj.obj["1"]]);
        // }
        if (defaultIsolatorResourceType === message.name) {
          console.log("Isolator type is", message.name, "and current resource is an isolator");
          handleIsolator(resultObj.obj["1"]);
        }
        newRows = resultObj?.row;
      });

      return newRows;
    });

    // messageBatch = [];
    isProcessingBatch = false;
  };

  useEffect(() => {
    if (websocketURL === "") {
      return
    }
    // Make a websocket connection
    const ws = new WebSocket(websocketURL)

    ws.onclose = () => {
      console.log("Websocket connection closed");
    }

    ws.onerror = (e) => {
      console.log("Websocket connection error", e);
    }

    ws.onopen = () => {
      console.log("Websocket connection established");
      handleCloseGlobalLoading()
    }

    ws.onmessage = (e => {
      const message: WebsocketData = JSON.parse(e.data);
      // New Code

      if (isProcessingBatch) {
        messageBuffer.push(message);
      } else {
        messageBatch.push(message);
      }

      if (messageBatch.length >= BATCH_SIZE) {
        clearTimeout(batchProcessingTimeout);
        // console.log("Batch size reached");
        processMessageBatch();
      } else {
        clearTimeout(batchProcessingTimeout);
        batchProcessingTimeout = setTimeout(() => {
          // console.log("Time out process invoked");
          processMessageBatch();
        }, PROCESSING_INTERVAL);
      }

      // Old Code
      const headerRowData = message.data[0].data;

      const columns = headerRowData.map((title, index) => ({
        title,
        dataIndex: index,
        key: index,
        defaultSortOrder: 'descend',
        sorter: (a, b) => a[index].length - b[index].length,
        // filters: () => {
        //   return []
        // },
        // filterMode: 'tree',
        // filterSearch: true,
        // onFilter: (value: string, record) => record[index].startsWith(value),
        ellipsis: {
          showTitle: false,
        },
        render: (address) => (
          <Tooltip placement="topLeft" title={address}>
            {address}
          </Tooltip>
        ),
      }

      ));

      if (wsTableName !== message.name) {
        count = 0;
        setDataRows(() => []);
        console.log("Table name changed", message.name);
        wsTableName = message.name
        setTableName(message.name)
        setSpecificActions(message.specificActions === undefined ? [] : message.specificActions)
        handleResourceSpecificAction(message.specificActions === undefined ? [] : message.specificActions)
      }


      setHeaderRow(columns)
    });

    return () => {
      console.log("!!!!!!!!! Closing websocket connection !!!!!!!!");
      ws.close()
    }
  }, [websocketURL])

  const calculateDataSource = (message: WebsocketData, currentRows, objKey) => {
    const normalRowData = message.data[1].data;
    const obj = normalRowData.reduce((obj, val, i) => ({ ...obj, [i]: val }), { key: objKey, "color": message.data[1].color })

    // console.log("Event type", message.eventType, "length of currentRows is", currentRows.length);
    switch (message.eventType) {

      case "added":
        // console.log("what is current rows", currentRows);
        return {
          "obj": obj,
          "row": [...currentRows, obj]
        }

      case "updated":
      case "modified":
        const name = normalRowData[1];
        // console.log("Name is", name, currentRows);
        const indexToUpdate = currentRows.findIndex((record) => record[1] === name);
        // If not found add it
        if (indexToUpdate === -1) {
          // console.log("Index to update is -1", name);

          return {
            "obj": obj,
            "row": [...currentRows, obj]
          }

        }
        // console.log("Index to update is", indexToUpdate);
        const rowFunc = [
          ...currentRows.slice(0, indexToUpdate), // copy elements before the updated index
          obj, // replace the element at the updated index with the new value
          ...currentRows.slice(indexToUpdate + 1) // copy elements after the updated index
        ]
        return {
          "obj": obj,
          "row": rowFunc
        }

      case "deleted":
        const delteName = normalRowData[1];
        const deleteIndexToUpdate = currentRows.findIndex((record) => record[1] === delteName);
        if (deleteIndexToUpdate === -1) {
          // console.log("delete Index to update is -1", delteName);
          return
        }
        // console.log("delete index to update is", deleteIndexToUpdate);
        const deleterowFunc = [
          ...currentRows.slice(0, deleteIndexToUpdate), // copy elements before the updated index
          ...currentRows.slice(deleteIndexToUpdate + 1) // copy elements after the updated index
        ]
        return {
          "obj": obj,
          "row": deleterowFunc
        }

      default:
        console.log("Default case", message.eventType);
    }
  }

  const columnActions = {
    title: 'Actions',
    dataIndex: 'actions',
    key: 'actions',
    render: (_, record: any) => (
      <>
        <Row gutter={8}>
          {itableName.toLowerCase() === "authentication" ?
            (
              <Col>
                <ReadOutlined
                  style={{ color: "blue" }}
                  onClick={() => handleActionClick("connect", "normal-action", record)}
                />
              </Col>
            ) : (
              <>
                <Col>
                  <ReadOutlined
                    style={{ color: "blue" }}
                    onClick={() => handleActionClick("read", "normal-action", record)}
                  />
                </Col>
                <Col>
                  <EditOutlined
                    style={{ color: "blue" }}
                    onClick={() => handleActionClick("edit", "normal-action", record)}
                  />
                </Col>
                <Col>
                  <DeleteOutlined
                    style={{ color: "red" }}
                    onClick={() => handleActionClick("delete", "normal-action", record)}
                  />
                </Col>

                {itableName === defaultIsolatorResourceType &&
                  <Col>
                    <FileAddOutlined
                      style={{ color: "gray" }}
                      onClick={() => handleActionClick("use", "normal-action", record)}
                    />
                  </Col>
                }
                <Col>
                  <Dropdown
                    overlay={() => specificActionMenu(record)}
                    placement="bottomRight">
                    <MoreOutlined />
                  </Dropdown>
                </Col>
              </>
            )
          }
        </Row>
      </>
    )
  }

  const handleActionClick = (eventName, eventType, record) => {
    onEvent({
      eventType: eventType,
      name: eventName,
      isolatorName: record["0"],
      pluginName: pluginName,
      resourceName: record["1"],
      resourceType: itableName.toLowerCase(),
      args: {},
    })
  }

  const specificActionMenu = (record: object) => (
    <Menu
      onClick={(event) => {
        const { key } = event
        handleActionClick(key, "specfic-action", record)
      }}
      items={
        ispecificActions.map((action) => {
          return {
            key: action.name,
            label: action.name,
          }
        })
      }
    >
    </Menu>
  );

  // Add a new state variable to store the currently selected row index
  const [selectedRowIndex, setSelectedRowIndex] = useState(-1);

  // Add an onKeyDown event handler to the table element
  const handleKeyDown = (event) => {
    if (event.key === 'ArrowDown') {
      setSelectedRowIndex((prevIndex) => Math.min(prevIndex + 1, dataRows.length - 1));
    } else if (event.key === 'ArrowUp') {
      setSelectedRowIndex((prevIndex) => Math.max(prevIndex - 1, 0));
    }
  };

  // Update the rowClassName function to highlight the selected row
  const getRowClassName = (record, index) => {
    const baseClassName = record["color"] === "" ? "white" : record["color"];
    return index === selectedRowIndex ? `${baseClassName} highlighted` : baseClassName;
  };

  const handleTableFocus = () => {
    console.log("Table focused");
    if (selectedRowIndex === -1) {
      setSelectedRowIndex(0);
    }
  };

  interface ExpandedDataType {
    key: React.Key;
    date: string;
    name: string;
    upgradeNum: string;
  }

  const items = [
    { key: '1', label: 'Action 1' },
    { key: '2', label: 'Action 2' },
  ];

  const expandedRowRender = () => {
    const actions = {
      title: 'Action',
      dataIndex: 'operation',
      key: 'operation',
      render: (_, record: object) => (
        <Space size="middle">
          <a
            onClick={(event) => {
              // const { key } = event
              console.log("Delete clicked", record);
              // return
              handleActionClick("delete-long-running", "specfic-action", { "0": record["id"], "1": record["id"] })
              setCurrentInternalTable(p => {
                return {
                  ...p,
                  dataRows: p.dataRows.filter((row) => row["id"] !== record["id"])
                }
              })
            }}
          >Delete</a>
          {/* <a>Stop</a>
          <Dropdown menu={{ items }}>
            <a>
              More <DownOutlined />
            </a>
          </Dropdown> */}
        </Space>
      ),
    }
    return <Table columns={[...currentInternalTable.headerRow, actions]} dataSource={currentInternalTable.dataRows} pagination={false} />;
  };

  // Function to handle row expansion and collapse
  const handleExpandCollapse = (expanded: boolean, record: DataType) => {
    if (expanded) {
      setExpandedRowKeys((prevKeys) => [...prevKeys, record.key]);
    } else {
      setExpandedRowKeys((prevKeys) => prevKeys.filter((key) => key !== record.key));
    }
  };


  return (
    itableName && itableName !== "" && (
      <Table
        title={() => (
          <center>
            <Typography.Title style={{ padding: 0, margin: 0 }} level={5}>
              {itableName.charAt(0).toUpperCase() + itableName.slice(1)} ({idataRows.length})
            </Typography.Title>
          </center>
        )}
        expandable={{
          expandedRowRender: expandedRowRender,
          expandedRowKeys: expandedRowKeys,
          onExpand: handleExpandCollapse,
          rowExpandable: (record: object) => {
            return expandedRowKeys.includes(record["key"]);
          },
        }}
        size="small"
        sticky={true}
        rowClassName={getRowClassName} // Update rowClassName to use the new getRowClassName function
        onRow={(record, index) => {
          return {
            onClick: (event) => {
              handleTableRowClick(record);
            }, // click row
            tabIndex: 0, // Make the table focusable
            onDoubleClick: (event) => { }, // double click row
            onContextMenu: (event) => { }, // right button click row
            onMouseEnter: (event) => { }, // mouse enter row
            onMouseLeave: (event) => { }, // mouse leave row
          };
        }}
        // scroll={{ y: 1000 }}
        scroll={{ x: 'max-content' }}
        // tableLayout='auto'
        pagination={false}
        // pagination={{
        //   pageSize: 100,
        // }}
        columns={[...iheaderRow, columnActions]}
        dataSource={idataRows}
        // rowClassName={(record, index) => record["color"] === "" ? "white" : record["color"]}
        bordered
      />
    ) || <></>
  );
}

export default ResourceTable
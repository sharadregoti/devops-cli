import React from "react";
import { Card, Descriptions, Typography } from "antd";
import { Action, General, Plugins } from "../../types/InfoCardTypes";
import './InfoCard.css'

type InfoCardPropsTypes = {
  title: string;
  content: Action | Plugins | General | { 0: string } | {};
}

const InfoCard: React.FC<InfoCardPropsTypes> = ({ title, content }) => {
  const keyValueStyle = {
    display: "block",
    maxWidth: "40ch",
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
    lineHeight: "1.5",
    paddingRight: "1ch",
    // display: "block",
  };

  return (
    <>
      <Card
        title={title}
        size="small"
      >
        {content &&
          Object.entries(content).map(([key, value]: [string, string]) => (
            // Object.entries(content).slice(0, 4).map(([key, value]: [string, string]) => (
            <span style={keyValueStyle} key={key} title={`${key}: ${value}`}>
              {/* {key}: {value} */}
              {value}
            </span>
          ))}
      </Card>
    </>
  );
}


export default InfoCard;
import React, { useState } from "react";
import { Card, Descriptions, Space, Typography } from "antd";
import { Action, General, Plugins } from "../../types/InfoCardTypes";
// import './InfoCard.css'
import CheckableTag from "antd/es/tag/CheckableTag";

type RecentlyUsedPropsTypes = {
    title: string;
    recentlyUsedItems: string[];
    // Define an async function
    onSearch: (value: string) => void
    // onSearch: async (value: string) => {
    // content: Action | Plugins | General | { 0: string } | {};
}

const RecentlyUsed: React.FC<RecentlyUsedPropsTypes> = ({ title, recentlyUsedItems, onSearch }) => {
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

    const [selectedTags, setSelectedTags] = useState<string[]>(['Books']);

    const handleChange = (tag: string, checked: boolean) => {
        const nextSelectedTags = checked
            ? [tag]
            : selectedTags.filter((t) => t !== tag);
        console.log('You are interested in: ', nextSelectedTags);
        setSelectedTags(nextSelectedTags);

        if (checked) {
            onSearch(tag);
        }
    };

    // const tagsData = ['Movies', 'Books', 'Music', 'Sports'];

    return (
        <>
            <Card
                title={title}
                size="small"
            >
                <>
                    {/* <span style={{ marginRight: 8 }}>Categories:</span> */}
                    {/* <Space size={[0, 8]} wrap> */}
                    {recentlyUsedItems.map((tag) => (
                        <CheckableTag
                            key={tag}
                            checked={selectedTags.includes(tag)}
                            onChange={(checked) => handleChange(tag, checked)}
                        >
                            {tag}
                        </CheckableTag>
                    ))}
                    {/* </Space> */}
                </>
            </Card>
        </>
    );
}


export default RecentlyUsed;
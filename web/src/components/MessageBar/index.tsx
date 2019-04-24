import React, { Component } from 'react'

export enum MessageType {
    Success,
    Failure,
}

type Props = {
    title: string,
    message: string,
    type: MessageType,
}

const MessageBar: React.FunctionComponent<Props> = (props) => {
    // Common styles
    let styles = "border px-4 py-3 rounded relative"

    // Apply specific styles
    if (props.type === MessageType.Failure) {
        styles += " bg-red-lightest border-red-light text-red-dark"
    } else if (props.type === MessageType.Success) {
        styles += " bg-green-lightest border-green-light text-green-dark"
    }

    return (
        <div className={styles} role="alert">
            <strong className="font-bold">{props.title}</strong>
            <span className="block sm:inline"> {props.message}</span>
        </div>
    )
}

export default MessageBar
import * as React from "react"
import styled from "styled-components";
import { Paragraph } from './text';
import Modal from 'react-modal'
import { TextBox } from './input';
import { FaCheck, FaWindowClose } from "react-icons/fa";
import { Button, Typography } from "@material-ui/core";

export interface OkCancelModalProps {
    isOpen: boolean
    heading: string
    text: string
    okText: string
    okHandler: () => void
    cancelText: string
    cancelHandler: () => void
}

/**
 * A modal dialog that contains a heading, text and cancel, ok button.
 * @param props 
 */
export const OkCancelModal: React.FunctionComponent<OkCancelModalProps> = (props) => {
    const onKeyDown = (e: React.KeyboardEvent) => {
        if (e.key == "Enter") {
            props.okHandler()
        }

        if (e.key == "Escape") {
            props.cancelHandler()
        }
    }

    return <Modal isOpen={props.isOpen}>
        <div onKeyDown={onKeyDown}>
            <Typography variant="h4">{props.heading}</Typography>
            <Typography variant="body1">{props.text}</Typography>
            {
                props.children
            }
            <div>
                <Button variant="outlined" color="secondary" onClick={() => props.cancelHandler()}>{props.cancelText}</Button>
                <Button variant="contained" color="primary" onClick={() => props.okHandler()}>{props.okText}</Button>
            </div>
        </div>
    </Modal>
}
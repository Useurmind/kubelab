import * as React from "react"
import styled from "styled-components";
import { H3 } from './headings';
import { Paragraph } from './text';
import Modal from 'react-modal'
import { TextBox } from './input';
import { Button } from './button';

export const ModalHeading = H3
export const ModalText = Paragraph
export const ModalButtonBar = styled.div``

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
            <ModalHeading>{props.heading}</ModalHeading>
            <ModalText>{props.text}</ModalText>
            {
                props.children
            }
            <ModalButtonBar>
                <Button onClick={() => props.cancelHandler()}>{props.cancelText}</Button>
                <Button onClick={() => props.okHandler()}>{props.okText}</Button>
            </ModalButtonBar>
        </div>
    </Modal>
}
import * as React from "react"
import { useState } from 'react'
import Modal from 'react-modal';
import { useStoreStateFromContainerContext } from 'rfluxx-react'
import { Button } from '../../components/button'
import { H3 } from '../../components/headings'
import { TextBox } from '../../components/input';
import { ModalButtonBar, ModalHeading, ModalText, OkCancelModal } from '../../components/modal';
import { IGroupListStore, IGroupListStoreState } from './group_list_store'

export const GroupList = () => {
    const [isModalOpen, setIsModalOpen] = useState(false)
    const [groupName, setGroupName ] = useState("")
    const [state, store] = useStoreStateFromContainerContext<IGroupListStore, IGroupListStoreState>({ storeRegistrationKey: "GroupListStore" })

    if (!state) {
        return null
    }

    const createGroup = () => {
        store.createGroup(groupName)
        setIsModalOpen(false)
        setGroupName("")
    }

    return <div>
        <H3>Group List</H3>
        <Button onClick={() => setIsModalOpen(true)}>Create Group</Button>
        <ul>
            {state.groups && state.groups.map(group => {
                const groupId = group.id
                return <li>
                    {group.name}<Button onClick={() => store.deleteGroup(groupId)}>Delete</Button>
                </li>
            })}
        </ul>
        <OkCancelModal isOpen={isModalOpen}
            heading="Create Group"
            text="Enter the name of the group to create"
            cancelHandler={() => setIsModalOpen(false)}
            cancelText="Cancel"
            okHandler={() => createGroup()}
            okText="Create">
            <TextBox value={groupName} autoFocus onChange={e => setGroupName(e.target.value)}></TextBox>
        </OkCancelModal>
    </div>
}

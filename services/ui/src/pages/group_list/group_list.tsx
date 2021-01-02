import * as React from "react"
import { useState } from 'react'
import Modal from 'react-modal';
import { useStoreStateFromContainerContext } from 'rfluxx-react'
import { root } from 'rxjs/internal/util/root';
import { Button } from '../../components/button'
import { H3 } from '../../components/headings'
import { TextBox } from '../../components/input';
import { ModalButtonBar, ModalHeading, ModalText, OkCancelModal } from '../../components/modal';
import { IGroup } from '../../models/project';
import { IGroupListStore, IGroupListStoreState } from './group_list_store'

export const GroupList = () => {
    const [modalState, setModalState] = useState({
        isOpen: false,
        rootGroupId: null as number,
        parentGroupId: null as number
    })
    const [groupName, setGroupName ] = useState("")
    const [state, store] = useStoreStateFromContainerContext<IGroupListStore, IGroupListStoreState>({ storeRegistrationKey: "GroupListStore" })

    if (!state) {
        return null
    }

    const createGroup = () => {
        if (modalState.rootGroupId !== null) {
            store.createSubgroup(modalState.rootGroupId, modalState.parentGroupId, groupName)
        } else {
            store.createGroup(groupName)
        }
        setModalState({
            isOpen: false,
            parentGroupId: null,
            rootGroupId: null
        })
        setGroupName("")
    }

    const startCreateGroup = (rootGroupId: number, parentGroupId: number) => {
        setModalState({
            isOpen: true,
            rootGroupId,
            parentGroupId
        })
    }

    const deleteGroup = (rootGroupId: number, groupId: number) => {
        if(!rootGroupId) {
            store.deleteGroup(groupId)
        } else {
            store.deleteSubgroup(rootGroupId, groupId)
        }
    }

    const groupList = (rootGroupId: number, parentGroupId: number, groups: IGroup[]) => {
        if (!groups || groups.length === 0) {
            return null
        }

        return <ul>
            {groups.map(group => {
                const groupId = group.id
                const currentRootId = rootGroupId
                const nextRootId = rootGroupId ? rootGroupId : groupId

                return <li key={groupId}>
                    {group.name}
                    <Button onClick={() => deleteGroup(currentRootId, groupId)}>Delete</Button>
                    <Button onClick={() => startCreateGroup(nextRootId, !rootGroupId ? null : groupId)}>Create Subgroup</Button>
                    { groupList(nextRootId, groupId, group.subgroups)}
                </li>
            })}
        </ul>
    }

    return <div>
        <H3>Group List</H3>
        <Button onClick={() => startCreateGroup(null, null)}>Create Group</Button>
        { groupList(null, null, state.groups) }
        <OkCancelModal isOpen={modalState.isOpen}
            heading="Create Group"
            text={`Enter the name of the ${modalState.parentGroupId !== null ? "sub" : ""}group to create`}
            cancelHandler={() => setModalState({isOpen: false, parentGroupId: null, rootGroupId: null})}
            cancelText="Cancel"
            okHandler={() => createGroup()}
            okText="Create">
            <TextBox value={groupName} autoFocus onChange={e => setGroupName(e.target.value)}></TextBox>
        </OkCancelModal>
    </div>
}


import * as React from "react"
import { useState } from 'react'
import Modal from 'react-modal';
import { useStoreStateFromContainerContext } from 'rfluxx-react'
import { root } from 'rxjs/internal/util/root';
import { TextBox } from '../../components/input';
import { OkCancelModal } from '../../components/modal';
import { IGroup } from '../../models/project';
import { IGroupListStore, IGroupListStoreState } from './group_list_store'
import { FaFolderPlus, FaTrash } from "react-icons/fa";
import { Button, Card, CardContent, List, ListItem, Typography } from "@material-ui/core";


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

    const groupCard = (group: IGroup) => {
        return <Card>
            <CardContent>
                <Typography color="textSecondary" gutterBottom>
                    group
                </Typography>
                <Typography variant="h5" component="h2">
                    {group.name}
                </Typography>
                { groupList(group.id, group.id, group.subgroups) }
            </CardContent>
        </Card>
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
                    <Button variant="outlined" color="secondary" onClick={() => startCreateGroup(nextRootId, !rootGroupId ? null : groupId)}>Create subgroup</Button>
                    <Button variant="outlined" color="secondary" onClick={() => deleteGroup(currentRootId, groupId)}>Delete</Button>
                    { groupList(nextRootId, groupId, group.subgroups)}
                </li>
            })}
        </ul>
    }

    return <div>
        <Typography variant="h3">Group List</Typography>
        <Button variant="outlined" color="secondary" onClick={() => startCreateGroup(null, null)}>Create Group</Button>
        <List>
        { state.groups && state.groups.map(g => {
            return <ListItem key={g.id}>{groupCard(g)}</ListItem>
        })}
        </List>
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


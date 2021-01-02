import * as rfluxx from "rfluxx"
import { reduceAction, useStore } from 'rfluxx';
import { take, takeLast } from 'rxjs/operators';
import { IConfigStore } from '../../config/config_store';
import { IGroup } from '../../models/project';
import { createGroup, deleteGroup, listGroups, updateGroup } from '../../services/projects';

export interface IGroupListStoreState {
    groups: IGroup[]
}

export const GroupListStore = (configStore: IConfigStore) => {
    const [state, setState, base] = useStore<IGroupListStoreState>({ groups: [] })

    const store = {
        ...base,
        loadGroups: () => {
            configStore.observeConfig().subscribe(config => {
                listGroups(config).then(groups => setState({ ...state.value, groups }))
            })
        },
        createGroup: (groupName: string) => {
            configStore.observeConfig().subscribe(config => {
                const group: IGroup = {
                    name: groupName,
                }
                createGroup(config, group).then(newGroup => setState({ ...state.value, groups: [...state.value.groups, newGroup] }))
            })
        },
        createSubgroup: (groupId: number, parentSubgroupId: number, subGroupName: string) => {
            var parentGroup = state.value.groups.find(g => g.id == groupId)
            const rootGroup = parentGroup

            if (!parentGroup.subgroups) {
                parentGroup.subgroups = []
            }

            if (parentSubgroupId !== null) {
                parentGroup = getSubgroupById(parentGroup.subgroups, parentSubgroupId)
                if (!parentGroup.subgroups) {
                    parentGroup.subgroups = []
                }
            }

            parentGroup.subgroups = [ ...parentGroup.subgroups, { name: subGroupName }]

            configStore.observeConfig().subscribe(config => {
                updateGroup(config, rootGroup).then(newGroup => setState(replaceGroupById(state.value, newGroup)))
            })
        },
        deleteGroup: (groupId: number) => {
            configStore.observeConfig().subscribe(config => {
                deleteGroup(config, groupId).then(deleted => {
                    const newGroups = state.value.groups.filter(x => x.id !== groupId)
                    setState({ ...state.value, groups: newGroups })
                })
            })
        },
        deleteSubgroup: (groupId: number, subgroupId: number) => {
            const parentGroup = state.value.groups.find(x => x.id === groupId)
            const deletedSubgroup = deleteSubgroupById(parentGroup, subgroupId)
            if (!deletedSubgroup) {
                return
            }
            configStore.observeConfig().subscribe(config => {
                updateGroup(config, parentGroup).then(newGroup => setState(replaceGroupById(state.value, newGroup)))
            })
        }
    }

    store.loadGroups()

    return store
}
export type IGroupListStore = ReturnType<typeof GroupListStore>


function replaceGroupById(state: IGroupListStoreState, group: IGroup): IGroupListStoreState {
    return {
        ...state,
        groups: [...state.groups.filter(x => x.id != group.id), group]
    }
}

function getSubgroupById(subgroups: IGroup[], subgroupId: number): IGroup {
    if (!subgroups) {
        return null
    }

    for (const subgroup of subgroups) {
        if (subgroup.id === subgroupId) {
            return subgroup
        }

        const found = getSubgroupById(subgroup.subgroups, subgroupId)
        if (found) {
            return found
        }
    }

    return null
}

function deleteSubgroupById(group: IGroup, subgroupId: number): boolean {
    if (!group.subgroups) {
        return false
    }

    const index = group.subgroups.findIndex(x => x.id === subgroupId)
    if (index >= 0) {
        group.subgroups.splice(index, 1)
        return true
    } else {
        for(const subgroup of group.subgroups) {
            const result = deleteSubgroupById(subgroup, subgroupId)
            if (result) {
                return true
            }
        }
    }

    return false
}
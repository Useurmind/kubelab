import * as rfluxx from "rfluxx"
import { reduceAction, useStore } from 'rfluxx';
import { take, takeLast } from 'rxjs/operators';
import { IConfigStore } from '../../config/config_store';
import { IGroup } from '../../models/project';
import { createGroup, deleteGroup, listGroups } from '../../services/projects';

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
        deleteGroup: (groupId: number) => {
            configStore.observeConfig().subscribe(config => {
                deleteGroup(config, groupId).then(deleted => {
                    const newGroups = state.value.groups.filter(x => x.id !== groupId)
                    setState({ ...state.value, groups: newGroups })
                })
            })
        }
    }

    store.loadGroups()

    return store
}
export type IGroupListStore = ReturnType<typeof GroupListStore>

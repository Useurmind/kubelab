import * as rfluxx from "rfluxx"
import { reduceAction, useStore } from 'rfluxx';
import { take, takeLast } from 'rxjs/operators';
import { IConfigStore } from '../../config/config_store';
import { IGroup } from '../../models/project';
import { listGroups } from '../../services/projects';

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
        }
    }

    store.loadGroups()

    return store
}
export type IGroupListStore = ReturnType<typeof GroupListStore>

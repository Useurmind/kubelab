import * as React from "react"
import { useStoreStateFromContainerContext } from 'rfluxx-react'
import { H3 } from '../../components/headings'
import { IGroupListStore, IGroupListStoreState } from './group_list_store'

export const GroupList = () => {
    const [state, store ] = useStoreStateFromContainerContext<IGroupListStore, IGroupListStoreState>({ storeRegistrationKey: "GroupListStore" })

    if (!state) {
        return null
    }

    return <div>
        <H3>Group List</H3>
        <ul>
            { state.groups && state.groups.map(group => 
                <li>
                    {group.name}
                </li>)}
        </ul>
    </div>
}


import * as rfluxx from "rfluxx"
import { reduceAction, useStore } from 'rfluxx';
import { Observable } from 'rxjs';
import { first, map } from 'rxjs/operators';
import { IUIConfig } from '../models/config';

export interface IConfigStoreState {
    config: IUIConfig
}

export const ConfigStore = () => {
    const [state, setState, base] = useStore<IConfigStoreState>({ config: null })

    const store = {
        ...base,
        setConfig: reduceAction(state, (s, config: IUIConfig) => ({ ...s, config })),
        loadConfig: () => {
            fetch("/config.json").then(r => r.json()).then(config => store.setConfig(config))
        },
        observeConfig: (): Observable<IUIConfig> => {
            return state.pipe(first(s => s != null && s.config != null), map(s => s.config))
        }
    }

    return store
}

export type IConfigStore = ReturnType<typeof ConfigStore>
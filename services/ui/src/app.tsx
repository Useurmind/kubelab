import * as React from "react";
import { Header } from './components/header';
import { GroupList } from './pages/group_list/group_list';

export const App = () => {
    return <div>
        <Header></Header>
        Hello kubelab
        <GroupList></GroupList>
    </div>
}
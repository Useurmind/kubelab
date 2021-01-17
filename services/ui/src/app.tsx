import * as React from "react";
import { BrowserRouter, Route, Switch } from "react-router-dom";
import { Header } from './components/header';
import { GroupList } from './pages/group_list/group_list';
import { Home } from "./pages/home/home";

export const App = () => {
    return <div>
        <BrowserRouter>
            <Header></Header>
            <Switch>
                <Route path="/groups" >
                    <GroupList></GroupList>
                </Route>
                <Route path="/" >
                    <Home />
                </Route>
            </Switch>
        </BrowserRouter>
    </div>
}
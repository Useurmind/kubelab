import * as React from "react";
import { BrowserRouter, Route, Switch } from "react-router-dom";
import styled from "styled-components";
import { Header } from './components/header';
import { GroupList } from './pages/group_list/group_list';
import { Home } from "./pages/home/home";

const Root = styled.div`
background-color: #f5f5f5;
`

export const App = () => {
    return <Root>
        <BrowserRouter>
            <Header></Header>
            <Switch>
                <Route path="/ui/groups" >
                    <GroupList></GroupList>
                </Route>
                <Route path="/ui" >
                    <Home />
                </Route>
            </Switch>
        </BrowserRouter>
    </Root>
}
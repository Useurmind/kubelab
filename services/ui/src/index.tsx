import * as React from "react";
import * as ReactDom from "react-dom";
import { SimpleContainerBuilder } from 'rfluxx';
import { ContainerContextProvider } from 'rfluxx-react';
import { App } from "./app";
import { ConfigStore, IConfigStore } from './config/config_store';
import { GroupListStore } from './pages/group_list/group_list_store';

const builder = new SimpleContainerBuilder()

builder.register(c => ConfigStore()).as("ConfigStore")
builder.register(c => GroupListStore(c.resolve("ConfigStore"))).as("GroupListStore")

const container = builder.build()

const configStore = container.resolve<IConfigStore>("ConfigStore")

configStore.loadConfig()

document.addEventListener("DOMContentLoaded", event =>
{
    const root = document.getElementById("root");
    ReactDom.render(
        <ContainerContextProvider container={container}>
            <App />
        </ContainerContextProvider>,
        root);
});
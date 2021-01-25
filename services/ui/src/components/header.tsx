import { AppBar, Link, Toolbar, Typography } from "@material-ui/core";
import * as React from "react";
import HomeIcon from '@material-ui/icons/Home';


export const Header = () => {
    return <AppBar position="static">
        <Toolbar>
            <a href="/ui/home"><HomeIcon></HomeIcon></a>
            <Typography variant="h6">
                Kubelab
            </Typography>
        </Toolbar>
    </AppBar>
}
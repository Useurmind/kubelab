import { Card, CardActions, CardContent, Link, Typography } from "@material-ui/core";
import * as React from "react"
export const Home = () => {

    return <div>
        <Typography variant="h3">Welcome to kubelab</Typography>

        <Card>
            <CardContent>
                <Typography variant="h6">Groups</Typography>
                <Typography variant="body1">Browse available groups.</Typography>
            </CardContent>
            <CardActions>
                <Link variant="button" href="/ui/groups">Go to groups</Link>
            </CardActions>
        </Card>
    </div>
}


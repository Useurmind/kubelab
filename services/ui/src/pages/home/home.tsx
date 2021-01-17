import * as React from "react"
import { Link } from "react-router-dom";
import { H3 } from '../../components/headings'
export const Home = () => {

    return <div>
        <H3>Welcome to kubelab</H3>

        <Link to="/groups">Groups</Link>
    </div>
}


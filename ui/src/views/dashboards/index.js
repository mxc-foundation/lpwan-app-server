import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import SessionStore from "../../stores/SessionStore";

import AdminDashboard from "./Admin";
import UserDashboard from "./User";


class Dashboard extends Component {

    render() {
        const user = SessionStore.getUser();

        return (<React.Fragment>
            {user.isAdmin ? <AdminDashboard user={user} /> : <UserDashboard user={user} />}
        </React.Fragment>
        );
    }
}

export default withRouter(Dashboard);

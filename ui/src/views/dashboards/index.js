import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import SessionStore from "../../stores/SessionStore";

import AdminDashboard from "./Admin";
import UserDashboard from "./User";
import Feature2FA from "./Feature2FA";


class Dashboard extends Component {
    constructor(props) {
        super(props);


        this.state = {
            show2FaFeature: false
        }
    }

    componentDidMount() {
        // TODO - api call to check if user has not enabled the feature
        this.setState({ show2FaFeature: true });
    }

    render() {
        const user = SessionStore.getUser();

        return (<React.Fragment>
            {user.isAdmin ? <AdminDashboard user={user} /> : <UserDashboard user={user} />}
            {this.state.show2FaFeature ? <Feature2FA /> : null}
        </React.Fragment>
        );
    }
}

export default withRouter(Dashboard);

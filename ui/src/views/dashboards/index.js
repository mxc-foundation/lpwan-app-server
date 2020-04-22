import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import i18n, { packageNS } from '../../i18n';
import SessionStore from "../../stores/SessionStore";
import AdminDashboard from "./Admin";
import Feature2FA from "./Feature2FA";
import UserDashboard from "./User";



class Dashboard extends Component {
    constructor(props) {
        super(props);


        this.state = {
            show2FaFeature: false
        }
    }

    componentDidMount() {
        // TODO - api call to check if user has not enabled the feature
        this.setState({ show2FaFeature: false });//edited 2020-03-11 by Namgyeong
    }

    render() {
        const user = SessionStore.getUser();

        return (<React.Fragment>
            <div style={{ padding: 30 }}>
                {user.isAdmin ? <AdminDashboard user={user} /> : <UserDashboard user={user} />}
            </div>
        </React.Fragment>
        );
    }
}

export default withRouter(Dashboard);
import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import SessionStore from "../../stores/SessionStore";

import i18n, { packageNS } from '../../i18n';
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
            <div className="position-relative">
                <div className="card-coming-soon-2">
                    <h1 className="title">{i18n.t(`${packageNS}:menu.dashboard.coming_soon`)}</h1>
                </div>
                {user.isAdmin ? <AdminDashboard user={user} /> : <UserDashboard user={user} />}

                {/* in order to disable popup - simple comment following line */}
                {this.state.show2FaFeature ? <Feature2FA /> : null}
            </div>
        </React.Fragment>
        );
    }
}

export default withRouter(Dashboard);

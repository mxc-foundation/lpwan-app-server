import React, { Component } from "react";
import { withRouter, Redirect } from "react-router-dom";

import SessionStore from "../stores/SessionStore";

class HomeComponent extends Component {
    constructor(props) {
        super(props);

        this.redirect = this.redirect.bind(this);
    }

    redirect = () => {
        const user = SessionStore.getUser();
        if (user) {
            const orgs = SessionStore.getOrganizations();
            if (SessionStore.getToken() && orgs.length > 0) {
                if(user.isAdmin){
                    return <Redirect to={`/users`}></Redirect>;
                }else{
                    return <Redirect to={`/stake/${orgs[0].organizationID}`}></Redirect>;
                }
            } else {
                console.log('User has no organisations. Redirecting to login');
                return <Redirect to={"/logout"}></Redirect>;
            }
        } else {
            return <Redirect to={"/logout"}></Redirect>;
        }
    }

    render() {
        return <React.Fragment>
            {this.redirect()}
        </React.Fragment>
    }
}

export default withRouter(HomeComponent);
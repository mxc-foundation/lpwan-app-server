import React, {Component} from "react";
import {withRouter} from "react-router-dom";


class TestAllAPIs extends Component {
    constructor(props) {
        super(props);

        this.redirect = this.redirect.bind(this);
    }

    redirect = () => {
        console.log("hey")
    }

    render() {
        return <React.Fragment>
            {this.redirect()}
        </React.Fragment>
    }
}

export default withRouter(TestAllAPIs);
import React from "react";
import { Container } from 'reactstrap';


const AuthLayout = (props) => {
    // get the child view which we would like to render
    const children = props.children || null;

    return (<div className="app">
        <div id="wrapper">
            {props.topBar}
            {props.topBanner}
            {props.sideNav}

            <div className="content-page">
                <div className="content">
                    <Container fluid>
                        {children}
                    </Container>
                </div>
            </div>
        </div>
    </div>);
}

export default AuthLayout;

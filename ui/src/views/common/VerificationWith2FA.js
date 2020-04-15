import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import NumberFormat from 'react-number-format';
import { Link, withRouter } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Card, Col, Row, Container } from 'reactstrap';
import TitleBar from "../../components/TitleBar";
import localStyles from "./Style";
import i18n, { packageNS } from "../../i18n";

const styles = {
    ...localStyles
};

class VerificationWith2FA extends Component {
    constructor(props) {
        super(props);
        this.state = {}
    }

    loadData = () => {
    }

    componentDidMount() {
        //this.loadData();
    }

    componentDidUpdate(prevProps, prevState) {
        if (prevState !== this.state && prevState.data !== this.state.data) {

        }
    }

    render() {
        const { classes } = this.props;

        return (
            <React.Fragment>
                <TitleBar>

                </TitleBar>
                <Container>
                    <Card className="card-box shadow-sm">
                        <Row>
                            <div className="text-center" style={{ width: '100%' }}>
                                <i className="mdi mdi-shield-check-outline text-primary display-3"></i>

                                <h3>Enter the one-time password{/* {i18n.t(`${packageNS}:menu.dashboard_2fa.title`)} */}</h3>
                                <p>We sent to your E-mail.{/* {i18n.t(`${packageNS}:menu.dashboard_2fa.description`)} */}</p>
                            </div>
                        </Row>
                        <Row className={classes.numLayout}>
                            <NumberFormat id="amount" className={classes.num} value={this.state.num_0} />
                            <NumberFormat id="amount" className={classes.num} value={this.state.num_1} />
                            <NumberFormat id="amount" className={classes.num} value={this.state.num_2} />
                            <NumberFormat id="amount" className={classes.num} value={this.state.num_3} />
                            <NumberFormat id="amount" className={classes.num} value={this.state.num_4} />
                            <NumberFormat id="amount" className={classes.num} value={this.state.num_5} />
                        </Row>
                        <Row>
                            <div className="text-center" style={{ width: '100%' }}>
                                <p>If you didn't receive it, please checkout your spam folder.</p>
                            </div>
                        </Row>
                        <Row className={classes.numLayout}>
                            <div class="base" style={{ display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                                <a href={'#'} style={{ color: 'white' }}>Restart</a>
                            </div>
                            <div class="baseRight" style={{ display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                                <a href={'#'} style={{ color: 'white' }}>next</a>
                            </div>
                        </Row>
                    </Card>
                </Container>
            </React.Fragment>
        );
    }
}

export default withStyles(styles)(withRouter(VerificationWith2FA));


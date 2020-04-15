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
        this.state = {
            isVerified: false
        }
    }


    componentDidMount() {
        //this.loadData();
    }

    componentDidUpdate(prevProps, prevState) {
        if (prevState !== this.state && prevState.data !== this.state.data) {

        }
    }

    restart = () => {
        this.props.history.push(this.props.restart);
    }

    next = () => {
        this.props.history.push(this.props.next);;
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

                                <h3>Enter the one-time password</h3>
                                <p>We've sent an E-mail to you.</p>
                            </div>
                        </Row>
                        <Row className={classes.numLayout}>
                            <NumberFormat id="amount" format="#" className={classes.num} value={this.state.num_0} />
                            <NumberFormat id="amount" format="#" className={classes.num} value={this.state.num_1} />
                            <NumberFormat id="amount" format="#" className={classes.num} value={this.state.num_2} />
                            <NumberFormat id="amount" format="#" className={classes.num} value={this.state.num_3} />
                            <NumberFormat id="amount" format="#" className={classes.num} value={this.state.num_4} />
                            <NumberFormat id="amount" format="#" className={classes.num} value={this.state.num_5} />
                        </Row>
                        <Row>
                            <div className="text-center" style={{ width: '100%' }}>
                                <p>If you didn't receive it, please checkout your spam folder.</p>
                            </div>
                        </Row>
                        <Row className={classes.numLayout}>
                            <div class="base" style={{ display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                                <Link to={this.props.restart}><span style={{ color: 'white' }}>restart</span></Link>
                            </div>
                            <div class="baseRight" style={{ display: 'flex', justifyContent: 'center', alignItems: 'center'}} onClick={this.next}>
                                <span  style={{ color: 'white', cursor: 'pointer' }}>next</span>
                            </div>
                        </Row>
                    </Card>
                </Container>
            </React.Fragment>
        );
    }
}

export default withStyles(styles)(withRouter(VerificationWith2FA));


import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import NumberFormat from 'react-number-format';
import { Link, withRouter } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Card, Col, Row, Container, Alert } from 'reactstrap';
import Modal from "../common/Modal";
import TitleBar from "../../components/TitleBar";
import localStyles from "./Style";
import SessionStore from "../../stores/SessionStore";
import i18n, { packageNS } from "../../i18n";

const styles = {
    ...localStyles
};

class VerificationWith2FA extends Component {
    constructor(props) {
        super(props);
        this.state = {
            isVerified: false,
            nsDialog: false,
            token:[0,0,0,0,0,0]
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
        if(SessionStore.getOTPToken('otp') == this.state.token.join("")){
            this.props.history.push(`/registration-confirm-steptwo/${this.state.token.join("")}`);
        }else{
            this.state.isVerified = false;
            let object = this.state;
            object.token = [0,0,0,0,0,0];
            object.isVerified = false;
            object.nsDialog = true;
            this.setState({object});
            //alert('Incorrect OTP code. Please, try again.');
        }
    }

    handleChange = (e) => {
        let token = this.state.token;
        token[e.target.id] = e.target.value;
        this.setState({token});
    }

    close = () => {
        let object = this.state;
        object.nsDialog = false;
        this.setState({object});
    }

    render() {
        const { classes } = this.props;

        return (
            <React.Fragment>
                {this.state.nsDialog && <Modal
                    title={""}
                    context={i18n.t(`${packageNS}:menu.common.code_unmatch`)}
                    closeModal={() => this.setState({ nsDialog: false })}
                    callback={this.close}
                />}
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
                            <NumberFormat id="0" format="#" onChange={this.handleChange} className={classes.num} value={this.state.token['0']} />
                            <NumberFormat id="1" format="#" onChange={this.handleChange} className={classes.num} value={this.state.token['1']} />
                            <NumberFormat id="2" format="#" onChange={this.handleChange} className={classes.num} value={this.state.token['2']} />
                            <NumberFormat id="3" format="#" onChange={this.handleChange} className={classes.num} value={this.state.token['3']} />
                            <NumberFormat id="4" format="#" onChange={this.handleChange} className={classes.num} value={this.state.token['4']} />
                            <NumberFormat id="5" format="#" onChange={this.handleChange} className={classes.num} value={this.state.token['5']} />
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


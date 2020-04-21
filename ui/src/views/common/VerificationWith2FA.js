import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import NumberFormat from 'react-number-format';
import { Link, withRouter } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Card, Col, Row, Container, Alert } from 'reactstrap';
import Modal from "../common/Modal";
import OtpInput from 'react-otp-input';
import TitleBar from "../../components/TitleBar";
import localStyles from "./Style";
import SessionStore from "../../stores/SessionStore";
import UserStore from "../../stores/UserStore";
import i18n, { packageNS } from "../../i18n";

const styles = {
    ...localStyles
};

class VerificationWith2FA extends Component {
    constructor(props) {
        super(props);
        this.state = {
            isVerified: false,
            modalOpen: false,
            token: null,
            focus: [false, false, false, true, false, false]
        }
    }


    componentDidMount() {
        //this.loadData();
    }

    loadData = async () => {
    }

    componentDidUpdate(prevProps, prevState) {
        /* if (prevState !== this.state && prevState.data !== this.state.data) {

        } */
    }

    restart = () => {
        this.props.history.push(this.props.restart);
    }

    next = async () => {
        const username = await SessionStore.getUsernameTemp();
        const res = await UserStore.getOTPCode(username);

        if (res !== undefined) {
            if (res.otpCode == this.state.token) {
                this.props.history.push(`/registration-confirm-steptwo/${this.state.token}`);
            } else {
                this.state.isVerified = false;
                let object = this.state;
                object.modalOpen = true;
                this.setState({ object });
            }
        } else {
            alert('OTPcode is undefined!');
        }
    }

    handleChange = (n) => {
        /* if(n.length >= 6){
            console.log('over six',n);
        }else{
            console.log('n',n);
            let object = this.state;
            object.token[0] = n;
            this.setState({object});
        } */

        let object = this.state;
        object.token = n;
        this.setState({ object });

        /* console.log('e.target.id', e.target.id);
        let object = this.state;
        object.token[e.target.id] = e.target.value;
        let focus = [false, false, false, false, false, false];
        if(e.target.id != 5){
            focus[Number.parseInt(e.target.id) + 1] = true; 
        }else{
            focus[0] = true; 
        }
        object.focus = focus;
        this.setState({object});
         */
    }

    closeModal = () => {
        let object = this.state;
        object.modalOpen = false;
        this.setState({ object });
    }

    render() {
        const { classes } = this.props;

        return (
            <React.Fragment>
                {this.state.modalOpen && <Modal
                    title={i18n.t(`${packageNS}:menu.topup.notice`)}
                    context={i18n.t(`${packageNS}:menu.registration.otp_unmatched`)}
                    callback={this.closeModal}
                />}
                <TitleBar>

                </TitleBar>
                <Container>
                    <Card className="card-box shadow-sm">
                        <Row>
                            <div className="text-center" style={{ width: '100%' }}>
                                <i className="mdi mdi-shield-check-outline text-primary display-3"></i>

                                <h3>{i18n.t(`${packageNS}:menu.registration.enter_pw`)}</h3>
                                <p>{i18n.t(`${packageNS}:menu.registration.notice_sent_email`)}</p>
                            </div>
                        </Row>
                        <Row className={classes.numLayout}>
                            <OtpInput
                                onChange={otp => this.handleChange(otp)}
                                value={this.state.token}
                                inputStyle={classes.num}

                                numInputs={6}
                                separator={<span>-</span>}
                            />
                        </Row>
                        <Row>
                            <div className="text-center" style={{ width: '100%' }}>
                                <p>{i18n.t(`${packageNS}:menu.registration.notice_check_spam`)}</p>
                            </div>
                        </Row>
                        <Row className={classes.numLayout}>
                            <div class="base" style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                                <Link to={this.props.restart}><span style={{ color: 'white' }}>{i18n.t(`${packageNS}:menu.registration.restart`)}</span></Link>
                            </div>
                            <div class="baseRight" style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }} onClick={this.next}>
                                <span style={{ color: 'white', cursor: 'pointer' }}>{i18n.t(`${packageNS}:menu.registration.next`)}</span>
                            </div>
                        </Row>
                    </Card>
                </Container>
            </React.Fragment>
        );
    }
}

export default withStyles(styles)(withRouter(VerificationWith2FA));


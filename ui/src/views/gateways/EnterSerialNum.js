import React, { Component, useState } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Button, Card, CardBody, Row, Col, CardHeader, Alert } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';
import { ReactstrapInput } from '../../components/FormInputs';
import Tooltips from "./Tooltips";
import { withStyles } from "@material-ui/core/styles";
import i18n, { packageNS } from '../../i18n';
import Spinner from "../../components/ScaleLoader";
import logo from '../../assets/images/MATCHX-SUPERNODE2.png';
import GatewayStore from "../../stores/GatewayStore";
import QReaderModal from './QReaderModal';
import ServerInfoStore from "../../stores/ServerInfoStore";
import Modal from "../common/Modal";

const styles = {
    center: {
        display: "flex",
        justifyContent: "center"
    },
    between: {
        display: "flex",
        justifyContent: "space-between"
    }
};

class EnterSerialNum extends Component {
    constructor(props) {
        super(props);
        this.state = {
            stage: 0,
            openQR: false,
            loading: false,
            modalOpen: false,
            isRegionCorrect: true,
            object: {
                serial: ''
            }
        };
    }

    componentDidMount() {
        this.loadData();
      }
    
    loadData = async () => {
        try {
            const res = await ServerInfoStore.getServerRegion();
            
            this.setState({
            serverRegion: res.serverRegion
            });
        } catch (error) {
            console.error(error);
            this.setState({ error });
        }
    }

    onSubmit = (serial) => {
        const object = this.state;
        object.loading = true;
        object.stage = 1;
        this.setState({ object });
        console.log('serial', serial);
        if (serial.serial.length === 0) {
            return false;
        }

        if (serial.serial.substring(0, 2).trim() !== 'MX' || serial.serial.substring(0, 2).trim() !== 'M2X') {
            this.props.history.push(
                `/organizations/${this.props.match.params.organizationID}/gateways/create`
            );
        } else {
            let gateway = {};
            gateway.organizationId = this.props.match.params.organizationID;
            gateway.sn = serial
            GatewayStore.register(gateway, resp => {
                this.props.history.push(
                    `/organizations/${this.props.match.params.organizationID}/gateways`
                );
            });
        }
    }

    back = () => {
        this.props.history.push(
            `/organizations/${this.props.match.params.organizationID}/gateways`
        );
    }

    readQR = (data) => {
        let isRegionCorrect = true;
        let QRCodeArray = '';

        if (typeof data === 'object') {
            return;
        }
        //data = "S/N: M2XO7FQOYGD, time: 3235, ID: MX1903, version: 1.0, MAC: 70:B3:D5:1C:B0:00}";
        if (data) {
            QRCodeArray = data.split(",");
        }

        /* 0: "S/N: M2XXXXXXXXX"
        1: " ID: MX1903"
        2: " 1.0"
        3: " 1320"
        4: " MAC: 70:B3:D5:00:00:00" */


        let modalOpen = false;
        if (QRCodeArray.length > 0) {
            const json = JSON.stringify(QRCodeArray);

            const serial = QRCodeArray[0].split(':')[1].trim();
            let time = '';
            let model = '';
            let version = ''; 
            let mac = ''; 
            if (serial.substring(0, 2).trim() !== 'M2X') {
                model = QRCodeArray[1].split(':')[1].trim();
                version = QRCodeArray[2].trim();
                time = QRCodeArray[3].trim(); 
                mac = QRCodeArray[4].substring(5, QRCodeArray[4].length).trim(); 
            } else {
                time = QRCodeArray[1].split(':')[1];
                model = QRCodeArray[2].split(':')[1].trim();
                version = QRCodeArray[3].split(':')[1]; 
            }

            if(this.state.serverRegion === 'RESTRICTED' ){
                if(model !== 'MX1903'){
                    isRegionCorrect = false;
                    modalOpen = true;
                }
            }else{
                if(model === 'MX1903'){
                    isRegionCorrect = false;
                    modalOpen = true;
                }
            }

            const object = this.state;
            object.isRegionCorrect = isRegionCorrect;
            object.modalOpen = modalOpen;
            object.object.serial = isRegionCorrect?serial:'';
            object.object.time = isRegionCorrect?time:'';
            object.object.model = isRegionCorrect?model:'';
            object.object.version = isRegionCorrect?version:'';
            if (serial.substring(0, 2).trim() !== 'M2X') {
                object.object.mac = isRegionCorrect?mac:'';
            }
            this.setState({ object });
        }
    }

    close = () => {
        const object = this.state;
        object.modalOpen = false;
        this.setState({object});
    }

    render() {
        const { classes } = this.props;
        const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

        let fieldsSchema = {
            serial: Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`)),
        }

        const formSchema = Yup.object().shape(fieldsSchema);
        const label = <span>
            {i18n.t(`${packageNS}:menu.gateways.enter_your_gateway_serial_number`)}{' '}
            <Link>
                <i id={'helper'} className="mdi mdi-help-circle-outline"></i>
            </Link>
        </span>

        return (<React.Fragment>
            {this.state.loading && <Spinner />}
            {this.state.openQR && <QReaderModal
                buttonLabel={i18n.t(`${packageNS}:tr000277`)}
                callback={this.handleLink} />}
            {this.state.modalOpen && <Modal
                buttonLabel={i18n.t(`${packageNS}:tr000277`)}
                title={i18n.t(`${packageNS}:menu.topup.notice`)}
                context={i18n.t(`${packageNS}:menu.gateways.restricted_region`)}
                callback={this.close} />}

            <Card>
                <Col xs={12}>
                    <Card className={classes.center} >
                        <CardHeader className={classes.center} style={{ marginTop: 100 }}>
                            <img src={logo} alt="" height="53" />
                        </CardHeader>
                        {this.state.stage === 0 && <CardBody className={classes.center} >
                            <Formik
                                enableReinitialize
                                initialValues={this.state.object}
                                validationSchema={formSchema}
                                onSubmit={(values) => {
                                    const castValues = formSchema.cast(values);
                                    this.onSubmit({ ...castValues })
                                }}>
                                {({
                                    handleSubmit,
                                    handleChange,
                                    setFieldValue,
                                    values,
                                    handleBlur,
                                }) => (
                                        <Form onSubmit={handleSubmit} noValidate>
                                            <div style={{ position: 'relative', display: 'inline-block' }}>
                                                <QReaderModal
                                                    buttonLabel={i18n.t(`${packageNS}:tr000277`)}
                                                    callback={this.readQR} />
                                                <Field
                                                    type="text"
                                                    label={label}
                                                    name="serial"
                                                    id="serial"
                                                    value={this.state.object.serial || ""}
                                                    autoComplete='off'
                                                    component={ReactstrapInput}
                                                    onBlur={handleBlur}
                                                    onChange={handleChange}

                                                    inputProps={{
                                                        clearable: true,
                                                        cache: false,
                                                    }}
                                                />
                                            </div>
                                            <Row>
                                                <Col className={classes.between}>
                                                    <Link to={`/organizations/${currentOrgID}/gateways`}><Button color="secondary" onClick={this.back}>{i18n.t(`${packageNS}:menu.common.back`)}</Button></Link>
                                                    <Button type="submit" color="secondary" className="btn" >{i18n.t(`${packageNS}:menu.common.submit`)}</Button>
                                                </Col>
                                            </Row>
                                        </Form>
                                    )}
                            </Formik>
                        </CardBody>}
                        {this.state.stage === 1 && <CardBody className={classes.center} style={{ height: '25vw' }}>
                            <span style={{ fontSize: 30, fontWeight: 400 }}>{i18n.t(`${packageNS}:menu.gateways.were_searching_with_your_gateway_please_wait`)}</span>
                        </CardBody>}
                        {this.state.stage === 2 && <CardBody className={classes.center} style={{ height: '25vw' }}>
                        </CardBody>}
                    </Card>
                </Col>
            </Card>
        </React.Fragment>
        );
    }
}

export default withStyles(styles)(withRouter(EnterSerialNum));
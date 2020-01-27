import React, { Component, useState } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Button, Card, CardBody, Row, Col, CardHeader } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';
import { ReactstrapInput } from '../../components/FormInputs';
import Tooltips from "./Tooltips";
import { withStyles } from "@material-ui/core/styles";
import i18n, { packageNS } from '../../i18n';
import Spinner from "../../components/ScaleLoader";
import logo from '../../assets/images/MATCHX-SUPERNODE2.png';
import { Divider } from "@material-ui/core";

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
            loading: false,
            object: {
                serial: ''
            }
        };
    }


    onSubmit = (serial) => {
        const object = this.state;
        object.loading = true;
        object.stage = 1;
        this.setState({ object });
        console.log('serial', serial);
    }

    back = () => {
        this.props.history.push(`/organizations/${this.props.match.params.organizationID}/gateways`);
    }

    readQR = () => {
        alert(1);
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
            <Card>
                <Row>
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
                                                    <i id={'helper'} onClick={this.readQR} className="mdi mdi-qrcode-scan" style={{ position: 'absolute', right: 0, top: '50%', transform: 'translate(-50%, -25%)', width: 20, height: 20 }}></i>
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
                                                        <Link to={`/organizations/${currentOrgID}/gateways/brand`}><Button color="secondary" onClick={this.back}>{i18n.t(`${packageNS}:menu.common.back`)}</Button></Link>
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
                </Row>
            </Card>
        </React.Fragment>
        );
    }
}

export default withStyles(styles)(withRouter(EnterSerialNum));

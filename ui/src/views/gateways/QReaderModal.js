import React, { useState } from 'react';
import i18n, { packageNS } from '../../i18n';
import { withRouter, Link } from "react-router-dom";
import { withStyles } from "@material-ui/core/styles";
import { Button, Modal, ModalHeader, ModalBody, ModalFooter, Card, CardBody, Row, Col } from 'reactstrap';
import logo from '../../assets/images/MATCHX-SUPERNODE2.png';
import QRCodeReader from '../../components/QRCodeReader';

const styles = {
    center: {
        display: "flex", 
        justifyContent: "center",
        height: 'auto'
    },
    modalWidth: {
        width: 800,
        height: 800
    }
  };

const QReaderModal = (props) => {
    const {
        buttonLabel,
        className
    } = props;

    const [modal, setModal] = useState(props.click);

    const toggle = (code) => {
        setModal(!modal)
        props.callback(code);
    };

    const {classes} = props;
    return (
        <div>
            <i id={'helper'} onClick={toggle} className="mdi mdi-qrcode-scan" style={{ position: 'absolute', cursor: 'pointer', right: 0, top: '50%', transform: 'translate(-50%, -25%)', width: 20, height: 20 }}></i>
            <Modal isOpen={modal} toggle={toggle} className={'center'} >
                <ModalHeader toggle={toggle}>{i18n.t(`${packageNS}:menu.gateways.qr_code_scan`)}</ModalHeader>
                <ModalBody>
                    <Row>
                        <Col>
                            <Card>
                                <CardBody className={classes.center}>
                                    <QRCodeReader toggle={toggle}/>
                                </CardBody>
                            </Card>
                        </Col>
                    </Row>
                </ModalBody>
            </Modal>
        </div>
    );
}

export default withStyles(styles)(withRouter(QReaderModal));
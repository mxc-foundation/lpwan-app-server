import { withStyles } from "@material-ui/core/styles";
import React, { useState } from 'react';
import { Link, withRouter } from "react-router-dom";
import { Button, Card, CardBody, Col, Modal, ModalBody, ModalFooter, ModalHeader, Row } from 'reactstrap';
import logo from '../../assets/images/MATCHX-SUPERNODE2.png';
import i18n, { packageNS } from '../../i18n';

const styles = {
    center: {
        display: "flex", 
        justifyContent: "center",
        height: '10vw'
    },
    modalWidth: {
        width: 800,
        height: 800
    }
  };

const BrandSelectModal = (props) => {
    const {
        buttonLabel,
        className
    } = props;

    const [modal, setModal] = useState(props.click);

    const toggle = () => {
        setModal(!modal)
    };

    const {classes} = props;
    return (
        <div>
            <Button color="primary" onClick={toggle}>{buttonLabel}</Button>
            <Modal isOpen={modal} toggle={toggle} className={'center'} >
                <ModalHeader toggle={toggle}>{i18n.t(`${packageNS}:menu.gateways.choose_a_manufacturer`)}</ModalHeader>
                <ModalBody>
                    <Row>
                        <Col>
                            <Card>
                                <CardBody className={classes.center}>
                                    <Link to=''><img src={logo} alt="" height="53" /></Link>
                                </CardBody>
                            </Card>
                        </Col>
                    </Row>
                    <Row>
                        <Col>
                            <Card>
                                <CardBody className={classes.center} >
                                    <Link to=''><b>{i18n.t(`${packageNS}:menu.gateways.other`)}</b></Link>
                                </CardBody>
                            </Card>
                        </Col>
                    </Row>
                </ModalBody>
                <ModalFooter  className={classes.center}>
                    {/* <Button color="primary" onClick={toggle}>Do Something</Button>{' '} */}
                    {/* <Button color="secondary" onClick={toggle}>Cancel</Button> */}
                </ModalFooter>
            </Modal>
        </div>
    );
}

export default withStyles(styles)(withRouter(BrandSelectModal));
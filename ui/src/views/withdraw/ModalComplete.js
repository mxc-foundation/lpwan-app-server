import { withStyles } from "@material-ui/core/styles";
import React, { useState } from 'react';
import { withRouter } from "react-router-dom";
import { Button, Col, FormGroup, Modal, ModalBody, ModalFooter, ModalHeader } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';
import localStyles from "./WithdrawStyle";


const styles = {
    ...localStyles
};

const ModalComplete = (props) => {
    const {
        className,
        closeModal,
        showCloseButton = true,
        show = true,
    } = props;

    const [modal, setModal] = useState(show);

    const toggle = () => {
        setModal(!modal);
        if (closeModal)
            closeModal();
    }

    return (
        <div>
            {/* {buttonLabel && <Button color={buttonColor} onClick={toggle}>{icon}{buttonLabel}</Button>} */}
            <Modal isOpen={modal} toggle={toggle} className={className} centered={true}>
                <ModalHeader toggle={toggle} className="border-0">

                </ModalHeader>
                <ModalBody>
                    <div style={{ display: 'flex', flexDirection: 'column' }}>
                        <FormGroup row style={{ display: 'flex', justifyContent: 'center', marginBottom: 0, height: 200 }}>
                            <i className="mdi mdi-check-circle-outline" style={{ color: '#10C469', fontSize: '150px' }}></i>
                        </FormGroup>
                        <FormGroup row style={{ display: 'flex', justifyContent: 'center' }}>
                            <span style={{ fontSize: '26px' }}>{i18n.t(`${packageNS}:menu.withdraw.confirmed`)}</span>
                        </FormGroup>
                    </div>
                    <div className="text-center">
                        <p>{i18n.t(`${packageNS}:menu.withdraw.request_withdraw_text_com`)}</p>
                        <FormGroup row>
                            <Col sm={12}>
                            </Col>
                        </FormGroup>
                    </div>
                </ModalBody>
                <ModalFooter className="border-0" style={{ display: 'flex', flexWrap: 'wrap', alignItems: 'center', justifyContent: 'center' }}>
                    {showCloseButton && <Button color="primary" onClick={toggle}>{i18n.t(`${packageNS}:menu.common.done`)}</Button>}
                </ModalFooter>
            </Modal>
        </div>
    );
}

export default withStyles(styles)(withRouter(ModalComplete));

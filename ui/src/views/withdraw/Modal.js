import React, { useState } from 'react';

import { withRouter } from "react-router-dom";
import { withStyles } from "@material-ui/core/styles";
import i18n, { packageNS } from '../../i18n';
import localStyles from "./WithdrawStyle"
import { Button, Modal, ModalHeader, ModalBody, ModalFooter, FormGroup, Label, Input, Col } from 'reactstrap';

const styles = {
    ...localStyles
};

const CommonModal = (props) => {
    const {
        className,
        closeModal,
        showCloseButton = true,
        showConfirmButton = true,
        show = true,
    } = props;

    const [modal, setModal] = useState(show);

    const toggle = () => {
        setModal(!modal);
        if (closeModal)
            closeModal();
    }

    const proc = () => {
        setModal(!modal);
        if (closeModal)
            closeModal();
        props.callback();
    }
    
    return (
        <div>
            {/* {buttonLabel && <Button color={buttonColor} onClick={toggle}>{icon}{buttonLabel}</Button>} */}
            <Modal isOpen={modal} toggle={toggle} className={className} centered={true}>
                <ModalHeader toggle={toggle}>{i18n.t(`${packageNS}:menu.withdraw.request_withdraw_title`)}</ModalHeader>
                <ModalBody>
                    <FormGroup row>
                        <Col sm={12}>
                            <span>{i18n.t(`${packageNS}:menu.withdraw.request_withdraw_text_0`)}</span>
                        </Col>
                    </FormGroup>
                    <FormGroup row>
                        <Col sm={12}>
                            <span>{i18n.t(`${packageNS}:menu.withdraw.request_withdraw_text_1`)}</span>
                        </Col>
                    </FormGroup>
                    <FormGroup row>
                        <Col sm={12}>
                            <span>{i18n.t(`${packageNS}:menu.withdraw.request_withdraw_text_2`)}</span>
                        </Col>
                    </FormGroup>
                    <FormGroup row>
                    </FormGroup>
                    <FormGroup row>
                        <Label for="amount" sm={2}>{i18n.t(`${packageNS}:menu.withdraw.amount`)}</Label>
                        <Col sm={10}>
                            <Input type="text" name="amount" id="amount" value={props.amount} readOnly />
                        </Col>
                    </FormGroup>
                </ModalBody>
                <ModalFooter>
                    {showCloseButton && <Button color="secondary" onClick={toggle}>{props.left !== undefined ? props.left : i18n.t(`${packageNS}:tr000424`)}</Button>}{' '}
                    {showConfirmButton && <Button color="primary" onClick={proc}>{props.right !== undefined ? props.right : i18n.t(`${packageNS}:tr000425`)}</Button>}
                </ModalFooter>
            </Modal>
        </div>
    );
}

export default withStyles(styles)(withRouter(CommonModal));

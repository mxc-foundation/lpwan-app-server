import { withStyles } from "@material-ui/core/styles";
import React, { useState } from 'react';
import { withRouter } from "react-router-dom";
import { Button, Col, FormGroup, Input, Label, Modal, ModalBody, ModalFooter, ModalHeader } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';
import localStyles from "./WithdrawStyle";


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
                {props.title ? <ModalHeader className="border-0" toggle={toggle}></ModalHeader> : null}
                <ModalBody className="pb-4">
                    <div className="text-center">

                        <h3>{i18n.t(`${packageNS}:menu.withdraw.request_withdraw_title`)}</h3>
                        <p>{i18n.t(`${packageNS}:menu.withdraw.request_withdraw_text_0`)}
                            {i18n.t(`${packageNS}:menu.withdraw.request_withdraw_text_1`)}
                            {i18n.t(`${packageNS}:menu.withdraw.request_withdraw_text_2`)}
                        </p>

                        <FormGroup row>
                            <Label for="amount" sm={2}>{i18n.t(`${packageNS}:menu.withdraw.amount`)}</Label>
                            <Col sm={10}>
                                <Input type="text" name="amount" id="amount" value={props.amount} readOnly />
                            </Col>
                        </FormGroup>
                        <ModalFooter className="border-0" style={{ display: 'flex', justifyContent: 'center' }}>
                            {showCloseButton && <Button color="secondary" onClick={toggle}>{props.left !== undefined ? props.left : i18n.t(`${packageNS}:tr000424`)}</Button>}{' '}
                            {showConfirmButton && <Button color="primary" onClick={proc}>{props.right !== undefined ? props.right : i18n.t(`${packageNS}:tr000425`)}</Button>}
                        </ModalFooter>

                    </div>
                </ModalBody>

            </Modal>
        </div>
    );
}

export default withStyles(styles)(withRouter(CommonModal));

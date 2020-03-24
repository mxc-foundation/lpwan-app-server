import React, { useState } from 'react';
import { Button, Col, FormGroup, Input, Label, Modal, ModalBody, ModalFooter, ModalHeader } from 'reactstrap';
import i18n, { packageNS } from '../../../i18n';

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
            <Modal isOpen={modal} toggle={toggle} className={className} centered={true}>
                {props.title ? <ModalHeader className="border-0" toggle={toggle}></ModalHeader> : null}
                <ModalBody className="pb-4">
                <div className="text-center">

                <h3>{i18n.t(`${packageNS}:menu.withdraw.confirm_modal_title`)}</h3>
                <p>{props.context}</p>

                <FormGroup row>
                    </FormGroup>
                    <FormGroup row>
                        <Label for="username" sm={2}>{i18n.t(`${packageNS}:menu.withdraw.requester`)}</Label>
                        <Col sm={10}>
                            <Input type="text" name="username" id="username" value={props.row.userName} readOnly />
                        </Col>
                    </FormGroup>
                    <FormGroup row>
                        <Label for="amount" sm={2}>{i18n.t(`${packageNS}:menu.withdraw.amount`)}</Label>
                        <Col sm={10}>
                            <Input type="text" name="amount" id="amount" value={props.row.amount} readOnly />
                        </Col>
                    </FormGroup>
                    <FormGroup row>
                        <Label for="balance" sm={2}>{i18n.t(`${packageNS}:menu.withdraw.balance`)}</Label>
                        <Col sm={10}>
                            <Input type="text" name="balance" id="balance" value={props.row.availableToken} readOnly />
                        </Col>
                    </FormGroup>
                    <FormGroup row></FormGroup>
                    {!props.status && <FormGroup>
                        <Input type="textarea" name="comment" id="comment" onChange={props.handleChange} placeholder={i18n.t(`${packageNS}:menu.withdraw.deny_reason`)} />
                    </FormGroup>}
                </div>
                </ModalBody>
                <ModalFooter className="border-0" style={{ display: 'flex', justifyContent: 'center' }}>
                    {showCloseButton && <Button color="secondary" onClick={toggle}>{props.left !== undefined ? props.left : i18n.t(`${packageNS}:tr000424`)}</Button>}{' '}
                    {showConfirmButton && <Button color="primary" onClick={proc}>{props.right !== undefined ? props.right : i18n.t(`${packageNS}:tr000425`)}</Button>}
                </ModalFooter>
            </Modal>
        </div>
    );
}

export default CommonModal;
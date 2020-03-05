import React, { useState } from 'react';
import i18n, { packageNS } from '../../../i18n';
import { Button, Modal, ModalHeader, ModalBody, ModalFooter, FormGroup, Label, Input, Col } from 'reactstrap';

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
                {props.title ? <ModalHeader toggle={toggle}>{props.title}</ModalHeader> : null}
                <ModalBody>
                    {props.context}
                    <FormGroup row>
                    </FormGroup>
                    <FormGroup row>
                        <Label for="username" sm={2}>{i18n.t(`${packageNS}:menu.withdraw.requester`)}</Label>
                        <Col sm={10}>
                        <Input type="text" name="username" id="username" value={props.row.userName} readOnly/>
                        </Col>
                    </FormGroup>
                    <FormGroup row>
                        <Label for="amount" sm={2}>{i18n.t(`${packageNS}:menu.withdraw.amount`)}</Label>
                        <Col sm={10}>
                        <Input type="text" name="amount" id="amount" value={props.row.amount} readOnly/>
                        </Col>
                    </FormGroup>
                    <FormGroup row>
                        <Label for="balance" sm={2}>{i18n.t(`${packageNS}:menu.withdraw.balance`)}</Label>
                        <Col sm={10}>
                        <Input type="text" name="balance" id="balance" value={props.row.availableToken} readOnly/>
                        </Col>
                    </FormGroup>
                    {!props.status && <FormGroup>
                        <Label for="exampleText">Comment</Label>
                        <Input type="textarea" name="comment" id="comment"  onChange={props.handleChange} placeholder={i18n.t(`${packageNS}:menu.withdraw.deny_reason`)}/>
                    </FormGroup>}
                </ModalBody>
                <ModalFooter>
                    {showCloseButton && <Button color="secondary" onClick={toggle}>{props.left !== undefined ? props.left : i18n.t(`${packageNS}:tr000424`)}</Button>}{' '}
                    {showConfirmButton && <Button color="primary" onClick={proc}>{props.right !== undefined ? props.right : i18n.t(`${packageNS}:tr000425`)}</Button>}
                </ModalFooter>
            </Modal>
        </div>
    );
}

export default CommonModal;
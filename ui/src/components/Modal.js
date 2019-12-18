import React, { useState } from 'react';
import i18n, { packageNS } from '../i18n';
import { Button, Modal, ModalHeader, ModalBody, ModalFooter } from 'reactstrap';

const CommonModal = (props) => {
    const {
        buttonLabel,
        className,
        callback
    } = props;

    const [modal, setModal] = useState(false);

    const toggle = () => setModal(!modal);
    const proc = () => {
        setModal(!modal);
        callback();
    }
    return (
        <div>
            <Button color="danger" onClick={toggle}>{buttonLabel}</Button>
            <Modal isOpen={modal} toggle={toggle} className={className}>
                <ModalHeader toggle={toggle}>{props.title}</ModalHeader>
                <ModalBody>
                    {props.context}
                </ModalBody>
                <ModalFooter>
                    <Button color="secondary" onClick={toggle}>{props.left!==undefined?props.left:i18n.t(`${packageNS}:tr000424`)}</Button>{' '}
                    <Button color="primary" onClick={proc}>{props.right!==undefined?props.right:i18n.t(`${packageNS}:tr000425`)}</Button>
                </ModalFooter>
            </Modal>
        </div>
    );
}

export default CommonModal;
import React, { useState } from 'react';
import i18n, { packageNS } from '../i18n';
import { Button, Modal, ModalHeader, ModalBody, ModalFooter } from 'reactstrap';

const CommonModal = (props) => {
    const {
        className,
        callback,
        showCloseButton = true,
        showConfirmButton = true,
    } = props;

    const [modal, setModal] = useState(true);

    const toggle = () => setModal(!modal);
    const proc = () => {
        setModal(!modal);
        props.callback();
    }
    /* const buttonColor = props.buttonColor === undefined
        ? 'primary'
        : props.buttonColor; 

    const icon = props.icon === undefined
        ? null
        : props.icon;*/

    return (
        <div>
            {/* {buttonLabel && <Button color={buttonColor} onClick={toggle}>{icon}{buttonLabel}</Button>} */}
            <Modal isOpen={modal} toggle={toggle} className={className} centered={true}>
                <ModalHeader toggle={toggle}>{props.title}</ModalHeader>
                <ModalBody>
                    {props.context}
                </ModalBody>
                <ModalFooter>
                    <Button color="secondary" onClick={toggle}>{props.left !== undefined ? props.left : i18n.t(`${packageNS}:tr000424`)}</Button>{' '}
                    <Button color="primary" onClick={proc}>{props.right !== undefined ? props.right : i18n.t(`${packageNS}:tr000425`)}</Button>
                </ModalFooter>
            </Modal>
        </div>
    );
}

export default CommonModal;
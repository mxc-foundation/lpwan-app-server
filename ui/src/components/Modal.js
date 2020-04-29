import React, { useState } from 'react';
import { Button, Modal, ModalBody, ModalFooter, ModalHeader } from 'reactstrap';
import i18n, { packageNS } from '../i18n';

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
            closeModal(1);
    }
    
    const proc = () => {
        setModal(!modal);
        if (closeModal)
            closeModal(1);
        props.callback();
    }

    return (
        <div>
            {/* {buttonLabel && <Button color={buttonColor} onClick={toggle}>{icon}{buttonLabel}</Button>} */}
            <Modal isOpen={modal} toggle={toggle} className={className} centered={true}>
                {props.title ? <ModalHeader toggle={toggle}>{props.title}</ModalHeader>: null}
                <ModalBody>
                    {props.context}
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
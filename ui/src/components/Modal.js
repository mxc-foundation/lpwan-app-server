import React, { useState } from 'react';
import i18n, { packageNS } from '../i18n';
import { Button, Modal, ModalHeader, ModalBody, ModalFooter } from 'reactstrap';

const CommonModal = (props) => {
    const {
        className,
        showCloseButton = true,
        showConfirmButton = true,
        show = true,
    } = props;

    const [modal, setModal] = useState(show);

    const toggle = () => setModal(!modal);
    const proc = () => {
        setModal(!modal);
        props.callback();
    }


    return (
        <div>
            {/* {buttonLabel && <Button color={buttonColor} onClick={toggle}>{icon}{buttonLabel}</Button>} */}
            <Modal isOpen={modal} toggle={toggle} className={className} centered={true}>
                <ModalHeader toggle={toggle}>{props.title}</ModalHeader>
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
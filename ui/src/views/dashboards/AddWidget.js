import React, { useState } from 'react';
import { Button, Modal, ModalHeader, ModalBody, Row, Col, Media } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';

const AddWidget = (props) => {
    const { closeModal, show = true, availableWidgets = [] } = props;
    const [modal, setModal] = useState(show);

    const toggle = () => {
        setModal(!modal);
        if (closeModal) closeModal();
    }

    return (<>
        <Modal isOpen={modal} toggle={toggle} centered={true} size="lg" className="modal-dialog-scrollable">
            <ModalHeader toggle={toggle}>{i18n.t(`${packageNS}:menu.dashboard.addWidget.title`)}</ModalHeader>
            <ModalBody>
                {availableWidgets.map((widget, idx) => {
                    return <div className="p-2 border-bottom" key={idx}>
                        <Row>
                            <Col className="mb-0">
                                <Media>
                                    <Media left className="align-self-center">
                                        <div className="avatar-lg pr-3">
                                            <img src={widget.avatar} alt="" className="img-fluid" />
                                        </div>
                                    </Media>
                                    <Media body className="pr-2">
                                        <Media heading tag="h5">{widget.label}</Media>
                                        <p>{widget.description}</p>
                                    </Media>
                                    <Media right className="align-self-center">
                                        <Button color="primary" size="sm">{i18n.t(`${packageNS}:menu.dashboard.addWidget.add`)}</Button>
                                    </Media>
                                </Media>
                            </Col>
                        </Row>
                    </div>
                })}
            </ModalBody>
        </Modal>
    </>
    );
}

export default AddWidget;
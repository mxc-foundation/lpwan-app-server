import classNames from "classnames";
import React, { useState } from 'react';
import { Button, Col, Media, Modal, ModalBody, ModalHeader, Row } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';


const AddWidget = (props) => {
    const { addWidget, closeModal, show = true, availableWidgets = [], addedWidgets = [] } = props;
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
                    return <div className={classNames("p-2", {"border-bottom": idx + 1 < availableWidgets.length})} key={idx}>
                        <Row>
                            <Col className="mb-0">
                                <Media className="align-items-center">
                                    <img src={widget.avatar} alt="" className="img-fluid border rounded avatar-xl mr-3" />
                                    <Media body className="mr-2">
                                        <Media heading tag="h5">{widget.label}</Media>
                                        <p>{widget.description}</p>
                                    </Media>
                                    <Button color="primary" size="sm"
                                        onClick={() => addWidget(widget)}
                                        disabled={addedWidgets.findIndex(w => w.name === widget.name) !== -1}>
                                        {i18n.t(`${packageNS}:menu.dashboard.addWidget.add`)}</Button>
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
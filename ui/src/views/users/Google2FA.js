import React, { useState } from "react";
import { Row, Col, Button, Input } from 'reactstrap';

import QRCode from "qrcode.react";
import i18n, { packageNS } from '../../i18n';


const Google = ({ title, code, confirm, skip, titleClass = "" }) => {
    const [confirmCode, setconfirmCode] = useState("");

    return <React.Fragment>
        <Row className="text-center">
            <Col className="mb-0">
                <h5 className={titleClass}>{title}</h5>

                <Row className="mt-3 text-center">
                    <Col>
                        <QRCode value={code} size={256} level={'H'} />

                        <Input type="text" name="confirm-code" value={confirmCode} onChange={(e) => setconfirmCode(e.target.value)} className="mt-2" />
                    </Col>
                </Row>

                <Button color="primary" className="btn-block mt-2" onClick={(e) => confirm(confirmCode)}
                    disabled={!confirmCode}>
                    {i18n.t(`${packageNS}:menu.google_2fa.confirm_button`)}</Button>
                <Button color="link" className="btn-block" onClick={skip}>{i18n.t(`${packageNS}:menu.google_2fa.skip_button`)}</Button>
            </Col>
        </Row>
    </React.Fragment>
}

export default Google;

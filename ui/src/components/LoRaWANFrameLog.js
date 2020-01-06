import React, { useState } from "react";
import { Card, CardBody, Row, Col, Collapse, NavLink } from 'reactstrap';
import moment from "moment";

import JSONTree from "./JSONTree";


const LoRaWANFrameLog = (props) => {
  const [isOpen, setIsOpen] = useState(false);
  const toggle = () => setIsOpen(!isOpen);

  let dir = "";
  let devID = "";

  if (props.frame.uplinkMetaData !== undefined) {
    dir = "UPLINK";
  } else {
    dir = "DOWNLINK";
  }

  const receivedAt = moment(props.frame.receivedAt).format("LTS");
  const mType = props.frame.phyPayload.mhdr.mType;

  if (props.frame.phyPayload.macPayload !== undefined) {
    if (props.frame.phyPayload.macPayload.devEUI !== undefined) {
      devID = props.frame.phyPayload.macPayload.devEUI;
    }

    if (props.frame.phyPayload.macPayload.fhdr !== undefined) {
      devID = props.frame.phyPayload.macPayload.fhdr.devAddr;
    }
  }

  return (<React.Fragment>

    <Card className="border shadow-none mb-1">
      <NavLink className="d-block pt-2 pb-2 text-dark" href="#" onClick={toggle}>
        <div className="d-flex">
          <div className="d-flex flex-row align-items-center">
            <div className="px-2">{dir}</div>
            <div className="px-2">{receivedAt}</div>
            <div className="px-2">{mType}</div>
            <div className="px-2">{devID}</div>
          </div>
          <div className="ml-auto">
            {!isOpen && <i className="mdi mdi-chevron-down font-20"></i>}
            {isOpen && <i className="mdi mdi-chevron-up font-20"></i>}
          </div>
        </div>
      </NavLink>

      <Collapse isOpen={isOpen}>
        <CardBody>
          <Row>
            <Col>
              <div className="p-2 border">
                {props.frame.uplinkMetaData && <JSONTree data={props.frame.uplinkMetaData} />}
                {props.frame.downlinkMetaData && <JSONTree data={props.frame.downlinkMetaData} />}
              </div>
            </Col>
            <Col>
              <div className="p-2 border">
                <JSONTree data={{ phyPayload: props.frame.phyPayload }} />
              </div>
            </Col>
          </Row>
        </CardBody>
      </Collapse>
    </Card>
  </React.Fragment>
  );
}

export default LoRaWANFrameLog;

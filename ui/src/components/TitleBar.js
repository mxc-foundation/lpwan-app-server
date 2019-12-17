import React, { Component } from "react";

import { Row, Col } from 'reactstrap';


class TitleBar extends Component {
  render() {
    return (
      <Row className="pt-3 pb-2">
        <Col lg={6}>
          {this.props.children}
        </Col>

        <Col lg={6} className="text-right">
          {this.props.buttons}
        </Col>
      </Row>
    );
  }
}

export default TitleBar;

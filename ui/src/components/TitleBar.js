import React, { Component } from "react";

import { Row, Col } from 'reactstrap';


class TitleBar extends Component {
  render() {
    const { noButtons } = this.props;

    return (
      <Row className="pt-3 pb-2">
        <Col xs={noButtons ? 12 : 8} style={{ display: "flex" }}>
          {this.props.children}
        </Col>
        {noButtons ? null
          : (
            <Col xs={4} className="text-right">
              {this.props.buttons}
            </Col>
          )
        }
      </Row>
    );
  }
}

export default TitleBar;

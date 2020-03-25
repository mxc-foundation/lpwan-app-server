import React, { Component } from "react";
import { Col, Row } from 'reactstrap';



class TitleBar extends Component {
  render() {
    const { noButtons } = this.props;

    return (
      <Row className="pt-3 pb-2">
        <Col md={noButtons ? 12 : 8} style={{ display: "flex" }}>
          {this.props.children}
        </Col>
        {noButtons ? null
          : (
            <Col md={4} className="text-md-right">
              {this.props.buttons}
            </Col>
          )
        }
      </Row>
    );
  }
}

export default TitleBar;

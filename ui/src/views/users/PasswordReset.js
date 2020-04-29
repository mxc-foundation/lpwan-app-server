import React, { Component } from "react";
import { Button, Card, CardBody, Col, Row } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';




class PasswordReset extends Component {
  constructor(props) {
    super(props);
    this.state = {
      object: props.object || {},
    };
  }

  render() {
    const { object } = this.state;

    if (object === undefined) {
      return null;
    }

    return (
      <React.Fragment>
        <Card className="h-auto">
          <CardBody className="pb-0">
            <div className="card-coming-soon">
              <h5>{i18n.t(`${packageNS}:menu.profile_password_reset.coming_soon`)}</h5>
            </div>

            <h5>{i18n.t(`${packageNS}:menu.profile_password_reset.title`)}</h5>
            <Row>
              <Col>
                <p className="mt-2">{i18n.t(`${packageNS}:menu.profile_password_reset.account_password.title`)}</p>
              </Col>
              <Col className="text-right">
                  <Button color="primary" outline>{i18n.t(`${packageNS}:menu.profile_password_reset.account_password.reset_button`)}</Button>
              </Col>
            </Row>
          </CardBody>
        </Card>
      </React.Fragment>
    );
  }
}

export default PasswordReset;

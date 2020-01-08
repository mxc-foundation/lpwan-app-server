import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import { Container, Row, Card, CardBody } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import UserStore from "../../stores/UserStore";
import UserForm from "./UserForm";

class UpdateUser extends Component {
  constructor() {
    super();
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(user) {
    UserStore.update(user, resp => {
      this.props.history.push("/users");
    });
  }

  render() {
    return (
      <Container fluid>
        <Row xs="1" lg="1">
          <Card>
            <CardBody>
              <UserForm
                submitLabel={i18n.t(`${packageNS}:tr000066`)}
                object={this.props.user}
                onSubmit={this.onSubmit}
              />
            </CardBody>
          </Card>
        </Row>
      </Container>
    );
  }
}

export default withRouter(UpdateUser);

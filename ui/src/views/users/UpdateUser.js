import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import { Container, Row, Card, CardBody } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import UserStore from "../../stores/UserStore";
import UserForm from "./UserForm";

class UpdateUser extends Component {
  onSubmit = (user) => {
    UserStore.update(user, resp => {
      this.props.history.push("/users");
    });
  }

  render() {
    const { loading, user } = this.props;

    return (
      <React.Fragment>
        <Container>
          <Row xs="1" lg="1">
            <Card>
              <CardBody>
                <UserForm
                  submitLabel={i18n.t(`${packageNS}:tr000066`)}
                  loading={loading}
                  object={user}
                  onSubmit={this.onSubmit}
                  update={true}
                />
              </CardBody>
            </Card>
          </Row>
        </Container>
      </React.Fragment>
    );
  }
}

export default withRouter(UpdateUser);

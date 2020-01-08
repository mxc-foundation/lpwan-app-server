import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Breadcrumb, BreadcrumbItem, Form, FormGroup, Label, Input, Button, Container, Row, Col, Card, CardBody } from 'reactstrap';

import TitleBar from "../../components/TitleBar";
import UserStore from "../../stores/UserStore";
import FormComponent from "../../classes/FormComponent";
import i18n, { packageNS } from '../../i18n';


class PasswordForm  extends FormComponent {
  render() {
    if (this.state.object === undefined) {
      return(<div></div>);
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
      >
        <FormGroup row>
          <Label for="password" sm={2}>{i18n.t(`${packageNS}:tr000004`)}</Label>
          <Col sm={10}>
            <Input type="password" name="password" id="password" value={this.state.object.password || ""} onChange={this.onChange} />
          </Col>
        </FormGroup>
        {this.props.submitLabel && <Button color="primary"
          onClick={this.onSubmit}
          disabled={this.props.disabled}
          className="btn-block">{this.props.submitLabel}
        </Button>}
      </Form>
    );
  }
}


class ChangeUserPassword extends Component {
  constructor() {
    super();
    this.state = {};

    this.onSubmit = this.onSubmit.bind(this);
  }

  componentDidMount() {
    UserStore.get(this.props.match.params.userID, resp => {
      this.setState({
        user: resp,
      });
    });
  }

  onSubmit(password) {
    UserStore.updatePassword(this.props.match.params.userID, password.password, resp => {
      this.props.history.push("/");
    });
  }

  render() {
    if (this.state.user === undefined) {
      return(<div></div>);
    }

    return(
      <React.Fragment>
        <TitleBar>
          <Breadcrumb>
          <BreadcrumbItem><Link to={`/users`}>{i18n.t(`${packageNS}:tr000036`)}</Link></BreadcrumbItem>
            <BreadcrumbItem><Link to={`/users/${this.state.user.user.id}`}>{this.state.user.user.username}</Link></BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000038`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Container fluid>
          <Row xs="1" lg="1">
            <Card>
              <CardBody>
              <PasswordForm
                submitLabel={i18n.t(`${packageNS}:tr000022`)}
                onSubmit={this.onSubmit}
              />
              </CardBody>
            </Card>
          </Row>
        </Container>
      </React.Fragment>
    );
  }
}

export default withRouter(ChangeUserPassword);

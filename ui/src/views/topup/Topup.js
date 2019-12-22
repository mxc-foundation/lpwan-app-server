import React, { Component } from "react";
import { withRouter } from 'react-router-dom';


import { Breadcrumb, BreadcrumbItem, Row, Col, Card, CardBody } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";

import SessionStorage from "../../stores/SessionStore";
import TopupForm from "./TopupForm";
import InfoCard from "./InfoCard";

class Topup extends Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  componentDidUpdate(oldProps) {
    if (this.props === oldProps) {
      return;
    }
  }

  onSubmit = () => {
    if (SessionStorage.getUser().isAdmin) {
      this.props.history.push(`/control-panel/modify-account`);
    } else {
      this.props.history.push(`/modify-account/${this.props.match.params.organizationID}`);
    }
  }

  render() {
    const path = `/modify-account/${this.props.match.params.organizationID}`;

    return (<React.Fragment>
      <TitleBar>
        <Breadcrumb>
          <BreadcrumbItem active>{i18n.t(`${packageNS}:menu.topup.topup`)}</BreadcrumbItem>
        </Breadcrumb>
      </TitleBar>

      <Row>
        <Col>
          <Card>
            <CardBody>
              <TopupForm
                reps={this.state.accounts} {...this.props}
                orgId={this.props.match.params.organizationID}
                path={path}
              />

            </CardBody>
          </Card>
        </Col>
        <Col>
          <InfoCard path={path} />
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default withRouter(Topup);
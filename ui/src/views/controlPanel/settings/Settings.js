import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem, Card, CardBody, Col, Row } from 'reactstrap';
import TitleBar from "../../../components/TitleBar";
import i18n, { packageNS } from '../../../i18n';
import SettingsForm from "./SettingsForm";



class Settings extends Component {
  constructor(props) {
    super(props);

    this.state = {};
  }

  onSubmit = (e, data) => {
    e.preventDefault();
  }

  render() {

    return (
      <React.Fragment>
        <TitleBar>
          <Breadcrumb>
            <BreadcrumbItem>
              <Link
                to={`/organizations`}
                onClick={() => {
                  // Change the sidebar content
                  this.props.switchToSidebarId('DEFAULT');
                }}
              >
                {i18n.t(`${packageNS}:menu.control_panel`)}
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem>{i18n.t(`${packageNS}:tr000451`)}</BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:menu.settings.system_settings`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          <Col>
            <Card>
              <CardBody>
                <SettingsForm
                  submitLabel={i18n.t(`${packageNS}:menu.withdraw.confirm`)}
                  onSubmit={this.onSubmit}
                />
              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withRouter(Settings);

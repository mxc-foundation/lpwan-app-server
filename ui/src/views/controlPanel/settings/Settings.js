import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Breadcrumb, BreadcrumbItem, Row } from 'reactstrap';
import i18n, { packageNS } from '../../../i18n';
import TitleBar from "../../../components/TitleBar";
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
            <BreadcrumbItem active>{i18n.t(`${packageNS}:menu.settings.system_settings`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          <SettingsForm
            submitLabel={i18n.t(`${packageNS}:menu.withdraw.confirm`)}
            onSubmit={this.onSubmit}
          />
        </Row>
      </React.Fragment>
    );
  }
}

export default withRouter(Settings);
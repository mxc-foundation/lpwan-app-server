import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Breadcrumb, BreadcrumbItem, Row } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import i18n, { packageNS } from '../../../i18n';
import TitleBar from "../../../components/TitleBar";
import SettingsForm from "./SettingsForm";

import breadcrumbStyles from "../../common/BreadcrumbStyles";
import Admin from "../../../components/Admin";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class Settings extends Component {
  constructor(props) {
    super(props);

    this.state = {};
  }

  onSubmit = (e, data) => {
    e.preventDefault();
  }

  render() {
    const { classes } = this.props;

    return (
      <React.Fragment>
        <TitleBar>
          <Breadcrumb className={classes.breadcrumb}>
            <Admin>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations`}
                onClick={() => {
                  // Change the sidebar content
                  this.props.switchToSidebarId('DEFAULT');
                }}
              >
                Control Panel
              </Link>
            </BreadcrumbItem>
            </Admin>
            <BreadcrumbItem className={classes.breadcrumbItem}>{i18n.t(`${packageNS}:tr000451`)}</BreadcrumbItem>
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

export default withStyles(styles)(withRouter(Settings));

import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { withStyles } from "@material-ui/core/styles";
import { Breadcrumb, BreadcrumbItem, Row, Col, Card, CardBody } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";

import SessionStorage from "../../stores/SessionStore";
import TopupForm from "./TopupForm";
import InfoCard from "./InfoCard";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

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
    const { classes } = this.props;

    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const path = `/modify-account/${this.props.match.params.organizationID}`;

    return (<React.Fragment>
      <TitleBar>
        <Breadcrumb className={classes.breadcrumb}>
          <BreadcrumbItem>
            <Link
              className={classes.breadcrumbItemLink}
              to={`/organizations`}
              onClick={() => { this.props.switchToSidebarId('DEFAULT'); }}
            >
                Organizations
            </Link>
          </BreadcrumbItem>
          <BreadcrumbItem>
            <Link
              className={classes.breadcrumbItemLink}
              to={`/organizations/${currentOrgID}`}
              onClick={() => { this.props.switchToSidebarId('DEFAULT'); }}
            >
              {currentOrgID}
            </Link>
          </BreadcrumbItem>
          <BreadcrumbItem className={classes.breadcrumbItem}>Wallet</BreadcrumbItem>
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

export default withStyles(styles)(withRouter(Topup));
import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import { Card, CardBody, Col, Row } from 'reactstrap';
import Loader from "../../components/Loader";
import OrgBreadCumb from '../../components/OrgBreadcrumb';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import NetworkServerStore from "../../stores/NetworkServerStore";
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import ServiceProfileForm from "./ServiceProfileForm";





class CreateServiceProfile extends Component {
  constructor() {
    super();
    this.state = {
      nsDialog: false,
    };
    this.onSubmit = this.onSubmit.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
  }

  componentDidMount() {
    this.loadData();
  }

  loadData = async () => {
    const res = await NetworkServerStore.list(0, 10, 0);
    
    if (res.totalCount === "0") {
      this.setState({
        nsDialog: true,
      });
    }
  }

  closeDialog() {
    this.setState({
      nsDialog: false,
    });
  }

  onSubmit = async (serviceProfile) => {
    let sp = serviceProfile;
    sp.organizationID = this.props.match.params.organizationID;

    this.setState({ loading: true });
    const res = await ServiceProfileStore.create(sp);
    
    this.setState({ loading: false });
    this.props.history.push(`/organizations/${this.props.match.params.organizationID}/service-profiles`);
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (
      <React.Fragment>
        <TitleBar>
          <OrgBreadCumb organizationID={currentOrgID} items={[
              { label: i18n.t(`${packageNS}:tr000078`), active: false, to: `/organizations/${currentOrgID}/service-profiles` },
              { label: i18n.t(`${packageNS}:tr000277`), active: false }]}></OrgBreadCumb>
        </TitleBar>
        <Row>
          <Col>
            <Card>
              <CardBody>
              {this.state.loading && <Loader />}

                <ServiceProfileForm
                  match={this.props.match}
                  submitLabel={i18n.t(`${packageNS}:tr000277`)}
                  onSubmit={this.onSubmit}
                  object={{}}
                />
              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withRouter(CreateServiceProfile);

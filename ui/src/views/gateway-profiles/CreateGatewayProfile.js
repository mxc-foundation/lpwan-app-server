import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem, Card, CardBody, Col, Form, Row } from 'reactstrap';
import Loader from "../../components/Loader";
import Modal from '../../components/Modal';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import GatewayProfileStore from "../../stores/GatewayProfileStore";
import NetworkServerStore from "../../stores/NetworkServerStore";
import GatewayProfileForm from "./GatewayProfileForm";







class CreateGatewayProfile extends Component {
  constructor() {
    super();
    this.state = {
      nsDialog: false,
      loading: false,
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

  closeDialog = () => {
    this.setState({
      nsDialog: false,
    });
  }

  onSubmit(gatewayProfile) {
    this.setState({ loading: true });
    GatewayProfileStore.create(gatewayProfile, resp => {
      this.setState({ loading: false });
      this.props.history.push("/gateway-profiles");
    }, error => { this.setState({ loading: false }) });
  }

  render() {

    return (<>
      <TitleBar>
        <Breadcrumb>
          <BreadcrumbItem>{i18n.t(`${packageNS}:menu.control_panel`)}</BreadcrumbItem>
          <BreadcrumbItem><Link to={`/gateway-profiles`}>{i18n.t(`${packageNS}:tr000046`)}</Link></BreadcrumbItem>
          <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000277`)}</BreadcrumbItem>
        </Breadcrumb>
      </TitleBar>

      <Form>
        {this.state.nsDialog && <Modal
          title={""}
          left={"DISMISS"}
          right={"ADD"}
          closeModal={this.closeDialog}
          context={i18n.t(`${packageNS}:tr000377`)}
          callback={this.deleteGatewayProfile} />}
        <Row>
          <Col>
            <Card>
              <CardBody>
                {this.state.loading && <Loader />}
                <GatewayProfileForm
                  submitLabel={i18n.t(`${packageNS}:tr000277`)}
                  onSubmit={this.onSubmit}
                />
              </CardBody>
            </Card>
          </Col>
        </Row>
      </Form>
    </>
    );
  }
}

export default withRouter(CreateGatewayProfile);

import React from "react";
import { withRouter } from "react-router-dom";

import { Row, Col, Card, CardBody, CardFooter } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import TitleBar from "../../components/TitleBar";

import ExtLink from "../../components/ExtLink";
import OrgBreadCumb from '../../components/OrgBreadcrumb';
import Typography from '@material-ui/core/Typography';
import StakeForm from "./StakeForm";
import StakeStore from "../../stores/StakeStore";

//import Modal from "../common/Modal";
import ModalTimer from "../common/ModalTimer";
//import Button from "@material-ui/core/Button";
import Spinner from "../../components/ScaleLoader";
import { EXT_URL_STAKE } from "../../util/Data"
import InfoCard from "../topup/InfoCard";
import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class SetStake extends FormComponent {

  state = {
    title: i18n.t(`${packageNS}:menu.messages.set_stake`)
  }

  componentDidMount() {
    this.loadStakeTextTranslation();
  }

  loadStakeTextTranslation = () => {
    this.setState({
      info: i18n.t(`${packageNS}:menu.messages.staking_enhances`)
    })
  }

  onChange = (event, name) => {
    this.setState({
      [name]: event.target.value
    });
  }

  setTitle = (isUnstake) => {
    const title = isUnstake ? i18n.t(`${packageNS}:menu.messages.unstake`) : i18n.t(`${packageNS}:menu.messages.set_stake`);
    this.setState({
      title
    })
    this.forceUpdate();
  }

  setInfo = (info) => {
    const object = this.state;
    object.info = info.text;
    object.infoStatus = info.status;
    this.setState({
      object
    })
  }

  render() {
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    let path = null;
    if (currentOrgID === process.env.REACT_APP_SUPER_ADMIN_LPWAN) {
      path = '/control-panel/modify-account/';
    } else {
      path = `/modify-account/${currentOrgID}`;
    }

    return (
      <>
        <TitleBar>
        <OrgBreadCumb organizationID={currentOrgID} items={[
            { label: i18n.t(`${packageNS}:menu.staking.set_stake`), active: false },
            { label: this.state.title, active: true }]}></OrgBreadCumb>
        </TitleBar>

        <Row xs="1">
          <Col>
            <Card>
              <CardBody>
                {this.state.info}
              </CardBody>
              <CardFooter>
                <ExtLink to={EXT_URL_STAKE} context={i18n.t(`${packageNS}:menu.common.learn_more`)} />
              </CardFooter>
            </Card>
          </Col>
        </Row>
        <Row xs="1" lg="2">
          <Col>
            <Card>
              <CardBody>
                <StakeForm setTitle={this.setTitle} />
              </CardBody>
            </Card>
          </Col>
          <Col><InfoCard path={path} /></Col>
        </Row>
      </>
    );
  }
}

export default withStyles(styles)(withRouter(SetStake));

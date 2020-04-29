import { withStyles } from "@material-ui/core/styles";
import React from "react";
import { withRouter } from "react-router-dom";
import { Card, CardBody, CardFooter, Col, Row } from 'reactstrap';
import FormComponent from "../../classes/FormComponent";
import ExtLink from "../../components/ExtLink";
import OrgBreadCumb from '../../components/OrgBreadcrumb';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import { EXT_URL_STAKE } from "../../util/Data";
import breadcrumbStyles from "../common/BreadcrumbStyles";
import StakeForm from "./StakeForm";
import StakeHistory from "./StakeHistory";





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
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

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
        <Row xs="1" lg="1">
          <Col>
            <Card>
              <CardBody>
                <StakeForm setTitle={this.setTitle} />
              </CardBody>
            </Card>
          </Col>
          {/* <Col><InfoCard path={path} /></Col> */}
        </Row>
        <Row xs="1" lg="1">
          <Col>
            <Card>
              <CardBody>
                <StakeHistory />
              </CardBody>
            </Card>
          </Col>
        </Row>
      </>
    );
  }
}

export default withStyles(styles)(withRouter(SetStake));

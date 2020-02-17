import React, { Component } from "react";
import { withRouter } from "react-router-dom";

import { Row, Col } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';
import { ReactstrapInput } from '../../components/FormInputs';

import Modal from "../../components/Modal";
import Spinner from "../../components/ScaleLoader";
import TopupStore from "../../stores/TopupStore";
import MoneyStore from "../../stores/MoneyStore";
import SessionStorage from "../../stores/SessionStore";
import { ETHER } from "../../util/CoinType"

import i18n, { packageNS } from '../../i18n';

function loadSuperNodeActiveMoneyAccount(organizationID) {
  return new Promise((resolve, reject) => {
    TopupStore.getTopUpDestination(ETHER, organizationID, resp => {
      resolve(resp.activeAccount);
    });

  });
}

function loadActiveMoneyAccount(organizationID) {
  return new Promise((resolve, reject) => {
    MoneyStore.getActiveMoneyAccount(ETHER, organizationID, resp => {
      resolve(resp.activeAccount);
    });
  });
}

class TopupForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      nsDialog: false,
      description: null,
      object: this.props.object || {},
    };

  }

  componentDidMount() {
    this.loadData();
  }

  componentDidUpdate(oldProps) {
    if (this.props === oldProps) {
      return;
    }

    this.loadData();
  }

  loadData = async () => {
    try {
      const organizationID = this.props.match.params.organizationID;

      this.setState({ loading: true })
      var superNodeAccount = await loadSuperNodeActiveMoneyAccount(organizationID);
      var account = await loadActiveMoneyAccount(organizationID);

      const accounts = {};
      accounts.superNodeAccount = superNodeAccount;
      accounts.account = account;


      let object = this.state.object;
      object.accounts = {
        superNodeAccount: superNodeAccount,
        account: account,
      }

      if (SessionStorage.getUser().isAdmin && !superNodeAccount) {
        this.showModal(true);
      }

      if (!accounts.account && !SessionStorage.getUser().isAdmin) {
        this.showModal(true);
      }

      let description = '';
      if (SessionStorage.getUser().isAdmin) {
        description = i18n.t(`${packageNS}:menu.topup.notice001`) + " " + i18n.t(`${packageNS}:menu.topup.notice003`);
      } else {
        description = i18n.t(`${packageNS}:menu.topup.notice002`) + " " + i18n.t(`${packageNS}:menu.topup.notice003`);
      }
      
      this.setState({
        object: object,
        description
      });

      this.setState({ loading: false })
    } catch (error) {
      this.setState({ loading: false })
      console.error(error);
      this.setState({ error });
    }
  }

  showModal = (nsDialog) => {
    this.setState({ nsDialog });
  }

  handleLink = () => {
    //window.location.replace(`http://wallet.mxc.org/`);
    this.props.history.push(this.props.path);
  }

  render() {
    let fieldsSchema = {
      accounts: Yup.object().shape({
        superNodeAccount: Yup.string().trim(),
        account: Yup.string().trim()
      })
    }

    const formSchema = Yup.object().shape(fieldsSchema);

    if (this.state.object === undefined) {
      return (<div></div>);
    }

    let accounts = {};
    if (this.state.object.accounts !== undefined) {
      if (this.state.object.accounts.superNodeAccount !== undefined) {
        accounts.superNodeAccount = this.state.object.accounts.superNodeAccount;
      } else {
        accounts.superNodeAccount = '';
      }
    }

    if (this.state.object.accounts !== undefined) {
      if (this.state.object.accounts.account !== undefined) {
        accounts.account = this.state.object.accounts.account;
      } else {
        accounts.account = '';
      }
    }

    return (<React.Fragment>
      <Row>
        <Col>
        {this.state.nsDialog && <Modal
          title={i18n.t(`${packageNS}:menu.topup.notice`)}
          left={"DISMISS"}
          right={"ADD ETH ACCOUNT"}
          context={this.state.description}
          callback={this.handleLink} />}
          <Formik
            enableReinitialize
            initialValues={this.state.object}
            validationSchema={formSchema}
            onSubmit={(values) => {
              const castValues = formSchema.cast(values);
              this.props.onSubmit({ ...castValues })
            }}>
            {({
              handleSubmit,
              setFieldValue,
              values,
              handleBlur,
            }) => (
                <Form onSubmit={handleSubmit} noValidate>
                  <Field
                    type="text"
                    label={i18n.t(`${packageNS}:menu.topup.from_eth_account`)}
                    name="accounts.superNodeAccount"
                    id="accounts.superNodeAccount"
                    value={accounts.account || `${i18n.t(`${packageNS}:menu.topup.can_not_find_any_account`)}`}
                    //helpText={i18n.t(`${packageNS}:tr000062`)}
                    component={ReactstrapInput}
                    onBlur={handleBlur}
                    inputProps={{
                      clearable: true,
                      cache: false,
                    }}
                  />
                  <Field
                    type="text"
                    label={i18n.t(`${packageNS}:menu.topup.to_eth_account`)}
                    name="accounts.account"
                    id="accounts.account"
                    value={accounts.superNodeAccount || `${i18n.t(`${packageNS}:menu.topup.can_not_find_any_account`)}`}
                    //helpText={i18n.t(`${packageNS}:tr000062`)}
                    component={ReactstrapInput}
                    onBlur={handleBlur}
                    inputProps={{
                      clearable: true,
                      cache: false,

                    }}
                  />
                </Form>
              )}
          </Formik>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default withRouter(TopupForm);

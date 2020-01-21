import React, { Component } from "react";

import { Button } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import { ReactstrapInput } from '../../../components/FormInputs';
import * as Yup from 'yup';

import i18n, { packageNS } from '../../../i18n';

import WithdrawStore from '../../../stores/WithdrawStore';
import { ETHER } from '../../../util/CoinType';
import SettingsStore from '../../../stores/SettingsStore';

class SettingsForm extends Component {

  constructor(props) {
    super(props);

    this.state = {
      object: {},
    };
  }

  componentDidMount() {
    this.loadSettings();
  }

  loadSettings = async () => {
    try {
      const organizationID = 0;
      //this.setState({loading: true})

      WithdrawStore.getWithdrawFee(ETHER, organizationID, (resp) => {
        this.setState({ object: { withdrawFee: resp.withdrawFee } });
      });

      SettingsStore.getSystemSettings((resp) => {
        this.setState({
          object: {
            downlinkPrice: resp.downlinkFee,
            percentageShare: resp.transactionPercentageShare,
            lbWarning: resp.lowBalanceWarning
          }
        });
      });
    } catch (e) {
      console.log("Error", e)
    }
  };

  saveSettings = async () => {
    try {
      let bodyWF = {
        moneyAbbr: 'Ether',
        orgId: '0',
        withdrawFee: this.state.object.withdrawFee
      };

      let bodySettings = {
        downlinkFee: this.state.object.downlinkPrice,
        lowBalanceWarning: this.state.object.lbWarning,
        transactionPercentageShare: this.state.object.percentageShare
      };

      WithdrawStore.setWithdrawFee(ETHER, 0, bodyWF, (resp) => { });

      SettingsStore.setSystemSettings(bodySettings, (resp) => { });
    } catch (e) {
      console.log("Error", e)
    }
  };

  reset = () => {
    this.loadSettings();
  }

  handleChange = (name, event) => {
    this.setState({
      [name]: event.target.value
    });
  };

  render() {
    
    if (this.state.object === undefined) {
      return (<div></div>);
    }

    let fieldsSchema = {
      withdrawFee: Yup.string().trim(),
      downlinkPrice: Yup.string().trim(),
      percentageShare: Yup.string().trim(),
      lbWarning: Yup.string().trim(),
    }

    const formSchema = Yup.object().shape(fieldsSchema);

    return (
      <React.Fragment>
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
            handleChange,
            setFieldValue,
            values,
            handleBlur,
          }) => (
              <Form onSubmit={handleSubmit} noValidate>
                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:menu.settings.withdraw_fee`)}
                  name="withdrawFee"
                  id="withdrawFee"
                  value={this.state.object.withdrawFee || ""}
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                  readOnly
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />

                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:menu.settings.downlink_price`)}
                  name="downlinkPrice"
                  id="downlinkPrice"
                  value={this.state.object.downlinkPrice || ""}
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                  readOnly
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />

                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:menu.withdraw.transaction_fee`)}
                  name="percentageShare"
                  id="percentageShare"
                  value={this.state.object.percentageShare || ""}
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                  readOnly
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />

                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:menu.settings.low_balance`)}
                  name="lbWarning"
                  id="lbWarning"
                  value={this.state.object.lbWarning || ""}
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                  readOnly
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />

                {/* <Button className="btn-block" onClick={this.reset}>{i18n.t(`${packageNS}:common.reset`)}</Button>
                <Button type="submit" className="btn-block" color="primary">{this.props.submitLabel || i18n.t(`${packageNS}:tr000066`)}</Button> */}
              </Form>
            )}
        </Formik>
      </React.Fragment>
    );
  }
}

export default SettingsForm;

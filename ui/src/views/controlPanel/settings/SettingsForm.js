import { Field, Form, Formik } from 'formik';
import React, { Component } from "react";
import * as Yup from 'yup';
import { ReactstrapInput } from '../../../components/FormInputs';
import i18n, { packageNS } from '../../../i18n';
import SettingsStore from '../../../stores/SettingsStore';
import WithdrawStore from '../../../stores/WithdrawStore';
import { ETHER } from '../../../util/CoinType';




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
      //this.setState({loading: true})

      WithdrawStore.getWithdrawFee(ETHER, (resp) => {
        this.setState({ object: { withdrawFee: resp.withdrawFee } });
      });

      SettingsStore.getSystemSettings((resp) => {
        this.setState({
          object: {
            lowBalanceWarning: resp.lowBalanceWarning,
            downlinkPrice: resp.downlinkPrice,
            supernodeIncomeRatio: resp.supernodeIncomeRatio,
            stakingPercentage: resp.stakingPercentage,
            stakingExpectedRevenuePercentage: resp.stakingExpectedRevenuePercentage,
            compensation: resp.compensation
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
        lowBalanceWarning: this.state.object.lowBalanceWarning,
        downlinkPrice: this.state.object.downlinkPrice,
        supernodeIncomeRatio: this.state.object.supernodeIncomeRatio,
        stakingPercentage: this.state.object.stakingPercentage,
        stakingExpectedRevenuePercentage: this.state.object.stakingExpectedRevenuePercentage,
        compensation: this.state.object.compensation
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
      lowBalanceWarning: Yup.string().trim(),
      downlinkPrice: Yup.string().trim(),
      supernodeIncomeRatio: Yup.string().trim(),
      stakingPercentage: Yup.string().trim(),
      stakingExpectedRevenuePercentage: Yup.string().trim(),
      compensation: Yup.string().trim(),
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
                    label={i18n.t(`${packageNS}:menu.settings.low_balance`)}
                    name="lowBalanceWarning"
                    id="lowBalanceWarning"
                    value={this.state.object.lowBalanceWarning || ""}
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
                    label={i18n.t(`${packageNS}:menu.settings.compensation`)}
                    name="compensation"
                    id="compensation"
                    value={this.state.object.compensation || ""}
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
                  label={i18n.t(`${packageNS}:menu.settings.supernode_income_ratio`)}
                  name="supernodeIncomeRatio"
                  id="supernodeIncomeRatio"
                  value={this.state.object.supernodeIncomeRatio || ""}
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
                    label={i18n.t(`${packageNS}:menu.settings.percentage_share`)}
                    name="stakingPercentage"
                    id="stakingPercentage"
                    value={this.state.object.stakingPercentage || ""}
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
                  label={i18n.t(`${packageNS}:menu.settings.expected_percentage`)}
                  name="stakingExpectedRevenuePercentage"
                  id="stakingExpectedRevenuePercentage"
                  value={this.state.object.stakingExpectedRevenuePercentage || ""}
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

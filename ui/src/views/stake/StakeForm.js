import React, { Component } from "react";

import { Button, FormGroup, Label, FormText, Card, CardBody } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';
import { ReactstrapInput } from '../../components/FormInputs';
import i18n, { packageNS } from '../../i18n';

import StakeStore from "../../stores/StakeStore";
//import Spinner from "../../components/ScaleLoader"
import { withRouter } from "react-router-dom";
import { withStyles } from "@material-ui/core/styles";

import NumberFormat from 'react-number-format';
import styles from "./StakeStyle"

const NumberFormatMXC = (props) => {
  const { inputRef, onChange, ...other } = props;

  return (
    <NumberFormat
      {...other}
      getInputRef={inputRef}
      onValueChange={(values) => {
        onChange({
          target: {
            value: values.value
          }
        });
      }}
      suffix=" MXC"
    />
  );
}

class StakeForm extends Component {

  constructor(props) {
    super(props);

    this.state = {
      object: {
        amount: 0,
        revRate: 0
      }
    };
  }

  componentWillReceiveProps(nextProps) {
    this.setState({ amount: nextProps.amount });  
  }

  componentDidMount() {
    if (this.props.amount > 0) {
      this.state.object.amount = this.props.amount;
      return;
    }
    this.loadData();
  }

  componentDidUpdate(oldProps) {
    if (this.props.revRate === oldProps.revRate) {
      return;
    }

    this.loadData();
  }

  loadData = async () => {
    let res = await StakeStore.getActiveStakes(this.props.match.params.organizationID);
    let amount = 0;

    if (res.actStake !== null) {
      amount = res.actStake.Amount;
    }

    res = await StakeStore.getStakingPercentage(this.props.match.params.organizationID);
    let revRate = 0;
    revRate = res.stakingPercentage;
    
    this.setState({
      object: {
        amount,
        revRate,
      }
    });
  }

  onChange = (event) => {
    const { id, value } = event.target;

    this.setState({
      object: {
        [id]: value
      }
    });
  }

  reset = () => {
    this.props.reset();
  }

  fix = () => {
    
  }
  render() {
    let fieldsSchema = {
      amount: Yup.number(),
      revRate: Yup.number(),
    }

    const formSchema = Yup.object().shape(fieldsSchema);
    
    return (
      <React.Fragment>
        <Formik
          enableReinitialize
          initialValues={this.state.object}
          validationSchema={formSchema}
          onSubmit={this.props.confirm}>
          {({
            handleSubmit,
            handleChange,
            setFieldValue,
            values,
            handleBlur,
          }) => (
              <Form onSubmit={handleSubmit} noValidate>
                <Field
                  type="number"
                  label={i18n.t(`${packageNS}:menu.common.amount`)}
                  name="amount"
                  id="amount"
                  value={this.state.object.amount || ""}
                  autoComplete='off'
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                  onChange={handleChange}
                  readOnly={this.props.isUnstake}
                  min={0}
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />

                <Field
                  type="number"
                  label={i18n.t(`${packageNS}:menu.messages.revenue_rate`)}
                  name="revRate"
                  id="revRate"
                  value={this.state.object.revRate || ""}
                  component={ReactstrapInput}
                  onChange={handleChange}
                  onBlur={handleBlur}
                  readOnly
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />

                <Button className="btn-block" onClick={this.reset}>{i18n.t(`${packageNS}:common.reset`)}</Button>
                <Button type="submit" className="btn-block" color="primary">{this.props.isUnstake ? i18n.t(`${packageNS}:menu.messages.confirm_unstake`) : i18n.t(`${packageNS}:menu.messages.confirm_stake`)}</Button>
              </Form>
            )}
        </Formik>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(StakeForm));

import React from "react";

import TextField from '@material-ui/core/TextField';
import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import TitleBarTitle from "../../components/TitleBarTitle";
import { withRouter, Link  } from "react-router-dom";
import Button from "@material-ui/core/Button";

class TopupForm extends FormComponent {

  handleOpenAXS = () => {
    window.location.replace(`http://wallet.mxc.org/`);
  } 

  render() {
    /* const extraButtons = <>
      <Button color="primary.main" onClick={this.handleOpenAXS} type="button" disabled={false}>{i18n.t(`${packageNS}:menu.topup.use_axs_wallet`)}</Button>
    </>; */
    
    if (this.props.reps === undefined) {
      return(
        <Form
          submitLabel={this.props.submitLabel}
          /* extraButtons={extraButtons} */
          onSubmit={this.onSubmit}
        >
          <TitleBarTitle component={Link} to={'#'} title={i18n.t(`${packageNS}:menu.topup.no_data_to_display`)} />
        </Form>
      );
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        /* extraButtons={extraButtons} */
        onSubmit={this.onSubmit}
      >
        <TitleBarTitle title={i18n.t(`${packageNS}:menu.topup.send_tokens`)} />
        <TextField
          id="to"
          label={i18n.t(`${packageNS}:menu.topup.from_eth_account`)}
          margin="normal"
          value={this.props.reps.account || `${i18n.t(`${packageNS}:menu.topup.can_not_find_any_account`)}` }
          InputProps={{
            readOnly: true,
          }}
          fullWidth
        />
        <TextField
          id="to"
          label={i18n.t(`${packageNS}:menu.topup.to_eth_account`)}
          margin="normal"
          value={this.props.reps.superNodeAccount || `${i18n.t(`${packageNS}:menu.topup.can_not_find_any_account`)}` }
          InputProps={{
            readOnly: true,
          }}
          fullWidth
        />
      </Form>
    );
  }
}

export default withRouter(TopupForm);

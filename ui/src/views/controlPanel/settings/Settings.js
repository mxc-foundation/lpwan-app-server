import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';
import { withStyles } from "@material-ui/core/styles";

import Grid from '@material-ui/core/Grid';
import i18n, { packageNS } from '../../../i18n';
import TitleBar from "../../../components/TitleBar";
import TitleBarTitle from "../../../components/TitleBarTitle";
import MoneyStore from "../../../stores/MoneyStore";
import SessionStore from "../../../stores/SessionStore";
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import SettingsForm from "./SettingsForm";
import styles from "./SettingsStyle";
import { ETHER } from "../../../util/CoinType";
import { SUPER_ADMIN } from "../../../util/M2mUtil";


class Settings extends Component {
    constructor(props) {
      super(props);
  
      this.state = {};
    }
  
    onSubmit = (e, data) => {
      e.preventDefault();
    }

  render() {
    return(
      <Grid container spacing={24}>
        <Grid item xs={12} md={12} lg={12} className={this.props.classes.divider}>
          <div className={this.props.classes.TitleBar}>
              <TitleBar className={this.props.classes.padding}>
                <TitleBarTitle title={i18n.t(`${packageNS}:menu.settings.system_settings`)} />
              </TitleBar>    
              {/* <Divider light={true}/> */}
              {/* <div className={this.props.classes.between}>
              <TitleBar>
                <TitleBarTitle component={Link} to="#" title="M2M Wallet" className={this.props.classes.link}/> 
                <TitleBarTitle component={Link} to="#" title="/" className={this.props.classes.link}/>
                <TitleBarTitle component={Link} to="#" title={i18n.t(`${packageNS}:menu.settings.system_settings`)} className={this.props.classes.link}/>
              </TitleBar>
              </div> */}
          </div>
        </Grid>
        <Grid item xs={12} md={12} lg={6} className={this.props.classes.column}>
          
                <SettingsForm
                  submitLabel={i18n.t(`${packageNS}:menu.withdraw.confirm`)}
                  onSubmit={this.onSubmit}
                />
              
          
        </Grid>
        <Grid item xs={6}>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(Settings));
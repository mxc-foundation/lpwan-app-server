import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import { withStyles } from "@material-ui/core/styles";
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import FormControl from "@material-ui/core/FormControl";
import FormLabel from "@material-ui/core/FormLabel";
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';
import TextField from "@material-ui/core/TextField";
import CardContent from "@material-ui/core/CardContent";
import Typography from "@material-ui/core/Typography";
import FormHelperText from "@material-ui/core/FormHelperText";

import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import AutocompleteSelect from "../../components/AutocompleteSelect";
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import UserStore from "../../stores/UserStore";
import OrganizationStore from "../../stores/OrganizationStore";
import SessionStore from "../../stores/SessionStore";
import theme from "../../theme";
import i18n, { packageNS } from '../../i18n';


const styles = {
  card: {
    overflow: "visible",
  },
  tabs: {
    borderBottom: "1px solid " + theme.palette.divider,
    height: "48px",
    overflow: "visible",
  },
  formLabel: {
    fontSize: 12,
  },
};



class AssignUserForm extends FormComponent {
  constructor() {
    super();

    // we need combo box
    // this.getUserOption = this.getUserOption.bind(this);
    this.getUserOptions = this.getUserOptions.bind(this);
  }

  getUserOptions(search, callbackFunc) {
    UserStore.list(search, 10, 0, resp => {
      const options = resp.result.map((u, i) => {return {label: u.username, value: u.id}});
      callbackFunc(options);
    });
  }

  render() {
    if (this.state.object === undefined) {
      return(<div></div>);
    }

    return(
      <Form
        submitLabel={i18n.t(`${packageNS}:tr000041`)}
        onSubmit={this.onSubmit}
      >
        <FormControl margin="normal" fullWidth>
          <FormLabel className={this.props.classes.formLabel} required>{i18n.t(`${packageNS}:tr000056`)}</FormLabel>
          <AutocompleteSelect
            id="userID"
            label={i18n.t(`${packageNS}:tr000137`)}
            value={this.state.object.userID || null}
            onChange={this.onChange}
            getOptions={this.getUserOptions}
          />
        </FormControl>
        <Typography variant="body1">
          {i18n.t(`${packageNS}:tr000138`)}
        </Typography>
        <FormControl fullWidth margin="normal">
          <FormControlLabel
            label={i18n.t(`${packageNS}:tr000139`)}
            control={
              <Checkbox
                id="isAdmin"
                checked={!!this.state.object.isAdmin}
                onChange={this.onChange}
                color="primary"
              />
            }
          />
          <FormHelperText>{i18n.t(`${packageNS}:tr000140`)}</FormHelperText>
        </FormControl>
        {!!!this.state.object.isAdmin && <FormControl fullWidth margin="normal">
          <FormControlLabel
            label={i18n.t(`${packageNS}:tr000141`)}
            control={
              <Checkbox
                id="isDeviceAdmin"
                checked={!!this.state.object.isDeviceAdmin}
                onChange={this.onChange}
                color="primary"
              />
            }
          />
          <FormHelperText>{i18n.t(`${packageNS}:tr000142`)}</FormHelperText>
        </FormControl>}
        {!!!this.state.object.isAdmin && <FormControl fullWidth margin="normal">
          <FormControlLabel
            label={i18n.t(`${packageNS}:tr000143`)}
            control={
              <Checkbox
                id="isGatewayAdmin"
                checked={!!this.state.object.isGatewayAdmin}
                onChange={this.onChange}
                color="primary"
              />
            }
          />
          <FormHelperText>{i18n.t(`${packageNS}:tr000144`)}</FormHelperText>
        </FormControl>}
      </Form>
    );
  };
}

AssignUserForm = withStyles(styles)(AssignUserForm);


class CreateUserForm extends FormComponent {
  render() {
    if (this.state.object === undefined) {
      return(<div></div>);
    }

    return(
      <Form
        submitLabel={i18n.t(`${packageNS}:tr000277`)}
        onSubmit={this.onSubmit}
      >
        <TextField
          id="username"
          label={i18n.t(`${packageNS}:tr000056`)}
          margin="normal"
          value={this.state.object.username || ""}
          onChange={this.onChange}
          required
          fullWidth
        />
        <TextField
          id="email"
          label={i18n.t(`${packageNS}:tr000147`)}
          margin="normal"
          value={this.state.object.email || ""}
          onChange={this.onChange}
          required
          fullWidth
        />
        <TextField
          id="note"
          label={i18n.t(`${packageNS}:tr000129`)}
          helperText={i18n.t(`${packageNS}:tr000130`)}
          margin="normal"
          value={this.state.object.note || ""}
          onChange={this.onChange}
          rows={4}
          fullWidth
          multiline
        />
        <TextField
          id="password"
          label={i18n.t(`${packageNS}:tr000004`)}
          type="password"
          margin="normal"
          value={this.state.object.password || ""}
          onChange={this.onChange}
          required
          fullWidth
        />
        <Typography variant="body1">
          {i18n.t(`${packageNS}:tr000138`)}
        </Typography>
        <FormControl fullWidth margin="normal">
          <FormControlLabel
            label={i18n.t(`${packageNS}:tr000139`)}
            control={
              <Checkbox
                id="isAdmin"
                checked={!!this.state.object.isAdmin}
                onChange={this.onChange}
                color="primary"
              />
            }
          />
          <FormHelperText>{i18n.t(`${packageNS}:tr000140`)}</FormHelperText>
        </FormControl>
        {!!!this.state.object.isAdmin && <FormControl fullWidth margin="normal">
          <FormControlLabel
            label={i18n.t(`${packageNS}:tr000141`)}
            control={
              <Checkbox
                id="isDeviceAdmin"
                checked={!!this.state.object.isDeviceAdmin}
                onChange={this.onChange}
                color="primary"
              />
            }
          />
          <FormHelperText>{i18n.t(`${packageNS}:tr000142`)}</FormHelperText>
        </FormControl>}
        {!!!this.state.object.isAdmin && <FormControl fullWidth margin="normal">
          <FormControlLabel
            label={i18n.t(`${packageNS}:tr000143`)}
            control={
              <Checkbox
                id="isGatewayAdmin"
                checked={!!this.state.object.isGatewayAdmin}
                onChange={this.onChange}
                color="primary"
              />
            }
          />
          <FormHelperText>{i18n.t(`${packageNS}:tr000144`)}</FormHelperText>
        </FormControl>}
      </Form>
    );
  }
}


class CreateOrganizationUser extends Component {
  constructor() {
    super();

    this.state = {
      tab: 0,
      assignUser: false,
    };

    this.onChangeTab = this.onChangeTab.bind(this);
    this.onAssignUser = this.onAssignUser.bind(this);
    this.onCreateUser = this.onCreateUser.bind(this);
    this.setAssignUser = this.setAssignUser.bind(this);
  }

  componentDidMount() {
    this.setAssignUser();

    SessionStore.on("change", this.setAssignUser);
  }

  comomentWillUnmount() {
    SessionStore.removeListener("change", this.setAssignUser);
  }

  setAssignUser() {
    const settings = SessionStore.getSettings();
    this.setState({
      assignUser: !settings.disableAssignExistingUsers || SessionStore.isAdmin(),
    });
  }

  onChangeTab(e, v) {
    this.setState({
      tab: v,
    });
  }

  onAssignUser(user) {
    OrganizationStore.addUser(this.props.match.params.organizationID, user, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/users`);
    });
  };

  onCreateUser(user) {
    const orgs = [
      {isAdmin: user.isAdmin, isDeviceAdmin: user.isDeviceAdmin, isGatewayAdmin: user.isGatewayAdmin, organizationID: this.props.match.params.organizationID},
    ];

    let u = user;
    u.isActive = true;

    delete u.isAdmin;
    delete u.isDeviceAdmin;
    delete u.isGatewayAdmin;

    UserStore.create(u, user.password, orgs, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/users`);
    });
  };

  render() {
    return(
      <Grid container spacing={4}>
        <TitleBar>
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000068`)} to={`/organizations/${this.props.match.params.organizationID}/users`} />
          <TitleBarTitle title="/" />
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000277`)} />
        </TitleBar>

        <Grid item xs={12}>
          <Tabs value={this.state.tab} onChange={this.onChangeTab} indicatorColor="primary" className={this.props.classes.tabs}>
            {this.state.assignUser && <Tab label={i18n.t(`${packageNS}:tr000136`)} />}
            <Tab label={i18n.t(`${packageNS}:tr000146`)} />
          </Tabs>
        </Grid>

        <Grid item xs={12}>
          <Card className={this.props.classes.card}>
            <CardContent>
              {(this.state.tab === 0 && this.state.assignUser) && <AssignUserForm onSubmit={this.onAssignUser} />}
              {(this.state.tab === 1 || !this.state.assignUser) && <CreateUserForm onSubmit={this.onCreateUser} />}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(CreateOrganizationUser));

import React from "react";

import { withStyles } from "@material-ui/core/styles";
import TextField from '@material-ui/core/TextField';
import FormControl from "@material-ui/core/FormControl";
import FormControlLabel from "@material-ui/core/FormControlLabel";
import FormLabel from "@material-ui/core/FormLabel";
import FormHelperText from "@material-ui/core/FormHelperText";
import Checkbox from "@material-ui/core/Checkbox";
import FormGroup from "@material-ui/core/FormGroup";
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import Button from "@material-ui/core/Button";
import Grid from "@material-ui/core/Grid";
import IconButton from '@material-ui/core/IconButton';
import Typography from "@material-ui/core/Typography";

import Delete from "mdi-material-ui/Delete";

import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import EUI64Field from "../../components/EUI64Field";
import AutocompleteSelect from "../../components/AutocompleteSelect";
import DeviceProfileStore from "../../stores/DeviceProfileStore";

import theme from "../../theme";


const styles = {
  formLabel: {
    fontSize: 12,
  },
  delete: {
    marginTop: 3 * theme.spacing(1),
  },
};

class DeviceKVForm extends FormComponent {
  onChange(e) {
    super.onChange(e);

    this.props.onChange(this.props.index, this.state.object);
  }

  onDelete = (e) => {
    e.preventDefault();
    this.props.onDelete(this.props.index);
  }

  render() {
    if (this.state.object === undefined) {
      return null;
    }

    return(
      <Grid container spacing={4}>
        <Grid item xs={4}>
          <TextField
            id="key"
            label={i18n.t(`${packageNS}:tr000042`)}
            margin="normal"
            value={this.state.object.key || ""}
            onChange={this.onChange}
            fullWidth
          />
        </Grid>
        <Grid item xs={7}>
          <TextField
            id="value"
            label="Value"
            margin="normal"
            value={this.state.object.value || ""}
            onChange={this.onChange}
            fullWidth
          />
        </Grid>
        <Grid item xs={1} className={this.props.classes.delete}>
          <IconButton aria-label="delete" onClick={this.onDelete}>
            <Delete />
          </IconButton>
        </Grid>
      </Grid>
    );
  }
}

DeviceKVForm = withStyles(styles)(DeviceKVForm);


class DeviceForm extends FormComponent {
  constructor() {
    super();
    this.getDeviceProfileOption = this.getDeviceProfileOption.bind(this);
    this.getDeviceProfileOptions = this.getDeviceProfileOptions.bind(this);

    this.state = {
      tab: 0,
      variables: [],
      tags: [],
    };
  }

  componentDidMount() {
    super.componentDidMount();

    this.setKVArrays(this.props.object || {});
  }

  componentDidUpdate(prevProps) {
    super.componentDidUpdate(prevProps);

    if (prevProps.object !== this.props.object) {
      this.setKVArrays(this.props.object || {});
    }
  }

  setKVArrays = (props) => {
    let variables = [];
    let tags = [];

    if (props.variables !== undefined) {
      for (let key in props.variables) {
        variables.push({key: key, value: props.variables[key]});
      }
    }

    if (props.tags !== undefined) {
      for (let key in props.tags) {
        tags.push({key: key, value: props.tags[key]});
      }
    }

    this.setState({
      variables: variables,
      tags: tags,
    });
  }

  getDeviceProfileOption(id, callbackFunc) {
    DeviceProfileStore.get(id, resp => {
      callbackFunc({label: resp.deviceProfile.name, value: resp.deviceProfile.id});
    });
  }

  getDeviceProfileOptions(search, callbackFunc) {
    DeviceProfileStore.list(0, this.props.match.params.applicationID, 999, 0, resp => {
      const options = resp.result.map((dp, i) => {return {label: dp.name, value: dp.id}});
      callbackFunc(options);
    });
  }

  onTabChange = (e, v) => {
    this.setState({
      tab: v,
    });
  }

  addKV = (name) => {
    return (e) => {
      e.preventDefault();

      let kvs = this.state[name];
      kvs.push({});

      let obj = {};
      obj[name] = kvs;

      this.setState(obj);
    };
  }

  onChangeKV = (name) => {
    return (index, obj) => {
      let kvs = this.state[name];
      let object = this.state.object;

      kvs[index] = obj;

      object[name] = {};
      kvs.forEach((obj, i) => {
        object[name][obj.key] = obj.value;
      });

      let ss = {
        object: object,
      };
      ss[name] = kvs;

      this.setState(ss);
    };
  }

  onDeleteKV = (name) => {
    return (index) => {
      let kvs = this.state[name];
      let object = this.state.object;

      kvs.splice(index, 1);

      object[name] = {};
      kvs.forEach((obj, i) => {
        object[name][obj.key] = obj.value;
      });

      let ss = {
        object: object,
      };
      ss[name] = kvs;

      this.setState(ss);
    };
  }

  render() {
    if (this.state.object === undefined) {
      return null;
    }

    const variables = this.state.variables.map((obj, i) => <DeviceKVForm key={i} index={i} object={obj} onChange={this.onChangeKV("variables")} onDelete={this.onDeleteKV("variables")} />);
    const tags = this.state.tags.map((obj, i) => <DeviceKVForm key={i} index={i} object={obj} onChange={this.onChangeKV("tags")} onDelete={this.onDeleteKV("tags")} />);

    return(
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
        disabled={this.props.disabled}
      >
        <Tabs value={this.state.tab} onChange={this.onTabChange} indicatorColor="primary">
          <Tab label={i18n.t(`${packageNS}:tr000167`)} />
          <Tab label={i18n.t(`${packageNS}:tr000305`)} />
          <Tab label={i18n.t(`${packageNS}:tr000308`)} />
        </Tabs>

        {this.state.tab === 0 && <div>
          <TextField
            id="name"
            label={i18n.t(`${packageNS}:tr000300`)}
            helperText={i18n.t(`${packageNS}:tr000062`)}
            margin="normal"
            value={this.state.object.name || ""}
            onChange={this.onChange}
            inputProps={{
              pattern: "[\\w-]+",
            }}
            fullWidth
            required
          />
          <TextField
            id="description"
            label={i18n.t(`${packageNS}:tr000301`)}
            margin="normal"
            value={this.state.object.description || ""}
            onChange={this.onChange}
            fullWidth
            required
          />
          {!this.props.update && <EUI64Field
            margin="normal"
            id="devEUI"
            label={i18n.t(`${packageNS}:tr000371`)}
            onChange={this.onChange}
            value={this.state.object.devEUI || ""}
            fullWidth
            required
            random
          />}
          <FormControl fullWidth margin="normal">
            <FormLabel className={this.props.classes.formLabel} required>{i18n.t(`${packageNS}:tr000281`)}</FormLabel>
            <AutocompleteSelect
              id="deviceProfileID"
              label={i18n.t(`${packageNS}:tr000281`)}
              value={this.state.object.deviceProfileID}
              onChange={this.onChange}
              getOption={this.getDeviceProfileOption}
              getOptions={this.getDeviceProfileOptions}
            />
          </FormControl>
          <FormControl margin="normal">
            <FormGroup>
              <FormControlLabel
                label={i18n.t(`${packageNS}:tr000303`)}
                control={
                  <Checkbox
                    id="skipFCntCheck"
                    checked={!!this.state.object.skipFCntCheck}
                    onChange={this.onChange}
                    color="primary"
                  />
                }
              />
            </FormGroup>
            <FormHelperText>
              {i18n.t(`${packageNS}:tr000304`)}
            </FormHelperText>
          </FormControl>
        </div>}

        {this.state.tab === 1 && <div>
          <FormControl fullWidth margin="normal">
            <Typography variant="body1">
              {i18n.t(`${packageNS}:tr000306`)}
            </Typography>
            {variables}
          </FormControl>
          <Button variant="outlined" onClick={this.addKV("variables")}>{i18n.t(`${packageNS}:tr000307`)}</Button>
        </div>}

        {this.state.tab === 2 && <div>
          <FormControl fullWidth margin="normal">
            <Typography variant="body1">
              {i18n.t(`${packageNS}:tr000309`)}
            </Typography>
            {tags}
          </FormControl>
          <Button variant="outlined" onClick={this.addKV("tags")}>Add tag</Button>
        </div>}
      </Form>
    );
  }
}

export default withStyles(styles)(DeviceForm);

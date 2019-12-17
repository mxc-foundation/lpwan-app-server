import React, { Component } from "react";

import { Row, Col, Button as RButton } from 'reactstrap';
import { Formik, Form as BForm, Field } from 'formik';
import * as Yup from 'yup';

import { withStyles } from "@material-ui/core/styles";
import TextField from '@material-ui/core/TextField';
import FormControl from "@material-ui/core/FormControl";
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormHelperText from "@material-ui/core/FormHelperText";
import FormGroup from "@material-ui/core/FormGroup";
import FormLabel from "@material-ui/core/FormLabel";
import Checkbox from '@material-ui/core/Checkbox';
import Button from "@material-ui/core/Button";

import { Map, Marker } from 'react-leaflet';

import { ReactstrapInput, ReactstrapCheckbox, AsyncAutoComplete } from '../../components/FormInputs';
import i18n, { packageNS } from '../../i18n';
import Form from "../../components/Form";

import AutocompleteSelect from "../../components/AutocompleteSelect";
import NetworkServerStore from "../../stores/NetworkServerStore";
import GatewayProfileStore from "../../stores/GatewayProfileStore";
import LocationStore from "../../stores/LocationStore";
import MapTileLayer from "../../components/MapTileLayer";
import EUI64Field from "../../components/EUI64Field";
import AESKeyField from "../../components/AESKeyField";
import theme from "../../theme";


const boardStyles = {
  formLabel: {
    color: theme.palette.primary.main,
  },
  a: {
    color: theme.palette.primary.main,
  },
};

class GatewayBoardForm extends Component {
  constructor() {
    super();

    this.onChange = this.onChange.bind(this);
    this.onDelete = this.onDelete.bind(this);
  }

  onChange(e) {
    let board = this.props.board;
    const field = e.target.id;

    board[field] = e.target.value;
    this.props.onChange(board);
  }

  onDelete(e) {
    e.preventDefault();
    this.props.onDelete();
  }

  render() {
    return (
      <FormControl fullWidth margin="normal">
        <FormLabel className={this.props.classes.formLabel}>{i18n.t(`${packageNS}:tr000400`)} #{this.props.i} (<a href="#delete" onClick={this.onDelete} className={this.props.classes.a}>{i18n.t(`${packageNS}:tr000401`)}</a>)</FormLabel>
        <EUI64Field
          id="fpgaID"
          label={i18n.t(`${packageNS}:tr000236`)}
          margin="normal"
          value={this.props.board.fpgaID || ""}
          onChange={this.onChange}
          helperText={i18n.t(`${packageNS}:tr000237`)}
          fullWidth
        />
        <AESKeyField
          id="fineTimestampKey"
          label={i18n.t(`${packageNS}:tr000238`)}
          margin="normal"
          value={this.props.board.fineTimestampKey || ""}
          onChange={this.onChange}
          helperText={i18n.t(`${packageNS}:tr000239`)}
          fullWidth
        />
      </FormControl>
    );
  }
}

GatewayBoardForm = withStyles(boardStyles)(GatewayBoardForm);


const styles = {
  mapLabel: {
    marginBottom: theme.spacing(1),
  },
  link: {
    color: theme.palette.primary.main,
  },
  formLabel: {
    fontSize: 12,
  },
};

class GatewayForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      mapZoom: 15,
    };

    this.getNetworkServerOption = this.getNetworkServerOption.bind(this);
    this.getNetworkServerOptions = this.getNetworkServerOptions.bind(this);
    this.getGatewayProfileOption = this.getGatewayProfileOption.bind(this);
    this.getGatewayProfileOptions = this.getGatewayProfileOptions.bind(this);
    this.setCurrentPosition = this.setCurrentPosition.bind(this);
    this.updatePosition = this.updatePosition.bind(this);
    this.updateZoom = this.updateZoom.bind(this);
    this.addGatewayBoard = this.addGatewayBoard.bind(this);
  }

  componentDidMount() {
    this.setState({
      object: this.props.object || {},
    });

    if (!this.props.update) {
      this.setCurrentPosition();
    }
  }

  // onChange(e) {
  //   if (e.target.id === "networkServerID" && e.target.value !== this.state.object.networkServerID) {
  //     let object = this.state.object;
  //     object.gatewayProfileID = null;
  //     this.setState({
  //       object: object,
  //     });
  //   }

  //   super.onChange(e);
  // }

  setCurrentPosition(e) {
    if (e !== undefined) {
      e.preventDefault();
    }

    LocationStore.getLocation(position => {
      let object = this.state.object;
      object.location = {
        latitude: position.coords.latitude,
        longitude: position.coords.longitude,
      }
      this.setState({
        object: object,
      });
    });
  }

  updatePosition() {
    const position = this.refs.marker.leafletElement.getLatLng();
    let object = this.state.object;
    object.location = {
      latitude: position.lat,
      longitude: position.lng,
    }
    this.setState({
      object: object,
    });
  }

  updateZoom(e) {
    this.setState({
      mapZoom: e.target.getZoom(),
    });
  }

  getNetworkServerOption(id, callbackFunc) {
    NetworkServerStore.get(id, resp => {
      callbackFunc({ label: resp.networkServer.name, value: resp.networkServer.id });
    });
  }

  getNetworkServerOptions(search, callbackFunc) {
    NetworkServerStore.list(this.props.match.params.organizationID, 999, 0, resp => {
      const options = resp.result.map((ns, i) => { return { label: ns.name, value: ns.id } });
      callbackFunc(options);
    });
  }

  getGatewayProfileOption(id, callbackFunc) {
    GatewayProfileStore.get(id, resp => {
      callbackFunc({ label: resp.gatewayProfile.name, value: resp.gatewayProfile.id });
    });
  }

  getGatewayProfileOptions(search, callbackFunc) {
    if (this.state.object === undefined || this.state.object.networkServerID === undefined) {
      callbackFunc([]);
      return;
    }

    GatewayProfileStore.list(this.state.object.networkServerID, 999, 0, resp => {
      const options = resp.result.map((gp, i) => { return { label: gp.name, value: gp.id } });
      callbackFunc(options);
    });
  }

  addGatewayBoard() {
    let object = this.state.object;
    if (object.boards === undefined) {
      object.boards = [{}];
    } else {
      object.boards.push({});
    }

    this.setState({
      object: object,
    });
  }

  deleteGatewayBoard(i) {
    let object = this.state.object;
    object.boards.splice(i, 1);
    this.setState({
      object: object,
    });
  }

  updateGatewayBoard(i, board) {
    let object = this.state.object;
    object.boards[i] = board;
    this.setState({
      object: object,
    });
  }

  render() {
    if (this.state.object === undefined) {
      return (<div></div>);
    }

    const style = {
      height: 400,
    };

    let position = [];
    if (this.state.object.location.latitude !== undefined && this.state.object.location.longitude !== undefined) {
      position = [this.state.object.location.latitude, this.state.object.location.longitude];
    } else {
      position = [0, 0];
    }

    let boards = [];
    if (this.state.object.boards !== undefined) {
      boards = this.state.object.boards.map((b, i) => <GatewayBoardForm key={i} i={i} board={b} onDelete={() => this.deleteGatewayBoard(i)} onChange={board => this.updateGatewayBoard(i, board)} />);
    }


    let fieldsSchema = {
      name: Yup.string()
        .required(i18n.t(`${packageNS}:tr000425`)),
      description: Yup.string()
        .required(i18n.t(`${packageNS}:tr000425`)),
      gatewayProfileID: Yup.string(),
      discoveryEnabled: Yup.bool(),
      'location.altitude': Yup.number().required(i18n.t(`${packageNS}:tr000425`)),
    }

    if (!this.props.update) {
      fieldsSchema['id'] = Yup.string().required(i18n.t(`${packageNS}:tr000425`));
      fieldsSchema['networkServerID'] =  Yup.string();
    }
    const formSchema = Yup.object().shape(fieldsSchema);

    return (<React.Fragment>
      <Row>
        <Col>
          <Formik
            initialValues={this.state.object}
            validationSchema={formSchema}
            onSubmit={values => {
              // same shape as initial values
              console.log(values);
            }}>
            {({
              handleSubmit,
              setFieldValue
            }) => (
                <BForm noValidate>
                  <Field
                    type="text"
                    label={i18n.t(`${packageNS}:tr000218`)}
                    name="name"
                    id="name"
                    helpText={i18n.t(`${packageNS}:tr000062`)}
                    component={ReactstrapInput}
                  />

                  <Field
                    type="textarea"
                    label={i18n.t(`${packageNS}:tr000219`)}
                    name="description"
                    id="description"
                    component={ReactstrapInput}
                  />

                  {!this.props.update && <EUI64Field
                    id="id"
                    label={i18n.t(`${packageNS}:tr000074`)}
                    name="id"
                    value={this.state.object.id || ""}
                    required
                    random
                  />}

                  {!this.props.update && <Field
                    type="text"
                    label={i18n.t(`${packageNS}:tr000047`)}
                    name="networkServerID"
                    id="networkServerID"
                    value={this.state.object.networkServerID || ""}
                    getOption={this.getNetworkServerOption}
                    getOptions={this.getNetworkServerOptions}
                    setFieldValue={setFieldValue}
                    helpText={i18n.t(`${packageNS}:tr000223`)}
                    inputProps={{
                      clearable: true,
                      cache: false,
                    }}
                    component={AsyncAutoComplete}
                  />}
                  

                  <Field
                    type="text"
                    label={i18n.t(`${packageNS}:tr000224`)}
                    name="gatewayProfileID"
                    id="gatewayProfileID"
                    value={this.state.object.gatewayProfileID || ""}
                    getOption={this.getGatewayProfileOption}
                    getOptions={this.getGatewayProfileOptions}
                    setFieldValue={setFieldValue}
                    inputProps={{
                      clearable: true,
                      cache: false,
                    }}
                    component={AsyncAutoComplete}
                  />

                  <Field
                    type="checkbox"
                    label={i18n.t(`${packageNS}:tr000228`)}
                    name="discoveryEnabled"
                    id="discoveryEnabled"
                    component={ReactstrapCheckbox}
                    helpText={i18n.t(`${packageNS}:tr000229`)}
                  />

                  <Field
                    type="number"
                    label={i18n.t(`${packageNS}:tr000230`)}
                    name="location.altitude"
                    id="location.altitude"
                    value={this.state.object.location.altitude || 0}
                    component={ReactstrapInput}
                    helpText={i18n.t(`${packageNS}:tr000231`)}
                  />

                  <RButton type="submit" color="primary">Submit</RButton>
                  {/* <Button onClick={this.addGatewayBoard}>{i18n.t(`${packageNS}:tr000234`)}</Button>   */}
                </BForm>
              )}
          </Formik>
        </Col>
      </Row>

      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
      >
        {/* <TextField
          id="name"
          label={i18n.t(`${packageNS}:tr000218`)}
          margin="normal"
          value={this.state.object.name || ""}
          onChange={this.onChange}
          inputProps={{
            pattern: "[\\w-]+",
          }}
          helperText={i18n.t(`${packageNS}:tr000062`)}
          required
          fullWidth
        />
        <TextField
          id="description"
          label={i18n.t(`${packageNS}:tr000219`)}
          margin="normal"
          value={this.state.object.description || ""}
          onChange={this.onChange}
          rows={4}
          multiline
          required
          fullWidth
        /> */}
        {/* {!this.props.update && <EUI64Field
          id="id"
          label={i18n.t(`${packageNS}:tr000074`)}
          margin="normal"
          value={this.state.object.id || ""}
          onChange={this.onChange}
          required
          fullWidth
          random
        />} */}
        {/* {!this.props.update && <FormControl fullWidth margin="normal">
          <FormLabel className={this.props.classes.formLabel} required>{i18n.t(`${packageNS}:tr000047`)}</FormLabel>
          <AutocompleteSelect
            id="networkServerID"
            label={i18n.t(`${packageNS}:tr000115`)}
            value={this.state.object.networkServerID || ""}
            onChange={(e) => { console.log(e) }}
            getOption={this.getNetworkServerOption}
            getOptions={this.getNetworkServerOptions}
          />
          <FormHelperText>
            {i18n.t(`${packageNS}:tr000223`)}
          </FormHelperText>
        </FormControl>}
        <FormControl fullWidth margin="normal">
          <FormLabel className={this.props.classes.formLabel}>{i18n.t(`${packageNS}:tr000224`)}</FormLabel>
          <AutocompleteSelect
            id="gatewayProfileID"
            label={i18n.t(`${packageNS}:tr000225`)}
            value={this.state.object.gatewayProfileID || ""}
            triggerReload={this.state.object.networkServerID || ""}
            onChange={(e) => { console.log(e) }}
            getOption={this.getGatewayProfileOption}
            getOptions={this.getGatewayProfileOptions}
            inputProps={{
              clearable: true,
              cache: false,
            }}
          />
          <FormHelperText>
            {i18n.t(`${packageNS}:tr000227`)}
          </FormHelperText>
        </FormControl> */}
        {/* <FormGroup>
          <FormControlLabel
            label={i18n.t(`${packageNS}:tr000228`)}
            control={
              <Checkbox
                id="discoveryEnabled"
                checked={!!this.state.object.discoveryEnabled}
                onChange={this.onChange}
                color="primary"
              />
            }
          />
          <FormHelperText>
            {i18n.t(`${packageNS}:tr000229`)}
          </FormHelperText>
        </FormGroup> */}
        {/* <TextField
          id="location.altitude"
          label={i18n.t(`${packageNS}:tr000230`)}
          margin="normal"
          type="number"
          value={this.state.object.location.altitude || 0}
          onChange={this.onChange}
          helperText={i18n.t(`${packageNS}:tr000231`)}
          required
          fullWidth
        /> */}
        <FormControl fullWidth margin="normal">
          <FormLabel className={this.props.classes.mapLabel}>{i18n.t(`${packageNS}:tr000232`)} (<a onClick={this.setCurrentPosition} href="#getlocation" className={this.props.classes.link}>{i18n.t(`${packageNS}:tr000328`)}</a>)</FormLabel>
          <Map
            center={position}
            zoom={this.state.mapZoom}
            style={style}
            animate={true}
            scrollWheelZoom={false}
            onZoomend={this.updateZoom}
          >
            <MapTileLayer />
            <Marker position={position} draggable={true} onDragend={this.updatePosition} ref="marker" />
          </Map>
          <FormHelperText>
            {i18n.t(`${packageNS}:tr000233`)}
          </FormHelperText>
        </FormControl>
        {/* {boards} */}
      </Form>
    </React.Fragment>
    );
  }
}

export default withStyles(styles)(GatewayForm);

import React, { Component } from "react";
import { withRouter, Route, Switch, Link } from "react-router-dom";

import BrandSelectModal from "./BrandSelectModal";
import moment from "moment";
import { Bar } from "react-chartjs-2";
import { Map, Marker, Popup } from 'react-leaflet';
import MarkerClusterGroup from "react-leaflet-markercluster";
import L from "leaflet";
import "leaflet.awesome-markers";
import { Row, Col, Card, CardBody } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import { MAX_DATA_LIMIT } from '../../util/pagination';
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import OrgBreadCumb from '../../components/OrgBreadcrumb';
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import GatewayAdmin from "../../components/GatewayAdmin";
import GatewayStore from "../../stores/GatewayStore";
import MapTileLayer from "../../components/MapTileLayer";


const GatewayActivityColumn = (cell, row, index, extraData) => {
  const stats = extraData['stats'];

  let rowStats = stats && stats[row.id] ? stats[row.id] : null;

  let dataTotal = 0;

  if (rowStats) {
    for (const row of rowStats) {
      dataTotal += parseFloat(row.rxPacketsReceivedOK + row.txPacketsEmitted);
    }
  }

  return (
    rowStats ? parseFloat(dataTotal) : <React.Fragment>0</React.Fragment>
  );
}

const GatewayColumn = (cell, row, index, extraData) => {
  const organizationId = extraData['organizationId'];
  return <Link to={`/organizations/${organizationId}/gateways/${row.id}`}>{row.name}</Link>;
}

const LastSeenAtColumn = (cell, row, index, extraData) => {
  return (row.lastSeenAt) ? row.lastSeenAt : '--:--';
}
const getColumns = (organizationId, stats) => (
  [{
    dataField: 'name',
    text: i18n.t(`${packageNS}:tr000042`),
    sort: false,
    formatter: GatewayColumn,
    formatExtraData: { organizationId: organizationId }
  }, {
    dataField: 'id',
    text: i18n.t(`${packageNS}:tr000074`),
    sort: false,
  }, {
    dataField: 'lastSeenAt',
    text: i18n.t(`${packageNS}:tr000075`),
    formatter: GatewayActivityColumn,
    formatExtraData: { stats },
    sort: false,
  }, {
    dataField: 'lastSeenAt',
    text: i18n.t(`${packageNS}:tr000283`),
    formatter: LastSeenAtColumn,
    sort: false,
  }, {
    dataField: 'downlink_price',
    text: i18n.t(`${packageNS}:tr000421`),
    sort: false,
  }, {
    dataField: 'mode',
    text: i18n.t(`${packageNS}:tr000422`),
    sort: false,
  }]
);

class ListGatewaysTable extends Component {
  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);
    this.getGateWayStats = this.getGateWayStats.bind(this);
    this.state = {
      data: [],
      stats: {},
      totalSize: 0
    }
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, searchText, sortField, sortOrder, searchField }) => {
    const offset = (page - 1) * sizePerPage;

    let searchQuery = null;
    if (type === 'search' && searchText && searchText.length) {
      searchQuery = searchText;
    }
    // TODO - how can I pass search query to server?
    this.getPage(sizePerPage, offset);
  }

  /**
   * Fetches data from server
   */
  getPage = (limit, offset) => {
    limit = MAX_DATA_LIMIT;
    this.setState({ loading: true });
    GatewayStore.list("", this.props.organizationID, limit, offset, (res) => {
      const object = this.state;
      object.totalSize = Number(res.totalCount);
      object.data = res.result;
      object.loading = false;
      this.setState({ object });
    });
  }

  /**
   * Gets the stats from server
   */
  getGateWayStats = (gatewayId) => {
    const start = moment().subtract(29, "days").toISOString();
    const end = moment().toISOString();
    GatewayStore.getStats(gatewayId, start, end, resp => {
      let stats = { ...this.state.stats };
      stats[gatewayId] = resp.result;
      this.setState({ stats: stats });
    });
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  componentDidUpdate(prevProps, prevState) {
    if (prevState !== this.state && prevState.data !== this.state.data) {
      for (const item of this.state.data) {
        this.getGateWayStats(item.id);
      }
    }
  }

  render() {
    return (
      <div className="position-relative">
        {this.state.loading && <Loader />}
        <AdvancedTable data={this.state.data} columns={getColumns(this.props.organizationID, this.state.stats)}
          keyField="id" onTableChange={this.handleTableChange} searchEnabled={false} totalSize={this.state.totalSize} rowsPerPage={10}></AdvancedTable>
      </div>
    );
  }
}


class ListGatewaysMap extends Component {
  constructor() {
    super();

    this.state = {
      items: null,
    };
  }

  componentDidMount() {
    GatewayStore.list("", this.props.organizationID, 9999, 0, resp => {
      this.setState({
        items: resp.result,
      });
    });
  }

  render() {
    if (this.state.items === null) {
      return null;
    }

    const style = {
      height: 800,
    };


    let bounds = [];
    let markers = [];

    const greenMarker = L.AwesomeMarkers.icon({
      icon: "wifi",
      prefix: "fa",
      markerColor: "green",
    });

    const grayMarker = L.AwesomeMarkers.icon({
      icon: "wifi",
      prefix: "fa",
      markerColor: "gray",
    });

    const redMarker = L.AwesomeMarkers.icon({
      icon: "wifi",
      prefix: "fa",
      markerColor: "red",
    });

    for (const item of this.state.items) {
      const position = [item.location.latitude, item.location.longitude];

      bounds.push(position);

      let marker = greenMarker;
      let lastSeen = "";

      if (item.lastSeenAt === undefined || item.lastSeenAt === null) {
        marker = grayMarker;
        lastSeen = "Never seen online";
      } else {
        const ts = moment(item.lastSeenAt);
        if (ts.isBefore(moment().subtract(5, 'minutes'))) {
          marker = redMarker;
        }

        lastSeen = ts.fromNow();
      }

      markers.push(
        <Marker position={position} key={`gw-${item.id}`} icon={marker}>
          <Popup>
            <Link to={`/organizations/${this.props.organizationID}/gateways/${item.id}`}>{item.name}</Link><br />
            {item.id}<br /><br />
            {lastSeen}
          </Popup>
        </Marker>
      );
    }

    return (<React.Fragment>
      {bounds && <Map bounds={bounds} maxZoom={19} style={style} animate={true} scrollWheelZoom={false}>
        <MapTileLayer />
        <MarkerClusterGroup>
          {markers}
        </MarkerClusterGroup>
      </Map>}
    </React.Fragment>
    );
  }
}


class ListGateways extends Component {
  constructor() {
    super();

    this.switchToList = this.switchToList.bind(this);
    this.locationToTab = this.locationToTab.bind(this);
    this.state = {
      viewMode: 'list',
      nsDialog: false,
      setModal: false
    };
  }

  componentDidMount() {
    this.locationToTab();
  }

  locationToTab = () => {
    if (window.location.href.endsWith("/map")) {
      this.setState({ viewMode: 'map' });
    }
  }

  showModal = () => {
    const object = this.state;
    object.nsDialog = true;
    this.setState({ object });
  }

  toggle = () => {
    const object = this.state;
    object.nsDialog = !object.nsDialog;
    this.setState({ object });
  }

  handleCloseModal = () => {
    const object = this.state;
    object.nsDialog = null;
    this.setState({
      object
    })
  }
  /**
   * Switch to list
   */
  switchToList() {
    this.setState({ viewMode: 'list' });
  }

  render() {
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (<React.Fragment>

      <TitleBar
        buttons={<GatewayAdmin organizationID={this.props.match.params.organizationID}>

          {/* <BrandSelectModal
          buttonLabel={i18n.t(`${packageNS}:tr000277`)}
          callback={this.handleLink} /> */}
          <TitleBarButton
            key={1}
            label={i18n.t(`${packageNS}:tr000277`)}
            icon={<i className="mdi mdi-plus mr-1 align-middle"></i>}
            onClick={this.toggle}
            //to={`/organizations/${this.props.match.params.organizationID}/gateways/brand`}
            //to={`/organizations/${currentOrgID}/gateways/create`}
            to={`/organizations/${currentOrgID}/gateways/input-serial`}
          />
        </GatewayAdmin>}
      >
        <OrgBreadCumb organizationID={currentOrgID} items={[
          { label: i18n.t(`${packageNS}:tr000063`), active: true }]}></OrgBreadCumb>
      </TitleBar>

      <Row>
        <Col>
          <Card className="card-box shadow-sm">
            {this.state.viewMode === 'map' &&
              <Link to={`/organizations/${this.props.match.params.organizationID}/gateways`} className="btn btn-primary mb-3" onClick={this.switchToList}>Show List</Link>}

            <Switch>
              <Route exact path={this.props.match.path} render={props => <ListGatewaysTable {...props} organizationID={this.props.match.params.organizationID} />} />
              <Route exact path={`${this.props.match.path}/map`} render={props => <ListGatewaysMap {...props} organizationID={this.props.match.params.organizationID} />} />
            </Switch>
          </Card>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default withRouter(ListGateways);

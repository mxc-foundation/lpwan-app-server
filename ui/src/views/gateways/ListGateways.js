import React, { Component } from "react";
import { Route, Switch, Link } from "react-router-dom";

import { withStyles } from "@material-ui/core/styles";

import moment from "moment";
import { Bar } from "react-chartjs-2";
import { Map, Marker, Popup } from 'react-leaflet';
import MarkerClusterGroup from "react-leaflet-markercluster";
import L from "leaflet";
import "leaflet.awesome-markers";
import { Row, Col, Card, CardBody } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import GatewayAdmin from "../../components/GatewayAdmin";
import GatewayStore from "../../stores/GatewayStore";
import MapTileLayer from "../../components/MapTileLayer";



// class GatewayRow extends Component {
//   constructor() {
//     super();

//     this.state = {};
//   }

//   componentWillMount() {
//     const start = moment().subtract(29, "days").toISOString();
//     const end = moment().toISOString();

//     GatewayStore.getStats(this.props.gateway.id, start, end, resp => {
//       let stats = {
//         labels: [],
//         datasets: [
//           {
//             data: [],
//             fillColor: "rgba(33, 150, 243, 1)",
//           },
//         ],
//       };

//       for (const row of resp.result) {
//         stats.labels.push(row.timestamp);
//         stats.datasets[0].data.push(row.rxPacketsReceivedOK + row.txPacketsEmitted);
//       }

//       this.setState({
//         stats: stats,
//       });
//     });
//   }

//   render() {
//     const options = {
//       elements: {
//         rectangle: {
//           backgroundColor: 'rgb(0, 255, 217)',
//         }
//       },
//       scales: {
//         xAxes: [{display: false}],
//         yAxes: [{display: false}],
//       },
//       tooltips: {
//         enabled: false,
//       },
//       legend: {
//         display: false,
//       },
//       responsive: false,
//       animation: {
//         duration: 0,
//       },
//     };

//     return(
//       <TableRow>
//           <TableCellLink to={`/organizations/${this.props.gateway.organizationID}/gateways/${this.props.gateway.id}`}>{this.props.gateway.name}</TableCellLink>
//           <TableCell>{this.props.gateway.id}</TableCell>
//           <TableCell>
//             {this.state.stats && <Bar
//               width={380}
//               height={23}
//               data={this.state.stats}
//               options={options}
//             />}
//           </TableCell>
//       </TableRow>
//     );
//   }
// }

const GatewayColumn = (cell, row, index, extraData) => {
  const organizationId = extraData['organizationId'];
  return <Link to={`/organizations/${organizationId}/gateways/${row.id}`}>{row.name}</Link>;
}

const getColumns = (organizationId) => (
  [{
    dataField: 'test_gateway_profile',
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
    sort: false,
  }, {
    dataField: 'status',
    text: i18n.t(`${packageNS}:tr000282`),
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
    this.state = {
      data: []
    }
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, searchText, sortField, sortOrder, searchField }) => {
    const offset = (page - 1) * sizePerPage + 1;

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
    this.setState({loading: true});
    GatewayStore.list("", this.props.organizationID, limit, offset, (res) => {
      this.setState({ data: res.result, loading: false });
    });
  }

  componentDidMount() {
    this.getPage(10);
  }

  render() {
    return (
      <div className="position-relative">
        {this.state.loading && <Loader />}
        <AdvancedTable data={this.state.data} columns={getColumns(this.props.organizationID)}
          keyField="id" onTableChange={this.handleTableChange} searchEnabled={true} rowsPerPage={10}></AdvancedTable>
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
      viewMode: 'list'
    };
  }

  componentDidMount() {
    this.locationToTab();
  }

  locationToTab = () => {
    if (window.location.href.endsWith("/map")) {
      this.setState({viewMode: 'map'});
    }
  }

  /**
   * Switch to list
   */
  switchToList() {
    this.setState({viewMode: 'list'});
  }

  render() {
    return (<React.Fragment>

      <TitleBar
        buttons={<GatewayAdmin organizationID={this.props.match.params.organizationID}>
          <TitleBarButton
            key={1}
            label={i18n.t(`${packageNS}:tr000277`)}
            icon={<i className="mdi mdi-plus mr-1 align-middle"></i>}
            to={`/organizations/${this.props.match.params.organizationID}/gateways/create`}
          />
        </GatewayAdmin>}
      >
        <TitleBarTitle title={i18n.t(`${packageNS}:tr000063`)} />
      </TitleBar>

      <Row>
        <Col>
          <Card>
            <CardBody>
              {this.state.viewMode === 'map' && 
                <Link to={`/organizations/${this.props.match.params.organizationID}/gateways`} className="btn btn-primary mb-3" onClick={this.switchToList}>Show List</Link>}

              <Switch>
                <Route exact path={this.props.match.path} render={props => <ListGatewaysTable {...props} organizationID={this.props.match.params.organizationID} />} />
                <Route exact path={`${this.props.match.path}/map`} render={props => <ListGatewaysMap {...props} organizationID={this.props.match.params.organizationID} />} />
              </Switch>

            </CardBody>
          </Card>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default ListGateways;

import React, { Component } from "react";
import { Link } from "react-router-dom";

import moment from "moment";
import { Bar } from "react-chartjs-2";

import i18n, { packageNS } from "../../i18n";
import { MAX_DATA_LIMIT } from '../../util/pagination';
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import DeviceStore from "../../stores/DeviceStore";

const OPTIONS_CONFIG = {
  elements: {
    rectangle: {
      backgroundColor: 'rgb(0, 255, 217)',
    }
  },
  scales: {
    xAxes: [{ display: false }],
    yAxes: [{ display: false }],
  },
  tooltips: {
    enabled: false,
  },
  legend: {
    display: false,
  },
  responsive: false,
  animation: {
    duration: 0,
  },
};

const CHART_DATA_CONFIG = {
  labels: [],
  datasets: [
    {
      data: [],
      fillColor: "rgba(33, 150, 243, 1)",
    },
  ],
};

/*
FIXME - this link is shown at http://localhost:3000/#/organizations/2/applications/3,
so shouldn't it have url /organizations/2/applications/3/devices 
FIXME - make the link dynamic depending on whether or not the device is associated with
an application or not.
*/
const DeviceNameColumn = (cell, row, index, extraData) => {
  const currentOrgID = extraData['currentOrgID'];
  const currentApplicationID = extraData['currentApplicationID'];
  return currentApplicationID 
    ? <Link to={`/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${row.devEUI}`}>{row.name}</Link>
    : <Link to={`/organizations/${currentOrgID}/devices/${row.devEUI}`}>{row.name}</Link>;
}

const DeviceLastSeenAtColumn = (cell, row, index, extraData) => {
  let lastseen;
  if (row.lastSeenAt !== undefined && row.lastSeenAt !== null) {
    lastseen = moment(row.lastSeenAt).fromNow();
  }
  return (
    lastseen ? lastseen : <React.Fragment>-</React.Fragment>
  );
}

const DeviceMarginColumn = (cell, row, index, extraData) => {
  let lastseen, margin;
  if (row.lastSeenAt !== undefined && row.lastSeenAt !== null) {
    lastseen = moment(row.lastSeenAt).fromNow();
  }
  if (row.deviceStatusMargin !== undefined && row.deviceStatusMargin !== 256) {
    margin = `${row.deviceStatusMargin} dB`;
  }
  return (
    margin ? (
      <React.Fragment>{margin}</React.Fragment>
    ) : (
      <React.Fragment>-</React.Fragment>
    )
  );
}

const DeviceExternalPowerSourceColumn = (cell, row, index, extraData) => {
  let isExternalPowerSource = row.deviceStatusExternalPowerSource ? true : false;
  return (
    isExternalPowerSource ? (
      <React.Fragment><i className="mdi mdi-power-plug mdi-24px align-middle"></i></React.Fragment>
    ) : (
      <React.Fragment>-</React.Fragment>
    )
  );
}

const DeviceBatteryLevelColumn = (cell, row, index, extraData) => {
  const stats = extraData['stats'];
  const options = OPTIONS_CONFIG;
  let rowStats = stats && stats[row.devEUI] ? stats[row.devEUI]: null;
  let chartData = CHART_DATA_CONFIG;
  let lastseen, batteryLevel;

  if (rowStats) {
    for (const row of rowStats) {
      if (!row.deviceStatusExternalPowerSource && !row.deviceStatusBatteryLevelUnavailable) {
        batteryLevel = `${row.deviceStatusBatteryLevel}%`
      }
      if (row.lastSeenAt !== undefined && row.lastSeenAt !== null) {
        lastseen = moment(row.lastSeenAt).fromNow();
        chartData.labels.push(lastseen);
        chartData.datasets[0].data.push(batteryLevel);
      }
    }
  }
  return (
    rowStats ? (
      <Bar
        width={380}
        height={23}
        data={chartData}
        options={options}
      />
    ) : <React.Fragment>-</React.Fragment>
  );
}

const getColumns = (currentOrgID, currentApplicationID, stats) => (
  [
    {
      dataField: 'name',
      text: i18n.t(`${packageNS}:tr000300`),
      sort: false,
      formatter: DeviceNameColumn,
      formatExtraData: {
        currentOrgID: currentOrgID,
        currentApplicationID: currentApplicationID
      }
    }, {
      dataField: 'devEUI',
      text: i18n.t(`${packageNS}:tr000371`),
      sort: false,
    }, {
      dataField: 'lastSeenAt',
      text: i18n.t(`${packageNS}:tr000242`),
      formatter: DeviceLastSeenAtColumn,
      sort: false,
    }, {
      dataField: 'deviceStatusMargin',
      text: i18n.t(`${packageNS}:tr000382`),
      sort: false,
      formatter: DeviceMarginColumn
    }, {
      dataField: 'deviceStatusExternalPowerSource',
      text: i18n.t(`${packageNS}:tr000500`),
      sort: false,
      formatter: DeviceExternalPowerSourceColumn
    }, {
      dataField: 'deviceStatusBatteryLevel',
      text: i18n.t(`${packageNS}:tr000383`),

      // TODO - find out whether we'll should implement

      // `listLocations` and `getStats` in DeviceStore.js
      // and the backend (similar to how ListGateways.js and GatewayStore.js)
      // formatter: DeviceBatteryLevelColumn,
      // formatExtraData: { stats },
      sort: false,
    }
  ]
);

class ListDevicesTable extends Component {
  constructor(props) {
    super(props);

    this.state = {
      data: [],
      stats: {}
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
    const currentApplicationID = this.props.applicationID;

    this.setState({ loading: true });
    const filters = {
      limit: limit,
      offset: offset,
      organizationID: this.props.organizationID,
      search: "",
    };

    DeviceStore.list(filters, (res) => {
      // Since this page is only shown when an Application ID is in the URL parameters
      // we need to filter so we only list devices that are part of the current Application
      const getOnlyDevicesFromCurrentApplication = (devices) =>
        devices.filter(device => device.applicationID === currentApplicationID)
      this.setState({ data: getOnlyDevicesFromCurrentApplication(res.result), loading: false });
    });
  }

  /**
   * Gets the stats from server
   */

  // TODO - see comment about wheter to implement `listLocation` and `getStats` in DeviceStore.js

  // getDeviceStats = (devEUI) => {
  //   const start = moment().subtract(29, "days").toISOString();
  //   const end = moment().toISOString();
  //   DeviceStore.getStats(devEUI, start, end, resp => {
  //     let stats = { ...this.state.stats };
  //     stats[devEUI] = resp.result;
  //     this.setState({ stats: stats });
  //   });
  // }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  componentDidUpdate(prevProps, prevState) {
    // TODO - see comment about wheter to implement `listLocation` and `getStats` in DeviceStore.js

    // if (prevState !== this.state && prevState.data !== this.state.data) {
    //   for (const item of this.state.data) {
    //     this.getDeviceStats(item.devEUI);
    //   }
    // }
  }

  render() {
    const currentOrgID = this.props.organizationID;
    const currentApplicationID = this.props.applicationID;

    return (
      <div className="position-relative">
        {this.state.loading && <Loader />}
        <AdvancedTable
          data={this.state.data}
          columns={getColumns(currentOrgID, currentApplicationID, this.state.stats)}
          keyField="devEUI"
          onTableChange={this.handleTableChange}
          rowsPerPage={10}
          searchEnabled={false}
        />
      </div>
    );
  }
}

export default ListDevicesTable;

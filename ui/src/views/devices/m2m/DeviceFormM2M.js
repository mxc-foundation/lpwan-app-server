import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { Link, withRouter } from "react-router-dom";
import AdvancedTable from "../../../components/AdvancedTable";
import Loader from "../../../components/Loader";
import NativeSelects from "../../../components/m2m/NativeSelectsM2M";
import SwitchLabels from "../../../components/m2m/SwitchM2M";
import i18n, { packageNS } from '../../../i18n';
import DeviceStore from "../../../stores/DeviceStore.js";
//import Wallet from "mdi-material-ui/OpenInNew";
//import Typography from '@material-ui/core/Typography';
import { DV_INACTIVE, DV_MODE_OPTION } from "../../../util/Data";
import { MAX_DATA_LIMIT } from '../../../util/pagination';




const styles = {
  flex: {
    display: 'flex',
    alignItems: 'center',
  },
  flex2: {
    left: 'calc(100%/3)',
  },
  maxW: {
    maxWidth: 120
  }
};

const DeviceM2MNameColumn = (cell, row, index, extraData) => {
  const currentOrgID = extraData['currentOrgID'];
  const currentApplicationID = extraData['currentApplicationID'] || row.applicationId; // Here we get the applicationId from the fetched data (elsewhere we get it from URL params or props)

  // return <Link to={`/organizations/${currentOrgID}/applications/${currentApplicationID ? currentApplicationID : 0}/devices/${row.devEui}`}>{row.name}</Link>;
  return currentApplicationID
    ? <Link to={`/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${row.devEui}`}>{row.name}</Link>
    : <Link to={`/organizations/${currentOrgID}/devices/${row.devEui}`}>{row.name}</Link>;
}

const DeviceM2MStatusColumn = (cell, row, index, extraData) => {
  return row.lastSeenAt.substring(0, 19);
}

class DeviceFormM2M extends Component {
  constructor(props) {
    super(props);

    this._isMounted = false;

    this.state = {
      data: [],
      loading: false,
      totalSize: 0
    }
  }

  componentDidMount() {
    this._isMounted = true;

    DeviceStore.on('update', () => {
      // re-render the table.
      this.forceUpdate();
    });

    this.getPage(MAX_DATA_LIMIT);
  }

  componentWillUnmount() {
    this._isMounted = false;
  }

  DeviceM2MAvailableColumn = (cell, row, index, extraData) => {
    let on = (row.mode !== DV_INACTIVE) ? true : false;

    return (
      <span className={this.props.classes.flex}>
        <SwitchLabels
          on={on}
          dvId={row.id}
          onSwitchChange={this.onSwitchChange}
        />
      </span>
    );
  }

  DeviceM2MModeColumn = (cell, row, index, extraData) => {
    /*
    let dValue = null;
    const options = DV_MODE_OPTION;
    
    switch(row.mode) {
        case options[1].value:
        dValue = options[1];
        break;
        case options[2].value:
        dValue = options[2];
        break;
        default:
        dValue = options[0];
        break;
    }
    */
    const options = DV_MODE_OPTION;
    const dValue = options[1];

    let on = (row.mode !== DV_INACTIVE) ? true : false;
    const isDisabled = on ? false : true;
    return (
      <span className={this.props.classes.maxW}>
        <NativeSelects
          options={options}
          isDisabled={isDisabled}
          defaultValue={dValue}
          haveGateway={this.props.haveGateway}
          mode={row.mode}
          dvId={row.id}
          onSelectChange={this.onSelectChange}
        />
      </span>
    );
  }

  getColumns = (currentOrgID, currentApplicationID) => (
    [
      {
        dataField: 'name',
        text: i18n.t(`${packageNS}:tr000300`),
        sort: false,
        formatter: DeviceM2MNameColumn,
        formatExtraData: {
          currentOrgID: currentOrgID,
          currentApplicationID: currentApplicationID
        }
      }, {
        dataField: 'lastSeenAt',
        text: "Status / Last Seen At", // i18n.t(`${packageNS}:menu.devices.status`)
        sort: false,
        formatter: DeviceM2MStatusColumn,
      }, {
        dataField: 'available',
        text: i18n.t(`${packageNS}:menu.devices.available`),
        sort: false,
        formatter: this.DeviceM2MAvailableColumn,
      },
/*      {
        dataField: 'mode',
        text: i18n.t(`${packageNS}:menu.devices.mode`),
        sort: false,
        formatter: this.DeviceM2MModeColumn,
        // Configure the CSS of the column so when switch to mobile viewport the selection box is not malformed
        headerStyle: (colum, colIndex) => {
          return { minWidth: '150px', textAlign: 'left' };
        }
      }*/
    ]
  )

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, searchText, sortField, sortOrder, searchField }) => {
    const offset = (page - 1) * sizePerPage ;

    /* let searchQuery = null;
    if (type === 'search' && searchText && searchText.length) {
      searchQuery = searchText;
    } */

    this.getPage(sizePerPage, offset);
  }

  /**
   * Fetches data from server
   */
  getPage = (limit, offset) => {
    limit = MAX_DATA_LIMIT;
    if (this._isMounted) {
      this.setState({ loading: true });
    }

    DeviceStore.getDeviceList(this.props.match.params.organizationID, offset, limit, (res) => {
      if (this._isMounted) {
        const object = this.state;
        object.totalSize = res.count;
        object.data = res.devProfile;
        object.loading = false;
        this.setState({ object });
      }
    }, (error) => {
      this.setState({ loading: false });
    });
  }

  onSelectChange = (e) => {
    const device = {
      dvId: e.dvId,
      dvMode: e.target.value
    }

    this.props.onSelectChange(device);
  }

  onSwitchChange = (dvId, available, e) => {
    const device = {
      dvId,
      available
    }

    this.props.onSwitchChange(device, e);
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;
    const { loading: loadingState } = this.state;
    const isLoading = loadingState;

    return (
      <React.Fragment>
        <div className="position-relative">
          {isLoading && <Loader />}
          <AdvancedTable
            data={this.state.data}
            columns={this.getColumns(currentOrgID, currentApplicationID)}
            keyField="id"
            onTableChange={this.handleTableChange}
            rowsPerPage={10}
            totalSize={this.state.totalSize}
            searchEnabled={false}
          />
        </div>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(DeviceFormM2M));

import React, { Component } from "react";
import { Link } from "react-router-dom";

import { Button } from "reactstrap";
import Grid from "@material-ui/core/Grid";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
// import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';

import moment from "moment";

import i18n, { packageNS } from '../../i18n';
import { MAX_DATA_LIMIT } from '../../util/pagination';
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import TableCellLink from "../../components/TableCellLink";
import DataTable from "../../components/DataTable";

import FUOTADeploymentStore from "../../stores/FUOTADeploymentStore";

const FUOTADeploymentDeviceNameColumn = (cell, row, index, extraData) => {
  const currentOrgID = extraData['currentOrgID'];
  const applicationId = extraData['applicationId'];
  return <Link to={
    applicationId
      ? `/organizations/${currentOrgID}/applications/${applicationId}/devices/${row.devEUI}`
      : `/organizations/${currentOrgID}/devices/${row.devEUI}`
  }>{row.deviceName}</Link>;
}

const FUOTADeploymentCreatedAtColumn = (cell, row, index) => {
  return moment(row.createdAt).format('lll');
}

const FUOTADeploymentUpdatedAtColumn = (cell, row, index) => {
  return moment(row.updatedAt).format('lll');
}

class FUOTADeploymentDevices extends Component {
  constructor() {
    super();

    this.state = {
      detailDialog: false,
      data: [],
      loading: true,
      totalSize: 0
    };
  }

  FUOTADeploymentDevicesShowDetailsColumn = (cell, row, index) => {
    return <Button onClick={() => this.showDetails(row.devEUI)}>Show Details</Button>;
  }

  getColumns = (currentOrgID, applicationId) => (
    [
      {
        dataField: 'devEUI',
        text: i18n.t(`${packageNS}:tr000371`),
        sort: false
      },
      {
        dataField: 'deviceName',
        text: i18n.t(`${packageNS}:tr000300`),
        sort: false,
        formatter: FUOTADeploymentDeviceNameColumn,
        formatExtraData: {
          currentOrgID: currentOrgID,
          applicationId: applicationId
        }
      }, {
        dataField: 'state',
        text: i18n.t(`${packageNS}:tr000350`),
        sort: false,
      }, {
        dataField: 'errorMessage',
        text: "Error Message",
        sort: false,
      }, {
        dataField: 'createdAt',
        text: i18n.t(`${packageNS}:tr000321`),
        sort: false,
        formatter: FUOTADeploymentCreatedAtColumn,
      }, {
        dataField: 'updatedAt',
        text: i18n.t(`${packageNS}:tr000322`),
        sort: false,
        formatter: FUOTADeploymentUpdatedAtColumn,
      }, {
        dataField: 'showDetails',
        text: "",
        sort: false,
        formatter: this.FUOTADeploymentDevicesShowDetailsColumn
      }
    ]
  );

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, searchText, sortField, sortOrder, searchField }) => {
    const offset = (page - 1) * sizePerPage ;

    let searchQuery = null;
    if (type === 'search' && searchText && searchText.length) {
      searchQuery = searchText;
    }

    this.getPage(sizePerPage, offset);
  }

  /**
   * Fetches data from server
   */
  getPage = (limit, offset) => {
    this.setState({ loading: true });

    FUOTADeploymentStore.listDeploymentDevices({
      fuota_deployment_id: this.props.match.params.fuotaDeploymentID,
      limit: limit,
      offset: offset,
    }, (res) => {
      const object = this.state;
      object.totalSize = res.totalCount;
      object.data = res.result;
      object.loading = false;
      this.setState({ object });
    });
  }

  showDetails = (devEUI) => {
    FUOTADeploymentStore.getDeploymentDevice(this.props.match.params.fuotaDeploymentID, devEUI, resp => {
      this.setState({
        data: new Array(Object.assign({}, this.state.data[0], resp.deploymentDevice)),
        detailDialog: true,
      });
    });
  }

  onCloseDialog = () => {
    this.setState({
      detailDialog: false,
    });
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  render() {
    const { data } = this.state;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;

    let fddUpdatedAt = "";
    if (data[0] !== undefined) {
      fddUpdatedAt = moment(data[0].updatedAt).format('lll');
    }

    return (
      <React.Fragment>
        {data[0] && <Dialog
          open={this.state.detailDialog}
          onClose={this.onCloseDialog}
        >
          <DialogTitle>{i18n.t(`${packageNS}:tr000339`)}</DialogTitle>
          <DialogContent>
            <Table>
              <TableBody>
                <TableRow>
                  <TableCell>{i18n.t(`${packageNS}:tr000340`)}</TableCell>
                  <TableCell>{fddUpdatedAt}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell>{i18n.t(`${packageNS}:tr000324`)}</TableCell>
                  <TableCell>{data[0].state}</TableCell>
                </TableRow>
                {data[0].state === "ERROR" && <TableRow>
                  <TableCell>{i18n.t(`${packageNS}:tr000341`)}</TableCell>
                  <TableCell>{data[0].errorMessage}</TableCell>
                </TableRow>}
              </TableBody>
            </Table>
          </DialogContent>
          <DialogActions>
            <Button color="primary" onClick={this.onCloseDialog}>{i18n.t(`${packageNS}:tr000166`)}</Button>
          </DialogActions>
        </Dialog>}

        <div className="position-relative">
          {this.state.loading && <Loader />}
          <AdvancedTable
            data={data}
            columns={this.getColumns(currentOrgID, currentApplicationID)}
            keyField="devEUI"
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

export default FUOTADeploymentDevices;

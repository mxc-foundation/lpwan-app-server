import Grid from "@material-ui/core/Grid";
import { withStyles } from "@material-ui/core/styles";
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
import moment from "moment";
import React, { Component } from "react";
import { Link } from "react-router-dom";
import { Button, Modal, ModalBody, ModalFooter, ModalHeader, NavLink } from 'reactstrap';
import AdvancedTable from "../../components/AdvancedTable";
import DeviceAdmin from "../../components/DeviceAdmin";
import Loader from "../../components/Loader";
import i18n, { packageNS } from '../../i18n';
import FUOTADeploymentStore from "../../stores/FUOTADeploymentStore";
import theme from "../../theme";
import { MAX_DATA_LIMIT } from '../../util/pagination';






const styles = {
  buttons: {
    textAlign: "right",
  },
  button: {
    marginLeft: 2 * theme.spacing(1),
  }
};

const FUOTADeploymentNameColumn = (cell, row, index, extraData) => {
  const currentOrgID = extraData['currentOrgID'];
  const currentApplicationID = extraData['currentApplicationID'];
  return <Link to={
    currentApplicationID
    ? `/organizations/${currentOrgID}/applications/${currentApplicationID}/fuota-deployments/${row.id}`
    : `/organizations/${currentOrgID}/fuota-deployments/${row.id}`
  }>{row.name}</Link>;
}

const FUOTADeploymentCreatedAtColumn = (cell, row, index) => {
  return moment(row.createdAt).format('lll');
}

const FUOTADeploymentUpdatedAtColumn = (cell, row, index) => {
  return moment(row.updatedAt).format('lll');
}

class ListFUOTADeploymentsForDevice extends Component {
  constructor() {
    super();

    this.state = {
      data: [],
      loading: true,
      detailDialog: false,
      totalSize: 0
    };
  }

  fUOTADeploymentShowDetailsColumn = (cell, row, index) => {
    return <Button size="small" onClick={() => this.showDetails(row.id)}>Show</Button>;
  }

  showDetails = (fuotaDeploymentID) => {
    const devEUI = this.props.devEUI || this.props.match.params.devEUI;
    FUOTADeploymentStore.getDeploymentDevice(fuotaDeploymentID, devEUI, resp => {
      this.setState({
        data: new Array(Object.assign({}, this.state.data[0], resp.deploymentDevice)),
        fuotaDeploymentID: fuotaDeploymentID,
        detailDialog: true,
      });
    });
  }

  toggleDialog = () => {
    this.setState({
      detailDialog: !this.state.detailDialog,
    });
  }

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

  getColumns = (currentOrgID, currentApplicationID) => (
    [
      {
        dataField: 'name',
        text: i18n.t(`${packageNS}:tr000320`),
        sort: false,
        formatter: FUOTADeploymentNameColumn,
        formatExtraData: {
          currentOrgID: currentOrgID,
          currentApplicationID: currentApplicationID
        }
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
        dataField: 'state',
        text: i18n.t(`${packageNS}:tr000323`),
        sort: false,
      }, {
        dataField: 'showDetails',
        text: "Show Details", // tr000324
        sort: false,
        formatter: this.fUOTADeploymentShowDetailsColumn,
      }
    ]
  );

  /**
   * Fetches data from server
   */
  getPage = (limit, offset) => {
    const devEUI = this.props.devEUI || this.props.match.params.devEUI;
    this.setState({ loading: true });

    FUOTADeploymentStore.list({
      devEUI: devEUI,
      limit: limit,
      offset: offset,
    }, (res) => {
      const object = this.state;
      object.totalSize = Number(res.totalCount);
      object.data = res.result;
      object.loading = false;
      this.setState({ object });
    });
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  render() {
    const { data } = this.state;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;
    const devEUI = this.props.devEUI || this.props.match.params.devEUI;

    let fddUpdatedAt = "";
    if (data[0] !== undefined) {
      fddUpdatedAt = moment(data[0].updatedAt).format('lll');
    }

    const closeBtn = <button className="close" onClick={this.toggleDialog}>&times;</button>;

    return(
      <React.Fragment>
        {data[0] &&
          <Modal
            isOpen={this.state.detailDialog}
            toggle={this.toggleDialog}
            aria-labelledby="help-dialog-title"
            aria-describedby="help-dialog-description"
          >
            <ModalHeader
              toggle={this.toggleDialog}
              close={closeBtn}
              id="help-dialog-title"
            >
              {i18n.t(`${packageNS}:tr000339`)}
            </ModalHeader>
            <ModalBody id="help-dialog-description">
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
            </ModalBody>
            <ModalFooter>
              <Button color="primary" onClick={this.toggleDialog}>{i18n.t(`${packageNS}:tr000166`)}</Button>{' '}
            </ModalFooter>
          </Modal>
        }

        <DeviceAdmin organizationID={currentOrgID}>
          <Grid item xs={12} className={this.props.classes.buttons}>
          <Button variant="outlined" style={{ marginBottom: "1em" }}>
            <NavLink
              style={{ color: "#fff", padding: "0" }}
              tag={Link}
              to={
                currentApplicationID
                ? `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${devEUI}/fuota-deployments/create`
                : `/organizations/${currentOrgID}/devices/${devEUI}/fuota-deployments/create`
              }
            >
              <i className="mdi mdi-cloud-upload"></i>&nbsp;
              {/* Create */} {i18n.t(`${packageNS}:tr000342`)} {/* Job */}
              {/* TODO - either Create or Add depending on whether any firmware update jobs already exist or not for the device */}
            </NavLink>
          </Button>
          </Grid>
        </DeviceAdmin>
          <div className="position-relative">
            {this.state.loading && <Loader />}
            <AdvancedTable
              data={data}
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

export default withStyles(styles)(ListFUOTADeploymentsForDevice);

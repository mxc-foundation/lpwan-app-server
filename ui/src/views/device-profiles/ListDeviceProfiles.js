import React, { Component } from "react";
import { Link } from "react-router-dom";

import { Card } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import Grid from '@material-ui/core/Grid';

import i18n, { packageNS } from '../../i18n';
import { MAX_DATA_LIMIT } from '../../util/pagination';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";
import DeviceAdmin from "../../components/DeviceAdmin";
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import DeviceProfileStore from "../../stores/DeviceProfileStore";
import OrganizationDevices from "../devices/OrganizationDevices";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

const DeviceProfileNameColumn = (cell, row, index, extraData) => {
  const currentOrgID = extraData['currentOrgID'];

  return <Link to={`/organizations/${currentOrgID}/device-profiles/${row.id}`}>{row.name}</Link>;
}

const getColumns = (currentOrgID) => (
  [
    {
      dataField: 'name',
      text: i18n.t(`${packageNS}:tr000042`),
      sort: false,
      formatter: DeviceProfileNameColumn,
      formatExtraData: {
        currentOrgID: currentOrgID
      }
    }
  ]
);

class ListDeviceProfiles extends Component {
  constructor(props) {
    super(props);

    this.state = {
      data: [],
      loading: true,
      totalSize: 0
    }
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, searchText, sortField, sortOrder, searchField }) => {
    const offset = (page - 1) * sizePerPage ;

    let searchQuery = null;
    if (type === 'search' && searchText && searchText.length) {
      searchQuery = searchText;
    }
    // TODO - how can I pass search query to server?
    this.getPage(sizePerPage, offset);
  }

  getPage = (limit, offset) => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    this.setState({ loading: true });

    // FIXME - should we be associating the Device Profile optionally with an Application ID?
    DeviceProfileStore.list(currentOrgID, 0, limit, offset, (res) => {
      const object = this.state;
      object.totalSize = res.totalCount;
      object.data = res.result;
      object.loading = false;
      this.setState({object});
    });
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  render() {
    const { classes } = this.props;
    // TODO - refactor this into a method or store in state on page load (apply to all components where this rushed approach used)
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return(
      <Grid container spacing={4}>
        <OrganizationDevices
          mainTabIndex={2}
          organizationID={currentOrgID}
        >
          <TitleBar
            buttons={
              <DeviceAdmin organizationID={currentOrgID}>
                <TitleBarButton
                  key={1}
                  label={i18n.t(`${packageNS}:tr000277`)}
                  icon={<i className="mdi mdi-plus mr-1 align-middle"></i>}
                  color="primary"
                  to={`/organizations/${currentOrgID}/device-profiles/create`}
                />,
              </DeviceAdmin>
            }
          >
            <TitleBarTitle title={i18n.t(`${packageNS}:tr000070`)} />
          </TitleBar>
          <Grid item xs={12}>
            <Card className="card-box shadow-sm">
              {this.state.loading && <Loader />}
              <AdvancedTable
                data={this.state.data}
                columns={getColumns(currentOrgID)}
                keyField="id"
                onTableChange={this.handleTableChange}
                rowsPerPage={10}
                totalSize={this.state.totalSize}
                searchEnabled={false}
              />
            </Card>
          </Grid>
        </OrganizationDevices>
      </Grid>
    );
  }
}

export default withStyles(styles)(ListDeviceProfiles);

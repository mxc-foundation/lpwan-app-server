import React, { Component } from "react";
import { Link } from "react-router-dom";

import { Breadcrumb, BreadcrumbItem, Card } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import Grid from '@material-ui/core/Grid';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import DeviceAdmin from "../../components/DeviceAdmin";
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import DeviceProfileStore from "../../stores/DeviceProfileStore";
import OrganizationDevices from "../devices/OrganizationDevices";

const styles = theme => ({
  [theme.breakpoints.down('sm')]: {
    breadcrumb: {
      fontSize: "1.1rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  [theme.breakpoints.up('sm')]: {
    breadcrumb: {
      fontSize: "1.25rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  breadcrumbItemLink: {
    color: "#71b6f9 !important"
  }
});

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
      loading: true
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

  getPage = (limit, offset) => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    this.setState({ loading: true });

    // FIXME - should we be associating the Device Profile optionally with an Application ID?
    DeviceProfileStore.list(currentOrgID, 0, limit, offset, (res) => {
      this.setState({
        data: res.result,
        loading: false
      });
    });
  }

  componentDidMount() {
    this.getPage(10);
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
            <Breadcrumb className={classes.breadcrumb}>
              <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000070`)}</BreadcrumbItem>
            </Breadcrumb>
          </TitleBar>
          <Grid item xs={12}>
            <Card className="card-box shadow-sm" style={{ minWidth: "25rem" }}>
              {this.state.loading && <Loader />}
              <AdvancedTable
                data={this.state.data}
                columns={getColumns(currentOrgID)}
                keyField="id"
                onTableChange={this.handleTableChange}
                rowsPerPage={10}
                searchEnabled={true}
              />
            </Card>
          </Grid>
        </OrganizationDevices>
      </Grid>
    );
  }
}

export default withStyles(styles)(ListDeviceProfiles);

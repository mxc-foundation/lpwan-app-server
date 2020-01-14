import React, { Component } from "react";
import { Link } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Button, Card, CardBody,
  CardSubtitle, CardTitle, Col, Container, Row, Spinner } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import { MAX_DATA_LIMIT } from '../../util/pagination';
import Loader from "../../components/Loader";
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";

import AdvancedTable from "../../components/AdvancedTable";

import NetworkServerStore from "../../stores/NetworkServerStore";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

const NetworkServerColumn = (cell, row, index, extraData) => {
  return <Link to={`/network-servers/${row.id}`}>{row.name}</Link>;
}

const NetworkServerAddressColumn = (cell, row, index, extraData) => {
  return <div>{row.server}</div>;
}

const columns = [{
  dataField: 'networkServerName',
  text: i18n.t(`${packageNS}:tr000042`),
  sort: false,
  formatter: NetworkServerColumn
}, {
  dataField: 'networkServerAddress',
  text: i18n.t(`${packageNS}:tr000043`),
  sort: false,
  formatter: NetworkServerAddressColumn
}];

class ListNetworkServers extends Component {
  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);
    this.state = {
      data: [],
      loading: true,
      totalSize: 0
    }
  }

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
    limit = MAX_DATA_LIMIT;
    const defaultOrgId = 0;
    this.setState({ loading: true });
    NetworkServerStore.list(defaultOrgId, limit, offset, (res) => {
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

    return(
      <React.Fragment>
        <TitleBar
          buttons={[
            <TitleBarButton
              aria-label={i18n.t(`${packageNS}:tr000277`)}
              icon={<i className="mdi mdi-plus mr-1 align-middle"></i>}
              label={i18n.t(`${packageNS}:tr000277`)}
              key={'b-1'}
              to={`/network-servers/create`}
              className="btn btn-primary">{i18n.t(`${packageNS}:tr000277`)}
            </TitleBarButton>,
          ]}
        >
          <Breadcrumb className={classes.breadcrumb}>
            <BreadcrumbItem className={classes.breadcrumbItem}>Control Panel</BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000040`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          <Col>
            <Card className="card-box shadow-sm">
              {/* <CardBody className="position-relative"> */}
                {this.state.loading && <Loader />}
                <AdvancedTable
                  data={this.state.data}
                  columns={columns}
                  keyField="id"
                  onTableChange={this.handleTableChange}
                  rowsPerPage={10}
                  totalSize={this.state.totalSize}
                  searchEnabled={false}
                />
              {/* </CardBody> */}
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(ListNetworkServers);

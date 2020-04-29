import React, { Component } from "react";
import { Link } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Card, Col, Row } from 'reactstrap';
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import i18n, { packageNS } from '../../i18n';
import GatewayProfileStore from "../../stores/GatewayProfileStore";
import SessionStore from '../../stores/SessionStore';
import { MAX_DATA_LIMIT } from '../../util/pagination';




const GatewayColumn = (cell, row, index, extraData) => {
  return SessionStore.isAdmin() ? <Link to={`/gateway-profiles/${row.id}`}>{row.name}</Link> : <span>{row.name}</span>;
}

const NetworkColumn = (cell, row, index, extraData) => {
  return SessionStore.isAdmin() ? <Link to={`/network-servers/${row.networkServerID}`}>{row.networkServerName}</Link> : <span>{row.networkServerName}</span>;
}

const getColumns = () => (
  [{
    dataField: 'name',
    text: i18n.t(`${packageNS}:tr000042`),
    sort: false,
    formatter: GatewayColumn
  }, {
    dataField: 'networkServerName',
    text: i18n.t(`${packageNS}:tr000047`),
    sort: false,
    formatter: NetworkColumn,
  }]
);

class ListGatewayProfiles extends Component {

  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);
    this.state = {
      data: [],
      totalSize: 0,
      loading: false,
    }
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, sortField, sortOrder }) => {
    const offset = (page - 1) * sizePerPage;
    this.getPage(sizePerPage, offset);
  }

  /**
   * Fetches data from server
   */
  getPage = (limit, offset) => {
    limit = MAX_DATA_LIMIT;
    this.setState({ loading: true });
    GatewayProfileStore.list(0, limit, offset, (res) => {
      const object = this.state;
      object.totalSize = Number(res.totalCount);
      object.data = res.result;
      object.loading = false;
      this.setState({ object });
    }, error => { this.setState({ loading: false }) });
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  render() {

    const buttons = SessionStore.isAdmin() ? [<TitleBarButton
      aria-label={i18n.t(`${packageNS}:tr000277`)}
      icon={<i className="mdi mdi-plus mr-1 align-middle"></i>}
      label={i18n.t(`${packageNS}:tr000277`)}
      key={'b-1'}
      to={`/gateway-profiles/create`}
      className="btn btn-primary">{i18n.t(`${packageNS}:tr000277`)}
    </TitleBarButton>,] : [];

    return (<React.Fragment>
      <TitleBar
        buttons={buttons}
      >
        <Breadcrumb>
          <BreadcrumbItem>{i18n.t(`${packageNS}:menu.control_panel`)}</BreadcrumbItem>
          <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000046`)}</BreadcrumbItem>
        </Breadcrumb>
      </TitleBar>

      <Row>
        <Col>
          <Card className="card-box shadow-sm" >
            {this.state.loading && <Loader />}
            <AdvancedTable data={this.state.data} columns={getColumns()} keyField="id" totalSize={this.state.totalSize} onTableChange={this.handleTableChange}></AdvancedTable>
          </Card>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default ListGatewayProfiles;

import React, { Component } from "react";
import { Link, withRouter } from "react-router-dom";
import { Card, Col, Row } from 'reactstrap';
import Admin from '../../components/Admin';
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import OrgBreadCumb from '../../components/OrgBreadcrumb';
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import i18n, { packageNS } from '../../i18n';
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import { MAX_DATA_LIMIT } from '../../util/pagination';





class ListServiceProfiles extends Component {
  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);
    this.serviceProfileColumn = this.serviceProfileColumn.bind(this)
    this.state = {
      data: [],
      loading: false,
      columns: [{
        dataField: 'name',
        text: i18n.t(`${packageNS}:tr000042`),
        sort: false,
        formatter: this.serviceProfileColumn,
      }],
      totalSize: 0
    }
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, sortField, sortOrder }) => {
    const offset = (page - 1) * sizePerPage;
    this.getPage(this.props.match.params.organizationID, sizePerPage, offset);
  }

  /**
   * Fetches data from server
   */
  getPage = async (organizationID, limit, offset) => {
    this.setState({ loading: true });

    const res = await ServiceProfileStore.list(organizationID, limit=10, offset=0);
    const object = this.state;
    object.totalSize = Number(res.totalCount);
    object.data = res.result;
    object.loading = false;
    this.setState({ object });
  }

  serviceProfileColumn = (cell, row, index, extraData) => {
    return <Link to={`/organizations/${this.props.match.params.organizationID}/service-profiles/${row.id}`}>{row.name}</Link>;
  }

  componentDidMount() {
    this.getPage(this.props.match.params.organizationID, MAX_DATA_LIMIT);
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (
      <React.Fragment>
        <Admin>
          <TitleBar
            buttons={
              <TitleBarButton
                key={1}
                label={i18n.t(`${packageNS}:tr000277`)}
                icon={<i className="mdi mdi-plus mr-1 align-middle"></i>}
                onClick={this.toggle}
                to={`/organizations/${currentOrgID}/service-profiles/create`}
              />
            }
          >
            <OrgBreadCumb organizationID={currentOrgID} items={[
              { label: i18n.t(`${packageNS}:tr000069`), active: false }]}></OrgBreadCumb>
          </TitleBar>
        </Admin>

        <Row>
          <Col>
            <Card className="card-box shadow-sm">
            {this.state.loading && <Loader />}
              <AdvancedTable data={this.state.data} columns={this.state.columns} keyField="id" totalSize={this.state.totalSize} onTableChange={this.handleTableChange}></AdvancedTable>
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withRouter(ListServiceProfiles);


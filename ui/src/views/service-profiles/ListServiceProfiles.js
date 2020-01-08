import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import { Breadcrumb, BreadcrumbItem, Row, Col, Card, CardBody } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import AdvancedTable from "../../components/AdvancedTable";

import ServiceProfileStore from "../../stores/ServiceProfileStore";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class ListServiceProfiles extends Component {
  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);
    this.serviceProfileColumn = this.serviceProfileColumn.bind(this)
    this.state = {
      data: [],
      columns: [{
        dataField: 'name',
        text: i18n.t(`${packageNS}:tr000042`),
        sort: false,
        formatter: this.serviceProfileColumn,
      }]
    }
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, sortField, sortOrder }) => {
    const offset = (page - 1) * sizePerPage + 1;
    this.getPage(this.props.match.params.organizationID, sizePerPage, offset);
  }

  /**
   * Fetches data from server
   */
  getPage = (organizationID, limit, offset) => {
    ServiceProfileStore.list(organizationID, limit, offset, (res) => {
      this.setState({ data: res.result });
    });
  }

  serviceProfileColumn = (cell, row, index, extraData) => {
    return <Link to={`/organizations/${this.props.match.params.organizationID}/service-profiles/${row.id}`}>{row.name}</Link>;
  }

  componentDidMount() {
    this.getPage(this.props.match.params.organizationID, 10);
  }

  render() {
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (
      <React.Fragment>
        <TitleBar>
          <Breadcrumb className={classes.breadcrumb}>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations`}
              >
                  Organizations
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations/${currentOrgID}`}
              >
                {currentOrgID}
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000069`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>

        <Row>
          <Col>
            <Card>
              <CardBody>
                <AdvancedTable data={this.state.data} columns={this.state.columns} keyField="id" onTableChange={this.handleTableChange}></AdvancedTable>
              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(ListServiceProfiles));


import React, { Component } from "react";
import { Link } from "react-router-dom";

import Check from "mdi-material-ui/Check";
import Close from "mdi-material-ui/Close";
import AdvancedTable from "../../components/AdvancedTable";
import { Button, Breadcrumb, BreadcrumbItem, Row, Col, Card, CardBody } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import Loader from "../../components/Loader";
import UserStore from "../../stores/UserStore";

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

const GatewayColumn = (cell, row, index, extraData) => {
  return <Link to={`/users/${row.id}`}>{row.username}</Link>;
}

const ActiveColumn = (cell, row, index, extraData) => {
  if (row.isActive) {
    return <Check/>;
  } else {
    return <Close />;
  }
}

const AdminColumn = (cell, row, index, extraData) => {
  if (row.isAdmin) {
    return <Check />;
  } else {
    return <Close />;
  }
}

const getColumns = () => (
  [{
    dataField: 'username',
    text: i18n.t(`${packageNS}:tr000056`),
    sort: false,
    formatter: GatewayColumn
  }, {
    dataField: 'isActive',
    text: i18n.t(`${packageNS}:tr000057`),
    sort: false,
    formatter: ActiveColumn,
  }, {
    dataField: 'isAdmin',
    text: i18n.t(`${packageNS}:tr000058`),
    sort: false,
    formatter: AdminColumn,
  }]
);

class ListUsers extends Component {
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
    this.setState({ loading: true });

    UserStore.list("", limit, offset, (res) => {
      const object = this.state;
      object.totalSize = res.totalCount;
      object.data = res.result;
      object.loading = false;
      this.setState({object});
    });
  }

  componentDidMount() {
    this.getPage(10);
  }

  render() {
    const { classes } = this.props;
    
    return (
      <React.Fragment>
        <TitleBar
          buttons={[
            <TitleBarButton
              key={1}
              label={i18n.t(`${packageNS}:tr000277`)}
              icon={<i className="mdi mdi-account-multiple-plus mr-1 align-middle"></i>}
              to={`/users/create`}
            />
          ]}
        >
          <Breadcrumb className={classes.breadcrumb}>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/users`}
              >
                {i18n.t(`${packageNS}:tr000036`)}
              </Link>
            </BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          <Col>
            <Card className="card-box shadow-sm" style={{ minWidth: "25rem" }}>
              {this.state.loading && <Loader />}
              <AdvancedTable
                data={this.state.data}
                columns={getColumns()}
                keyField="id"
                onTableChange={this.handleTableChange}
                searchEnabled={false}
                rowsPerPage={10}
                totalSize={this.state.totalSize}
              />
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(ListUsers);

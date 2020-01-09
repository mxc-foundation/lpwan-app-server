import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import { Collapse, Button, ButtonGroup } from 'reactstrap';

import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";

import Refresh from "mdi-material-ui/Refresh";
import Delete from "mdi-material-ui/Delete";

import i18n, { packageNS } from '../../../../i18n';
import { MAX_DATA_LIMIT } from '../../../../util/pagination';
import AdvancedTable from "../../../../components/AdvancedTable";
import Loader from "../../../../components/Loader";
import DeviceQueueStore from "../../../../stores/DeviceQueueStore";

const CURRENT_CARD = "queueCard";

const DeviceQueueFConfirmedColumn = (cell, row, index, extraData) => {
  return (
    <React.Fragment>{row.confirmed
      ? <i className="mdi mdi-checkbox-marked-circle" style={{ color: "green", fontSize: "1.5em" }}></i>
      : <i className="mdi mdi-close-circle" style={{ color: "red", fontSize: "1.5em" }}></i>
    }</React.Fragment>
  );
}

const getColumns = (currentOrgID, currentApplicationID) => (
  [
    {
      dataField: 'devEUI',
      text: i18n.t(`${packageNS}:tr000371`),
      sort: false,
    }, {
      dataField: 'fCnt',
      text: i18n.t(`${packageNS}:tr000294`),
      sort: false,
    }, {
      dataField: 'fPort',
      text: i18n.t(`${packageNS}:tr000295`),
      sort: false,
    }, {
      dataField: 'confirmed',
      text: i18n.t(`${packageNS}:tr000296`),
      sort: false,
      formatter: DeviceQueueFConfirmedColumn,
      formatExtraData: {
        currentOrgID: currentOrgID,
        currentApplicationID: currentApplicationID
      }
    }, {
      dataField: 'data',
      text: i18n.t(`${packageNS}:tr000297`),
      sort: false,
    }
  ]
);

class QueueCard extends Component {
  constructor() {
    super();

    this.state = {
      data: [],
    };
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);

    DeviceQueueStore.on("enqueue", this.getQueue);
  }

  componentWillUnmount() {
    DeviceQueueStore.removeListener("enqueue", this.getQueue);
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
    const filters = {
      limit: limit,
      offset: offset,
      organizationID: this.props.organizationID,
      search: "",
    };
    this.getQueue();
  }

  getQueue = () => {
    this.setState({ loading: true });

    DeviceQueueStore.list(this.props.match.params.devEUI, resp => {
      this.setState({
        data: resp.deviceQueueItems,
        loading: false
      });
    });
  }

  flushQueue = () => {
    if (window.confirm("Are you sure you want to flush the device queue?")) {
      DeviceQueueStore.flush(this.props.match.params.devEUI, resp => {
        this.getQueue();
      });
    }
  }

  render() {
    const { collapseCard, setCollapseCard } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;

    return(
      <Card>
        <ButtonGroup>
          <Button color="secondary" onClick={() => setCollapseCard(CURRENT_CARD)}>
            <i className={`mdi mdi-arrow-${collapseCard[CURRENT_CARD] ? 'up' : 'down'}`}></i>
            &nbsp;&nbsp;
            <h5 style={{ color: "#fff", display: "inline" }}>
              {i18n.t(`${packageNS}:tr000293`)}
            </h5>
          </Button>
          <Button
            onClick={this.getQueue}
            color="primary"
            style={{ borderRadius: "2.5px" }}
          >
            <Refresh style={{ fill: "white" }} />
          </Button>
          <Button
            onClick={this.flushQueue}
            color="danger"
            style={{ borderRadius: "2.5px" }}
          >
            <Delete style={{ fill: "white" }} />
          </Button>
          <br />
        </ButtonGroup>
        <Collapse isOpen={collapseCard[CURRENT_CARD]} style={{ marginTop: "10px" }}>
          <CardContent>
            {this.state.loading && <Loader />}
            <AdvancedTable
              data={this.state.data}
              columns={getColumns(currentOrgID, currentApplicationID)}
              keyField="devEUI"
              onTableChange={this.handleTableChange}
              rowsPerPage={10}
              searchEnabled={false}
            />
          </CardContent>
        </Collapse>
      </Card>
    );
  }
}

export default withRouter(QueueCard);

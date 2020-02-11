import React, { Component } from "react";

import { Row, Col } from "reactstrap";
import { EditingState } from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  TableHeaderRow,
  TableInlineCellEditing
} from "@devexpress/dx-react-grid-bootstrap4";

import i18n, { packageNS } from "../../i18n";

const getRowId = row => row.id;

const FocusableCell = ({ onClick, ...restProps }) => (
  <Table.Cell {...restProps} tabIndex={0} onFocus={onClick} />
);

class GatewayFormLBT extends Component {
  constructor(props) {
    super(props);

    this.state = {
      columns: [
        { name: "name", title: "Name" },
        { name: "gender", title: "Gender" },
        { name: "city", title: "City" }
      ],
      rows: [
        {
          id: 1,
          name: "abc",
          gender: "m",
          city: "c1"
        }
      ]
    };
  }

  componentDidMount() {}

  componentDidUpdate(prevProps) {}

  onCommitChanges = ({ added, changed, deleted }) => {
    let changedRows;
    const rows = this.state.rows;

    if (added) {
      const startingAddedId =
        rows.length > 0 ? rows[rows.length - 1].id + 1 : 0;
      changedRows = [
        ...rows,
        ...added.map((row, index) => ({
          id: startingAddedId + index,
          ...row
        }))
      ];
    }
    if (changed) {
      changedRows = rows.map(row =>
        changed[row.id] ? { ...row, ...changed[row.id] } : row
      );
    }
    if (deleted) {
      const deletedSet = new Set(deleted);
      changedRows = rows.filter(row => !deletedSet.has(row.id));
    }
    this.setState({ rows: changedRows });
  };

  render() {
    return (
      <React.Fragment>
        <Row>
          <Col>
            <h5>{i18n.t(`${packageNS}:tr000598`)}</h5>
            <div className="card">
              <Grid
                rows={this.state.rows}
                columns={this.state.columns}
                getRowId={getRowId}
              >
                <EditingState onCommitChanges={this.onCommitChanges} />
                <Table cellComponent={FocusableCell} />
                <TableHeaderRow />

                <TableInlineCellEditing
                  startEditAction="click"
                  selectTextOnEditStart={true}
                />
              </Grid>
            </div>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default GatewayFormLBT;

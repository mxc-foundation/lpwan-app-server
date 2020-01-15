import React, { Component } from "react";
import { Collapse, Button } from 'reactstrap';

import Card from "@material-ui/core/Card";
import CardHeader from "@material-ui/core/CardHeader";
import CardContent from "@material-ui/core/CardContent";
import Table from "@material-ui/core/Table";
import TableRow from "@material-ui/core/TableRow";
import TableCell from "@material-ui/core/TableCell";
import TableBody from "@material-ui/core/TableBody";

import moment from "moment";

import i18n, { packageNS } from '../../../../i18n';

const CURRENT_CARD = "statusCard";

class StatusCard extends Component {
  render() {
    const { collapseCard, setCollapseCard } = this.props;
    let lastSeenAt = i18n.t(`${packageNS}:tr000372`);

    if (this.props.device.lastSeenAt !== null) {
      lastSeenAt = moment(this.props.device.lastSeenAt).format("lll");
    }

    return(
      <Card>
        <Button color="secondary" onClick={() => setCollapseCard(CURRENT_CARD)}>
          <i className={`mdi mdi-arrow-${collapseCard[CURRENT_CARD] ? 'up' : 'down'}`}></i>
          &nbsp;&nbsp;
          <h5 style={{ color: "#fff", display: "inline" }}>
            {i18n.t(`${packageNS}:tr000282`)}
          </h5>
        </Button>
        <Collapse isOpen={collapseCard[CURRENT_CARD]}>
          <CardContent>
            <Table>
              <TableBody>
                <TableRow>
                  <TableCell>{i18n.t(`${packageNS}:tr000283`)}</TableCell>
                  <TableCell>{lastSeenAt}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Collapse>
      </Card>
    );
  }
}

export default StatusCard;

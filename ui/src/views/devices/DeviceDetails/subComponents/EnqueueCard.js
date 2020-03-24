import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import { Button, Collapse } from 'reactstrap';
import i18n, { packageNS } from '../../../../i18n';
import DeviceQueueStore from "../../../../stores/DeviceQueueStore";
import DeviceQueueItemForm from "./DeviceQueueItemForm";



const CURRENT_CARD = "enqueueCard";

class EnqueueCard extends Component {
  constructor() {
    super();

    this.state = {
      object: {},
    };
  }

  onSubmit = (queueItem) => {
    let qi = queueItem;
    qi.devEUI = this.props.match.params.devEUI;

    DeviceQueueStore.enqueue(qi, resp => {
      this.setState({
        object: {},
      });
    });
  }

  render() {
    const { collapseCard, setCollapseCard } = this.props;

    return(
      <Card>
        <Button color="secondary" onClick={() => setCollapseCard(CURRENT_CARD)}>
          <i className={`mdi mdi-arrow-${collapseCard[CURRENT_CARD] ? 'up' : 'down'}`}></i>
          &nbsp;&nbsp;
          <h5 style={{ color: "#fff", display: "inline" }}>
            {i18n.t(`${packageNS}:tr000284`)}
          </h5>
        </Button>
        <Collapse isOpen={collapseCard[CURRENT_CARD]}>
          <CardContent>
            <DeviceQueueItemForm
              submitLabel={i18n.t(`${packageNS}:tr000292`)}
              onSubmit={this.onSubmit}
              object={this.state.object}
            />
          </CardContent>
        </Collapse>
      </Card>
    );
  }
}

export default withRouter(EnqueueCard);

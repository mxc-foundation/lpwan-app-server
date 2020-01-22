import React, { useState } from 'react';
import { Tooltip } from 'reactstrap';

const Tooltips = (props) => {
  const [tooltipOpen, setTooltipOpen] = useState(false);

  const toggle = () => setTooltipOpen(!tooltipOpen);

  return (
    <div>
      <Tooltip placement="right" isOpen={tooltipOpen} target={props.targeId} toggle={toggle}>
        Hello world!
      </Tooltip>
    </div>
  );
}

export default Tooltips;
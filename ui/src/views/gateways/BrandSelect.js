import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Button, Modal, ModalHeader, ModalBody, ModalFooter, Card, CardBody, Row, Col } from 'reactstrap';

import { withStyles } from "@material-ui/core/styles";
import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import Loader from "../../components/Loader";
import CommonModal from '../../components/Modal';
import logo from '../../assets/images/MATCHX-SUPERNODE2.png';
import { Divider } from "@material-ui/core";

const styles = {
    center: {
        display: "flex",
        justifyContent: "center"
    }
};

class BrandSelect extends Component {
    constructor(props) {
        super(props);
        this.state = {};
    }


    onSubmit = () => {

    }

    back = () => {
        this.props.history.push(`/organizations/${this.props.match.params.organizationID}/gateways`);
    }

    render() {
        const { classes } = this.props;

        const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

        return (<React.Fragment>
            <Card>
                <Row>
                    <Col>
                        <Card>
                            <CardBody className={classes.center}>
                                <span style={{ fontSize: '24px' }}>{i18n.t(`${packageNS}:menu.gateways.choose_a_manufacturer`)}</span>
                            </CardBody>
                        </Card>
                    </Col>
                </Row>
                <Row>
                    <Col>
                        <Card>
                            <CardBody className={classes.center}>
                                <Link to={`/organizations/${currentOrgID}/gateways/input-serial`}><img src={logo} alt="" height="136" /></Link>
                            </CardBody>
                        </Card>
                    </Col>
                </Row>
                <Row>
                    <Col>
                        <Card>
                            <CardBody className={classes.center} >
                                <Link to={`/organizations/${currentOrgID}/gateways/create`}><b>{i18n.t(`${packageNS}:menu.gateways.other`)}</b></Link>
                            </CardBody>
                        </Card>
                    </Col>
                </Row>
                <Divider />
                <Row>
                    <Col className={classes.center}>
                        <Button color="secondary" onClick={this.back}>{i18n.t(`${packageNS}:menu.common.cancel`)}</Button>
                    </Col>
                </Row>
            </Card>
        </React.Fragment>
        );
    }
}

export default withStyles(styles)(withRouter(BrandSelect));

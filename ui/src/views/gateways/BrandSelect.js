import { Divider } from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';
import { Button, Card, CardBody, Col, Row } from 'reactstrap';
import logo from '../../assets/images/matchx.png';
import i18n, { packageNS } from '../../i18n';



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

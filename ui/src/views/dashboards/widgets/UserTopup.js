import React from "react";
import { Link } from "react-router-dom";
import { Badge, Col, Media, Row } from 'reactstrap';
import defaultProfilePic from "../../../assets/images/users/profile-icon.png";
import i18n, { packageNS } from '../../../i18n';
import SessionStore from "../../../stores/SessionStore";
import WidgetActions from './WidgetActions';



/**
 * Topup
 * @param {*} props 
 */
const Topup = (props) => {
    const data = props.data || {};
    const { profilePic, username }  = SessionStore.getUser();
    const formattedVal = (data.amount || 0).toLocaleString(navigator.language, { minimumFractionDigits: 4 });

    const userRole = SessionStore.isOrganizationAdmin() ? i18n.t(`${packageNS}:menu.dashboard.roleOrgAdmin`) :
        SessionStore.isOrganizationDeviceAdmin() ? i18n.t(`${packageNS}:menu.dashboard.roleDeviceAdmin`) :
            SessionStore.isOrganizationGatewayAdmin() ? i18n.t(`${packageNS}:menu.dashboard.roleGatewayAdmin`) : "";

    const orgId = SessionStore.getOrganizationID();

    return <div className="card-box">
        <div className="float-right">
            <WidgetActions widget={props.widget} actionItems={[{ to: '#', label: 'Week' }]} onDelete={props.onDelete} />
        </div>

        <Media className="align-items-center">
            <Media left className="avatar-xl">
                <img src={profilePic || defaultProfilePic} className="img-fluid rounded-circle" alt="user" />
            </Media>
            <Media body>
                <div className="ml-2">
                    <h3 className="font-weight-normal mt-0">{username}</h3>
                    <h5 className="text-primary mb-0">{userRole}</h5>
                </div>
            </Media>
        </Media>

        <Row className="mt-3">
            <Col className="text-right mb-0">
                {data.growth ? <h5><Badge className="px-1">
                    {data.growth} <i className="mdi mdi-arrow-up"></i></Badge></h5> : null}
                <h2 className="my-2 font-2rem">{formattedVal} MXC</h2>
                <Link className="btn btn-primary" to={`/topup/${orgId}`}>{i18n.t(`${packageNS}:menu.dashboard.topupButton`)}</Link>
            </Col>
        </Row>
    </div>;
}

export default Topup;
import React, { Component } from "react";
import { Link } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem } from 'reactstrap';
import i18n, { packageNS } from '../i18n';
import OrganizationStore from "../stores/OrganizationStore";
import SessionStorage from "../stores/SessionStore";


class OrgBreadcrumb extends Component {
  constructor(props) {
    super(props);
    this.state = {
      currentOrg: null
    };
  }

  componentDidMount() {
    this.getOrganization();
  }
  
  getOrganization = async () => {
    if (!this.state.currentOrg) {
      const orgId = this.props.organizationID;
      const orgs = SessionStorage.getOrganizations();
      let org = null;
      if (orgs && orgs.length) {
        org = orgs.find(o => o.organizationID === orgId);

        if (!org) {
          let organization = await OrganizationStore.get(orgId);
          this.setState({ currentOrg: organization.organization });
        } else {
          this.setState({ currentOrg: org });
        }
      }
    }
  }

  render() {
    let currentOrgName = this.state.currentOrg ? (this.state.currentOrg.organizationName || this.state.currentOrg.name) : "";
    let currentOrgFullName = null;
    if (currentOrgName.length > 5) {
      currentOrgFullName = currentOrgName;
      currentOrgName = currentOrgName.slice(0, 5) + "...";
    }

    return (
      <Breadcrumb>
        <BreadcrumbItem>
          <Link to={`/organizations`} onClick={(e) => { if (this.props.orgListCallback) this.props.orgListCallback(e) }}>
            {i18n.t(`${packageNS}:tr000049`)}
          </Link>
        </BreadcrumbItem>
        <BreadcrumbItem>
          <Link to={`/organizations/${this.props.organizationID}`} title={currentOrgFullName}
            onClick={(e) => { if (this.props.orgNameCallback) this.props.orgNameCallback(e) }}
          >
            {currentOrgName}
          </Link>
        </BreadcrumbItem>

        {(this.props.items || []).map((item, idx) => {
          return <BreadcrumbItem key={idx} active={item.active}>
            {item.to ? <Link to={item.to} onClick={(e) => { if (item.onClick) item.onClick(e) }}>{item.label}</Link> : <React.Fragment>
              {item.label}</React.Fragment>}
          </BreadcrumbItem>
        })}
      </Breadcrumb>
    );
  }
}

export default OrgBreadcrumb;

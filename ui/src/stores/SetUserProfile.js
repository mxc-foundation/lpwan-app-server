import hash from "object-hash";
import sessionStore from "./SessionStore";

export default function updateOrganizations(response) {
    const organizationList = response.body.organizations;

    if(!organizationList){
        return false;
    }
    if(sessionStore.getOrganizationList() !== null){
        if(hash(sessionStore.getOrganizationList()) !== hash(organizationList)){
            
            sessionStore.setOrganizationList(organizationList);
        }
    }else{
        sessionStore.setOrganizationList(organizationList);
    }
     
    return response; 
};
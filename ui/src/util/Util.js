import SessionStore from "../stores/SessionStore";

export function getM2MLink() {
    let host = process.env.REACT_APP_MXPROTOCOL_SERVER;
    const origin = window.location.origin;
    
    if(origin.includes(process.env.REACT_APP_SUBDOM_LORA)){
        host = origin.replace(process.env.REACT_APP_SUBDOM_LORA, process.env.REACT_APP_SUBDOM_M2M);
    }
    
    return host;
}

export function openM2M(org, isBelongToOrg, path) {
    let orgName = org.name;
    let orgId = org.id;
    
    if(!orgId){
      return false;
    }
    const user = SessionStore.getUser();

    if(user.isAdmin && !isBelongToOrg){
      orgId = '0';
      orgName = 'Super_admin';
    }

    const data = {
      jwt: window.localStorage.getItem("jwt"),
      path: `${path}/${orgId}`,
      orgId,
      orgName,
      username: user.username,
      loraHostUrl: window.location.origin,
      language: SessionStore.getLanguage()
    };

    const dataString = encodeURIComponent(JSON.stringify(data));

    const host = getM2MLink();

    // for new tab, see: https://stackoverflow.com/questions/427479/programmatically-open-new-pages-on-tabs
    window.location.replace(host + `/#/j/${dataString}`);
  }
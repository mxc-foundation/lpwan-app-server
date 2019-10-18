export function getM2MLink() {
    let host = process.env.REACT_APP_MXPROTOCOL_SERVER;
    const origin = window.location.origin;
    
    if(origin.includes(process.env.REACT_APP_SUBDOM_LORA)){
        host = origin.replace(process.env.REACT_APP_SUBDOM_LORA, process.env.REACT_APP_SUBDOM_M2M);
    }
    
    return host;
}
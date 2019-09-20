export function getM2MLink() {
    let host = process.env.REACT_APP_M2M_LOCAL_SERVER;
    const origin = window.location.origin;
    
    if(origin.includes(process.env.REACT_APP_DEMO_HOST_SERVER)){
        host = process.env.REACT_APP_M2M_DEMO_SERVER;
    }
    if(origin.includes(process.env.REACT_APP_TEST_HOST_SERVER)){
        host = process.env.REACT_APP_M2M_TEST_SERVER;
    }
    
    return host;
}
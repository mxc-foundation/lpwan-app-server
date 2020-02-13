# basic
import sys

from RESTful_api.send_request import init_config
from RESTful_api.network_server import getNetworkServerID
from script.global_admin_user import user_login
from RESTful_api.gateway_profile import getGatewayProfileID
from RESTful_api.gateway import createGateway
from RESTful_api.service_profile import getServiceProfileID
from RESTful_api.application import getApplicationID
from RESTful_api.device_profile import getDeviceProfileID
from RESTful_api.device import createDevice

if __name__ == "__main__":
    if 1 == len(sys.argv):
        print("Invalid input")
        exit(1)

    # init
    srv_url = sys.argv[1]
    user_info = sys.argv[2]
    init_config(srv_url=srv_url)

    # user login
    if not user_login(user_info):
        print("User {} login failed!".format(user_info))
        exit(1)

    # get network server id
    networkServerID = getNetworkServerID()

    # get gateway profile id
    gatewayProfileID = getGatewayProfileID(networkServerID)

    # create gateways
    for i in range(10):
        createGateway(networkServerID, gatewayProfileID)

    # get service profile
    serviceProfileID = getServiceProfileID(networkServerID)

    # get application id
    applicationID = getApplicationID(serviceProfileID)

    # get device profile id
    deviceProfileID = getDeviceProfileID(applicationID, networkServerID)

    # create devices
    for i in range(10):
        createDevice(applicationID, deviceProfileID)

    print("end")





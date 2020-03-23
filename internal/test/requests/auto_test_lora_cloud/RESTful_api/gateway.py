import os
import binascii

from RESTful_api.send_request import get_request, post_request
from script.global_admin_user import all_users, user_login
from RESTful_api.network_server import getNetworkServerID
from RESTful_api.gateway_profile import getGatewayProfileID


def getGateway(user):
    response = get_request(api_url='api/gateways?limit=999&offset=0&organizationID=1', jwt=user.get_user_jwt(), data='')
    if response.status_code != 200:
        print("api/gateways?limit=999&offset=0&organizationID=1 response code {}".format(response.status_code))
        exit()

    tmpRes = response.text.replace('null', '123')
    result = eval(tmpRes)
    if eval(result['totalCount']) == 0:
        return "null"
    else:
        gatewayList = result['result']

    return gatewayList


def createGateway(networkServerID, gatewayProfileID):
    user = all_users.get_user_list()[0]
    mac = binascii.hexlify(os.urandom(8)).decode()
    gateway = '''
    {{
      "gateway": {{
        "boards": [],
        "description": "auto generated",
        "discoveryEnabled": false,
        "gatewayProfileID": "{}",
        "id": "{}",
        "location": {{
          "accuracy": 0,
          "altitude": 0,
          "latitude": 0,
          "longitude": 0,
          "source": "UNKNOWN"
        }},
        "name": "gateway_{}",
        "networkServerID": "{}",
        "organizationID": "1"
      }}
    }}
    '''.format(gatewayProfileID, mac, mac, networkServerID)

    response = post_request(api_url='api/gateways', jwt=user.get_user_jwt(), data=gateway)
    if response.status_code != 200:
        print("api/gateways response code {}".format(response.status_code))
        exit()

    print("gateway with mac {} created..".format(mac))

if __name__ == "__main__":
    # user login
    user_info = '{"password": "admin", "username": "admin"}'

    if not user_login(user_info):
        print("User {} login failed!".format(user_info))
        exit(1)

    user = all_users.get_user_list()[0]
    networkServerID = getNetworkServerID()
    gatewayProfileID = getGatewayProfileID(networkServerID)

    gatewayList = getGateway(user)
    createGateway(user, networkServerID, gatewayProfileID)



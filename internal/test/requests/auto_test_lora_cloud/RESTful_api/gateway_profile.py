from RESTful_api.send_request import get_request, post_request
from script.global_admin_user import all_users, user_login
from RESTful_api.network_server import getNetworkServerID

def getGatewayProfile(user, networkServerID):
    response = get_request(api_url='api/gateway-profiles?limit=999&offset=0&networkServerID={}'.format(networkServerID), jwt=user.get_user_jwt(), data='')
    if response.status_code != 200:
        print("api/gateway-profiles?limit=999&offset=0&networkServerID={} response code {}".format(networkServerID, response.status_code))
        exit()

    tmpRes = response.text.replace('null', '123')
    result = eval(tmpRes)
    if eval(result['totalCount']) == 0:
        return "null"
    else:
        gatewayProfileList = result['result']

    return gatewayProfileList


def getGatewayProfileByID(id):
    response = get_request(api_url='api/gateway-profiles/{}'.format(id), jwt=user.get_user_jwt(), data='')
    if response.status_code != 200:
        print("api/gateway-profiles/{} response code {}".format(id, response.status_code))
        exit()

    result = eval(response.text)
    return result


def createGatewayProfile(user, networkServerID):
    gatewayProfile = '''
    {{
      "gatewayProfile": {{
        "channels": [
          0,1,2,3,4,5,6,7,8
        ],
        "extraChannels": [],
        "id": "0",
        "name": "gatewayProfile",
        "networkServerID": "{}"
      }}
    }}
    '''.format(networkServerID)

    response = post_request(api_url='api/gateway-profiles', jwt=user.get_user_jwt(), data=gatewayProfile)
    if response.status_code != 200:
        print("api/gateway-profiles response code {}".format(response.status_code))
        exit()

    result = eval(response.text)
    return result['id']


def getGatewayProfileID(networkServerID):
    user = all_users.get_user_list()[0]
    gatewayProfileList = getGatewayProfile(user, networkServerID)
    if "null" == gatewayProfileList:
        gatewayProfileID = createGatewayProfile(user, networkServerID)
        return gatewayProfileID
    else:
        gatewayProfile = gatewayProfileList[0]
        return gatewayProfile['id']


if __name__ == "__main__":
    # user login
    user_info = '{"password": "admin", "username": "admin"}'

    if not user_login(user_info):
        print("User {} login failed!".format(user_info))
        exit(1)

    user = all_users.get_user_list()[0]
    networkServerID = getNetworkServerID()

    gatewayProfileList = getGatewayProfile(user, networkServerID)
    if "null" == gatewayProfileList:
        gatewayProfileID = createGatewayProfile(user, networkServerID)
        print(gatewayProfileID)
    else:
        gatewayProfile = gatewayProfileList[0]
        print(gatewayProfile['id'])

from RESTful_api.send_request import get_request, post_request
from script.global_admin_user import all_users, user_login


def get_network_server(user):
    response = get_request(api_url='api/network-servers?limit=999&offset=0&organizationID=0', jwt=user.get_user_jwt(), data='')
    if response.status_code != 200:
        print("api/network-servers?limit=999&offset=0&organizationID=0 response code {}".format(response.status_code))
        exit()

    result = eval(response.text)
    if eval(result['totalCount']) == 0:
        return "null"
    else:
        networkServerList = result['result']

    return networkServerList


def getNetworkServerByID(id):
    response = get_request(api_url='api/network-servers/{}'.format(id), jwt=user.get_user_jwt(), data='')
    if response.status_code != 200:
        print("api/network-servers/{} response code {}".format(id, response.status_code))
        exit()

    result = eval(response.text)
    return result


def create_network_server(user):
    network_server = '''
    {
        "networkServer": {
            "caCert": "",
            "gatewayDiscoveryDR": 0,
            "gatewayDiscoveryEnabled": false,
            "gatewayDiscoveryInterval": 0,
            "gatewayDiscoveryTXFrequency": 0,
            "id": "0",
            "name": "network-server",
            "routingProfileCACert": "",
            "routingProfileTLSCert": "",
            "routingProfileTLSKey": "",
            "server": "network-server:8000",
            "tlsCert": "",
            "tlsKey": ""
        }
    }
    '''

    response = post_request(api_url='api/network-servers', jwt=user.get_user_jwt(), data=network_server)
    if response.status_code != 200:
        print("api/network-servers response code {}".format(response.status_code))
        exit()

    result = eval(response.text)
    return eval(result['id'])


def getNetworkServerID():
    user = all_users.get_user_list()[0]
    networkServerList = get_network_server(user)
    if "null" == networkServerList:
        networkServerID = create_network_server(user)
        networkServer = getNetworkServerByID(networkServerID)
    else:
        networkServer = networkServerList[0]

    networkServerID = networkServer['id']
    print(networkServerID)

    return networkServerID


if __name__ == "__main__":
    # user login
    user_info = '{"password": "admin", "username": "admin"}'

    if not user_login(user_info):
        print("User {} login failed!".format(user_info))
        exit(1)

    user = all_users.get_user_list()[0]
    networkServerList = get_network_server(user)
    if "null" == networkServerList:
        networkServerID = create_network_server(user)
        networkServer = getNetworkServerByID(networkServerID)
    else:
        networkServer = networkServerList[0]

    print(networkServer['id'])
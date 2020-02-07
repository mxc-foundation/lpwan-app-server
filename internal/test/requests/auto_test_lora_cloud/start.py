# basic
import sys
# init
from RESTful_api.send_request import init_config
from RESTful_api.
# script
from script.global_admin_user import user_login, all_users


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
    networkServerID = networkServer.getNetworkServerID()

    # get gateway profile id
    gatewayProfileID =

    print("end")





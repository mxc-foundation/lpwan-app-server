from RESTful_api.send_request import get_request, post_request
from model.gateways import Gateway
false = False
true = True


# get gateway list from a particular organization
def gateway_list(jwt, oid, limit, offset):
    '''
    updated
    :param jwt:
    :param oid:
    :param limit:
    :param offset:
    :return: a list of gateway objects
    '''
    limit = 100
    offset = 0
    response = get_request(api_url='api/gateways?limit={}&offset={}&organizationID={}'.format(limit, offset, oid),
                           jwt=jwt, data='')
    if 200 != response.status_code:
        print("Gateway list request error: {} \n {}".format(response.status_code, response.text))
        return "null"

    try:
        result = eval(response.text)
        glist = result["result"]
    except TypeError:
        return "null"
    except KeyError:
        return "null"

    gateway_obj_list = []
    for item in glist:
        gateway = Gateway()
        gateway.init_gateway(item)
        gateway_obj_list.append(gateway)
    return gateway_obj_list



import requests

config = {"server": "http://localhost:8080/"}


def init_config(srv_url):
    config["server"] = srv_url


def get_request(api_url, jwt, data=''):
    headers = dict()
    headers['Accept'] = 'application/json'
    if '' != jwt:
        headers['Grpc-Metadata-Authorization'] = jwt
    url = config["server"] + api_url
    return requests.get(url=url, headers=headers, data=data)


def post_request(api_url, jwt, data=''):
    headers = dict()
    headers['Content-Type'] = 'application/json'
    headers['Accept'] = 'application/json'
    if '' != jwt:
        headers['Grpc-Metadata-Authorization'] = jwt
    url = config["server"] + api_url
    return requests.post(url=url, headers=headers, data=data)


if __name__ == "__main__":
    headers = dict()
    headers['Content-Type'] = 'application/json'
    headers['Accept'] = 'application/json'
    data = '{"password": "appuser1", "username": "appuser1"}'
    res = post_request('https://lora-test-srv.matchx.io/api/internal/login', '', data)
    print('')

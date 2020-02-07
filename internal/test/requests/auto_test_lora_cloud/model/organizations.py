class Organization:
    def __init__(self):
        self.id = ''
        self.name = ''
        self.displayName = ''
        self.canHaveGateways = ''
        self.createdAt = ''
        self.updatedAt = ''
        self.user_list = []
        self.gateway_list = []
        self.app_list = []
        self.admin_user = ''

    def init_organization(self, data):
        if self._init_from_organization_list(data):
            return
        if self._init_from_organization_id(data):
            return

    def _init_from_organization_list(self, data):
        try:
                data["id"]
                data["name"]
                data["displayName"]
                data["canHaveGateways"]
                data["createdAt"]
                data["updatedAt"]
        except KeyError:
            return False

        self.id = data["id"]
        self.name = data["name"]
        self.displayName = data["displayName"]
        self.canHaveGateways = data["canHaveGateways"]
        self.createdAt = data["createdAt"]
        self.updatedAt = data["updatedAt"]
        return True

    def _init_from_organization_id(self, data):
        try:
            data["organization"]["id"]
            data["organization"]["name"]
            data["organization"]["displayName"]
            data["organization"]["canHaveGateways"]
            data["createdAt"]
            data["updatedAt"]
        except KeyError:
            return False

        self.id = data["organization"]["id"]
        self.name = data["organization"]["name"]
        self.displayName = data["organization"]["displayName"]
        self.canHaveGateways = data["organization"]["canHaveGateways"]
        self.createdAt = data["createdAt"]
        self.updatedAt = data["updatedAt"]
        return True

    def add_admin_user(self, uid):
        self.admin_user = uid

    def get_admin_user_id(self):
        return self.admin_user

    def get_organization_id(self):
        return self.id

    def get_organization_dname(self):
        return self.displayName

    def can_have_gateways(self):
        return self.canHaveGateways

    def is_id(self, oid):
        return oid == self.id

    def add_user_to_org(self, user_id):
        self.user_list.append(user_id)

    def add_app_to_org(self, app_id):
        self.app_list.append(app_id)

    def add_gateway_to_org(self, gateway_id):
        self.gateway_list.append(gateway_id)

    def add_user_id_list(self, ulist):
        self.user_list = ulist

    def get_user_id_list(self):
        return self.user_list

    def add_app_id_list(self, applist):
        self.app_list = applist

    def get_app_id_list(self):
        return self.app_list

    def add_gateway_id_list(self, gateway_list):
        self.gateway_list = gateway_list

    def get_gateway_id_list(self):
        return self.gateway_list


class OrganizationList:
    def __init__(self):
        self.organization_list = []

    def insert_organization(self, org):
        self.organization_list.append(org)

    def append_organization_list(self, org_list):
        self.organization_list += org_list

    def get_organization_list(self):
        return self.organization_list

    def get_organization_by_id(self, oid):
        for item in self.organization_list:
            if item.is_id(oid):
                return item
        return "null"



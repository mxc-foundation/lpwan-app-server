class User:
    '''
    Should break down all elements into smallest unit
    '''
    def __init__(self):
        self.user = dict()
        self.user["id"] = ''
        self.user["username"] = ''
        self.user["sessionTTL"] = ''
        self.user["isAdmin"] = False
        self.user["isGlobalAdmin"] = False
        self.user["isActive"] = ''
        self.user["email"] = ''
        self.user["note"] = ''
        self.user["createdAt"] = ''
        self.user["updatedAt"] = ''

        self.organization_ids = []
        self.organization_ids_admin = []

        self.settings = dict()
        self.settings["disableAssignExistingUsers"] = ''
        self.jwt = ''

    def _init_from_internal_user_profile(self, data):
        try:
            # user info got from internal api
            data["user"]["id"]
            data["user"]["username"]
            data["user"]["sessionTTL"]
            data["user"]["isAdmin"]
            data["user"]["isActive"]
            data["user"]["email"]
            data["user"]["note"]

            if 0 != len(data["organizations"]):
                item = data["organizations"][0]
                item["organizationID"]
                item["organizationName"]
                item["isAdmin"]
                item["createdAt"]
                item["updatedAt"]

            data["settings"]["disableAssignExistingUsers"]
        except KeyError:
            return False

        self.user["id"] = data["user"]["id"]
        self.user["username"] = data["user"]["username"]
        self.user["sessionTTL"] = data["user"]["sessionTTL"]
        self.user["isActive"] = data["user"]["isActive"]
        self.user["email"] = data["user"]["email"]
        self.user["note"] = data["user"]["note"]
        if data["user"]["isAdmin"]:
            self.user["isGlobalAdmin"] = True

        for item in data["organizations"]:
            if item["isAdmin"]:
                self.organization_ids_admin.append(item["organizationID"])
            else:
                self.organization_ids.append(item["organizationID"])
        self.settings["disableAssignExistingUsers"] = data["settings"]["disableAssignExistingUsers"]
        return True

    def _init_from_organization_users(self, data):
        try:
            data["result"]["userID"]
            data["result"]["username"]
            data["result"]["isAdmin"]
            data["result"]["createdAt"]
            data["result"]["updatedAt"]
        except KeyError:
            return False

        self.user["id"] = data["result"]["userID"]
        self.user["username"] = data["result"]["username"]
        self.user["isAdmin"] = data["result"]["isAdmin"]
        self.user["createdAt"] = data["result"]["createdAt"]
        self.user["updatedAt"] = data["result"]["updatedAt"]
        return True

    def _init_from_organization_user_id(self, data):
        try:
            data["organizationUser"]["organizationID"]
            data["organizationUser"]["userID"]
            data["organizationUser"]["isAdmin"]
            data["organizationUser"]["username"]
            data["createdAt"]
            data["updatedAt"]
        except KeyError:
            return False

        self.user["id"] = data["organizationUser"]["userID"]
        self.user["username"] = data["organizationUser"]["username"]
        self.user["isAdmin"] = data["organizationUser"]["isAdmin"]
        self.user["createdAt"] = data["createdAt"]
        self.user["updatedAt"] = data["updatedAt"]
        if data["organizationUser"]["isAdmin"]:
            self.organization_ids_admin.append(data["organizationUser"]["organizationID"])
        else:
            self.organization_ids.append(data["organizationUser"]["organizationID"])
        return True

    def _init_from_user_user_id(self, data):
        try:
            data["user"]["id"]
            data["user"]["username"]
            data["user"]["sessionTTL"]
            data["user"]["isAdmin"]
            data["user"]["isActive"]
            data["user"]["email"]
            data["user"]["note"]
            data["createdAt"]
            data["updatedAt"]
        except KeyError:
            return False

        self.user["id"] = data["user"]["id"]
        self.user["username"] = data["user"]["username"]
        self.user["sessionTTL"] = data["user"]["sessionTTL"]
        self.user["isActive"] = data["user"]["isActive"]
        self.user["email"] = data["user"]["email"]
        self.user["note"] = data["user"]["note"]
        self.user["createdAt"] = data["createdAt"]
        self.user["updatedAt"] = data["updatedAt"]
        if data["user"]["isAdmin"]:
            self.user["isGlobalAdmin"] = True
        return True

    def init_user(self, data):
        if self._init_from_internal_user_profile(data):
            return
        if self._init_from_organization_users(data):
            return
        if self._init_from_organization_user_id(data):
            return
        if self._init_from_user_user_id(data):
            return
        print("Model.users init user object failed!")

    def set_user_jwt(self, jwt):
        self.jwt = jwt

    def get_user_jwt(self):
        return self.jwt

    def get_user_id(self):
        return self.user["id"]

    def get_username(self):
        return self.user["username"]

    def is_user_gadmin(self):
        return self.user["isGlobalAdminAdmin"]

    def is_user_oadmin(self):
        return self.user["isAdmin"]

    def is_user_active(self):
        return self.user["isActive"]

    def get_user_orglist(self):
        return self.organization_ids

    def is_id(self, uid):
        return uid == self.user["id"]

    def mark_user_in_organization(self, org_id, is_admin=False):
        if is_admin:
            self.organization_ids_admin.append(org_id)
        else:
            self.organization_ids.append(org_id)


class UserList:
    def __init__(self):
        self.user_list = []

    def insert_user(self, user):
        self.user_list.append(user)

    def append_user_list(self, ulist):
        self.user_list += ulist

    def get_user_list(self):
        return self.user_list

    def get_user_by_id(self, uid):
        for item in self.user_list:
            if item.is_id(uid):
                return item
        return "null"


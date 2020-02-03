/// <reference types="jest" />

import React from 'react';
import { mount, configure } from 'enzyme';
import { MemoryRouter, Route } from 'react-router-dom'
import CreateGateway from './CreateGateway';
import Adapter from 'enzyme-adapter-react-16'


configure({ adapter: new Adapter() })

it('CreateGateway', async () => {
    const orgId = "1";
    const wrapper = mount(
        <MemoryRouter initialEntries={[`/organizations/${orgId}/gateways/create`]}>
            <Route exact path="/organizations/:organizationID(\d+)/gateways/create" component={CreateGateway} />
        </MemoryRouter>
    );
    //   let tree = component.toJSON();
    //   expect(tree).toMatchSnapshot();

    await ServiceProfileStore.list(orgId, 0, 0);

    expect(wrapper.state().loading).toBe(true);

    expect(wrapper.find(`[to="/organizations/${orgId}"]`)).toHaveLength(1);

    expect(wrapper.state().loading).toBe(false);
});

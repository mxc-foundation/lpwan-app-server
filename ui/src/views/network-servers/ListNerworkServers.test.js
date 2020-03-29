import React from 'react';
import { mount } from 'enzyme';
import { MemoryRouter } from 'react-router-dom'
import ListNetworkServers from './ListNetworkServers';


it('ListNetworkServer list', () => {
  const wrapper = mount(
    <MemoryRouter>
      <ListNetworkServers ></ListNetworkServers>
    </MemoryRouter>
  );
  expect(wrapper).toMatchSnapshot();
});


import React from 'react';
import { mount, configure } from 'enzyme';
import { MemoryRouter } from 'react-router-dom'
import ListNetworkServers from './ListNetworkServers';
import Adapter from 'enzyme-adapter-react-16';


configure({ adapter: new Adapter() })

it('ListNetworkServer list', () => {
  const wrapper = mount(
    <MemoryRouter>
      <ListNetworkServers ></ListNetworkServers>
    </MemoryRouter>
  );
  expect(wrapper).toMatchSnapshot();
});


import React from 'react';
import { shallow } from 'enzyme';
import ListNetworkServers from './ListNetworkServers';

jest.mock('history',  () =>  ({
  createHashHistory: jest.fn(),
  createMemoryHistory: jest.fn(),
}))

it('ListNetworkServer list', () => {

  const wrapper = shallow(
    <ListNetworkServers />
  );
  expect(wrapper).toMatchSnapshot();
});


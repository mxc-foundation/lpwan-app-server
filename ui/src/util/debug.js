import debug from 'debug';

import { name } from '../../package.json';

const Debug = namespace => debug(`${name}:${namespace}`);

export default Debug;

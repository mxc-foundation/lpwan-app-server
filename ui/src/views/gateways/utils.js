/**
 * Extracts the channels
 * @param {*} config
 */
const getChannels = config => {
  const obj = config[Object.keys(config)[0]];
  let channels = [];
  for (const f in obj) {
    if (f.startsWith("chan_")) {
      channels.push({ channel: f, ...obj[f] });
    }
  }
  return channels;
};

/**
 * Extracts the lbt channels
 * @param {*} config
 */
const getLBTChannels = config => {
  let lbtData = [];
  const channels = getChannels(config);
  const obj = config[Object.keys(config)[0]];

  if (
    obj &&
    obj.hasOwnProperty("lbt_cfg") &&
    obj["lbt_cfg"].hasOwnProperty("chan_cfg")
  ) {
    const lbtRecords = obj["lbt_cfg"]["chan_cfg"];

    for (let idx = 0; idx < lbtRecords.length; idx++) {
      lbtData.push({ channel: channels[idx]["channel"], ...lbtRecords[idx] });
    }
  }

  return lbtData;
};

/**
 * Extracts channels with frequency
 * @param {*} config
 */
const getChannelsWithFrequency = config => {
  let channelFreqData = [];
  const channels = getChannels(config);
  const obj = config[Object.keys(config)[0]];
  for (const channel of channels) {
    const radioName = "radio_" + channel.radio;
    const radioObj = obj.hasOwnProperty(radioName)
      ? obj[radioName]
      : { freq: 0 };
    const freq_hz = radioObj.freq + (channel.if ? channel.if : 0);
    channelFreqData.push({ ...channel, freq_hz: freq_hz });
  }
  return channelFreqData;
};

/**
 * Get Antenna gain data from object
 * @param {*} config
 */
const getAntennaGain = config => {
  let obj = config[Object.keys(config)[0]];
  return obj ? obj["antenna_gain"] : null;
};

/**
 * Extracts the lbt channels to get status
 * @param {*} config
 */
const getLBTConfigStatus = config => {
  let obj = config[Object.keys(config)[0]];
  console.log('getLBTConfigStatus util', obj);
  return obj ? obj["lbt_cfg"]["enable"] : null;
};

export { getChannels, getLBTChannels, getChannelsWithFrequency, getAntennaGain, getLBTConfigStatus };


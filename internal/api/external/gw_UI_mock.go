package external

var GwConf = "{" +
	"\"SX1301_conf\": {" +
	"\"lorawan_public\": true," +
	"\"clksrc\": 1," +
	"\"lbt_cfg\": {" +
	"\"enable\": true," +
	"\"rssi_target\": -81," +
	"\"chan_cfg\":[ " +
	"{ \"freq_hz\": 868100000, \"scan_time_us\": 5000 }, " +
	"{ \"freq_hz\": 868300000, \"scan_time_us\": 5000 }, " +
	"{ \"freq_hz\": 868500000, \"scan_time_us\": 5000 }, " +
	"{ \"freq_hz\": 868800000, \"scan_time_us\": 5000 }, " +
	"{ \"freq_hz\": 864700000, \"scan_time_us\": 5000 }, " +
	"{ \"freq_hz\": 864900000, \"scan_time_us\": 5000 }, " +
	"{ \"freq_hz\": 865100000, \"scan_time_us\": 5000 }, " +
	"{ \"freq_hz\": 869525000, \"scan_time_us\": 5000 }" +
	"]," +
	"\"sx127x_rssi_offset\": -7" +
	"}," +
	"\"antenna_gain\": 2.5," +
	"\"radio_0\": {" +
	"\"enable\": true," +
	"\"type\": \"SX1257\"," +
	"\"freq\": 864900000," +
	"\"rssi_offset\": -166.0," +
	"\"tx_enable\": true," +
	"\"tx_notch_freq\": 129000," +
	"\"tx_freq_min\": 863000000," +
	"\"tx_freq_max\": 870000000" +
	"}," +
	"\"radio_1\": {" +
	"\"enable\": true," +
	"\"type\": \"SX1257\"," +
	"\"freq\": 868500000," +
	"\"rssi_offset\": -166.0," +
	"\"tx_enable\": false" +
	"}," +
	"\"chan_multiSF_0\": {" +
	"\"enable\": true," +
	"\"radio\": 1," +
	"\"if\": -400000" +
	"}," +
	"\"chan_multiSF_1\": {" +
	"\"enable\": true," +
	"\"radio\": 1," +
	"\"if\": -200000" +
	"}," +
	"\"chan_multiSF_2\": {" +
	"\"enable\": true," +
	"\"radio\": 1," +
	"\"if\": 0" +
	"}," +
	"\"chan_multiSF_3\": {" +
	"\"enable\": true," +
	"\"radio\": 1," +
	"\"if\": 300000" +
	"}," +
	"\"chan_multiSF_4\": {" +
	"\"enable\": true," +
	"\"radio\": 0," +
	"\"if\": -200000" +
	"}," +
	"\"chan_multiSF_5\": {" +
	"\"enable\": true," +
	"\"radio\": 0," +
	"\"if\": 0" +
	"}," +
	"\"chan_multiSF_6\": {" +
	"\"enable\": true," +
	"\"radio\": 0," +
	"\"if\": 200000" +
	"}," +
	"\"chan_multiSF_7\": {" +
	"\"enable\": true," +
	"\"radio\": 0," +
	"\"if\": 400000" +
	"}," +
	"\"chan_Lora_std\": {" +
	"\"enable\": true," +
	"\"radio\": 1," +
	"\"if\": -200000," +
	"\"bandwidth\": 250000," +
	"\"spread_factor\": 7" +
	"}," +
	"\"chan_FSK\": {" +
	"\"enable\": true," +
	"\"radio\": 1," +
	"\"if\": 300000," +
	"\"bandwidth\": 125000," +
	"\"datarate\": 50000" +
	"}," +
	"\"tx_lut_0\": {" +
	"\"pa_gain\": 0," +
	"\"mix_gain\": 8," +
	"\"rf_power\": -6," +
	"\"dig_gain\": 2" +
	"}," +
	"\"tx_lut_1\": {" +
	"\"pa_gain\": 0," +
	"\"mix_gain\": 11," +
	"\"rf_power\": -3," +
	"\"dig_gain\": 3" +
	"}," +
	"\"tx_lut_2\": {" +
	"\"pa_gain\": 0," +
	"\"mix_gain\": 11," +
	"\"rf_power\": 0," +
	"\"dig_gain\": 1" +
	"}," +
	"\"tx_lut_3\": {" +
	"\"pa_gain\": 0," +
	"\"mix_gain\": 14," +
	"\"rf_power\": 3," +
	"\"dig_gain\": 0" +
	"}," +
	"\"tx_lut_4\": {" +
	"\"pa_gain\": 1," +
	"\"mix_gain\": 11," +
	"\"rf_power\": 6," +
	"\"dig_gain\": 3" +
	"}," +
	"\"tx_lut_5\": {" +
	"\"pa_gain\": 1," +
	"\"mix_gain\": 11," +
	"\"rf_power\": 10," +
	"\"dig_gain\": 0" +
	"}," +
	"\"tx_lut_6\": {" +
	"\"pa_gain\": 1," +
	"\"mix_gain\": 13," +
	"\"rf_power\": 11," +
	"\"dig_gain\": 2" +
	"}," +
	"\"tx_lut_7\": {" +
	"\"pa_gain\": 1," +
	"\"mix_gain\": 13," +
	"\"rf_power\": 12," +
	"\"dig_gain\": 1" +
	"}," +
	"\"tx_lut_8\": {" +
	"\"pa_gain\": 1," +
	"\"mix_gain\": 14," +
	"\"rf_power\": 13," +
	"\"dig_gain\": 1" +
	"}," +
	"\"tx_lut_9\": {" +
	"\"pa_gain\": 1," +
	"\"mix_gain\": 14," +
	"\"rf_power\": 14," +
	"\"dig_gain\": 0" +
	"}," +
	"\"tx_lut_10\": {" +
	"\"pa_gain\": 2," +
	"\"mix_gain\": 9," +
	"\"rf_power\": 16," +
	"\"dig_gain\": 0" +
	"}," +
	"\"tx_lut_11\": {" +
	"\"pa_gain\": 2," +
	"\"mix_gain\": 12," +
	"\"rf_power\": 20," +
	"\"dig_gain\": 1" +
	"}," +
	"\"tx_lut_12\": {" +
	"\"pa_gain\": 2," +
	"\"mix_gain\": 13," +
	"\"rf_power\": 23," +
	"\"dig_gain\": 0" +
	"}," +
	"\"tx_lut_13\": {" +
	"\"pa_gain\": 1," +
	"\"mix_gain\": 10," +
	"\"rf_power\": 25," +
	"\"dig_gain\": 1" +
	"}," +
	"\"tx_lut_14\": {" +
	"\"pa_gain\": 3," +
	"\"mix_gain\": 12," +
	"\"rf_power\": 26," +
	"\"dig_gain\": 2" +
	"}," +
	"\"tx_lut_15\": {" +
	"\"pa_gain\": 3," +
	"\"mix_gain\": 14," +
	"\"rf_power\": 27," +
	"\"dig_gain\": 0" +
	"}" +
	"}," +
	"\"gateway_conf\": {" +
	"\"server_address\": \"192.168.0.124\"," +
	"\"serv_port_up\": 1700," +
	"\"serv_port_down\": 1700," +
	"\"keepalive_interval\": 10," +
	"\"stat_interval\": 30," +
	"\"push_timeout_ms\": 100," +
	"\"forward_crc_valid\": true," +
	"\"forward_crc_error\": false," +
	"\"forward_crc_disabled\": false," +
	"\"gps_tty_path\": \"/dev/ttyS1\"," +
	"\"ref_latitude\": 0.0," +
	"\"ref_longitude\": 0.0," +
	"\"ref_altitude\": 0" +
	"}" +
	"}"

//var ConfGw = strings.ReplaceAll(GwConf, "\\", "")
//var ConfGw = []byte(GwConf)
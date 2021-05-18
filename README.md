# README for dyndnscd

dyndnscd is the dyndns client daemon. It is a daemon that continually polls for
IP address changes an in the event of a change, triggers an IP address update.
It is somewhat configurable.

## Downloading

You can find the latest version and the git repository of dyndnscd under the 
following URL: http://github.com/akrennmair/dyndnscd/

## Installation

Run `go get github.com/akrennmair/dyndnscd`, then run it like this: `dyndnscd -f configfile.yaml`

## Configuration

The configuration is a YAML file, consisting of a list of configuration 
sections.  Every section defines the IP polling mechanism with a `type` 
(allowed values: `device`, `ipbouncer`). The `device` type regularly polls the 
IPv4 address of a network device (specified by the configuration key `device`), 
while the `ipbouncer` regularly polls the IPv4 address by calling a bouncer 
URL. A bouncer URL returns the client's IP address as the only content of the 
response body, and is configured with the configuration key `bouncer_url`.

The URL update is configured with the configuraton key `update_url`. Simply 
write `<ip>` (no quotes) where the client IP shall be inserted. dyndnscd will 
replace it with the IP and will do a `GET` request on the resulting URL.

To configure the polling interval, use the configuration key `interval` (defines 
the minimum amount of time between two polling attempts). If no interval is specified,
it defaults to 60 seconds.

### Example 1

	- name: dyndns-bouncer
	  type: ipbouncer
	  bouncer_url: http://example.com/ip-bouncer
	  update_url: http://username:password@members.dyndns.org/nic/update?hostname=example.dyndns.org&myip=<ip>&wildcard=NOCHG&mx=NOCHG&backmx=NOCHG

### Example 2

	- name: dyndns-eth0
	  type: device
	  device: eth0
	  update_url: http://example.com/?myip=<ip>

## Contact

Andreas Krennmair <ak@synflood.at>

## License

dyndnscd is licensed under the MIT License. See the file LICENSE for further 
details.

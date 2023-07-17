import axios from 'axios';

import * as util from './util';

const rpcAddr = '/rpc/network';

const methods = {
  getNetworkConfig: { method: 'getNetworkConfig' },
  setNetworkConfig: { method: 'setNetworkConfig' },
  authorizeMember: { method: 'authorizeMember' },
  removeMember: { method: 'removeMember' },
  deleteWaiting: { method: 'deleteWaiting' },
  setAcceptAddress: { method: 'setAcceptAddress' },
  sendToken: { method: 'sendToken' },
  nknPing: { method: 'nknPing' },
}

var rpc = {};
for (let method in methods) {
  if (methods.hasOwnProperty(method)) {
    rpc[method] = (addr, params) => {
      params = util.assignDefined({}, methods[method].defaultParams, params)
      return rpcCall(addr, methods[method].method, params);
    }
  }
}

async function rpcCall(addr, method, params = {}) {
  let headers;
  try {
    headers = await window.rpcHeaders;
  } catch (e) {
    console.error('Await rpc headers error:', e);
  }

  let response = await axios({
    url: addr,
    method: 'POST',
    timeout: 10000,
    headers,
    // withCredentials: true,
    data: {
      id: 'nConnect-web',
      jsonrpc: '2.0',
      method: method,
      params: params,
    },
  });

  let data = response.data;

  if (data.error) {
    throw data.error;
  }

  if (data.result !== undefined) {
    return data.result;
  }

  throw new Error('rpc response contains no result or error field');
}

export async function getNetworkConfig() {
  return rpc.getNetworkConfig(rpcAddr);
}

export async function setNetworkConfig(networkConfig) {
  return rpc.setNetworkConfig(rpcAddr, {domain: networkConfig.domain, ipStart: networkConfig.ipStart, ipEnd: networkConfig.ipEnd, 
    netmask: networkConfig.netmask, gateway: networkConfig.gateway, dns: networkConfig.dns});
}

export async function authorizeMember(address) {
    return rpc.authorizeMember(rpcAddr, {address});
}

export async function removeMember(address) {
    return rpc.removeMember(rpcAddr, {address});
}

export async function deleteWaiting(address) {
    return rpc.deleteWaiting(rpcAddr, {address});
}

export async function setAcceptAddress(address, acceptAddresses) {
    return rpc.setAcceptAddress(rpcAddr, {address: address, AcceptAddresses: acceptAddresses});
}

export async function sendToken(address, amount) {
  return rpc.sendToken(rpcAddr, {address, amount});
}

export async function nknPing(address) {
  return rpc.nknPing(rpcAddr, {address});
}